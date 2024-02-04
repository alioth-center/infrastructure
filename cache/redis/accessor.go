package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alioth-center/infrastructure/cache"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

type accessor struct {
	db *redis.Client
	kb keyBuilder
}

func (ra *accessor) copySenderToReceiver(key string, senderPtr, receiverPtr any) error {
	rv := reflect.ValueOf(receiverPtr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		// 如果接收者不是指针，或者是空指针，返回错误
		return ra.kb.BuildError("copy sender to receiver", errors.New("receiver is not pointer or is nil"), key)
	}

	rv.Elem().Set(reflect.ValueOf(senderPtr))
	return nil
}

func (ra *accessor) DriverName() string {
	return DriverName
}

func (ra *accessor) Increase(ctx context.Context, key string, delta uint64) (result cache.CounterResultEnum) {
	resultNum, executeErr := ra.db.IncrBy(ctx, ra.kb.BuildKey(key), int64(delta)).Result()
	if executeErr != nil && !errors.Is(executeErr, redis.Nil) {
		return cache.CounterResultEnumFailed
	} else if errors.Is(executeErr, redis.Nil) {
		return cache.CounterResultEnumNotEffective
	}

	return cache.CounterResultEnum(resultNum)
}

func (ra *accessor) IncreaseWithExpireWhenNotExist(ctx context.Context, key string, delta uint64, expire time.Duration) (result cache.CounterResultEnum) {
	exist, existErr := ra.ExistKey(ctx, key)
	if existErr != nil {
		// 无法获取key是否存在
		return cache.CounterResultEnumFailed
	}
	if exist {
		// key已存在，直接增加，不需要设置过期时间
		return ra.Increase(ctx, ra.kb.BuildKey(key), delta)
	}

	resultNum, incrErr := ra.db.IncrBy(ctx, ra.kb.BuildKey(key), int64(delta)).Result()
	if incrErr != nil {
		// 增加失败
		return cache.CounterResultEnumFailed
	}
	if ra.db.Expire(ctx, ra.kb.BuildKey(key), expire).Err() != nil {
		// 设置过期时间失败
		return cache.CounterResultEnumFailed
	}

	// 增加成功
	return cache.CounterResultEnum(resultNum)
}

func (ra *accessor) SetExpire(ctx context.Context, key string, expire time.Duration) (result cache.CounterResultEnum) {
	exist, existErr := ra.ExistKey(ctx, key)
	if existErr != nil {
		// 无法获取key是否存在
		return cache.CounterResultEnumFailed
	}
	if !exist {
		// key不存在，设置过期时间无效
		return cache.CounterResultEnumNotEffective
	}
	if ra.db.Expire(ctx, ra.kb.BuildKey(key), expire).Err() != nil {
		// 设置过期时间失败
		return cache.CounterResultEnumFailed
	}

	// 设置成功
	return cache.CounterResultEnumSuccess
}

func (ra *accessor) SetExpireWhenNotSet(ctx context.Context, key string, expire time.Duration) (result cache.CounterResultEnum) {
	exist, expiredAt, existErr := ra.GetExpiredTime(ctx, key)
	if existErr != nil {
		// 无法获取key是否存在
		return cache.CounterResultEnumFailed
	}
	if !exist {
		// key不存在，设置过期时间无效
		return cache.CounterResultEnumNotEffective
	}
	if !expiredAt.IsZero() {
		// key已设置过期时间，设置过期时间无效
		return cache.CounterResultEnumNotEffective
	}
	if ra.db.Expire(ctx, ra.kb.BuildKey(key), expire).Err() != nil {
		// 设置过期时间失败
		return cache.CounterResultEnumFailed
	}

	// 设置成功
	return cache.CounterResultEnumSuccess
}

func (ra *accessor) ExpireImmediately(ctx context.Context, key string) (result cache.CounterResultEnum) {
	exist, existErr := ra.ExistKey(ctx, key)
	if existErr != nil {
		// 无法获取key是否存在
		return cache.CounterResultEnumFailed
	}
	if !exist {
		// key不存在，设置过期时间无效
		return cache.CounterResultEnumNotEffective
	}
	if ra.Delete(ctx, key) != nil {
		// 删除失败
		return cache.CounterResultEnumFailed
	}

	// 设置成功
	return cache.CounterResultEnumSuccess
}

