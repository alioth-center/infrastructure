package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type accessor struct {
	mtx sync.RWMutex
	db  map[string]entry
	ec  chan struct{}
}

func (ca *accessor) Add(_ context.Context, key string, delta int64) (_ error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		entry = newCounterEntry(delta)
		ca.mtx.Lock()
		ca.db[key] = entry
		ca.mtx.Unlock()
	} else {
		entry.(*counterEntry).Add(delta)
	}

	if entry.IsExpired() {
		_ = ca.Delete(nil, key)
		return nil
	}

	return nil
}

func (ca *accessor) Sub(_ context.Context, key string, delta int64) (_ error) {
	return ca.Add(nil, key, -delta)
}

func (ca *accessor) Get(_ context.Context, key string) (value int64, _ error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		return 0, nil
	}

	if entry.IsExpired() {
		_ = ca.Delete(nil, key)
		return 0, nil
	}

	if entry.Type() != Int {
		return 0, NewValueTypeNotMatchError(Int, entry.Type())
	}

	return entry.(*counterEntry).Value(), nil
}

func (ca *accessor) ExistKey(_ context.Context, key string) (exist bool, err error) {
	ca.mtx.RLock()
	defer ca.mtx.RUnlock()
	_, exist = ca.db[key]
	return exist, nil
}

func (ca *accessor) GetExpiredTime(_ context.Context, key string) (exist bool, expiredAt time.Time, err error) {
	ca.mtx.RLock()
	defer ca.mtx.RUnlock()
	entry, isExist := ca.db[key]
	if !isExist {
		return false, time.Time{}, nil
	} else {
		return true, entry.GetExpiredAt(), nil
	}
}

func (ca *accessor) Load(_ context.Context, key string) (exist bool, value string, err error) {
	exist, _, value, _ = ca.LoadWithEX(nil, key)
	return exist, value, nil
}

func (ca *accessor) LoadWithEX(_ context.Context, key string) (loaded bool, expiredTime time.Duration, value string, err error) {
	ca.mtx.RLock()
	entry, isExist := ca.db[key]
	ca.mtx.RUnlock()
	if !isExist {
		// 如果不存在，直接返回
		return false, 0, "", nil
	}

	if entry.IsExpired() {
		// 如果已经过期，删除后返回
		_ = ca.Delete(nil, key)
		return false, 0, "", nil
	}

	if entry.Type() != String {
		// 如果类型不匹配，返回错误
		return false, 0, "", NewValueTypeNotMatchError(String, entry.Type())
	}

	// 如果存在，且没有过期，返回值
	return true, entry.GetExpiredAt().Sub(time.Now()), entry.(*stringEntry).Value(), nil
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
	} else {
		return true, exTime, nil
	}
}

func (ca *accessor) Store(_ context.Context, key string, value string) (err error) {
	ca.mtx.Lock()
	ca.db[key] = newStringEntry(value)
	ca.mtx.Unlock()
	return nil
}

func (ca *accessor) StoreEX(_ context.Context, key string, value string, expiration time.Duration) (err error) {
	ca.mtx.Lock()
	entry := newStringEntry(value)
	entry.SetExpireTime(expiration)
	ca.db[key] = entry
	ca.mtx.Unlock()
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
	ca.mtx.Lock()
	delete(ca.db, key)
	ca.mtx.Unlock()
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
	} else {
		return false, storeValue, ca.Store(nil, key, storeValue)
	}
}

func (ca *accessor) LoadOrStoreEx(_ context.Context, key string, storeValue string, expiration time.Duration) (loaded bool, value string, err error) {
	if exist, val, _ := ca.Load(nil, key); exist {
		return true, val, nil
	} else {
		return false, storeValue, ca.StoreEX(nil, key, storeValue, expiration)
	}
}

