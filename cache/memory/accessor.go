package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/cache"
	"github.com/alioth-center/infrastructure/utils/values"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

func getEntryWithType[T any](mp *accessor, wantType Type, key string) (result T, exist bool, expireTime time.Duration) {
	rawEntry, isExist := mp.getEntry(key)
	if !isExist {
		// 如果不存在，返回空指针
		return values.Nil[T](), false, 0
	}
	if rawEntry.Type() != wantType {
		// 如果类型不匹配，返回空指针
		return values.Nil[T](), true, rawEntry.GetExpireTime()
	}

	entry, convertSuccess := rawEntry.(T)
	if !convertSuccess {
		// 如果类型转换失败，返回空指针
		return values.Nil[T](), true, rawEntry.GetExpireTime()
	}

	// 如果类型匹配，返回对应类型
	return entry, true, rawEntry.GetExpireTime()
}

type accessor struct {
	mtx sync.RWMutex
	db  map[string]entry
	ec  chan struct{}
}

func (ca *accessor) delete(key string) {
	ca.mtx.Lock()
	delete(ca.db, key)
	ca.mtx.Unlock()
}

func (ca *accessor) get(key string) (result entry, exist bool) {
	ca.mtx.RLock()
	rawEntry, isExist := ca.db[key]
	ca.mtx.RUnlock()
	return rawEntry, isExist
}

func (ca *accessor) create(key string, value entry) {
	ca.mtx.Lock()
	ca.db[key] = value
	ca.mtx.Unlock()
}

func (ca *accessor) update(key string, value entry) {
	ca.mtx.RLock()
	defer ca.mtx.RUnlock()
	if _, exist := ca.get(key); !exist {
		// 如果不存在，不应该更新
		return
	}

	ca.db[key] = value
}

func (ca *accessor) getEntry(key string) (result entry, exist bool) {
	result, exist = ca.get(key)
	if !exist {
		return nil, false
	}
	if result.IsExpired() {
		ca.delete(key)
		return nil, false
	}
	return result, true
}

func (ca *accessor) getCounterEntry(key string) (result *counterEntry, exist bool, expireTime time.Duration) {
	return getEntryWithType[*counterEntry](ca, Int, key)
}

func (ca *accessor) getStringEntry(key string) (result *stringEntry, exist bool, expireTime time.Duration) {
	return getEntryWithType[*stringEntry](ca, String, key)
}

func (ca *accessor) getSetEntry(key string) (result *setEntry, exist bool, expireTime time.Duration) {
	return getEntryWithType[*setEntry](ca, Set, key)
}

func (ca *accessor) getHashEntry(key string) (result *hashEntry, exist bool, expireTime time.Duration) {
	return getEntryWithType[*hashEntry](ca, Hash, key)
}

func (ca *accessor) copySenderToReceiver(senderPtr, receiverPtr any) error {
	rv := reflect.ValueOf(receiverPtr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		// 如果接收者不是指针，或者是空指针，返回错误
		return fmt.Errorf("failed to copy sender to receiver: %w", NewReceiverTypeIncorrectError(reflect.TypeOf(receiverPtr).String(), rv.IsNil()))
	}

	rv.Elem().Set(reflect.ValueOf(senderPtr))
	return nil
}

func (ca *accessor) Increase(_ context.Context, key string, delta uint64) (result cache.CounterResultEnum) {
	entry, exist, _ := ca.getCounterEntry(key)
	if !exist {
		// 如果不存在，创建一个新的计数器
		newEntry := newCounterEntry(int64(delta))
		ca.create(key, newEntry)
		return cache.CounterResultEnum(delta)
	}
	if entry == nil {
		// 计数器类型不匹配或转换失败，返回错误
		return cache.CounterResultEnumFailed
	}

	// 如果存在，增加计数器的值
	entry.Add(int64(delta))
	ca.update(key, entry)
	return cache.CounterResultEnum(entry.Value())
}