func (ra *accessor) ExistKey(ctx context.Context, key string) (exist bool, err error) {
	count, executeRedisErr := ra.db.Exists(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil && !errors.Is(executeRedisErr, redis.Nil) {
		return false, ra.kb.BuildError("check key exist", err, key)
	}

	return count == 1, nil
}

func (ra *accessor) GetExpiredTime(ctx context.Context, key string) (exist bool, expiredAt time.Time, err error) {
	// 检查是否存在key
	existed, existErr := ra.ExistKey(ctx, key)
	if existErr != nil {
		return false, time.Time{}, existErr
	}
	if !existed {
		return false, time.Time{}, nil
	}

	// 获取过期时间
	result, executeRedisErr := ra.db.TTL(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		return false, time.Time{}, ra.kb.BuildError("get expired time", executeRedisErr, key)
	}
	if result == -1 {
		// 永不过期，返回空
		return true, time.Time{}, nil
	}
	if result == -2 {
		// 不存在key，返回空
		return false, time.Time{}, nil
	}

	// 返回过期时间
	return true, time.Now().Add(result), nil
}

func (ra *accessor) Load(ctx context.Context, key string) (exist bool, value string, err error) {
	result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, "", nil
		}

		return false, "", ra.kb.BuildError("load data", executeRedisErr, key)
	}

	return true, result, nil
}

func (ra *accessor) LoadWithEX(ctx context.Context, key string) (loaded bool, expiredTime time.Duration, value string, err error) {
	exist, result, executeRedisErr := ra.Load(ctx, key)
	if executeRedisErr != nil {
		return false, 0, "", executeRedisErr
	}
	if !exist {
		return false, 0, "", nil
	}

	hasTTL, ttl, getTTLErr := ra.GetExpiredTime(ctx, key)
	if getTTLErr != nil {
		return true, 0, "", getTTLErr
	}
	if !hasTTL {
		return true, 0, result, nil
	}
	if !ttl.IsZero() {
		return true, time.Until(ttl), result, nil
	}

	return true, 0, result, nil
}

func (ra *accessor) LoadJson(ctx context.Context, key string, receiverPtr any) (exist bool, err error) {
	existValue, value, loadValueErr := ra.Load(ctx, key)
	if loadValueErr != nil {
		return false, loadValueErr
	}
	if !existValue {
		return false, nil
	}

	unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr)
	if unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("load json data", unmarshalJsonErr, key)
	}

	return true, nil
}

func (ra *accessor) LoadJsonWithEX(ctx context.Context, key string, receiverPtr any) (exist bool, expiredTime time.Duration, err error) {
	existValue, expired, value, loadValueErr := ra.LoadWithEX(ctx, key)
	if loadValueErr != nil {
		return false, 0, loadValueErr
	}
	if !existValue {
		return false, 0, nil
	}

	unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr)
	if unmarshalJsonErr != nil {
		return true, 0, ra.kb.BuildError("load ex json data", unmarshalJsonErr, key)
	}

	return true, expired, nil
}