func (ca *accessor) LoadOrStoreJson(_ context.Context, key string, senderPtr any, receiverPtr any) (loaded bool, err error) {
	if exist, loadErr := ca.LoadJson(nil, key, receiverPtr); exist {
		return true, loadErr
	} else {
		receiverPtr = senderPtr
		return false, ca.StoreJson(nil, key, senderPtr)
	}
}

func (ca *accessor) LoadOrStoreJsonEx(_ context.Context, key string, senderPtr any, receiverPtr any, expiration time.Duration) (loaded bool, err error) {
	if exist, loadErr := ca.LoadJson(nil, key, receiverPtr); exist {
		return true, loadErr
	} else {
		receiverPtr = senderPtr
		return false, ca.StoreJsonEX(nil, key, senderPtr, expiration)
	}
}

func (ca *accessor) IsMember(_ context.Context, key string, member string) (isMember bool, err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		// 如果不存在，直接返回
		return false, nil
	}

	if entry.IsExpired() {
		// 如果已经过期，删除后返回
		_ = ca.Delete(nil, key)
		return false, nil
	}

	if entry.Type() != Set {
		// 类型不是set，返回错误
		return false, NewValueTypeNotMatchError(Set, entry.Type())
	}

	return entry.(*setEntry).IsMember(member), nil
}

func (ca *accessor) IsMembers(_ context.Context, key string, members ...string) (isMembers bool, err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		// 如果不存在，直接返回
		return false, nil
	}

	if entry.IsExpired() {
		// 如果已经过期，删除后返回
		_ = ca.Delete(nil, key)
		return false, nil
	}

	if entry.Type() != Set {
		// 类型不是set，返回错误
		return false, NewValueTypeNotMatchError(Set, entry.Type())
	}

	for _, member := range members {
		if !entry.(*setEntry).IsMember(member) {
			return false, nil
		}
	}
	return true, nil
}

func (ca *accessor) AddMember(_ context.Context, key string, member string) (err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		ca.mtx.Lock()
		entry = newSetEntry()
		ca.db[key] = entry
		ca.mtx.Unlock()
	}

	entry.(*setEntry).AddMember(member)
	return nil
}

func (ca *accessor) AddMembers(_ context.Context, key string, members ...string) (err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		entry = newSetEntry()
		ca.mtx.Lock()
		ca.db[key] = entry
		ca.mtx.Unlock()
	}

	entry.(*setEntry).AddMembers(members...)
	return nil
}

func (ca *accessor) RemoveMember(_ context.Context, key string, member string) (err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		return nil
	}

	entry.(*setEntry).RemoveMember(member)
	return nil
}

func (ca *accessor) GetMembers(_ context.Context, key string) (members []string, err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		return []string{}, nil
	}

	return entry.(*setEntry).Members(), nil
}

func (ca *accessor) GetRandomMember(_ context.Context, key string) (member string, err error) {
	members, _ := ca.GetMembers(nil, key)
	if len(members) > 0 {
		n := rand.Intn(len(members))
		return members[n], nil
	} else {
		return "", nil
	}
}

func (ca *accessor) GetRandomMembers(_ context.Context, key string, count int64) (members []string, err error) {
	existMembers, _ := ca.GetMembers(nil, key)
	length := int64(len(existMembers))
	if length <= count {
		// 如果得到的元素小于要求个数，全部返回
		return existMembers, nil
	}

	if length > count*10 {
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
	} else {
		// 获取大量元素采取乱序抽取的方式
		rand.Shuffle(int(length), func(i, j int) { existMembers[i], existMembers[j] = existMembers[j], existMembers[i] })
		return existMembers[:count], nil
	}
}

func (ca *accessor) HGetValue(_ context.Context, key string, field string) (exist bool, value string, err error) {
	ca.mtx.RLock()
	entry, isExist := ca.db[key]
	ca.mtx.RUnlock()
	if !isExist {
		// 不存在hash的key，直接返回
		return false, "", nil
	}

	if entry.IsExpired() {
		// 如果已经过期，删除后返回
		_ = ca.Delete(nil, key)
		return false, "", nil
	}

	if entry.Type() != Hash {
		// 如果entry的类型不正确，返回
		return false, "", NewValueTypeNotMatchError(Hash, entry.Type())
	}

	value, exist = entry.(*hashEntry).GetField(field)
	return exist, value, nil
}

