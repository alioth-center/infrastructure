package cache

import (
	"context"
	"time"
)

type Cache interface {
	ExistKey(ctx context.Context, key string) (exist bool, err error)
	GetExpiredTime(ctx context.Context, key string) (exist bool, expiredAt time.Time, err error)
	Load(ctx context.Context, key string) (exist bool, value string, err error)
	LoadWithEX(ctx context.Context, key string) (loaded bool, expiredTime time.Duration, value string, err error)
	LoadJson(ctx context.Context, key string, receiverPtr any) (exist bool, err error)
	LoadJsonWithEX(ctx context.Context, key string, receiverPtr any) (exist bool, expiredTime time.Duration, err error)
	Store(ctx context.Context, key string, value string) (err error)
	StoreEX(ctx context.Context, key string, value string, expiration time.Duration) (err error)
	StoreJson(ctx context.Context, key string, senderPtr any) (err error)
	StoreJsonEX(ctx context.Context, key string, senderPtr any, expiration time.Duration) (err error)
	Delete(ctx context.Context, key string) (err error)
	LoadAndDelete(ctx context.Context, key string) (loaded bool, value string, err error)
	LoadAndDeleteJson(ctx context.Context, key string, receivePtr any) (loaded bool, err error)
	LoadOrStore(ctx context.Context, key string, storeValue string) (loaded bool, value string, err error)
	LoadOrStoreEx(ctx context.Context, key string, storeValue string, expiration time.Duration) (loaded bool, value string, err error)
	LoadOrStoreJson(ctx context.Context, key string, senderPtr any, receiverPtr any) (loaded bool, err error)
	LoadOrStoreJsonEx(ctx context.Context, key string, senderPtr any, receiverPtr any, expiration time.Duration) (loaded bool, err error)
	IsMember(ctx context.Context, key string, member string) (isMember bool, err error)
	IsMembers(ctx context.Context, key string, members ...string) (isMembers bool, err error)
	AddMember(ctx context.Context, key string, member string) (err error)
	AddMembers(ctx context.Context, key string, members ...string) (err error)
	RemoveMember(ctx context.Context, key string, member string) (err error)
	GetMembers(ctx context.Context, key string) (members []string, err error)
	GetRandomMember(ctx context.Context, key string) (member string, err error)
	GetRandomMembers(ctx context.Context, key string, count int64) (members []string, err error)
	HGetValue(ctx context.Context, key string, field string) (exist bool, value string, err error)
	HGetValues(ctx context.Context, key string, fields ...string) (resultMap map[string]string, err error)
	HGetJson(ctx context.Context, key string, field string, receiverPtr any) (exist bool, err error)
	HGetAll(ctx context.Context, key string) (resultMap map[string]string, err error)
	HGetAllJson(ctx context.Context, key string, receiverPtr any) (err error)
	HSetValue(ctx context.Context, key string, field string, value string) (err error)
	HSetValues(ctx context.Context, key string, values map[string]string) (err error)
	HRemoveValue(ctx context.Context, key string, field string) (err error)
	HRemoveValues(ctx context.Context, key string, fields ...string) (err error)
	Expire(ctx context.Context, key string, expire time.Duration) (err error)
}