func (ra *accessor) Store(ctx context.Context, key string, value string) (err error) {
	setRedisErr := ra.db.Set(ctx, ra.kb.BuildKey(key), value, 0).Err()
	if setRedisErr != nil {
		if errors.Is(setRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("store data", setRedisErr, key)
	}

	return nil
}

func (ra *accessor) StoreEX(ctx context.Context, key string, value string, expiration time.Duration) (err error) {
	setRedisErr := ra.db.SetEX(ctx, ra.kb.BuildKey(key), value, expiration).Err()
	if setRedisErr != nil {
		if errors.Is(setRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("store ex data", setRedisErr, key)
	}

	return nil
}

func (ra *accessor) StoreJson(ctx context.Context, key string, senderPtr any) (err error) {
	payload, marshalJsonErr := json.Marshal(senderPtr)
	if marshalJsonErr != nil {
		return ra.kb.BuildError("marshal json data", marshalJsonErr, key)
	}

	storeRedisErr := ra.Store(ctx, key, string(payload))
	if storeRedisErr != nil {
		return storeRedisErr
	}

	return nil
}

func (ra *accessor) StoreJsonEX(ctx context.Context, key string, senderPtr any, expiration time.Duration) (err error) {
	payload, marshalJsonErr := json.Marshal(senderPtr)
	if marshalJsonErr != nil {
		return ra.kb.BuildError("marshal ex json data", marshalJsonErr, key)
	}

	storeRedisErr := ra.StoreEX(ctx, key, string(payload), expiration)
	if storeRedisErr != nil {
		return storeRedisErr
	}

	return nil
}

func (ra *accessor) Delete(ctx context.Context, key string) (err error) {
	deleteRedisErr := ra.db.Del(ctx, ra.kb.BuildKey(key)).Err()
	if deleteRedisErr != nil {
		if errors.Is(deleteRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("delete data", deleteRedisErr, key)
	}

	return nil
}

func (ra *accessor) LoadAndDelete(ctx context.Context, key string) (loaded bool, value string, err error) {
	result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, "", nil
		}

		return false, "", ra.kb.BuildError("load data", executeRedisErr, key)
	}

	if deleteRedisErr := ra.db.Del(ctx, ra.kb.BuildKey(key)).Err(); deleteRedisErr != nil {
		if errors.Is(deleteRedisErr, redis.Nil) {
			return true, result, nil
		}

		return true, result, ra.kb.BuildError("delete data", deleteRedisErr, key)
	}

	return true, result, nil
}

func (ra *accessor) LoadAndDeleteJson(ctx context.Context, key string, receivePtr any) (loaded bool, err error) {
	exist, value, loadValueErr := ra.LoadAndDelete(ctx, key)
	if loadValueErr != nil {
		return false, loadValueErr
	}
	if !exist {
		return false, nil
	}
	if unmarshalJsonErr := json.Unmarshal([]byte(value), receivePtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("unmarshal json data", unmarshalJsonErr, key)
	}

	return true, nil
}

func (ra *accessor) LoadOrStore(ctx context.Context, key string, storeValue string) (loaded bool, value string, err error) {
	result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if !errors.Is(executeRedisErr, redis.Nil) {
			return false, "", ra.kb.BuildError("load data", executeRedisErr, key)
		}

		// 如果不存在，存储
		if setRedisErr := ra.db.Set(ctx, ra.kb.BuildKey(key), storeValue, 0).Err(); setRedisErr != nil {
			if errors.Is(setRedisErr, redis.Nil) {
				return false, "", nil
			}

			return false, "", ra.kb.BuildError("store data", setRedisErr, key)
		}

		return false, storeValue, nil
	}

	// 如果存在，返回
	return true, result, nil
}

func (ra *accessor) LoadOrStoreEX(ctx context.Context, key string, storeValue string, expiration time.Duration) (loaded bool, value string, err error) {
	result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if !errors.Is(executeRedisErr, redis.Nil) {
			return false, "", ra.kb.BuildError("load ex data", executeRedisErr, key)
		}

		// 如果不存在，存储
		if setRedisErr := ra.db.SetEX(ctx, ra.kb.BuildKey(key), storeValue, expiration).Err(); setRedisErr != nil {
			if errors.Is(setRedisErr, redis.Nil) {
				return false, "", nil
			}

			return false, "", ra.kb.BuildError("store ex data", setRedisErr, key)
		}

		return false, storeValue, nil
	}

	// 如果存在，返回
	return true, result, nil
}

func (ra *accessor) LoadOrStoreJson(ctx context.Context, key string, senderPtr any, receiverPtr any) (loaded bool, err error) {
	payload, marshalJsonErr := json.Marshal(senderPtr)
	if marshalJsonErr != nil {
		return false, ra.kb.BuildError("marshal json data", marshalJsonErr, key)
	}

	exist, value, loadValueErr := ra.LoadOrStore(ctx, key, string(payload))
	if loadValueErr != nil {
		return false, loadValueErr
	}
	if !exist && value == "" {
		return false, nil
	}

	if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("unmarshal json data", unmarshalJsonErr, key)
	}

	return exist, nil
}