func (ca *accessor) HGetValues(_ context.Context, key string, fields ...string) (resultMap map[string]string, err error) {
	ca.mtx.RLock()
	entry, isExist := ca.db[key]
	ca.mtx.RUnlock()
	if !isExist {
		// 不存在hash的key，直接返回
		return map[string]string{}, nil
	}

	if entry.IsExpired() {
		// 如果已经过期，删除后返回
		_ = ca.Delete(nil, key)
		return map[string]string{}, nil
	}

	if entry.Type() != Hash {
		// 如果entry的类型不正确，返回
		return map[string]string{}, NewValueTypeNotMatchError(Hash, entry.Type())
	}

	return entry.(*hashEntry).GetFields(fields...), nil
}

func (ca *accessor) HGetJson(_ context.Context, key string, field string, receiverPtr any) (exist bool, err error) {
	isExist, value, _ := ca.HGetValue(nil, key, field)
	if !isExist {
		return false, nil
	}

	if unmarshalErr := json.Unmarshal([]byte(value), receiverPtr); unmarshalErr != nil {
		return true, fmt.Errorf("unmarshal hash key %s field %s error: %w", key, field, unmarshalErr)
	}

	return true, nil
}

func (ca *accessor) HGetAll(_ context.Context, key string) (resultMap map[string]string, err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		// 如果不存在，直接返回
		return map[string]string{}, nil
	}

	if entry.IsExpired() {
		// 如果已经过期，删除后返回
		_ = ca.Delete(nil, key)
		return map[string]string{}, nil
	}

	if entry.Type() != Hash {
		// 如果entry的类型不正确，返回
		return map[string]string{}, NewValueTypeNotMatchError(Hash, entry.Type())
	}

	return entry.(*hashEntry).GetAllFields(), nil
}

func (ca *accessor) HGetAllJson(_ context.Context, key string, receiverPtr any) (err error) {
	result, _ := ca.HGetAll(nil, key)
	if buffer, marshalErr := json.Marshal(&result); marshalErr != nil {
		return fmt.Errorf("marshal hash key %s json data: %w", key, marshalErr)
	} else if unmarshalErr := json.Unmarshal(buffer, receiverPtr); unmarshalErr != nil {
		return fmt.Errorf("unmarshal hash key %s json data: %w", key, unmarshalErr)
	}

	return nil
}

func (ca *accessor) HSetValue(_ context.Context, key string, field string, value string) (err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		entry = newHashEntry()
		ca.mtx.Lock()
		ca.db[key] = entry
		ca.mtx.Unlock()
	}

	if entry.IsExpired() {
		_ = ca.Delete(nil, key)
		return nil
	}

	if entry.Type() != Hash {
		return NewValueTypeNotMatchError(Hash, entry.Type())
	}

	entry.(*hashEntry).AddField(field, value)
	return nil
}

func (ca *accessor) HSetValues(_ context.Context, key string, values map[string]string) (err error) {
	ca.mtx.RLock()
	entry, exist := ca.db[key]
	ca.mtx.RUnlock()
	if !exist {
		entry = newHashEntry()
		ca.mtx.Lock()
		ca.db[key] = entry
		ca.mtx.Unlock()
	}

	if entry.IsExpired() {
		_ = ca.Delete(nil, key)
		return nil
	}

	if entry.Type() != Hash {
		return NewValueTypeNotMatchError(Hash, entry.Type())
	}

	entry.(*hashEntry).AddFields(values)
	return nil
}

func (ca *accessor) HRemoveValue(_ context.Context, key string, field string) (err error) {
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

	entry.(*hashEntry).RemoveField(field)
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

	entry.SetExpireTime(expire)
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
