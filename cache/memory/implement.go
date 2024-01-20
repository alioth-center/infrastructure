package memory

import (
	"sync"
	"time"
)

type Type string

const (
	Int    Type = "int"
	String Type = "string"
	Set    Type = "set"
	Hash   Type = "hash"
)

type trackable struct {
	createdAt   time.Time
	expiredTime time.Duration
}

func (e *trackable) GetExpiredAt() time.Time {
	if e.expiredTime < 0 || e.createdAt.IsZero() {
		return time.Time{}
	}

	return e.createdAt.Add(e.expiredTime)
}

func (e *trackable) IsExpired() bool {
	if e.expiredTime < 0 || e.createdAt.IsZero() {
		return false
	}

	return time.Now().After(e.createdAt.Add(e.expiredTime))
}

func (e *trackable) SetExpireTime(expiredTime time.Duration) {
	e.createdAt, e.expiredTime = time.Now(), expiredTime
}

func (e *trackable) GetExpireTime() time.Duration {
	if e.expiredTime < 0 || e.createdAt.IsZero() {
		// 如果没有设置过期时间，或者没有创建时间，则返回0
		return 0
	} else {
		// 返回剩余过期时间
		return time.Until(e.createdAt.Add(e.expiredTime))
	}
}

type entry interface {
	Type() Type
	GetExpiredAt() time.Time
	IsExpired() bool
	SetExpireTime(expiredTime time.Duration)
	GetExpireTime() time.Duration
}

type counterEntry struct {
	trackable
	mtx sync.Mutex
	val int64
}

func (e *counterEntry) Type() Type { return Int }

func (e *counterEntry) Value() int64 { e.mtx.Lock(); defer e.mtx.Unlock(); return e.val }

func (e *counterEntry) Add(delta int64) { e.mtx.Lock(); defer e.mtx.Unlock(); e.val += delta }

func (e *counterEntry) Sub(delta int64) { e.mtx.Lock(); defer e.mtx.Unlock(); e.val -= delta }

func (e *counterEntry) Set(val int64) { e.mtx.Lock(); defer e.mtx.Unlock(); e.val = val }

func newCounterEntry(val int64) entry { return &counterEntry{val: val} }

type stringEntry struct {
	trackable
	val string
}

func (e *stringEntry) Type() Type { return String }

func (e *stringEntry) Value() string { return e.val }

func newStringEntry(val string) entry { return &stringEntry{val: val} }

type setEntry struct {
	trackable
	mtx sync.RWMutex
	val map[string]struct{}
}

func (e *setEntry) Type() Type { return Set }

func (e *setEntry) AddMember(key string) {
	e.mtx.Lock()
	e.val[key] = struct{}{}
	e.mtx.Unlock()
}

func (e *setEntry) AddMembers(keys ...string) {
	e.mtx.Lock()
	for _, key := range keys {
		e.val[key] = struct{}{}
	}
	e.mtx.Unlock()
}

func (e *setEntry) RemoveMember(key string) {
	e.mtx.Lock()
	delete(e.val, key)
	e.mtx.Unlock()
}

func (e *setEntry) RemoveMembers(keys ...string) {
	e.mtx.Lock()
	for _, key := range keys {
		delete(e.val, key)
	}
	e.mtx.Unlock()
}

func (e *setEntry) IsMember(key string) bool {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	_, ok := e.val[key]
	return ok
}

func (e *setEntry) Members() (members []string) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	for key := range e.val {
		members = append(members, key)
	}
	return members
}

func newSetEntry() *setEntry { return &setEntry{mtx: sync.RWMutex{}, val: map[string]struct{}{}} }

type hashEntry struct {
	trackable
	mtx sync.RWMutex
	val map[string]string
}

func (e *hashEntry) Type() Type { return Hash }

func (e *hashEntry) AddField(field, value string) {
	e.mtx.Lock()
	e.val[field] = value
	e.mtx.Unlock()
}

func (e *hashEntry) AddFields(fields map[string]string) {
	e.mtx.Lock()
	for field, value := range fields {
		e.val[field] = value
	}
	e.mtx.Unlock()
}

func (e *hashEntry) RemoveField(field string) {
	e.mtx.Lock()
	delete(e.val, field)
	e.mtx.Unlock()
}

func (e *hashEntry) RemoveFields(fields ...string) {
	e.mtx.Lock()
	for _, field := range fields {
		delete(e.val, field)
	}
	e.mtx.Unlock()
}

func (e *hashEntry) GetField(field string) (value string, exist bool) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	value, exist = e.val[field]
	return value, exist
}

func (e *hashEntry) GetFields(fields ...string) (resultMap map[string]string) {
	if fields == nil || len(fields) == 0 {
		return map[string]string{}
	}

	e.mtx.RLock()
	defer e.mtx.RUnlock()
	resultMap = map[string]string{}
	for _, field := range fields {
		if value, exist := e.val[field]; exist {
			resultMap[field] = value
		}
	}

	return resultMap
}

func (e *hashEntry) GetAllFields() (resultMap map[string]string) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	resultMap = map[string]string{}
	for k, v := range e.val {
		resultMap[k] = v
	}
	return resultMap
}

func newHashEntry() *hashEntry { return &hashEntry{mtx: sync.RWMutex{}, val: map[string]string{}} }