func (ra *accessor) LoadOrStoreJsonEX(ctx context.Context, key string, senderPtr any, receiverPtr any, expiration time.Duration) (loaded bool, err error) {
	payload, marshalJsonErr := json.Marshal(senderPtr)
	if marshalJsonErr != nil {
		return false, ra.kb.BuildError("marshal ex json data", marshalJsonErr, key)
	}

	exist, value, loadValueErr := ra.LoadOrStoreEX(ctx, key, string(payload), expiration)
	if loadValueErr != nil {
		return false, loadValueErr
	}
	if !exist && value == "" {
		return false, nil
	}

	if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("marshal ex json data", marshalJsonErr, key)
	}

	return exist, nil
}

func (ra *accessor) IsMember(ctx context.Context, key string, member string) (isMember bool, err error) {
	result, executeRedisErr := ra.db.SIsMember(ctx, ra.kb.BuildKey(key), member).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, nil
		}

		return false, ra.kb.BuildError("check is member", executeRedisErr, key)
	}

	return result, nil
}

func (ra *accessor) IsMembers(ctx context.Context, key string, members ...string) (isMembers bool, err error) {
	if len(members) == 0 {
		return false, nil
	}

	membersInterfaces := make([]interface{}, len(members))
	for i, member := range members {
		membersInterfaces[i] = member
	}

	result, executeRedisErr := ra.db.SMIsMember(ctx, ra.kb.BuildKey(key), membersInterfaces...).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, nil
		}

		return false, ra.kb.BuildError("check is member", executeRedisErr, key)
	}

	// 如果有一个不是成员，返回false
	for _, b := range result {
		if !b {
			return false, nil
		}
	}

	return true, nil
}

func (ra *accessor) AddMember(ctx context.Context, key string, member string) (err error) {
	addRedisErr := ra.db.SAdd(ctx, ra.kb.BuildKey(key), member).Err()
	if addRedisErr != nil {
		if errors.Is(addRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("add member", addRedisErr, key)
	}

	return nil
}

func (ra *accessor) AddMembers(ctx context.Context, key string, members ...string) (err error) {
	if len(members) == 0 {
		return nil
	}

	membersInterfaces := make([]interface{}, len(members))
	for i, member := range members {
		membersInterfaces[i] = member
	}

	addRedisErr := ra.db.SAdd(ctx, ra.kb.BuildKey(key), membersInterfaces...).Err()
	if addRedisErr != nil {
		if errors.Is(addRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("add members", addRedisErr, key)
	}

	return nil
}

func (ra *accessor) RemoveMember(ctx context.Context, key string, member string) (err error) {
	removeRedisErr := ra.db.SRem(ctx, ra.kb.BuildKey(key), member).Err()
	if removeRedisErr != nil {
		if errors.Is(removeRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("remove member", removeRedisErr, key)
	}

	return nil
}

func (ra *accessor) GetMembers(ctx context.Context, key string) (members []string, err error) {
	result, executeRedisErr := ra.db.SMembers(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return []string{}, nil
		}

		return []string{}, ra.kb.BuildError("get members", executeRedisErr, key)
	}

	return result, nil
}

func (ra *accessor) GetRandomMember(ctx context.Context, key string) (member string, err error) {
	result, executeRedisErr := ra.db.SRandMember(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return "", nil
		}

		return "", ra.kb.BuildError("get random member", executeRedisErr, key)
	}

	return result, nil
}

func (ra *accessor) GetRandomMembers(ctx context.Context, key string, count int64) (members []string, err error) {
	result, executeRedisErr := ra.db.SRandMemberN(ctx, ra.kb.BuildKey(key), count).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return []string{}, nil
		}

		return []string{}, ra.kb.BuildError("get random members", executeRedisErr, key)
	}

	return result, nil
}

func (ra *accessor) HGetValue(ctx context.Context, key string, field string) (exist bool, value string, err error) {
	result, executeRedisErr := ra.db.HGet(ctx, ra.kb.BuildKey(key), field).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, "", nil
		}

		return false, "", ra.kb.BuildError("get hash value", executeRedisErr, key)
	}

	return true, result, nil
}