func (ca *accessor) IncreaseWithExpireWhenNotExist(_ context.Context, key string, delta uint64, expire time.Duration) (result cache.CounterResultEnum) {
	entry, exist, _ := ca.getCounterEntry(key)
	if !exist {
		// 如果不存在，创建一个新的计数器
		newEntry := newCounterEntry(int64(delta))
		newEntry.SetExpireTime(expire)
		ca.create(key, newEntry)
		return cache.CounterResultEnum(delta)
	}
	if entry == nil {
		// 计数器类型不匹配或转换失败，返回错误
		return cache.CounterResultEnumFailed
	}

	// 存在计数器，直接增加值，不需要修改过期时间
	entry.Add(int64(delta))
	ca.update(key, entry)
	return cache.CounterResultEnum(entry.Value())
}

func (ca *accessor) SetExpire(_ context.Context, key string, expire time.Duration) (result cache.CounterResultEnum) {
	entry, exist := ca.getEntry(key)
	if !exist {
		// 不存在计数器，返回修改无影响
		return cache.CounterResultEnumNotEffective
	}

	// 如果存在计数器，修改过期时间
	entry.SetExpireTime(expire)
	ca.update(key, entry)
	return cache.CounterResultEnumSuccess
}

func (ca *accessor) SetExpireWhenNotSet(_ context.Context, key string, expire time.Duration) (result cache.CounterResultEnum) {
	entry, exist := ca.getEntry(key)
	if !exist {
		// 不存在计数器，返回修改无影响
		return cache.CounterResultEnumNotEffective
	}
	if entry.GetExpireTime() == 0 {
		// 如果存在计数器，但是没有设置过期时间，修改过期时间
		ca.update(key, entry)
		return cache.CounterResultEnumSuccess
	}

	// 如果存在计数器，且已经设置过期时间，返回修改无影响
	return cache.CounterResultEnumNotEffective
}

func (ca *accessor) ExpireImmediately(_ context.Context, key string) (result cache.CounterResultEnum) {
	_, exist := ca.getEntry(key)
	if !exist {
		// 不存在计数器，返回修改无影响
		return cache.CounterResultEnumNotEffective
	}

	// 如果存在计数器，立即删除
	ca.delete(key)
	return cache.CounterResultEnumSuccess
}

func (ca *accessor) ExistKey(_ context.Context, key string) (exist bool, err error) {
	_, ext := ca.getEntry(key)
	if !ext {
		return false, nil
	}

	return true, nil
}

func (ca *accessor) GetExpiredTime(_ context.Context, key string) (exist bool, expiredAt time.Time, err error) {
	ca.mtx.RLock()
	entry, isExist := ca.db[key]
	ca.mtx.RUnlock()
	if !isExist {
		return false, time.Time{}, nil
	}
	if entry.IsExpired() {
		ca.mtx.Lock()
		delete(ca.db, key)
		ca.mtx.Unlock()
		return false, time.Time{}, nil
	}

	return true, entry.GetExpiredAt(), nil
}

func (ca *accessor) Load(_ context.Context, key string) (exist bool, value string, err error) {
	exist, _, value, err = ca.LoadWithEX(nil, key)
	return exist, value, err
}

func (ca *accessor) LoadWithEX(_ context.Context, key string) (loaded bool, expiredTime time.Duration, value string, err error) {
	resultEntry, exist, expireTime := ca.getStringEntry(key)
	if !exist {
		// 如果不存在，返回空值
		return false, 0, "", nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return true, expireTime, "", NewValueTypeNotMatchError(String, resultEntry.Type())
	}

	// 如果存在，且没有过期，返回值
	return true, resultEntry.GetExpireTime(), resultEntry.Value(), nil
}

func (ca *accessor) LoadJson(_ context.Context, key string, receiverPtr any) (exist bool, err error) {
	exist, _, err = ca.LoadJsonWithEX(nil, key, receiverPtr)
	return exist, err
}

