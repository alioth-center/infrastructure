package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type accessor struct {
	db *redis.Client
	kb keyBuilder
}

func (ra *accessor) ExistKey(ctx context.Context, key string) (exist bool, err error) {
	if count, executeRedisErr := ra.db.Exists(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		return false, ra.kb.BuildError("check key exist", err, key)
	} else {
		return count == 1, nil
	}
}

func (ra *accessor) GetExpiredTime(ctx context.Context, key string) (exist bool, expiredAt time.Time, err error) {
	// 检查是否存在key
	if exist, existErr := ra.ExistKey(ctx, key); existErr != nil {
		return false, time.Time{}, existErr
	} else if !exist {
		return false, time.Time{}, nil
	}

	// 获取过期时间
	if result, executeRedisErr := ra.db.TTL(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		return false, time.Time{}, ra.kb.BuildError("get expired time", executeRedisErr, key)
	} else if result == -1 {
		// 永不过期，返回空
		return true, time.Time{}, nil
	} else if result == -2 {
		// 不存在key，返回空
		return false, time.Time{}, nil
	} else {
		// 返回过期时间
		return true, time.Now().Add(result), nil
	}
}

func (ra *accessor) Load(ctx context.Context, key string) (exist bool, value string, err error) {
	if result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, "", nil
		} else {
			return false, "", ra.kb.BuildError("load data", executeRedisErr, key)
		}
	} else {
		return true, result, nil
	}
}

func (ra *accessor) LoadWithEX(ctx context.Context, key string) (loaded bool, expiredTime time.Duration, value string, err error) {
	if exist, result, executeRedisErr := ra.Load(ctx, key); executeRedisErr != nil {
		return false, 0, "", executeRedisErr
	} else if !exist {
		return false, 0, "", nil
	} else if hasTTL, ttl, getTTLErr := ra.GetExpiredTime(ctx, key); getTTLErr != nil {
		return true, 0, "", getTTLErr
	} else if !hasTTL {
		return true, 0, result, nil
	} else if !ttl.IsZero() {
		return true, time.Until(ttl), result, nil
	} else {
		return true, 0, result, nil
	}
}

func (ra *accessor) LoadJson(ctx context.Context, key string, receiverPtr any) (exist bool, err error) {
	if existValue, value, loadValueErr := ra.Load(ctx, key); loadValueErr != nil {
		return false, loadValueErr
	} else if !existValue {
		return false, nil
	} else if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("load json data", unmarshalJsonErr, key)
	} else {
		return true, nil
	}
}

func (ra *accessor) LoadJsonWithEX(ctx context.Context, key string, receiverPtr any) (exist bool, expiredTime time.Duration, err error) {
	if existValue, expiredTime, value, loadValueErr := ra.LoadWithEX(ctx, key); loadValueErr != nil {
		return false, 0, loadValueErr
	} else if !existValue {
		return false, 0, nil
	} else if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, 0, ra.kb.BuildError("load ex json data", unmarshalJsonErr, key)
	} else {
		return true, expiredTime, nil
	}
}