func (ra *accessor) HGetValues(ctx context.Context, key string, fields ...string) (resultMap map[string]string, err error) {
	if len(fields) == 0 {
		return map[string]string{}, nil
	}

	result, executeRedisErr := ra.db.HMGet(ctx, ra.kb.BuildKey(key), fields...).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return map[string]string{}, nil
		}

		return map[string]string{}, ra.kb.BuildError("get hash values", executeRedisErr, key)
	}

	resultMap = map[string]string{}
	for i, field := range fields {
		if result[i] == nil {
			resultMap[field] = ""
		} else {
			resultMap[field] = result[i].(string)
		}
	}

	return resultMap, nil
}

func (ra *accessor) HGetJson(ctx context.Context, key string, field string, receiverPtr any) (exist bool, err error) {
	existValue, value, loadValueErr := ra.HGetValue(ctx, key, field)
	if loadValueErr != nil {
		return false, loadValueErr
	}
	if !existValue {
		return false, nil
	}
	if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("unmarshal hash json data", unmarshalJsonErr, key)
	}

	return true, nil
}

func (ra *accessor) HGetAll(ctx context.Context, key string) (resultMap map[string]string, err error) {
	result, executeRedisErr := ra.db.HGetAll(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return map[string]string{}, nil
		}

		return map[string]string{}, ra.kb.BuildError("get all hash data", executeRedisErr, key)
	}

	return result, nil
}

func (ra *accessor) HGetAllJson(ctx context.Context, key string, receiverPtr any) (err error) {
	result, executeRedisErr := ra.db.HGetAll(ctx, ra.kb.BuildKey(key)).Result()
	if executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("get all hash json", executeRedisErr, key)
	}

	marshalBytes, marshalJsonErr := json.Marshal(result)
	if marshalJsonErr != nil {
		return ra.kb.BuildError("marshal hash json data", marshalJsonErr, key)
	}
	if unmarshalJsonErr := json.Unmarshal(marshalBytes, receiverPtr); unmarshalJsonErr != nil {
		return ra.kb.BuildError("unmarshal hash json data", unmarshalJsonErr, key)
	}

	return nil
}

func (ra *accessor) HSetValue(ctx context.Context, key string, field string, value string) (err error) {
	setRedisErr := ra.db.HSet(ctx, ra.kb.BuildKey(key), field, value).Err()
	if setRedisErr != nil {
		if errors.Is(setRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("set hash value", setRedisErr, key)
	}

	return nil
}

func (ra *accessor) HSetValues(ctx context.Context, key string, values map[string]string) (err error) {
	if len(values) == 0 {
		return nil
	}

	setRedisErr := ra.db.HSet(ctx, ra.kb.BuildKey(key), values).Err()
	if setRedisErr != nil {
		if errors.Is(setRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("set hash values", setRedisErr, key)
	}

	return nil
}

func (ra *accessor) HRemoveValue(ctx context.Context, key string, field string) (err error) {
	removeRedisErr := ra.db.HDel(ctx, ra.kb.BuildKey(key), field).Err()
	if removeRedisErr != nil {
		if errors.Is(removeRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("remove hash value", removeRedisErr, key)
	}

	return nil
}

func (ra *accessor) HRemoveValues(ctx context.Context, key string, fields ...string) (err error) {
	if len(fields) == 0 {
		return nil
	}

	removeRedisErr := ra.db.HDel(ctx, ra.kb.BuildKey(key), fields...).Err()
	if removeRedisErr != nil {
		if errors.Is(removeRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("remove hash values", removeRedisErr, key)
	}

	return nil
}

func (ra *accessor) Expire(ctx context.Context, key string, expire time.Duration) (err error) {
	expireRedisErr := ra.db.Expire(ctx, ra.kb.BuildKey(key), expire).Err()
	if expireRedisErr != nil {
		if errors.Is(expireRedisErr, redis.Nil) {
			return nil
		}

		return ra.kb.BuildError("expire key", expireRedisErr, key)
	}

	return nil
}