func (ca *accessor) LoadJsonWithEX(_ context.Context, key string, receiverPtr any) (exist bool, expiredTime time.Duration, err error) {
	isExist, exTime, value, _ := ca.LoadWithEX(nil, key)
	if !isExist {
		return false, 0, nil
	}

	unmarshalErr := json.Unmarshal([]byte(value), receiverPtr)
	if unmarshalErr != nil {
		return true, exTime, fmt.Errorf("load json failed for key %s: %w", key, unmarshalErr)
	}

	return true, exTime, nil
}

func (ca *accessor) Store(_ context.Context, key string, value string) (err error) {
	resultEntry, exist, _ := ca.getStringEntry(key)
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(String, resultEntry.Type())
	}
	if exist {
		// 如果存在，更新值
		ca.update(key, newStringEntry(value))
		return nil
	}

	// 如果不存在，创建新的值
	ca.create(key, newStringEntry(value))
	return nil
}

func (ca *accessor) StoreEX(_ context.Context, key string, value string, expiration time.Duration) (err error) {
	resultEntry, exist, _ := ca.getStringEntry(key)
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(String, resultEntry.Type())
	}
	if exist {
		// 如果存在，更新值
		entry := newStringEntry(value)
		entry.SetExpireTime(expiration)
		ca.update(key, entry)
		return nil
	}

	// 如果不存在，创建新的值
	entry := newStringEntry(value)
	entry.SetExpireTime(expiration)
	ca.create(key, entry)
	return nil
}

func (ca *accessor) StoreJson(_ context.Context, key string, senderPtr any) (err error) {
	marshaled, marshalErr := json.Marshal(senderPtr)
	if marshalErr != nil {
		return fmt.Errorf("marshal json failed for key %s: %w", key, marshalErr)
	}

	return ca.Store(nil, key, string(marshaled))
}

func (ca *accessor) StoreJsonEX(_ context.Context, key string, senderPtr any, expiration time.Duration) (err error) {
	marshaled, marshalErr := json.Marshal(senderPtr)
	if marshalErr != nil {
		return fmt.Errorf("marshal json failed for key %s: %w", key, marshalErr)
	}

	return ca.StoreEX(nil, key, string(marshaled), expiration)
}

func (ca *accessor) Delete(_ context.Context, key string) (err error) {
	ca.delete(key)
	return nil
}

func (ca *accessor) LoadAndDelete(_ context.Context, key string) (loaded bool, value string, err error) {
	if exist, val, _ := ca.Load(nil, key); exist {
		return true, val, ca.Delete(nil, key)
	} else {
		return false, "", nil
	}
}

func (ca *accessor) LoadAndDeleteJson(_ context.Context, key string, receivePtr any) (loaded bool, err error) {
	if exist, loadErr := ca.LoadJson(nil, key, receivePtr); exist && loadErr == nil {
		return true, ca.Delete(nil, key)
	} else if exist && loadErr != nil {
		// 即使没有序列化成功，也需要删除键值对
		_ = ca.Delete(nil, key)
		return true, loadErr
	} else {
		return false, nil
	}
}

func (ca *accessor) LoadOrStore(_ context.Context, key string, storeValue string) (loaded bool, value string, err error) {
	if exist, val, _ := ca.Load(nil, key); exist {
		return true, val, nil
	}

	return false, storeValue, ca.Store(nil, key, storeValue)
}

func (ca *accessor) LoadOrStoreEX(_ context.Context, key string, storeValue string, expiration time.Duration) (loaded bool, value string, err error) {
	if exist, val, _ := ca.Load(nil, key); exist {
		return true, val, nil
	}

	return false, storeValue, ca.StoreEX(nil, key, storeValue, expiration)
}

func (ca *accessor) LoadOrStoreJson(_ context.Context, key string, senderPtr any, receiverPtr any) (loaded bool, err error) {
	if exist, loadErr := ca.LoadJson(nil, key, receiverPtr); exist {
		return true, loadErr
	}

	copyErr := ca.copySenderToReceiver(senderPtr, receiverPtr)
	if copyErr != nil {
		// 如果拷贝失败，返回错误
		return false, copyErr
	}

	return false, ca.StoreJson(nil, key, senderPtr)
}