func (ra *accessor) Store(ctx context.Context, key string, value string) (err error) {
	if setRedisErr := ra.db.Set(ctx, ra.kb.BuildKey(key), value, 0).Err(); setRedisErr != nil {
		return ra.kb.BuildError("store data", setRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) StoreEX(ctx context.Context, key string, value string, expiration time.Duration) (err error) {
	if setRedisErr := ra.db.SetEX(ctx, ra.kb.BuildKey(key), value, expiration).Err(); setRedisErr != nil {
		return ra.kb.BuildError("store ex data", setRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) StoreJson(ctx context.Context, key string, senderPtr any) (err error) {
	if payload, marshalJsonErr := json.Marshal(senderPtr); marshalJsonErr != nil {
		return ra.kb.BuildError("marshal json data", marshalJsonErr, key)
	} else if storeRedisErr := ra.Store(ctx, key, string(payload)); storeRedisErr != nil {
		return storeRedisErr
	} else {
		return nil
	}
}

func (ra *accessor) StoreJsonEX(ctx context.Context, key string, senderPtr any, expiration time.Duration) (err error) {
	if payload, marshalJsonErr := json.Marshal(senderPtr); marshalJsonErr != nil {
		return ra.kb.BuildError("marshal ex json data", marshalJsonErr, key)
	} else if storeRedisErr := ra.StoreEX(ctx, key, string(payload), expiration); storeRedisErr != nil {
		return storeRedisErr
	} else {
		return nil
	}
}

func (ra *accessor) Delete(ctx context.Context, key string) (err error) {
	if deleteRedisErr := ra.db.Del(ctx, ra.kb.BuildKey(key)).Err(); deleteRedisErr != nil {
		return ra.kb.BuildError("delete data", deleteRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) LoadAndDelete(ctx context.Context, key string) (loaded bool, value string, err error) {
	if result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, "", nil
		} else {
			return false, "", ra.kb.BuildError("load data", executeRedisErr, key)
		}
	} else if deleteRedisErr := ra.db.Del(ctx, ra.kb.BuildKey(key)).Err(); deleteRedisErr != nil {
		return true, result, ra.kb.BuildError("delete data", deleteRedisErr, key)
	} else {
		return true, result, nil
	}
}

func (ra *accessor) LoadAndDeleteJson(ctx context.Context, key string, receivePtr any) (loaded bool, err error) {
	if loaded, value, loadValueErr := ra.LoadAndDelete(ctx, key); loadValueErr != nil {
		return false, loadValueErr
	} else if !loaded {
		return false, nil
	} else if unmarshalJsonErr := json.Unmarshal([]byte(value), receivePtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("unmarshal json data", unmarshalJsonErr, key)
	} else {
		return true, nil
	}
}

func (ra *accessor) LoadOrStore(ctx context.Context, key string, storeValue string) (loaded bool, value string, err error) {
	if result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		if !errors.Is(executeRedisErr, redis.Nil) {
			return false, "", ra.kb.BuildError("load data", executeRedisErr, key)
		} else {
			// 如果不存在，存储
			if setRedisErr := ra.db.Set(ctx, ra.kb.BuildKey(key), storeValue, 0).Err(); setRedisErr != nil {
				return false, "", ra.kb.BuildError("store data", setRedisErr, key)
			} else {
				return false, storeValue, nil
			}
		}
	} else {
		// 如果存在，返回
		return true, result, nil
	}
}

func (ra *accessor) LoadOrStoreEx(ctx context.Context, key string, storeValue string, expiration time.Duration) (loaded bool, value string, err error) {
	if result, executeRedisErr := ra.db.Get(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		if !errors.Is(executeRedisErr, redis.Nil) {
			return false, "", ra.kb.BuildError("load ex data", executeRedisErr, key)
		} else {
			// 如果不存在，存储
			if setRedisErr := ra.db.SetEX(ctx, ra.kb.BuildKey(key), storeValue, expiration).Err(); setRedisErr != nil {
				return false, "", ra.kb.BuildError("store ex data", setRedisErr, key)
			} else {
				return false, storeValue, nil
			}
		}
	} else {
		// 如果存在，返回
		return true, result, nil
	}
}

func (ra *accessor) LoadOrStoreJson(ctx context.Context, key string, senderPtr any, receiverPtr any) (loaded bool, err error) {
	if payload, marshalJsonErr := json.Marshal(senderPtr); marshalJsonErr != nil {
		return false, ra.kb.BuildError("marshal json data", marshalJsonErr, key)
	} else if loaded, value, loadValueErr := ra.LoadOrStore(ctx, key, string(payload)); loadValueErr != nil {
		return false, loadValueErr
	} else if !loaded {
		return false, nil
	} else if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("unmarshal json data", unmarshalJsonErr, key)
	} else {
		return true, nil
	}
}

func (ra *accessor) LoadOrStoreJsonEx(ctx context.Context, key string, senderPtr any, receiverPtr any, expiration time.Duration) (loaded bool, err error) {
	if payload, marshalJsonErr := json.Marshal(senderPtr); marshalJsonErr != nil {
		return false, ra.kb.BuildError("marshal ex json data", marshalJsonErr, key)
	} else if loaded, value, loadValueErr := ra.LoadOrStoreEx(ctx, key, string(payload), expiration); loadValueErr != nil {
		return false, loadValueErr
	} else if !loaded {
		return false, nil
	} else if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("marshal ex json data", marshalJsonErr, key)
	} else {
		return true, nil
	}
}

func (ra *accessor) IsMember(ctx context.Context, key string, member string) (isMember bool, err error) {
	if result, executeRedisErr := ra.db.SIsMember(ctx, ra.kb.BuildKey(key), member).Result(); executeRedisErr != nil {
		return false, ra.kb.BuildError("check is member", executeRedisErr, key)
	} else {
		return result, nil
	}
}

func (ra *accessor) IsMembers(ctx context.Context, key string, members ...string) (isMembers bool, err error) {
	if len(members) == 0 {
		return true, nil
	}

	membersInterfaces := make([]interface{}, len(members))
	for i, member := range members {
		membersInterfaces[i] = member
	}

	if result, executeRedisErr := ra.db.SMIsMember(ctx, ra.kb.BuildKey(key), membersInterfaces...).Result(); executeRedisErr != nil {
		return false, ra.kb.BuildError("check is member", executeRedisErr, key)
	} else {
		// 如果有一个不是成员，返回false
		for _, b := range result {
			if !b {
				return false, nil
			}
		}
		return true, nil
	}
}

func (ra *accessor) AddMember(ctx context.Context, key string, member string) (err error) {
	if addRedisErr := ra.db.SAdd(ctx, ra.kb.BuildKey(key), member).Err(); addRedisErr != nil {
		return ra.kb.BuildError("add member", addRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) AddMembers(ctx context.Context, key string, members ...string) (err error) {
	if len(members) == 0 {
		return nil
	}

	membersInterfaces := make([]interface{}, len(members))
	for i, member := range members {
		membersInterfaces[i] = member
	}

	if addRedisErr := ra.db.SAdd(ctx, ra.kb.BuildKey(key), membersInterfaces...).Err(); addRedisErr != nil {
		return ra.kb.BuildError("add members", addRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) RemoveMember(ctx context.Context, key string, member string) (err error) {
	if removeRedisErr := ra.db.SRem(ctx, ra.kb.BuildKey(key), member).Err(); removeRedisErr != nil {
		return ra.kb.BuildError("remove member", removeRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) GetMembers(ctx context.Context, key string) (members []string, err error) {
	if result, executeRedisErr := ra.db.SMembers(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		return nil, ra.kb.BuildError("get members", executeRedisErr, key)
	} else {
		return result, nil
	}
}

func (ra *accessor) GetRandomMember(ctx context.Context, key string) (member string, err error) {
	if result, executeRedisErr := ra.db.SRandMember(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		return "", ra.kb.BuildError("get random member", executeRedisErr, key)
	} else {
		return result, nil
	}
}

func (ra *accessor) GetRandomMembers(ctx context.Context, key string, count int64) (members []string, err error) {
	if result, executeRedisErr := ra.db.SRandMemberN(ctx, ra.kb.BuildKey(key), count).Result(); executeRedisErr != nil {
		return nil, ra.kb.BuildError("get random members", executeRedisErr, key)
	} else {
		return result, nil
	}
}

func (ra *accessor) HGetValue(ctx context.Context, key string, field string) (exist bool, value string, err error) {
	if result, executeRedisErr := ra.db.HGet(ctx, ra.kb.BuildKey(key), field).Result(); executeRedisErr != nil {
		if errors.Is(executeRedisErr, redis.Nil) {
			return false, "", nil
		} else {
			return false, "", ra.kb.BuildError("get hash value", executeRedisErr, key)
		}
	} else {
		return true, result, nil
	}
}

func (ra *accessor) HGetValues(ctx context.Context, key string, fields ...string) (resultMap map[string]string, err error) {
	if len(fields) == 0 {
		return map[string]string{}, nil
	}

	if result, executeRedisErr := ra.db.HMGet(ctx, ra.kb.BuildKey(key), fields...).Result(); executeRedisErr != nil {
		return nil, ra.kb.BuildError("get hash values", executeRedisErr, key)
	} else {
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
}

func (ra *accessor) HGetJson(ctx context.Context, key string, field string, receiverPtr any) (exist bool, err error) {
	if existValue, value, loadValueErr := ra.HGetValue(ctx, key, field); loadValueErr != nil {
		return false, loadValueErr
	} else if !existValue {
		return false, nil
	} else if unmarshalJsonErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalJsonErr != nil {
		return true, ra.kb.BuildError("unmarshal hash json data", unmarshalJsonErr, key)
	} else {
		return true, nil
	}
}

func (ra *accessor) HGetAll(ctx context.Context, key string) (resultMap map[string]string, err error) {
	if result, executeRedisErr := ra.db.HGetAll(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		return map[string]string{}, ra.kb.BuildError("get all hash data", executeRedisErr, key)
	} else {
		return result, nil
	}
}

func (ra *accessor) HGetAllJson(ctx context.Context, key string, receiverPtr any) (err error) {
	if result, executeRedisErr := ra.db.HGetAll(ctx, ra.kb.BuildKey(key)).Result(); executeRedisErr != nil {
		return ra.kb.BuildError("get all hash json", executeRedisErr, key)
	} else if marshalBytes, marshalJsonErr := json.Marshal(result); marshalJsonErr != nil {
		return ra.kb.BuildError("marshal hash json data", marshalJsonErr, key)
	} else if unmarshalJsonErr := json.Unmarshal(marshalBytes, receiverPtr); unmarshalJsonErr != nil {
		return ra.kb.BuildError("unmarshal hash json data", unmarshalJsonErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) HSetValue(ctx context.Context, key string, field string, value string) (err error) {
	if setRedisErr := ra.db.HSet(ctx, ra.kb.BuildKey(key), field, value).Err(); setRedisErr != nil {
		return ra.kb.BuildError("set hash value", setRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) HSetValues(ctx context.Context, key string, values map[string]string) (err error) {
	if len(values) == 0 {
		return nil
	}

	if setRedisErr := ra.db.HSet(ctx, ra.kb.BuildKey(key), values).Err(); setRedisErr != nil {
		return ra.kb.BuildError("set hash values", setRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) HRemoveValue(ctx context.Context, key string, field string) (err error) {
	if removeRedisErr := ra.db.HDel(ctx, ra.kb.BuildKey(key), field).Err(); removeRedisErr != nil {
		return ra.kb.BuildError("remove hash value", removeRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) HRemoveValues(ctx context.Context, key string, fields ...string) (err error) {
	if len(fields) == 0 {
		return nil
	}

	if removeRedisErr := ra.db.HDel(ctx, ra.kb.BuildKey(key), fields...).Err(); removeRedisErr != nil {
		return ra.kb.BuildError("remove hash values", removeRedisErr, key)
	} else {
		return nil
	}
}

func (ra *accessor) Expire(ctx context.Context, key string, expire time.Duration) (err error) {
	if expireRedisErr := ra.db.Expire(ctx, ra.kb.BuildKey(key), expire).Err(); expireRedisErr != nil {
		return ra.kb.BuildError("expire key", expireRedisErr, key)
	} else {
		return nil
	}
}