func (ca *accessor) LoadOrStoreJsonEX(_ context.Context, key string, senderPtr any, receiverPtr any, expiration time.Duration) (loaded bool, err error) {
	if exist, loadErr := ca.LoadJson(nil, key, receiverPtr); exist {
		return true, loadErr
	}

	copyErr := ca.copySenderToReceiver(senderPtr, receiverPtr)
	if copyErr != nil {
		// 如果拷贝失败，返回错误
		return false, copyErr
	}

	return false, ca.StoreJsonEX(nil, key, senderPtr, expiration)
}

func (ca *accessor) IsMember(_ context.Context, key string, member string) (isMember bool, err error) {
	resultEntry, exist, _ := ca.getSetEntry(key)
	if !exist {
		// 如果不存在，直接返回
		return false, nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return false, NewValueTypeNotMatchError(Set, resultEntry.Type())
	}

	return resultEntry.IsMember(member), nil
}

func (ca *accessor) IsMembers(_ context.Context, key string, members ...string) (isMembers bool, err error) {
	if members == nil || len(members) == 0 {
		// 如果没有元素，直接返回
		return false, nil
	}

	resultEntry, exist, _ := ca.getSetEntry(key)
	if !exist {
		// 如果不存在，直接返回
		return false, nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return false, NewValueTypeNotMatchError(Set, resultEntry.Type())
	}

	for _, member := range members {
		if !resultEntry.IsMember(member) {
			return false, nil
		}
	}

	return true, nil
}

func (ca *accessor) AddMember(_ context.Context, key string, member string) (err error) {
	resultEntry, exist, _ := ca.getSetEntry(key)
	if !exist {
		// 如果不存在，创建后添加
		resultEntry = newSetEntry()
		resultEntry.AddMember(member)
		ca.create(key, resultEntry)
		return nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(Set, resultEntry.Type())
	}

	// 如果存在，添加元素
	resultEntry.AddMember(member)
	ca.update(key, resultEntry)
	return nil
}

func (ca *accessor) AddMembers(_ context.Context, key string, members ...string) (err error) {
	if members == nil || len(members) == 0 {
		// 如果没有元素，直接返回
		return nil
	}

	resultEntry, exist, _ := ca.getSetEntry(key)
	if !exist {
		// 如果不存在，创建后添加
		resultEntry = newSetEntry()
		resultEntry.AddMembers(members...)
		ca.create(key, resultEntry)
		return nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(Set, resultEntry.Type())
	}

	resultEntry.AddMembers(members...)
	ca.update(key, resultEntry)
	return nil
}

func (ca *accessor) RemoveMember(_ context.Context, key string, member string) (err error) {
	resultEntry, exist, _ := ca.getSetEntry(key)
	if !exist {
		// 如果不存在，直接返回
		return nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(Set, resultEntry.Type())
	}

	resultEntry.RemoveMember(member)
	ca.update(key, resultEntry)
	return nil
}

func (ca *accessor) GetMembers(_ context.Context, key string) (members []string, err error) {
	resultEntry, exist, _ := ca.getSetEntry(key)
	if !exist {
		// 如果不存在，直接返回
		return []string{}, nil
	}
	if resultEntry == nil && exist {
		// 如果类型不匹配，返回错误
		return []string{}, NewValueTypeNotMatchError(Set, resultEntry.Type())
	}

	return resultEntry.Members(), nil
}

func (ca *accessor) GetRandomMember(_ context.Context, key string) (member string, err error) {
	members, getMembersErr := ca.GetMembers(nil, key)
	if getMembersErr != nil {
		return "", getMembersErr
	}
	if len(members) > 0 {
		return members[rand.Intn(len(members))], nil
	}

	return "", nil
}

func (ca *accessor) GetRandomMembers(_ context.Context, key string, count int64) (members []string, err error) {
	if count <= 0 {
		// 如果需要的元素个数小于等于0，直接返回
		return []string{}, nil
	}

	existMembers, getMembersErr := ca.GetMembers(nil, key)
	if getMembersErr != nil {
		return []string{}, getMembersErr
	}
	length := int64(len(existMembers))
	if length <= count {
		// 如果得到的元素小于要求个数，全部返回
		return existMembers, nil
	}

	if length > count*10 || length <= 10 {
		// 如果元素个数远大于需要的个数，采取随机取数的方法
		randomSet := map[int]string{}
		for i := int64(0); i < count; {
			idx := rand.Intn(int(length))
			if _, exist := randomSet[idx]; exist {
				// 如果元素已经取过，再取一次
				continue
			}

			// 如果是全新的元素，添加进结果中
			randomSet[idx] = existMembers[idx]
			i++
		}

		// 将结果导出
		for _, s := range randomSet {
			members = append(members, s)
		}
		return members, nil
	}

	// 获取大量元素采取乱序抽取的方式
	rand.Shuffle(int(length), func(i, j int) { existMembers[i], existMembers[j] = existMembers[j], existMembers[i] })
	return existMembers[:count], nil
}

func (ca *accessor) HGetValue(_ context.Context, key string, field string) (exist bool, value string, err error) {
	resultEntry, isExist, _ := ca.getHashEntry(key)
	if !isExist {
		// 如果不存在hash的key，直接返回
		return false, "", nil
	}
	if resultEntry == nil && isExist {
		// 如果类型不匹配，返回错误
		return true, "", NewValueTypeNotMatchError(Hash, resultEntry.Type())
	}

	value, _ = resultEntry.GetField(field)
	return true, value, nil
}

func (ca *accessor) HGetValues(_ context.Context, key string, fields ...string) (resultMap map[string]string, err error) {
	resultEntry, isExist, _ := ca.getHashEntry(key)
	if !isExist {
		// 如果不存在hash的key，直接返回
		return map[string]string{}, nil
	}
	if resultEntry == nil && isExist {
		// 如果类型不匹配，返回错误
		return map[string]string{}, NewValueTypeNotMatchError(Hash, resultEntry.Type())
	}

	return resultEntry.GetFields(fields...), nil
}

func (ca *accessor) HGetJson(_ context.Context, key string, field string, receiverPtr any) (exist bool, err error) {
	isExist, value, getValueErr := ca.HGetValue(nil, key, field)
	if !isExist {
		return false, nil
	}
	if getValueErr != nil {
		return true, getValueErr
	}

	if unmarshalErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalErr != nil {
		return true, fmt.Errorf("unmarshal hash key %s field %s error: %w", key, field, unmarshalErr)
	}

	return true, nil
}

func (ca *accessor) HGetAll(_ context.Context, key string) (resultMap map[string]string, err error) {
	resultEntry, isExist, _ := ca.getHashEntry(key)
	if !isExist {
		// 如果不存在hash的key，直接返回
		return map[string]string{}, nil
	}
	if resultEntry == nil && isExist {
		// 如果类型不匹配，返回错误
		return map[string]string{}, NewValueTypeNotMatchError(Hash, resultEntry.Type())
	}

	return resultEntry.GetAllFields(), nil
}

func (ca *accessor) HGetAllJson(_ context.Context, key string, receiverPtr any) (err error) {
	result, getErr := ca.HGetAll(nil, key)
	if getErr != nil {
		return getErr
	}
	buffer, marshalErr := json.Marshal(&result)
	if marshalErr != nil {
		return fmt.Errorf("marshal hash key %s json data: %w", key, marshalErr)
	}
	if unmarshalErr := json.Unmarshal(buffer, receiverPtr); unmarshalErr != nil {
		return fmt.Errorf("unmarshal hash key %s json data: %w", key, unmarshalErr)
	}

	return nil
}

func (ca *accessor) HSetValue(_ context.Context, key string, field string, value string) (err error) {
	resultEntry, isExist, _ := ca.getHashEntry(key)
	if !isExist {
		// 如果不存在hash的key，创建后添加
		resultEntry = newHashEntry()
		resultEntry.AddField(field, value)
		ca.create(key, resultEntry)
		return nil
	}
	if resultEntry == nil && isExist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(Hash, resultEntry.Type())
	}

	// 如果存在，添加元素
	resultEntry.AddField(field, value)
	ca.update(key, resultEntry)
	return nil
}

func (ca *accessor) HSetValues(_ context.Context, key string, values map[string]string) (err error) {
	resultEntry, isExist, _ := ca.getHashEntry(key)
	if !isExist {
		// 如果不存在hash的key，创建后添加
		resultEntry = newHashEntry()
		resultEntry.AddFields(values)
		ca.create(key, resultEntry)
		return nil
	}
	if resultEntry == nil && isExist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(Hash, resultEntry.Type())
	}

	// 如果存在，添加元素
	resultEntry.AddFields(values)
	ca.update(key, resultEntry)
	return nil
}

func (ca *accessor) HRemoveValue(_ context.Context, key string, field string) (err error) {
	resultEntry, isExist, _ := ca.getHashEntry(key)
	if !isExist {
		// 如果不存在hash的key，直接返回
		return nil
	}
	if resultEntry == nil && isExist {
		// 如果类型不匹配，返回错误
		return NewValueTypeNotMatchError(Hash, resultEntry.Type())
	}

	// 如果存在，移除元素
	resultEntry.RemoveField(field)
	ca.update(key, resultEntry)
	return nil
}

func (ca *accessor) HRemoveValues(_ context.Context, key string, fields ...string) (err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		return nil
	}

	if entry.IsExpired() {
		_ = ca.Delete(nil, key)
		return nil
	}

	if entry.Type() != Hash {
		return NewValueTypeNotMatchError(Hash, entry.Type())
	}

	entry.(*hashEntry).RemoveFields(fields...)
	return nil
}

func (ca *accessor) Expire(_ context.Context, key string, expire time.Duration) (err error) {
	entry, exist := ca.getEntry(key)
	if !exist {
		return nil
	}

	entry.SetExpireTime(expire)
	ca.update(key, entry)
	return nil
}

func (ca *accessor) cleanCache(interval time.Duration, maxExecutionTime time.Duration, maxExecutionPercentage int) {
	exitChan, pauseChan, resumeChan := make(chan struct{}, 1), make(chan struct{}, 1), make(chan struct{}, 1)

	cleanFunction := func(exit, pause, resume chan struct{}) {
		for {
			select {
			case <-exit:
				return
			case <-pause:
				// 阻塞，等待恢复信号
				<-resume
			default:
				// 没有阻塞时进入此分支，执行完后会重新进入阻塞状态
				ca.mtx.RLock()
				maxExecution := len(ca.db) * maxExecutionPercentage / 100
				ca.mtx.RUnlock()
				deleteList, endTimer := make([]string, 0, maxExecution), time.Now().Add(maxExecutionTime)

				// 统计需要删除的key
				ca.mtx.RLock()
				for k, v := range ca.db {
					if v.IsExpired() {
						deleteList = append(deleteList, k)
					}

					if len(deleteList) > maxExecution {
						break
					}

					if len(deleteList)%100 == 0 && time.Now().After(endTimer) {
						break
					}
				}
				ca.mtx.RUnlock()

				// 执行删除任务
				ca.mtx.Lock()
				for _, k := range deleteList {
					delete(ca.db, k)
				}
				ca.mtx.Unlock()

				// 将自己阻塞
				pause <- struct{}{}
			}
		}
	}

	go cleanFunction(exitChan, pauseChan, resumeChan)

	for {
		select {
		case <-time.After(interval):
			resumeChan <- struct{}{}
		case <-ca.ec:
			exitChan <- struct{}{}
		}
	}
}

func (ca *accessor) close() { ca.ec <- struct{}{} }
