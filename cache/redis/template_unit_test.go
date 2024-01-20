package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/cache"
	"testing"
	"time"
)

var (
	interval      = time.Second * 2
	expire        = time.Second
	sleepInterval = time.Millisecond * 1100

	BaseCacheUnitTestCaseList = []TestCase{
		{
			CaseName:     "ExistKey",
			TestFunction: ExistKeyTestFunction,
		},
		{
			CaseName:     "GetExpiredTime",
			TestFunction: GetExpiredTimeFunction,
		},
		{
			CaseName:     "Load",
			TestFunction: LoadFunction,
		},
		{
			CaseName:     "LoadWithEX",
			TestFunction: LoadWithEXFunction,
		},
		{
			CaseName:     "LoadJson",
			TestFunction: LoadJsonFunction,
		},
		{
			CaseName:     "LoadJsonWithEX",
			TestFunction: LoadJsonWithEXFunction,
		},
		{
			CaseName:     "Store",
			TestFunction: StoreFunction,
		},
		{
			CaseName:     "StoreEX",
			TestFunction: StoreEXFunction,
		},
		{
			CaseName:     "StoreJson",
			TestFunction: StoreJsonFunction,
		},
		{
			CaseName:     "StoreJsonEX",
			TestFunction: StoreJsonEXFunction,
		},
		{
			CaseName:     "Delete",
			TestFunction: DeleteFunction,
		},
		{
			CaseName:     "LoadAndDelete",
			TestFunction: LoadAndDeleteFunction,
		},
		{
			CaseName:     "LoadAndDeleteJson",
			TestFunction: LoadAndDeleteJsonFunction,
		},
		{
			CaseName:     "LoadOrStore",
			TestFunction: LoadOrStoreFunction,
		},
		{
			CaseName:     "LoadOrStoreEX",
			TestFunction: LoadOrStoreEXFunction,
		},
		{
			CaseName:     "LoadOrStoreJson",
			TestFunction: LoadOrStoreJsonFunction,
		},
		{
			CaseName:     "LoadOrStoreJsonEX",
			TestFunction: LoadOrStoreJsonEXFunction,
		},
		{
			CaseName:     "IsMember",
			TestFunction: IsMemberFunction,
		},
		{
			CaseName:     "IsMembers",
			TestFunction: IsMembersFunction,
		},
		{
			CaseName:     "AddMember",
			TestFunction: AddMemberFunction,
		},
		{
			CaseName:     "AddMembers",
			TestFunction: AddMembersFunction,
		},
		{
			CaseName:     "RemoveMember",
			TestFunction: RemoveMemberFunction,
		},
		{
			CaseName:     "GetMembers",
			TestFunction: GetMembersFunction,
		},
		{
			CaseName:     "GetRandomMember",
			TestFunction: GetRandomMemberFunction,
		},
		{
			CaseName:     "GetRandomMembers",
			TestFunction: GetRandomMembersFunction,
		},
		{
			CaseName:     "HGetValue",
			TestFunction: HGetValueFunction,
		},
		{
			CaseName:     "HGetValues",
			TestFunction: HGetValuesFunction,
		},
		{
			CaseName:     "HGetJson",
			TestFunction: HGetJsonFunction,
		},
		{
			CaseName:     "HGetAll",
			TestFunction: HGetAllFunction,
		},
		{
			CaseName:     "HGetAllJson",
			TestFunction: HGetAllJsonFunction,
		},
		{
			CaseName:     "HSetValue",
			TestFunction: HSetValueFunction,
		},
		{
			CaseName:     "HSetValues",
			TestFunction: HSetValuesFunction,
		},
		{
			CaseName:     "HRemoveValue",
			TestFunction: HRemoveValueFunction,
		},
		{
			CaseName:     "HRemoveValues",
			TestFunction: HRemoveValuesFunction,
		},
	}
)

func timeAbs(d time.Duration) time.Duration {
	if d > 0 {
		return d
	}

	return -d
}

func expiredTimeIsCorrect(expiredTime time.Duration, expect time.Duration) bool {
	return timeAbs(expiredTime-expect) < interval
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

type TestCase struct {
	CaseName     string
	TestFunction func(impl cache.Cache) func(t *testing.T)
}

type testStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ExistKeyTestFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 存在key的样例
		t.Run("ExistKey:Exist", func(t *testing.T) {
			key := "ExistKey:Exist"
			storeErr := impl.Store(context.Background(), key, "ExistKey:Exist")
			if storeErr != nil {
				t.Errorf("ExistKey:Exist case failed when storing key: %v", storeErr.Error())
			}

			exist, existErr := impl.ExistKey(context.Background(), key)
			if existErr != nil {
				t.Errorf("ExistKey:Exist case failed when checking key: %v", existErr.Error())
			}
			if !exist {
				t.Errorf("ExistKey:Exist case failed: key not exist, want exist")
			}
		})

		// 不存在key的样例
		t.Run("ExistKey:NotExist", func(t *testing.T) {
			key := "ExistKey:NotExist"
			exist, existErr := impl.ExistKey(context.Background(), key)
			if existErr != nil {
				t.Errorf("ExistKey:NotExist case failed when checking key: %v", existErr.Error())
			}
			if exist {
				t.Errorf("ExistKey:NotExist case failed: key exist, want not exist")
			}
		})

		// 存在key，设置了过期时间，但是没过期的样例
		t.Run("ExistKey:NotExpired", func(t *testing.T) {
			key := "ExistKey:NotExpired"
			storeErr := impl.StoreEX(context.Background(), key, "ExistKey:NotExpired", time.Hour)
			if storeErr != nil {
				t.Errorf("ExistKey:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			exist, existErr := impl.ExistKey(context.Background(), key)
			if existErr != nil {
				t.Errorf("ExistKey:NotExpired case failed when checking key: %v", existErr.Error())
			}
			if !exist {
				t.Errorf("ExistKey:NotExpired case failed: key not exist, want exist")
			}
		})

		// 存在key，但是过期了的样例
		t.Run("ExistKey:Expired", func(t *testing.T) {
			key := "ExistKey:Expired"
			storeErr := impl.StoreEX(context.Background(), key, "ExistKey:Expired", expire)
			if storeErr != nil {
				t.Errorf("ExistKey:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, existErr := impl.ExistKey(context.Background(), key)
			if existErr != nil {
				t.Errorf("ExistKey:Expired case failed when checking key: %v", existErr.Error())
			}
			if exist {
				t.Errorf("ExistKey:Expired case failed: key exist, want not exist")
			}
		})
	}
}

func GetExpiredTimeFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 设置了过期时间并且没到期的样例
		t.Run("GetExpiredTime:NotExpired", func(t *testing.T) {
			key := "GetExpiredTime:NotExpired"
			storeErr := impl.StoreEX(context.Background(), key, "GetExpiredTime:NotExpired", time.Minute)
			if storeErr != nil {
				t.Errorf("GetExpiredTime:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			exist, expiredAt, getErr := impl.GetExpiredTime(context.Background(), key)
			if getErr != nil {
				t.Errorf("GetExpiredTime:NotExpired case failed when getting expried time: %v", getErr.Error())
			}
			if !exist {
				t.Errorf("GetExpiredTime:NotExpired case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(time.Until(expiredAt), time.Minute) {
				t.Errorf("GetExpiredTime:NotExpired case failed: incorrect expired time")
			}
		})

		// 没设置过期时间的样例
		t.Run("GetExpiredTime:NotSet", func(t *testing.T) {
			key := "GetExpiredTime:NotSet"
			storeErr := impl.Store(context.Background(), key, "GetExpiredTime:NotSet")
			if storeErr != nil {
				t.Errorf("GetExpiredTime:NotSet case failed when storing key: %v", storeErr.Error())
			}

			exist, expiredAt, getErr := impl.GetExpiredTime(context.Background(), key)
			if getErr != nil {
				t.Errorf("GetExpiredTime:NotSet case failed when getting expried time: %v", getErr.Error())
			}
			if !exist {
				t.Errorf("GetExpiredTime:NotSet case failed: key not exist, want exist")
			}
			if !expiredAt.IsZero() {
				t.Errorf("GetExpiredTime:NotSet case failed: incorrect expired time")
			}
		})

		// 设置了过期时间，但是已经过期的样例
		t.Run("GetExpiredTime:Expired", func(t *testing.T) {
			key := "GetExpiredTime:Expired"
			storeErr := impl.StoreEX(context.Background(), key, "GetExpiredTime:Expired", expire)
			if storeErr != nil {
				t.Errorf("GetExpiredTime:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, expiredAt, getErr := impl.GetExpiredTime(context.Background(), key)
			if getErr != nil {
				t.Errorf("GetExpiredTime:Expired case failed when getting expried time: %v", getErr.Error())
			}
			if exist || !expiredAt.IsZero() {
				t.Errorf("GetExpiredTime:Expired case failed: incorrect expired time or exist status")
			}
		})

		// 不存在的样例
		t.Run("GetExpiredTime:NotExist", func(t *testing.T) {
			key := "GetExpiredTime:NotExist"
			exist, expiredAt, getErr := impl.GetExpiredTime(context.Background(), key)
			if getErr != nil {
				t.Errorf("GetExpiredTime:NotExist case failed when getting expried time: %v", getErr.Error())
			}
			if exist || !expiredAt.IsZero() {
				t.Errorf("GetExpiredTime:NotExist case failed: incorrect expired time or exist status")
			}
		})
	}
}

func LoadFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 正常存在的样例
		t.Run("Load:Exist", func(t *testing.T) {
			key := "Load:Exist"
			storeErr := impl.Store(context.Background(), key, "Load:Exist")
			if storeErr != nil {
				t.Errorf("Load:Exist case failed when storing key: %v", storeErr.Error())
			}

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Load:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("Load:Exist case failed: key not exist, want exist")
			}
			if value != "Load:Exist" {
				t.Errorf("Load:Exist case failed: incorrect value")
			}
		})

		// 不存在的样例
		t.Run("Load:NotExist", func(t *testing.T) {
			key := "Load:NotExist"
			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Load:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("Load:NotExist case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("Load:NotExist case failed: incorrect value")
			}
		})

		// 存在但是过期的样例
		t.Run("Load:Expired", func(t *testing.T) {
			key := "Load:Expired"
			storeErr := impl.StoreEX(context.Background(), key, "Load:Expired", expire)
			if storeErr != nil {
				t.Errorf("Load:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Load:Expired case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("Load:Expired case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("Load:Expired case failed: incorrect value")
			}
		})

		// 设置了过期时间且没有过期的样例
		t.Run("Load:NotExpired", func(t *testing.T) {
			key := "Load:NotExpired"
			storeErr := impl.StoreEX(context.Background(), key, "Load:NotExpired", time.Minute)
			if storeErr != nil {
				t.Errorf("Load:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Load:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("Load:NotExpired case failed: key not exist, want exist")
			}
			if value != "Load:NotExpired" {
				t.Errorf("Load:NotExpired case failed: incorrect value")
			}
		})

		// 存储了一个非字符串的样例
		t.Run("Load:NotString", func(t *testing.T) {
			key, value := "Load:NotString", "NotString"
			storeErr := impl.AddMember(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("Load:NotString case failed when storing key: %v", storeErr.Error())
			}

			_, loadedValue, loadErr := impl.Load(context.Background(), key)
			if loadErr == nil {
				t.Errorf("Load:NotString case failed: want error but not")
			}
			if loadedValue != "" {
				t.Errorf("Load:NotString case failed: incorrect value")
			}
		})
	}
}

func LoadWithEXFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 正常存在的样例
		t.Run("LoadWithEX:Exist", func(t *testing.T) {
			key := "LoadWithEX:Exist"
			storeErr := impl.Store(context.Background(), key, "LoadWithEX:Exist")
			if storeErr != nil {
				t.Errorf("LoadWithEX:Exist case failed when storing key: %v", storeErr.Error())
			}

			loaded, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("LoadWithEX:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadWithEX:Exist case failed: key not exist, want exist")
			}
			if expiredTime != 0 {
				t.Errorf("LoadWithEX:Exist case failed: incorrect expired time")
			}
			if value != "LoadWithEX:Exist" {
				t.Errorf("LoadWithEX:Exist case failed: incorrect value")
			}
		})

		// 不存在的样例
		t.Run("LoadWithEX:NotExist", func(t *testing.T) {
			key := "LoadWithEX:NotExist"
			loaded, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("LoadWithEX:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadWithEX:NotExist case failed: key exist, want not exist")
			}
			if expiredTime != 0 {
				t.Errorf("LoadWithEX:NotExist case failed: incorrect expired time")
			}
			if value != "" {
				t.Errorf("LoadWithEX:NotExist case failed: incorrect value")
			}
		})

		// 存在但是过期的样例
		t.Run("LoadWithEX:Expired", func(t *testing.T) {
			key := "LoadWithEX:Expired"
			storeErr := impl.StoreEX(context.Background(), key, "LoadWithEX:Expired", expire)
			if storeErr != nil {
				t.Errorf("LoadWithEX:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			loaded, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("LoadWithEX:Expired case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadWithEX:Expired case failed: key exist, want not exist")
			}
			if expiredTime != 0 {
				t.Errorf("LoadWithEX:Expired case failed: incorrect expired time")
			}
			if value != "" {
				t.Errorf("LoadWithEX:Expired case failed: incorrect value")
			}
		})

		// 设置了过期时间且没有过期的样例
		t.Run("LoadWithEX:NotExpired", func(t *testing.T) {
			key := "LoadWithEX:NotExpired"
			storeErr := impl.StoreEX(context.Background(), key, "LoadWithEX:NotExpired", time.Minute)
			if storeErr != nil {
				t.Errorf("LoadWithEX:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			loaded, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("LoadWithEX:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadWithEX:NotExpired case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("LoadWithEX:NotExpired case failed: incorrect expired time")
			}
			if value != "LoadWithEX:NotExpired" {
				t.Errorf("LoadWithEX:NotExpired case failed: incorrect value")
			}
		})
	}
}

func LoadJsonFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 正常存在的样例
		t.Run("LoadJson:Exist", func(t *testing.T) {
			key := "LoadJson:Exist"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("LoadJson:Exist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJson:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("LoadJson:Exist case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadJson:Exist case failed: incorrect value")
			}
		})

		// 不存在的样例
		t.Run("LoadJson:NotExist", func(t *testing.T) {
			key := "LoadJson:NotExist"

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJson:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("LoadJson:NotExist case failed: key exist, want not exist")
			}
			if receiver.Key != "" || receiver.Value != "" {
				t.Errorf("LoadJson:NotExist case failed: incorrect value")
			}
		})

		// 存在但是过期的样例
		t.Run("LoadJson:Expired", func(t *testing.T) {
			key := "LoadJson:Expired"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJsonEX(context.Background(), key, value, expire)
			if storeErr != nil {
				t.Errorf("LoadJson:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJson:Expired case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("LoadJson:Expired case failed: key exist, want not exist")
			}
			if receiver.Key != "" || receiver.Value != "" {
				t.Errorf("LoadJson:Expired case failed: incorrect value")
			}
		})

		// 设置了过期时间且没有过期的样例
		t.Run("LoadJson:NotExpired", func(t *testing.T) {
			key := "LoadJson:NotExpired"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJsonEX(context.Background(), key, value, time.Minute)
			if storeErr != nil {
				t.Errorf("LoadJson:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJson:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("LoadJson:NotExpired case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadJson:NotExpired case failed: incorrect value")
			}
		})
	}
}

func LoadJsonWithEXFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 正常存在的样例
		t.Run("LoadJsonWithEX:Exist", func(t *testing.T) {
			key := "LoadJsonWithEX:Exist"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("LoadJsonWithEX:Exist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJsonWithEX:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("LoadJsonWithEX:Exist case failed: key not exist, want exist")
			}
			if expiredTime != 0 {
				t.Errorf("LoadJsonWithEX:Exist case failed: incorrect expired time")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadJsonWithEX:Exist case failed: incorrect value")
			}
		})

		// 不存在的样例
		t.Run("LoadJsonWithEX:NotExist", func(t *testing.T) {
			key := "LoadJsonWithEX:NotExist"

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJsonWithEX:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("LoadJsonWithEX:NotExist case failed: key exist, want not exist")
			}
			if expiredTime != 0 {
				t.Errorf("LoadJsonWithEX:NotExist case failed: incorrect expired time")
			}
			if receiver.Key != "" || receiver.Value != "" {
				t.Errorf("LoadJsonWithEX:NotExist case failed: incorrect value")
			}
		})

		// 存在但是过期的样例
		t.Run("LoadJsonWithEX:Expired", func(t *testing.T) {
			key := "LoadJsonWithEX:Expired"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJsonEX(context.Background(), key, value, expire)
			if storeErr != nil {
				t.Errorf("LoadJsonWithEX:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJsonWithEX:Expired case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("LoadJsonWithEX:Expired case failed: key exist, want not exist")
			}
			if expiredTime != 0 {
				t.Errorf("LoadJsonWithEX:Expired case failed: incorrect expired time")
			}
			if receiver.Key != "" || receiver.Value != "" {
				t.Errorf("LoadJsonWithEX:Expired case failed: incorrect value")
			}
		})

		// 设置了过期时间且没有过期的样例
		t.Run("LoadJsonWithEX:NotExpired", func(t *testing.T) {
			key := "LoadJsonWithEX:NotExpired"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJsonEX(context.Background(), key, value, time.Minute)
			if storeErr != nil {
				t.Errorf("LoadJsonWithEX:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("LoadJsonWithEX:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("LoadJsonWithEX:NotExpired case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("LoadJsonWithEX:NotExpired case failed: incorrect expired time")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadJsonWithEX:NotExpired case failed: incorrect value")
			}
		})
	}
}

func StoreFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 存储一个不存在的样例
		t.Run("Store:NotExist", func(t *testing.T) {
			key := "Store:NotExist"
			storeErr := impl.Store(context.Background(), key, "Store:NotExist")
			if storeErr != nil {
				t.Errorf("Store:NotExist case failed when storing key: %v", storeErr.Error())
			}

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Store:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("Store:NotExist case failed: key not exist, want exist")
			}
			if value != "Store:NotExist" {
				t.Errorf("Store:NotExist case failed: incorrect value")
			}
		})

		// 存储一个已经存在的样例
		t.Run("Store:Exist", func(t *testing.T) {
			prevKey, key := "Store:Exist:Prev", "Store:Exist"
			storePrevErr := impl.Store(context.Background(), prevKey, "Store:Exist:Prev")
			if storePrevErr != nil {
				t.Errorf("Store:Exist case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.Store(context.Background(), key, "Store:Exist")
			if storeErr != nil {
				t.Errorf("Store:Exist case failed when storing key: %v", storeErr.Error())
			}

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Store:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("Store:Exist case failed: key not exist, want exist")
			}
			if value != "Store:Exist" {
				t.Errorf("Store:Exist case failed: incorrect value")
			}
		})

		// 存储一个已经存在但是过期的样例
		t.Run("Store:Expired", func(t *testing.T) {
			prevKey, key := "Store:Expired:Prev", "Store:Expired"
			storePrevErr := impl.StoreEX(context.Background(), prevKey, "Store:Expired:Prev", expire)
			if storePrevErr != nil {
				t.Errorf("Store:Expired case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.Store(context.Background(), key, "Store:Expired")
			if storeErr != nil {
				t.Errorf("Store:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Store:Expired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("Store:Expired case failed: key not exist, want exist")
			}
			if value != "Store:Expired" {
				t.Errorf("Store:Expired case failed: incorrect value")
			}
		})

		// 存储一个设置了过期时间且没过期的样例
		t.Run("Store:NotExpired", func(t *testing.T) {
			prevKey, key := "Store:NotExpired:Prev", "Store:NotExpired"
			storePrevErr := impl.StoreEX(context.Background(), prevKey, "Store:NotExpired:Prev", time.Minute)
			if storePrevErr != nil {
				t.Errorf("Store:NotExpired case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.Store(context.Background(), key, "Store:NotExpired")
			if storeErr != nil {
				t.Errorf("Store:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			exist, value, loadErr := impl.Load(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Store:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("Store:NotExpired case failed: key not exist, want exist")
			}
			if value != "Store:NotExpired" {
				t.Errorf("Store:NotExpired case failed: incorrect value")
			}
		})

		// 存储一个非字符串的样例
		t.Run("Store:NotString", func(t *testing.T) {
			key, value := "Store:NotString", "NotString"
			addErr := impl.AddMember(context.Background(), key, value)
			if addErr == nil {
				t.Errorf("Store:NotString case failed: want error but not")
			}
		})
	}
}

func StoreEXFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 存储一个不存在的样例
		t.Run("StoreEX:NotExist", func(t *testing.T) {
			key := "StoreEX:NotExist"
			storeErr := impl.StoreEX(context.Background(), key, "StoreEX:NotExist", time.Minute)
			if storeErr != nil {
				t.Errorf("StoreEX:NotExist case failed when storing key: %v", storeErr.Error())
			}

			exist, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("StoreEX:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreEX:NotExist case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("StoreEX:NotExist case failed: incorrect expired time")
			}
			if value != "StoreEX:NotExist" {
				t.Errorf("StoreEX:NotExist case failed: incorrect value")
			}
		})

		// 存储一个已经存在的样例
		t.Run("StoreEX:Exist", func(t *testing.T) {
			prevKey, key := "StoreEX:Exist:Prev", "StoreEX:Exist"
			storePrevErr := impl.StoreEX(context.Background(), prevKey, "StoreEX:Exist:Prev", time.Hour)
			if storePrevErr != nil {
				t.Errorf("StoreEX:Exist case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.StoreEX(context.Background(), key, "StoreEX:Exist", time.Minute)
			if storeErr != nil {
				t.Errorf("StoreEX:Exist case failed when storing key: %v", storeErr.Error())
			}

			exist, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("StoreEX:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreEX:Exist case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("StoreEX:Exist case failed: incorrect expired time")
			}
			if value != "StoreEX:Exist" {
				t.Errorf("StoreEX:Exist case failed: incorrect value")
			}
		})

		// 存储一个已经存在但是过期的样例
		t.Run("StoreEX:Expired", func(t *testing.T) {
			prevKey, key := "StoreEX:Expired:Prev", "StoreEX:Expired"
			storePrevErr := impl.StoreEX(context.Background(), prevKey, "StoreEX:Expired:Prev", expire)
			if storePrevErr != nil {
				t.Errorf("StoreEX:Expired case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.StoreEX(context.Background(), key, "StoreEX:Expired", time.Minute)
			if storeErr != nil {
				t.Errorf("StoreEX:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, expiredTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("StoreEX:Expired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreEX:Expired case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("StoreEX:Expired case failed: incorrect expired time")
			}
			if value != "StoreEX:Expired" {
				t.Errorf("StoreEX:Expired case failed: incorrect value")
			}
		})
	}
}

func StoreJsonFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 存储一个不存在的样例
		t.Run("StoreJson:NotExist", func(t *testing.T) {
			key := "StoreJson:NotExist"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("StoreJson:NotExist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJson:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJson:NotExist case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJson:NotExist case failed: incorrect value")
			}
		})

		// 存储一个已经存在的样例
		t.Run("StoreJson:Exist", func(t *testing.T) {
			prevKey, key := "StoreJson:Exist:Prev", "StoreJson:Exist"
			prevValue, value := testStruct{Key: prevKey, Value: prevKey}, testStruct{Key: key, Value: key}

			storePrevErr := impl.StoreJson(context.Background(), prevKey, prevValue)
			if storePrevErr != nil {
				t.Errorf("StoreJson:Exist case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("StoreJson:Exist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJson:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJson:Exist case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJson:Exist case failed: incorrect value")
			}
		})

		// 存储一个已经存在但是过期的样例
		t.Run("StoreJson:Expired", func(t *testing.T) {
			prevKey, key := "StoreJson:Expired:Prev", "StoreJson:Expired"
			prevValue, value := testStruct{Key: prevKey, Value: prevKey}, testStruct{Key: key, Value: key}

			storePrevErr := impl.StoreJsonEX(context.Background(), prevKey, prevValue, expire)
			if storePrevErr != nil {
				t.Errorf("StoreJson:Expired case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("StoreJson:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJson:Expired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJson:Expired case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJson:Expired case failed: incorrect value")
			}
		})

		// 存储一个设置了过期时间且没过期的样例
		t.Run("StoreJson:NotExpired", func(t *testing.T) {
			prevKey, key := "StoreJson:NotExpired:Prev", "StoreJson:NotExpired"
			prevValue, value := testStruct{Key: prevKey, Value: prevKey}, testStruct{Key: key, Value: key}

			storePrevErr := impl.StoreJsonEX(context.Background(), prevKey, prevValue, time.Hour)
			if storePrevErr != nil {
				t.Errorf("StoreJson:NotExpired case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("StoreJson:NotExpired case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, loadErr := impl.LoadJson(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJson:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJson:NotExpired case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJson:NotExpired case failed: incorrect value")
			}
		})
	}
}

func StoreJsonEXFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 存储一个不存在的样例
		t.Run("StoreJsonEX:NotExist", func(t *testing.T) {
			key := "StoreJsonEX:NotExist"
			value := testStruct{Key: key, Value: key}

			storeErr := impl.StoreJsonEX(context.Background(), key, value, time.Minute)
			if storeErr != nil {
				t.Errorf("StoreJsonEX:NotExist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJsonEX:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJsonEX:NotExist case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("StoreJsonEX:NotExist case failed: incorrect expired time")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJsonEX:NotExist case failed: incorrect value")
			}
		})

		// 存储一个已经存在的样例
		t.Run("StoreJsonEX:Exist", func(t *testing.T) {
			prevKey, key := "StoreJsonEX:Exist:Prev", "StoreJsonEX:Exist"
			prevValue, value := testStruct{Key: prevKey, Value: prevKey}, testStruct{Key: key, Value: key}

			storePrevErr := impl.StoreJsonEX(context.Background(), prevKey, prevValue, time.Hour)
			if storePrevErr != nil {
				t.Errorf("StoreJsonEX:Exist case failed when storing prev key: %v", storePrevErr.Error())
			}
			storeErr := impl.StoreJsonEX(context.Background(), key, value, time.Minute)
			if storeErr != nil {
				t.Errorf("StoreJsonEX:Exist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJsonEX:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJsonEX:Exist case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("StoreJsonEX:Exist case failed: incorrect expired time")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJsonEX:Exist case failed: incorrect value")
			}
		})

		// 存储一个已经存在但是过期的样例
		t.Run("StoreJsonEX:Expired", func(t *testing.T) {
			prevKey, key := "StoreJsonEX:Expired:Prev", "StoreJsonEX:Expired"
			prevValue, value := testStruct{Key: prevKey, Value: prevKey}, testStruct{Key: key, Value: key}

			storePrevErr := impl.StoreJsonEX(context.Background(), prevKey, prevValue, expire)
			if storePrevErr != nil {
				t.Errorf("StoreJsonEX:Expired case failed when storing prev key: %v", storePrevErr.Error())
			}

			time.Sleep(sleepInterval)

			storeErr := impl.StoreJsonEX(context.Background(), key, value, time.Minute)
			if storeErr != nil {
				t.Errorf("StoreJsonEX:Expired case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, expiredTime, loadErr := impl.LoadJsonWithEX(context.Background(), key, &receiver)
			if loadErr != nil {
				t.Errorf("StoreJsonEX:Expired case failed when loading key: %v", loadErr.Error())
			}
			if !exist {
				t.Errorf("StoreJsonEX:Expired case failed: key not exist, want exist")
			}
			if !expiredTimeIsCorrect(expiredTime, time.Minute) {
				t.Errorf("StoreJsonEX:Expired case failed: incorrect expired time")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("StoreJsonEX:Expired case failed: incorrect value")
			}
		})
	}
}

func DeleteFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 删除一个不存在的样例
		t.Run("Delete:NotExist", func(t *testing.T) {
			key := "Delete:NotExist"
			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("Delete:NotExist case failed when deleting key: %v", deleteErr.Error())
			}

			exist, expireTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Delete:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("Delete:NotExist case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("Delete:NotExist case failed: incorrect value")
			}
			if expireTime != 0 {
				t.Errorf("Delete:NotExist case failed: incorrect expired time")
			}
		})

		// 删除一个存在的样例
		t.Run("Delete:Exist", func(t *testing.T) {
			key := "Delete:Exist"
			storeErr := impl.Store(context.Background(), key, "Delete:Exist")
			if storeErr != nil {
				t.Errorf("Delete:Exist case failed when storing key: %v", storeErr.Error())
			}

			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("Delete:Exist case failed when deleting key: %v", deleteErr.Error())
			}

			exist, expireTime, value, loadErr := impl.LoadWithEX(context.Background(), key)
			if loadErr != nil {
				t.Errorf("Delete:Exist case failed when loading key: %v", loadErr.Error())
			}
			if exist {
				t.Errorf("Delete:Exist case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("Delete:Exist case failed: incorrect value")
			}
			if expireTime != 0 {
				t.Errorf("Delete:Exist case failed: incorrect expired time")
			}
		})

		// 删除一个集合类型的样例
		t.Run("Delete:Set", func(t *testing.T) {
			key, value := "Delete:Set", "Delete:Set"
			addErr := impl.AddMember(context.Background(), key, value)
			if addErr != nil {
				t.Errorf("Delete:Set case failed when storing key: %v", addErr.Error())
			}

			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("Delete:Set case failed when deleting key: %v", deleteErr.Error())
			}

			checkExist, checkErr := impl.IsMember(context.Background(), key, value)
			if checkErr != nil {
				t.Errorf("Delete:Set case failed when checking key: %v", checkErr.Error())
			}
			if checkExist {
				t.Errorf("Delete:Set case failed: key exist, want not exist")
			}
		})

		// 删除一个哈希类型的样例
		t.Run("Delete:Hash", func(t *testing.T) {
			key, field, value := "Delete:Hash", "Delete:Hash", "Delete:Hash"
			addErr := impl.HSetValue(context.Background(), key, field, value)
			if addErr != nil {
				t.Errorf("Delete:Hash case failed when storing key: %v", addErr.Error())
			}

			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("Delete:Hash case failed when deleting key: %v", deleteErr.Error())
			}

			checkExist, checkValue, checkErr := impl.HGetValue(context.Background(), key, field)
			if checkErr != nil {
				t.Errorf("Delete:Hash case failed when checking key: %v", checkErr.Error())
			}
			if checkExist {
				t.Errorf("Delete:Hash case failed: key exist, want not exist")
			}
			if checkValue != "" {
				t.Errorf("Delete:Hash case failed: incorrect value")
			}
		})
	}
}

func LoadAndDeleteFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取并删除一个不存在的样例
		t.Run("LoadAndDelete:NotExist", func(t *testing.T) {
			key := "LoadAndDelete:NotExist"
			exist, value, deleteErr := impl.LoadAndDelete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("LoadAndDelete:NotExist case failed when deleting key: %v", deleteErr.Error())
			}
			if exist {
				t.Errorf("LoadAndDelete:NotExist case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("LoadAndDelete:NotExist case failed: incorrect value")
			}
		})

		// 获取并删除一个存在的样例
		t.Run("LoadAndDelete:Exist", func(t *testing.T) {
			key := "LoadAndDelete:Exist"
			storeErr := impl.Store(context.Background(), key, "LoadAndDelete:Exist")
			if storeErr != nil {
				t.Errorf("LoadAndDelete:Exist case failed when storing key: %v", storeErr.Error())
			}

			exist, value, deleteErr := impl.LoadAndDelete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("LoadAndDelete:Exist case failed when deleting key: %v", deleteErr.Error())
			}
			if !exist {
				t.Errorf("LoadAndDelete:Exist case failed: key not exist, want exist")
			}
			if value != "LoadAndDelete:Exist" {
				t.Errorf("LoadAndDelete:Exist case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadAndDelete:Exist case failed when loading key: %v", checkErr.Error())
			}
			if checkExist {
				t.Errorf("LoadAndDelete:Exist case failed: key exist, want not exist")
			}
			if checkValue != "" {
				t.Errorf("LoadAndDelete:Exist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadAndDelete:Exist case failed: incorrect expired time")
			}
		})
	}
}

func LoadAndDeleteJsonFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取并删除一个不存在的样例
		t.Run("LoadAndDeleteJson:NotExist", func(t *testing.T) {
			key := "LoadAndDeleteJson:NotExist"
			var receiver testStruct
			exist, deleteErr := impl.LoadAndDeleteJson(context.Background(), key, &receiver)
			if deleteErr != nil {
				t.Errorf("LoadAndDeleteJson:NotExist case failed when deleting key: %v", deleteErr.Error())
			}
			if exist {
				t.Errorf("LoadAndDeleteJson:NotExist case failed: key exist, want not exist")
			}
			if receiver.Key != "" || receiver.Value != "" {
				t.Errorf("LoadAndDeleteJson:NotExist case failed: incorrect value")
			}
		})

		// 获取并删除一个存在的样例
		t.Run("LoadAndDeleteJson:Exist", func(t *testing.T) {
			key := "LoadAndDeleteJson:Exist"
			value := testStruct{Key: key, Value: key}
			storeErr := impl.StoreJson(context.Background(), key, value)
			if storeErr != nil {
				t.Errorf("LoadAndDeleteJson:Exist case failed when storing key: %v", storeErr.Error())
			}

			var receiver testStruct
			exist, deleteErr := impl.LoadAndDeleteJson(context.Background(), key, &receiver)
			if deleteErr != nil {
				t.Errorf("LoadAndDeleteJson:Exist case failed when deleting key: %v", deleteErr.Error())
			}
			if !exist {
				t.Errorf("LoadAndDeleteJson:Exist case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadAndDeleteJson:Exist case failed: incorrect value")
			}

			var checkReceiver testStruct
			checkExist, checkExpireTime, checkErr := impl.LoadJsonWithEX(context.Background(), key, &checkReceiver)
			if checkErr != nil {
				t.Errorf("LoadAndDeleteJson:Exist case failed when loading key: %v", checkErr.Error())
			}
			if checkExist {
				t.Errorf("LoadAndDeleteJson:Exist case failed: key exist, want not exist")
			}
			if checkReceiver.Key != "" || checkReceiver.Value != "" {
				t.Errorf("LoadAndDeleteJson:Exist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadAndDeleteJson:Exist case failed: incorrect expired time")
			}
		})
	}
}

func LoadOrStoreFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取或存储一个不存在的样例
		t.Run("LoadOrStore:NotExist", func(t *testing.T) {
			key := "LoadOrStore:NotExist"
			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("LoadOrStore:NotExist case failed when deleting key: %v", deleteErr.Error())
			}

			loaded, value, loadErr := impl.LoadOrStore(context.Background(), key, "LoadOrStore:NotExist")
			if loadErr != nil {
				t.Errorf("LoadOrStore:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStore:NotExist case failed: key exist, want not exist")
			}
			if value != "LoadOrStore:NotExist" {
				t.Errorf("LoadOrStore:NotExist case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStore:NotExist case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStore:NotExist case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStore:NotExist" {
				t.Errorf("LoadOrStore:NotExist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadOrStore:NotExist case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在的样例
		t.Run("LoadOrStore:Exist", func(t *testing.T) {
			key := "LoadOrStore:Exist"
			storeErr := impl.Store(context.Background(), key, "LoadOrStore:Exist:Prev")
			if storeErr != nil {
				t.Errorf("LoadOrStore:Exist case failed when storing key: %v", storeErr.Error())
			}

			loaded, value, loadErr := impl.LoadOrStore(context.Background(), key, "LoadOrStore:Exist")
			if loadErr != nil {
				t.Errorf("LoadOrStore:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadOrStore:Exist case failed: key not exist, want exist")
			}
			if value != "LoadOrStore:Exist:Prev" {
				t.Errorf("LoadOrStore:Exist case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStore:Exist case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStore:Exist case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStore:Exist:Prev" {
				t.Errorf("LoadOrStore:Exist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadOrStore:Exist case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在但是过期的样例
		t.Run("LoadOrStore:Expired", func(t *testing.T) {
			key := "LoadOrStore:Expired"
			storeErr := impl.StoreEX(context.Background(), key, "LoadOrStore:Expired:Prev", expire)
			if storeErr != nil {
				t.Errorf("LoadOrStore:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			loaded, value, loadErr := impl.LoadOrStore(context.Background(), key, "LoadOrStore:Expired")
			if loadErr != nil {
				t.Errorf("LoadOrStore:Expired case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStore:Expired case failed: key exist, want not exist")
			}
			if value != "LoadOrStore:Expired" {
				t.Errorf("LoadOrStore:Expired case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStore:Expired case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStore:Expired case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStore:Expired" {
				t.Errorf("LoadOrStore:Expired case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadOrStore:Expired case failed: incorrect expired time")
			}
		})

		// 获取或存储一个设置了过期时间且没过期的样例
		t.Run("LoadOrStore:NotExpired", func(t *testing.T) {
			key := "LoadOrStore:NotExpired"
			storePrevErr := impl.StoreEX(context.Background(), key, "LoadOrStore:NotExpired:Prev", time.Hour)
			if storePrevErr != nil {
				t.Errorf("LoadOrStore:NotExpired case failed when storing prev key: %v", storePrevErr.Error())
			}

			loaded, value, loadErr := impl.LoadOrStore(context.Background(), key, "LoadOrStore:NotExpired")
			if loadErr != nil {
				t.Errorf("LoadOrStore:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadOrStore:NotExpired case failed: key not exist, want exist")
			}
			if value != "LoadOrStore:NotExpired:Prev" {
				t.Errorf("LoadOrStore:NotExpired case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStore:NotExpired case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStore:NotExpired case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStore:NotExpired:Prev" {
				t.Errorf("LoadOrStore:NotExpired case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Hour) {
				t.Errorf("LoadOrStore:NotExpired case failed: incorrect expired time")
			}
		})
	}
}

func LoadOrStoreEXFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取或存储一个不存在的样例
		t.Run("LoadOrStoreEX:NotExist", func(t *testing.T) {
			key := "LoadOrStoreEX:NotExist"

			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("LoadOrStoreEX:NotExist case failed when deleting key: %v", deleteErr.Error())
			}

			loaded, value, loadErr := impl.LoadOrStoreEX(context.Background(), key, "LoadOrStoreEX:NotExist", time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreEX:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStoreEX:NotExist case failed: key exist, want not exist")
			}
			if value != "LoadOrStoreEX:NotExist" {
				t.Errorf("LoadOrStoreEX:NotExist case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStoreEX:NotExist case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStoreEX:NotExist case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStoreEX:NotExist" {
				t.Errorf("LoadOrStoreEX:NotExist case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Minute) {
				t.Errorf("LoadOrStoreEX:NotExist case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在的样例
		t.Run("LoadOrStoreEX:Exist", func(t *testing.T) {
			key := "LoadOrStoreEX:Exist"
			storeErr := impl.Store(context.Background(), key, "LoadOrStoreEX:Exist:Prev")
			if storeErr != nil {
				t.Errorf("LoadOrStoreEX:Exist case failed when storing key: %v", storeErr.Error())
			}

			loaded, value, loadErr := impl.LoadOrStoreEX(context.Background(), key, "LoadOrStoreEX:Exist", time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreEX:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadOrStoreEX:Exist case failed: key not exist, want exist")
			}
			if value != "LoadOrStoreEX:Exist:Prev" {
				t.Errorf("LoadOrStoreEX:Exist case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStoreEX:Exist case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStoreEX:Exist case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStoreEX:Exist:Prev" {
				t.Errorf("LoadOrStoreEX:Exist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadOrStoreEX:Exist case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在但是过期的样例
		t.Run("LoadOrStoreEX:Expired", func(t *testing.T) {
			key := "LoadOrStoreEX:Expired"
			storeErr := impl.StoreEX(context.Background(), key, "LoadOrStoreEX:Expired:Prev", expire)
			if storeErr != nil {
				t.Errorf("LoadOrStoreEX:Expired case failed when storing key: %v", storeErr.Error())
			}

			time.Sleep(sleepInterval)

			loaded, value, loadErr := impl.LoadOrStoreEX(context.Background(), key, "LoadOrStoreEX:Expired", time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreEX:Expired case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStoreEX:Expired case failed: key exist, want not exist")
			}
			if value != "LoadOrStoreEX:Expired" {
				t.Errorf("LoadOrStoreEX:Expired case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStoreEX:Expired case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStoreEX:Expired case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStoreEX:Expired" {
				t.Errorf("LoadOrStoreEX:Expired case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Minute) {
				t.Errorf("LoadOrStoreEX:Expired case failed: incorrect expired time")
			}
		})

		// 获取或存储一个设置了过期时间且没过期的样例
		t.Run("LoadOrStoreEX:NotExpired", func(t *testing.T) {
			key := "LoadOrStoreEX:NotExpired"
			storePrevErr := impl.StoreEX(context.Background(), key, "LoadOrStoreEX:NotExpired:Prev", time.Hour)
			if storePrevErr != nil {
				t.Errorf("LoadOrStoreEX:NotExpired case failed when storing prev key: %v", storePrevErr.Error())
			}

			loaded, value, loadErr := impl.LoadOrStoreEX(context.Background(), key, "LoadOrStoreEX:NotExpired", time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreEX:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadOrStoreEX:NotExpired case failed: key not exist, want exist")
			}
			if value != "LoadOrStoreEX:NotExpired:Prev" {
				t.Errorf("LoadOrStoreEX:NotExpired case failed: incorrect value")
			}

			checkExist, checkExpireTime, checkValue, checkErr := impl.LoadWithEX(context.Background(), key)
			if checkErr != nil {
				t.Errorf("LoadOrStoreEX:NotExpired case failed when loading key: %v", checkErr.Error())
			}
			if !checkExist {
				t.Errorf("LoadOrStoreEX:NotExpired case failed: key not exist, want exist")
			}
			if checkValue != "LoadOrStoreEX:NotExpired:Prev" {
				t.Errorf("LoadOrStoreEX:NotExpired case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Hour) {
				t.Errorf("LoadOrStoreEX:NotExpired case failed: incorrect expired time")
			}
		})
	}
}

func LoadOrStoreJsonFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取或存储一个不存在的样例
		t.Run("LoadOrStoreJson:NotExist", func(t *testing.T) {
			key := "LoadOrStoreJson:NotExist"
			value := testStruct{Key: key, Value: key}

			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("LoadOrStoreJson:NotExist case failed when deleting key: %v", deleteErr.Error())
			}

			var receiver testStruct
			loaded, loadErr := impl.LoadOrStoreJson(context.Background(), key, value, &receiver)
			if loadErr != nil {
				t.Errorf("LoadOrStoreJson:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStoreJson:NotExist case failed: key exist, want not exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadOrStoreJson:NotExist case failed: incorrect value: %+v", receiver)
			}

			var checkReceiver testStruct
			exist, checkExpireTime, checkErr := impl.LoadJsonWithEX(context.Background(), key, &checkReceiver)
			if checkErr != nil {
				t.Errorf("LoadOrStoreJson:NotExist case failed when loading key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("LoadOrStoreJson:NotExist case failed: key not exist, want exist")
			}
			if checkReceiver.Key != key || checkReceiver.Value != key {
				t.Errorf("LoadOrStoreJson:NotExist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadOrStoreJson:NotExist case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在的样例
		t.Run("LoadOrStoreJson:Exist", func(t *testing.T) {
			key := "LoadOrStoreJson:Exist"
			prevValue, value := testStruct{Key: key, Value: key + ":Prev"}, testStruct{Key: key, Value: key}
			storePrevErr := impl.StoreJson(context.Background(), key, prevValue)
			if storePrevErr != nil {
				t.Errorf("LoadOrStoreJson:Exist case failed when storing prev key: %v", storePrevErr.Error())
			}

			var receiver testStruct
			loaded, loadErr := impl.LoadOrStoreJson(context.Background(), key, &value, &receiver)
			if loadErr != nil {
				t.Errorf("LoadOrStoreJson:Exist case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadOrStoreJson:Exist case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key+":Prev" {
				t.Errorf("LoadOrStoreJson:Exist case failed: incorrect value")
			}

			var checkReceiver testStruct
			exist, checkExpireTime, checkErr := impl.LoadJsonWithEX(context.Background(), key, &checkReceiver)
			if checkErr != nil {
				t.Errorf("LoadOrStoreJson:Exist case failed when loading key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("LoadOrStoreJson:Exist case failed: key not exist, want exist")
			}
			if checkReceiver.Key != key || checkReceiver.Value != key+":Prev" {
				t.Errorf("LoadOrStoreJson:Exist case failed: incorrect value")
			}
			if checkExpireTime != 0 {
				t.Errorf("LoadOrStoreJson:Exist case failed: incorrect expired time")
			}
		})
	}
}

func LoadOrStoreJsonEXFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取或存储一个不存在的样例
		t.Run("LoadOrStoreJsonEX:NotExist", func(t *testing.T) {
			key := "LoadOrStoreJsonEX:NotExist"
			value := testStruct{Key: key, Value: key}

			deleteErr := impl.Delete(context.Background(), key)
			if deleteErr != nil {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed when deleting key: %v", deleteErr.Error())
			}

			var receiver testStruct
			loaded, loadErr := impl.LoadOrStoreJsonEX(context.Background(), key, &value, &receiver, time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed: key exist, want not exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed: incorrect value")
			}

			var checkReceiver testStruct
			exist, checkExpireTime, checkErr := impl.LoadJsonWithEX(context.Background(), key, &checkReceiver)
			if checkErr != nil {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed when loading key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed: key not exist, want exist")
			}
			if checkReceiver.Key != key || checkReceiver.Value != key {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Minute) {
				t.Errorf("LoadOrStoreJsonEX:NotExist case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在且未过期的样例
		t.Run("LoadOrStoreJsonEX:NotExpired", func(t *testing.T) {
			key := "LoadOrStoreJsonEX:NotExpired"
			prevValue, value := testStruct{Key: key, Value: key + ":Prev"}, testStruct{Key: key, Value: key}
			storePrevErr := impl.StoreJsonEX(context.Background(), key, prevValue, time.Hour)
			if storePrevErr != nil {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed when storing prev key: %v", storePrevErr.Error())
			}

			var receiver testStruct
			loaded, loadErr := impl.LoadOrStoreJsonEX(context.Background(), key, value, &receiver, time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed when loading key: %v", loadErr.Error())
			}
			if !loaded {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed: key not exist, want exist")
			}
			if receiver.Key != key || receiver.Value != key+":Prev" {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed: incorrect value")
			}

			var checkReceiver testStruct
			exist, checkExpireTime, checkErr := impl.LoadJsonWithEX(context.Background(), key, &checkReceiver)
			if checkErr != nil {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed when loading key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed: key not exist, want exist")
			}
			if checkReceiver.Key != key || checkReceiver.Value != key+":Prev" {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Hour) {
				t.Errorf("LoadOrStoreJsonEX:NotExpired case failed: incorrect expired time")
			}
		})

		// 获取或存储一个存在但是过期的样例
		t.Run("LoadOrStoreJsonEX:Expired", func(t *testing.T) {
			key := "LoadOrStoreJsonEX:Expired"
			prevValue, value := testStruct{Key: key, Value: key + ":Prev"}, testStruct{Key: key, Value: key}
			storePrevErr := impl.StoreJsonEX(context.Background(), key, prevValue, expire)
			if storePrevErr != nil {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed when storing prev key: %v", storePrevErr.Error())
			}

			time.Sleep(sleepInterval)

			var receiver testStruct
			loaded, loadErr := impl.LoadOrStoreJsonEX(context.Background(), key, &value, &receiver, time.Minute)
			if loadErr != nil {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed when loading key: %v", loadErr.Error())
			}
			if loaded {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed: key exist, want not exist")
			}
			if receiver.Key != key || receiver.Value != key {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed: incorrect value")
			}

			var checkReceiver testStruct
			exist, checkExpireTime, checkErr := impl.LoadJsonWithEX(context.Background(), key, &checkReceiver)
			if checkErr != nil {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed when loading key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed: key not exist, want exist")
			}
			if checkReceiver.Key != key || checkReceiver.Value != key {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed: incorrect value")
			}
			if !expiredTimeIsCorrect(checkExpireTime, time.Minute) {
				t.Errorf("LoadOrStoreJsonEX:Expired case failed: incorrect expired time")
			}
		})
	}
}

func IsMemberFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 判断一个不存在的样例
		t.Run("IsMember:NotExist", func(t *testing.T) {
			key, member := "IsMember:NotExist", "NotExistCase"
			exist, err := impl.IsMember(context.Background(), key, member)
			if err != nil {
				t.Errorf("IsMember:NotExist case failed when checking key: %v", err.Error())
			}
			if exist {
				t.Errorf("IsMember:NotExist case failed: key exist, want not exist")
			}
		})

		// 判断一个存在且未过期的样例
		t.Run("IsMember:NotExpired", func(t *testing.T) {
			key, member := "IsMember:NotExpired", "NotExpired"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("IsMember:NotExpired case failed when adding member: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("IsMember:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			exist, err := impl.IsMember(context.Background(), key, member)
			if err != nil {
				t.Errorf("IsMember:NotExpired case failed when checking key: %v", err.Error())
			}
			if !exist {
				t.Errorf("IsMember:NotExpired case failed: key not exist, want exist")
			}
		})

		// 判断一个存在但是过期的样例
		t.Run("IsMember:Expired", func(t *testing.T) {
			key, member := "IsMember:Expired", "Expired"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("IsMember:Expired case failed when adding member: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("IsMember:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, err := impl.IsMember(context.Background(), key, member)
			if err != nil {
				t.Errorf("IsMember:Expired case failed when checking key: %v", err.Error())
			}
			if exist {
				t.Errorf("IsMember:Expired case failed: key exist, want not exist")
			}
		})

		// 判断错误类型的样例
		t.Run("IsMember:WrongType", func(t *testing.T) {
			key, member := "IsMember:WrongType", "WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("IsMember:WrongType case failed when storing key: %v", storeErr.Error())
			}

			exist, err := impl.IsMember(context.Background(), key, member)
			if err == nil {
				t.Errorf("IsMember:WrongType case failed: no error")
			}
			if exist {
				t.Errorf("IsMember:WrongType case failed: key exist, want not exist")
			}
		})
	}
}

func IsMembersFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 判断多个不存在的样例
		t.Run("IsMembers:NotExist", func(t *testing.T) {
			key, members := "IsMembers:NotExist", []string{"NotExistCase1", "NotExistCase2"}
			exist, err := impl.IsMembers(context.Background(), key, members...)
			if err != nil {
				t.Errorf("IsMembers:NotExist case failed when checking key: %v", err.Error())
			}
			if exist {
				t.Errorf("IsMembers:NotExist case failed: key exist, want not exist")
			}
		})

		// 判断多个存在且未过期的样例
		t.Run("IsMembers:NotExpired", func(t *testing.T) {
			key, members := "IsMembers:NotExpired", []string{"NotExpired1", "NotExpired2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("IsMembers:NotExpired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("IsMembers:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			exist, err := impl.IsMembers(context.Background(), key, members...)
			if err != nil {
				t.Errorf("IsMembers:NotExpired case failed when checking key: %v", err.Error())
			}
			if !exist {
				t.Errorf("IsMembers:NotExpired case failed: key not exist, want exist")
			}
		})

		// 判断多个存在但是过期的样例
		t.Run("IsMembers:Expired", func(t *testing.T) {
			key, members := "IsMembers:Expired", []string{"Expired1", "Expired2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("IsMembers:Expired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("IsMembers:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, err := impl.IsMembers(context.Background(), key, members...)
			if err != nil {
				t.Errorf("IsMembers:Expired case failed when checking key: %v", err.Error())
			}
			if exist {
				t.Errorf("IsMembers:Expired case failed: key exist, want not exist")
			}
		})

		// 判断部分存在且未过期的样例
		t.Run("IsMembers:PartNotExpired", func(t *testing.T) {
			key, members := "IsMembers:PartNotExpired", []string{"PartNotExpired1", "PartNotExpired2"}
			addErr := impl.AddMembers(context.Background(), key, members[:1]...)
			if addErr != nil {
				t.Errorf("IsMembers:PartNotExpired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("IsMembers:PartNotExpired case failed when setting expire: %v", expireErr.Error())
			}

			exist, err := impl.IsMembers(context.Background(), key, members...)
			if err != nil {
				t.Errorf("IsMembers:PartNotExpired case failed when checking key: %v", err.Error())
			}
			if exist {
				t.Errorf("IsMembers:PartNotExpired case failed: key exist, want not exist")
			}
		})

		// 判断空样例
		t.Run("IsMembers:Empty", func(t *testing.T) {
			key, members := "IsMembers:Empty", []string{"Empty1", "Empty2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("IsMembers:Empty case failed when adding members: %v", addErr.Error())
			}

			exist, err := impl.IsMembers(context.Background(), key)
			if err != nil {
				t.Errorf("IsMembers:Empty case failed when checking key: %v", err.Error())
			}
			if exist {
				t.Errorf("IsMembers:Empty case failed: key exist, want not exist")
			}
		})
	}
}

func AddMemberFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 添加一个不存在的样例
		t.Run("AddMember:NotExist", func(t *testing.T) {
			key, member := "AddMember:NotExist", "NotExistCase"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("AddMember:NotExist case failed when adding member: %v", addErr.Error())
			}

			exist, checkErr := impl.IsMember(context.Background(), key, member)
			if checkErr != nil {
				t.Errorf("AddMember:NotExist case failed when checking key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("AddMember:NotExist case failed: key not exist, want exist")
			}
		})

		// 添加一个存在且未过期的样例
		t.Run("AddMember:NotExpired", func(t *testing.T) {
			key, member := "AddMember:NotExpired", "NotExpired"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("AddMember:NotExpired case failed when adding member: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("AddMember:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			executeErr := impl.AddMember(context.Background(), key, member)
			if executeErr != nil {
				t.Errorf("AddMember:NotExpired case failed when adding member: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMember(context.Background(), key, member)
			if checkErr != nil {
				t.Errorf("AddMember:NotExpired case failed when checking key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("AddMember:NotExpired case failed: key not exist, want exist")
			}
		})

		// 添加一个存在但是过期的样例
		t.Run("AddMember:Expired", func(t *testing.T) {
			key, member := "AddMember:Expired", "Expired"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("AddMember:Expired case failed when adding member: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("AddMember:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			executeErr := impl.AddMember(context.Background(), key, member)
			if executeErr != nil {
				t.Errorf("AddMember:Expired case failed when adding member: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMember(context.Background(), key, member)
			if checkErr != nil {
				t.Errorf("AddMember:Expired case failed when checking key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("AddMember:Expired case failed: key not exist, want exist")
			}
		})

		// 向已经添加的样例中添加一个不存在的样例
		t.Run("AddMember:NotExistInExist", func(t *testing.T) {
			key, member1, member2 := "AddMember:NotExistInExist", "NotExistInExist1", "NotExistInExist2"
			addErr := impl.AddMember(context.Background(), key, member1)
			if addErr != nil {
				t.Errorf("AddMember:NotExistInExist case failed when adding member: %v", addErr.Error())
			}

			executeErr := impl.AddMember(context.Background(), key, member2)
			if executeErr != nil {
				t.Errorf("AddMember:NotExistInExist case failed when adding member: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMembers(context.Background(), key, member1, member2)
			if checkErr != nil {
				t.Errorf("AddMember:NotExistInExist case failed when checking key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("AddMember:NotExistInExist case failed: key not exist, want exist")
			}
		})

		// 添加错误类型的样例
		t.Run("AddMember:WrongType", func(t *testing.T) {
			key, member := "AddMember:WrongType", "WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("AddMember:WrongType case failed when storing key: %v", storeErr.Error())
			}

			executeErr := impl.AddMember(context.Background(), key, member)
			if executeErr == nil {
				t.Errorf("AddMember:WrongType case failed: no error")
			}

		})
	}
}

func AddMembersFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 添加多个不存在的样例
		t.Run("AddMembers:NotExist", func(t *testing.T) {
			key, members := "AddMembers:NotExist", []string{"NotExistCase1", "NotExistCase2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("AddMembers:NotExist case failed when adding members: %v", addErr.Error())
			}

			exist, checkErr := impl.IsMembers(context.Background(), key, members...)
			if checkErr != nil {
				t.Errorf("AddMembers:NotExist case failed when checking key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("AddMembers:NotExist case failed: key not exist, want exist")
			}
		})

		// 添加多个存在且未过期的样例
		t.Run("AddMembers:NotExpired", func(t *testing.T) {
			key, members1, members2 := "AddMembers:NotExpired", []string{"NotExpired1", "NotExpired2"}, []string{"NotExpired3", "NotExpired4"}
			addErr := impl.AddMembers(context.Background(), key, members1...)
			if addErr != nil {
				t.Errorf("AddMembers:NotExpired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("AddMembers:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			executeErr := impl.AddMembers(context.Background(), key, members2...)
			if executeErr != nil {
				t.Errorf("AddMembers:NotExpired case failed when adding members: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMembers(context.Background(), key, append(members1, members2...)...)
			if checkErr != nil {
				t.Errorf("AddMembers:NotExpired case failed when checking key: %v", checkErr.Error())
			}
			if !exist {
				t.Errorf("AddMembers:NotExpired case failed: key not exist, want exist")
			}
		})

		// 添加多个存在但是过期的样例
		t.Run("AddMembers:Expired", func(t *testing.T) {
			key, members1, members2 := "AddMembers:Expired", []string{"Expired1", "Expired2"}, []string{"Expired3", "Expired4"}
			addErr := impl.AddMembers(context.Background(), key, members1...)
			if addErr != nil {
				t.Errorf("AddMembers:Expired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("AddMembers:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			executeErr := impl.AddMembers(context.Background(), key, members2...)
			if executeErr != nil {
				t.Errorf("AddMembers:Expired case failed when adding members: %v", executeErr.Error())
			}

			checkExist1, checkErr1 := impl.IsMembers(context.Background(), key, members1...)
			if checkErr1 != nil {
				t.Errorf("AddMembers:Expired case failed when checking key: %v", checkErr1.Error())
			}
			if checkExist1 {
				t.Errorf("AddMembers:Expired case failed: key exist, want not exist")
			}

			checkExist2, checkErr2 := impl.IsMembers(context.Background(), key, members2...)
			if checkErr2 != nil {
				t.Errorf("AddMembers:Expired case failed when checking key: %v", checkErr2.Error())
			}
			if !checkExist2 {
				t.Errorf("AddMembers:Expired case failed: key not exist, want exist")
			}
		})

		// 添加错误类型的样例
		t.Run("AddMembers:WrongType", func(t *testing.T) {
			key, members := "AddMembers:WrongType", []string{"WrongType1", "WrongType2"}
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("AddMembers:WrongType case failed when storing key: %v", storeErr.Error())
			}

			executeErr := impl.AddMembers(context.Background(), key, members...)
			if executeErr == nil {
				t.Errorf("AddMembers:WrongType case failed: no error")
			}
		})
	}
}

func RemoveMemberFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 移除一个不存在的样例
		t.Run("RemoveMember:NotExist", func(t *testing.T) {
			key := "RemoveMember:NotExist"
			executeErr := impl.RemoveMember(context.Background(), key, "NotExistCase")
			if executeErr != nil {
				t.Errorf("RemoveMember:NotExist case failed when removing member: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMember(context.Background(), key, "NotExistCase")
			if checkErr != nil {
				t.Errorf("RemoveMember:NotExist case failed when checking key: %v", checkErr.Error())
			}
			if exist {
				t.Errorf("RemoveMember:NotExist case failed: key exist, want not exist")
			}
		})

		// 移除一个存在且未过期的样例
		t.Run("RemoveMember:NotExpired", func(t *testing.T) {
			key, member := "RemoveMember:NotExpired", "NotExpired"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("RemoveMember:NotExpired case failed when adding member: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("RemoveMember:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			executeErr := impl.RemoveMember(context.Background(), key, member)
			if executeErr != nil {
				t.Errorf("RemoveMember:NotExpired case failed when removing member: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMember(context.Background(), key, member)
			if checkErr != nil {
				t.Errorf("RemoveMember:NotExpired case failed when checking key: %v", checkErr.Error())
			}
			if exist {
				t.Errorf("RemoveMember:NotExpired case failed: key exist, want not exist")
			}
		})

		// 移除一个存在但是过期的样例
		t.Run("RemoveMember:Expired", func(t *testing.T) {
			key, member := "RemoveMember:Expired", "Expired"
			addErr := impl.AddMember(context.Background(), key, member)
			if addErr != nil {
				t.Errorf("RemoveMember:Expired case failed when adding member: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("RemoveMember:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			executeErr := impl.RemoveMember(context.Background(), key, member)
			if executeErr != nil {
				t.Errorf("RemoveMember:Expired case failed when removing member: %v", executeErr.Error())
			}

			exist, checkErr := impl.IsMember(context.Background(), key, member)
			if checkErr != nil {
				t.Errorf("RemoveMember:Expired case failed when checking key: %v", checkErr.Error())
			}
			if exist {
				t.Errorf("RemoveMember:Expired case failed: key exist, want not exist")
			}
		})

		// 移除错误类型的样例
		t.Run("RemoveMember:WrongType", func(t *testing.T) {
			key, member := "RemoveMember:WrongType", "WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("RemoveMember:WrongType case failed when storing key: %v", storeErr.Error())
			}

			executeErr := impl.RemoveMember(context.Background(), key, member)
			if executeErr == nil {
				t.Errorf("RemoveMember:WrongType case failed: no error")
			}
		})
	}
}

func GetMembersFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("GetMembers:NotExist", func(t *testing.T) {
			key := "GetMembers:NotExist"
			members, getMembersErr := impl.GetMembers(context.Background(), key)
			if getMembersErr != nil {
				t.Errorf("GetMembers:NotExist case failed when getting members: %v", getMembersErr.Error())
			}
			if members == nil {
				t.Errorf("GetMembers:NotExist case failed: members is nil")
			}
			if len(members) != 0 {
				t.Errorf("GetMembers:NotExist case failed: incorrect members")
			}
		})

		// 获取一个存在且未过期的样例
		t.Run("GetMembers:NotExpired", func(t *testing.T) {
			key, members := "GetMembers:NotExpired", []string{"NotExpired1", "NotExpired2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetMembers:NotExpired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("GetMembers:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			getMembers, getMembersErr := impl.GetMembers(context.Background(), key)
			if getMembersErr != nil {
				t.Errorf("GetMembers:NotExpired case failed when getting members: %v", getMembersErr.Error())
			}
			if getMembers == nil {
				t.Errorf("GetMembers:NotExpired case failed: members is nil")
			}
			if len(getMembers) != len(members) {
				t.Errorf("GetMembers:NotExpired case failed: incorrect members")
			}
			for _, member := range members {
				if !containsString(getMembers, member) {
					t.Errorf("GetMembers:NotExpired case failed: incorrect members")
				}
			}
		})

		// 获取一个存在但是过期的样例
		t.Run("GetMembers:Expired", func(t *testing.T) {
			key, members := "GetMembers:Expired", []string{"Expired1", "Expired2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetMembers:Expired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("GetMembers:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			getMembers, getMembersErr := impl.GetMembers(context.Background(), key)
			if getMembersErr != nil {
				t.Errorf("GetMembers:Expired case failed when getting members: %v", getMembersErr.Error())
			}
			if getMembers == nil {
				t.Errorf("GetMembers:Expired case failed: members is nil")
			}
			if len(getMembers) != 0 {
				t.Errorf("GetMembers:Expired case failed: incorrect members")
			}
		})

		// 获取部分删除的样例
		t.Run("GetMembers:PartRemoved", func(t *testing.T) {
			key, members, deletedMember := "GetMembers:PartRemoved", []string{"PartRemoved1", "PartRemoved2"}, "PartRemoved1"
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetMembers:PartRemoved case failed when adding members: %v", addErr.Error())
			}

			removeErr := impl.RemoveMember(context.Background(), key, deletedMember)
			if removeErr != nil {
				t.Errorf("GetMembers:PartRemoved case failed when removing member: %v", removeErr.Error())
			}

			getMembers, getMembersErr := impl.GetMembers(context.Background(), key)
			if getMembersErr != nil {
				t.Errorf("GetMembers:PartRemoved case failed when getting members: %v", getMembersErr.Error())
			}
			if getMembers == nil {
				t.Errorf("GetMembers:PartRemoved case failed: members is nil")
			}
			if len(getMembers) != len(members)-1 {
				t.Errorf("GetMembers:PartRemoved case failed: incorrect members")
			}
			if containsString(getMembers, deletedMember) {
				t.Errorf("GetMembers:PartRemoved case failed: incorrect members")
			}
			if !containsString(getMembers, members[1]) {
				t.Errorf("GetMembers:PartRemoved case failed: incorrect members")
			}
		})

		// 获取错误类型的样例
		t.Run("GetMembers:WrongType", func(t *testing.T) {
			key := "GetMembers:WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("GetMembers:WrongType case failed when storing key: %v", storeErr.Error())
			}

			getMembers, getMembersErr := impl.GetMembers(context.Background(), key)
			if getMembersErr == nil {
				t.Errorf("GetMembers:WrongType case failed: no error")
			}
			if getMembers == nil {
				t.Errorf("GetMembers:WrongType case failed: members is nil")
			}
			if len(getMembers) != 0 {
				t.Errorf("GetMembers:WrongType case failed: incorrect members")
			}
		})
	}
}

func GetRandomMemberFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("GetRandomMember:NotExist", func(t *testing.T) {
			key := "GetRandomMember:NotExist"
			member, getRandomMemberErr := impl.GetRandomMember(context.Background(), key)
			if getRandomMemberErr != nil {
				t.Errorf("GetRandomMember:NotExist case failed when getting random member: %v", getRandomMemberErr.Error())
			}
			if member != "" {
				t.Errorf("GetRandomMember:NotExist case failed: incorrect member")
			}
		})

		// 获取一个存在且未过期的样例
		t.Run("GetRandomMember:NotExpired", func(t *testing.T) {
			key, members := "GetRandomMember:NotExpired", []string{"NotExpired1", "NotExpired2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMember:NotExpired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("GetRandomMember:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			member, getRandomMemberErr := impl.GetRandomMember(context.Background(), key)
			if getRandomMemberErr != nil {
				t.Errorf("GetRandomMember:NotExpired case failed when getting random member: %v", getRandomMemberErr.Error())
			}
			if member == "" {
				t.Errorf("GetRandomMember:NotExpired case failed: member is empty")
			}
			if !containsString(members, member) {
				t.Errorf("GetRandomMember:NotExpired case failed: incorrect member")
			}
		})

		// 获取一个存在但是过期的样例
		t.Run("GetRandomMember:Expired", func(t *testing.T) {
			key, members := "GetRandomMember:Expired", []string{"Expired1", "Expired2"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMember:Expired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("GetRandomMember:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			member, getRandomMemberErr := impl.GetRandomMember(context.Background(), key)
			if getRandomMemberErr != nil {
				t.Errorf("GetRandomMember:Expired case failed when getting random member: %v", getRandomMemberErr.Error())
			}
			if member != "" {
				t.Errorf("GetRandomMember:Expired case failed: incorrect member")
			}
		})
	}
}

func GetRandomMembersFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("GetRandomMembers:NotExist", func(t *testing.T) {
			key := "GetRandomMembers:NotExist"
			members, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 2)
			if getRandomMembersErr != nil {
				t.Errorf("GetRandomMembers:NotExist case failed when getting random members: %v", getRandomMembersErr.Error())
			}
			if members == nil {
				t.Errorf("GetRandomMembers:NotExist case failed: members is nil")
			}
			if len(members) != 0 {
				t.Errorf("GetRandomMembers:NotExist case failed: incorrect members")
			}
		})

		// 获取一个存在且未过期的样例
		t.Run("GetRandomMembers:NotExpired", func(t *testing.T) {
			key, members := "GetRandomMembers:NotExpired", []string{"NotExpired1", "NotExpired2", "NotExpired3", "NotExpired4"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMembers:NotExpired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("GetRandomMembers:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			gotMembers, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 2)
			if getRandomMembersErr != nil {
				t.Errorf("GetRandomMembers:NotExpired case failed when getting random members: %v", getRandomMembersErr.Error())
			}
			if gotMembers == nil {
				t.Errorf("GetRandomMembers:NotExpired case failed: members is nil")
			}
			if len(gotMembers) != 2 {
				t.Errorf("GetRandomMembers:NotExpired case failed: incorrect members")
			}
			for _, member := range gotMembers {
				if !containsString(members, member) {
					t.Errorf("GetRandomMembers:NotExpired case failed: incorrect members")
				}
			}
		})

		// 获取一个存在但是过期的样例
		t.Run("GetRandomMembers:Expired", func(t *testing.T) {
			key, members := "GetRandomMembers:Expired", []string{"Expired1", "Expired2", "Expired3", "Expired4"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMembers:Expired case failed when adding members: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("GetRandomMembers:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			gotMembers, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 2)
			if getRandomMembersErr != nil {
				t.Errorf("GetRandomMembers:Expired case failed when getting random members: %v", getRandomMembersErr.Error())
			}
			if gotMembers == nil {
				t.Errorf("GetRandomMembers:Expired case failed: members is nil")
			}
			if len(gotMembers) != 0 {
				t.Errorf("GetRandomMembers:Expired case failed: incorrect members")
			}
		})

		// 获取超额的样例
		t.Run("GetRandomMembers:OverLimit", func(t *testing.T) {
			key, members := "GetRandomMembers:OverLimit", []string{"OverLimit1", "OverLimit2", "OverLimit3", "OverLimit4"}
			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMembers:OverLimit case failed when adding members: %v", addErr.Error())
			}

			gotMembers, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 5)
			if getRandomMembersErr != nil {
				t.Errorf("GetRandomMembers:OverLimit case failed when getting random members: %v", getRandomMembersErr.Error())
			}
			if gotMembers == nil {
				t.Errorf("GetRandomMembers:OverLimit case failed: members is nil")
			}
			if len(gotMembers) != 4 {
				t.Errorf("GetRandomMembers:OverLimit case failed: incorrect members")
			}
			for _, member := range gotMembers {
				if !containsString(members, member) {
					t.Errorf("GetRandomMembers:OverLimit case failed: incorrect members")
				}
			}
		})

		// 在大量数据中获取少量数据的样例
		t.Run("GetRandomMembers:InLargeData", func(t *testing.T) {
			key := "GetRandomMembers:InLargeData"
			var members []string
			for i := 1; i <= 1010; i++ {
				members = append(members, fmt.Sprintf("InLargeData%d", i))
			}

			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMembers:InLargeData case failed when adding members: %v", addErr.Error())
			}

			gotMembers, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 10)
			if getRandomMembersErr != nil {
				t.Errorf("GetRandomMembers:InLargeData case failed when getting random members: %v", getRandomMembersErr.Error())
			}
			if gotMembers == nil {
				t.Errorf("GetRandomMembers:InLargeData case failed: members is nil")
			}
			if len(gotMembers) != 10 {
				t.Errorf("GetRandomMembers:InLargeData case failed: incorrect members")
			}
			for _, member := range gotMembers {
				if !containsString(members, member) {
					t.Errorf("GetRandomMembers:InLargeData case failed: incorrect members")
				}
			}
		})

		// 在大量数据中获取大量数据的样例
		t.Run("GetRandomMembers:OutLargeData", func(t *testing.T) {
			key := "GetRandomMembers:OutLargeData"
			var members []string
			for i := 1; i <= 1010; i++ {
				members = append(members, fmt.Sprintf("OutLargeData%d", i))
			}

			addErr := impl.AddMembers(context.Background(), key, members...)
			if addErr != nil {
				t.Errorf("GetRandomMembers:OutLargeData case failed when adding members: %v", addErr.Error())
			}

			gotMembers, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 999)
			if getRandomMembersErr != nil {
				t.Errorf("GetRandomMembers:OutLargeData case failed when getting random members: %v", getRandomMembersErr.Error())
			}
			if gotMembers == nil {
				t.Errorf("GetRandomMembers:OutLargeData case failed: members is nil")
			}
			if len(gotMembers) != 999 {
				t.Errorf("GetRandomMembers:OutLargeData case failed: incorrect members")
			}
			for _, member := range gotMembers {
				if !containsString(members, member) {
					t.Errorf("GetRandomMembers:OutLargeData case failed: incorrect members")
				}
			}
		})

		// 获取错误类型的样例
		t.Run("GetRandomMembers:WrongType", func(t *testing.T) {
			key := "GetRandomMembers:WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("GetRandomMembers:WrongType case failed when storing key: %v", storeErr.Error())
			}

			gotMembers, getRandomMembersErr := impl.GetRandomMembers(context.Background(), key, 2)
			if getRandomMembersErr == nil {
				t.Errorf("GetRandomMembers:WrongType case failed: no error")
			}
			if gotMembers == nil {
				t.Errorf("GetRandomMembers:WrongType case failed: members is nil")
			}
			if len(gotMembers) != 0 {
				t.Errorf("GetRandomMembers:WrongType case failed: incorrect members")
			}
		})
	}
}

func HGetValueFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("HGetValue:NotExist", func(t *testing.T) {
			key, field := "HGetValue:NotExist", "NotExist"
			exist, value, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HGetValue:NotExist case failed when getting value: %v", getValueErr.Error())
			}
			if exist {
				t.Errorf("HGetValue:NotExist case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("HGetValue:NotExist case failed: incorrect value")
			}
		})

		// 获取一个存在且未过期的样例
		t.Run("HGetValue:NotExpired", func(t *testing.T) {
			key, field, value := "HGetValue:NotExpired", "NotExpired", "NotExpired"
			addErr := impl.HSetValue(context.Background(), key, field, value)
			if addErr != nil {
				t.Errorf("HGetValue:NotExpired case failed when adding value: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("HGetValue:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HGetValue:NotExpired case failed when getting value: %v", getValueErr.Error())
			}
			if !exist {
				t.Errorf("HGetValue:NotExpired case failed: key not exist, want exist")
			}
			if getValue != value {
				t.Errorf("HGetValue:NotExpired case failed: incorrect value")
			}
		})

		// 获取一个存在但是过期的样例
		t.Run("HGetValue:Expired", func(t *testing.T) {
			key, field, value := "HGetValue:Expired", "Expired", "Expired"
			addErr := impl.HSetValue(context.Background(), key, field, value)
			if addErr != nil {
				t.Errorf("HGetValue:Expired case failed when adding value: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("HGetValue:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HGetValue:Expired case failed when getting value: %v", getValueErr.Error())
			}
			if exist {
				t.Errorf("HGetValue:Expired case failed: key exist, want not exist")
			}
			if getValue != "" {
				t.Errorf("HGetValue:Expired case failed: incorrect value")
			}
		})

		// 获取字段不存在的样例
		t.Run("HGetValue:FieldNotExist", func(t *testing.T) {
			key, field, value, notExistField := "HGetValue:FieldNotExist", "FieldNotExist", "FieldNotExist", "NotExist"
			addErr := impl.HSetValue(context.Background(), key, field, value)
			if addErr != nil {
				t.Errorf("HGetValue:FieldNotExist case failed when adding value: %v", addErr.Error())
			}

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, notExistField)
			if getValueErr != nil {
				t.Errorf("HGetValue:FieldNotExist case failed when getting value: %v", getValueErr.Error())
			}
			if exist {
				t.Errorf("HGetValue:FieldNotExist case failed: key exist, want not exist")
			}
			if getValue != "" {
				t.Errorf("HGetValue:FieldNotExist case failed: incorrect value")
			}
		})

		// 获取错误类型的样例
		t.Run("HGetValue:WrongType", func(t *testing.T) {
			key, field := "HGetValue:WrongType", "WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("HGetValue:WrongType case failed when storing key: %v", storeErr.Error())
			}

			exist, value, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr == nil {
				t.Errorf("HGetValue:WrongType case failed: no error")
			}
			if exist {
				t.Errorf("HGetValue:WrongType case failed: key exist, want not exist")
			}
			if value != "" {
				t.Errorf("HGetValue:WrongType case failed: incorrect value")
			}
		})
	}
}

func HGetValuesFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("HGetValues:NotExist", func(t *testing.T) {
			key := "HGetValues:NotExist"
			values, getValuesErr := impl.HGetValues(context.Background(), key)
			if getValuesErr != nil {
				t.Errorf("HGetValues:NotExist case failed when getting values: %v", getValuesErr.Error())
			}
			if values == nil {
				t.Errorf("HGetValues:NotExist case failed: values is nil")
			}
		})

		// 获取一个存在且未过期的样例
		t.Run("HGetValues:NotExpired", func(t *testing.T) {
			key, hashSet := "HGetValues:NotExpired", map[string]string{"NotExpired1": "NotExpired1", "NotExpired2": "NotExpired2"}
			addErr := impl.HSetValues(context.Background(), key, hashSet)
			if addErr != nil {
				t.Errorf("HGetValues:NotExpired case failed when adding values: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("HGetValues:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, []string{"NotExpired1", "NotExpired2"}...)
			if getValuesErr != nil {
				t.Errorf("HGetValues:NotExpired case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HGetValues:NotExpired case failed: values is nil")
			}
			if len(getValues) != len(hashSet) {
				t.Errorf("HGetValues:NotExpired case failed: incorrect values")
			}
			for field, value := range hashSet {
				if getValues[field] != value {
					t.Errorf("HGetValues:NotExpired case failed: incorrect values")
				}
			}
		})

		// 获取一个存在但是过期的样例
		t.Run("HGetValues:Expired", func(t *testing.T) {
			key, hashSet := "HGetValues:Expired", map[string]string{"Expired1": "Expired1", "Expired2": "Expired2"}
			addErr := impl.HSetValues(context.Background(), key, hashSet)
			if addErr != nil {
				t.Errorf("HGetValues:Expired case failed when adding values: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, expire)
			if expireErr != nil {
				t.Errorf("HGetValues:Expired case failed when setting expire: %v", expireErr.Error())
			}

			time.Sleep(sleepInterval)

			getValues, getValuesErr := impl.HGetValues(context.Background(), key)
			if getValuesErr != nil {
				t.Errorf("HGetValues:Expired case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HGetValues:Expired case failed: values is nil")
			}
			if len(getValues) != 0 {
				t.Errorf("HGetValues:Expired case failed: incorrect values")
			}
		})

		// 获取部分字段不存在的样例
		t.Run("HGetValues:PartFieldNotExist", func(t *testing.T) {
			key, hashSet, getFields := "HGetValues:PartFieldNotExist", map[string]string{"PartFieldNotExist1": "PartFieldNotExist1", "PartFieldNotExist2": "PartFieldNotExist2"}, []string{"PartFieldNotExist1", "PartFieldNotExist2", "NotExist"}
			addErr := impl.HSetValues(context.Background(), key, hashSet)
			if addErr != nil {
				t.Errorf("HGetValues:PartFieldNotExist case failed when adding values: %v", addErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, getFields...)
			t.Log(getValues)
			if getValuesErr != nil {
				t.Errorf("HGetValues:PartFieldNotExist case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HGetValues:PartFieldNotExist case failed: values is nil")
			}
			if len(getValues) != len(getFields) {
				t.Errorf("HGetValues:PartFieldNotExist case failed: incorrect values")
			}
			for field, value := range hashSet {
				if field == getFields[2] && value != "" {
					t.Errorf("HGetValues:PartFieldNotExist case failed: incorrect values")
				}
				if containsString(getFields, field) && getValues[field] != value {
					t.Errorf("HGetValues:PartFieldNotExist case failed: incorrect values")
				}
			}
		})

		// 获取错误类型的样例
		t.Run("HGetValues:WrongType", func(t *testing.T) {
			key := "HGetValues:WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("HGetValues:WrongType case failed when storing key: %v", storeErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, "WrongType")
			if getValuesErr == nil {
				t.Errorf("HGetValues:WrongType case failed: no error")
			}
			if getValues == nil {
				t.Errorf("HGetValues:WrongType case failed: values is nil")
			}
			if len(getValues) != 0 {
				t.Errorf("HGetValues:WrongType case failed: incorrect values")
			}
		})
	}
}

func HGetJsonFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("HGetJson:NotExist", func(t *testing.T) {
			key, field := "HGetJson:NotExist", "NotExist"
			var receiver testStruct
			exist, getJsonErr := impl.HGetJson(context.Background(), key, field, &receiver)
			if getJsonErr != nil {
				t.Errorf("HGetJson:NotExist case failed when getting json: %v", getJsonErr.Error())
			}
			if exist {
				t.Errorf("HGetJson:NotExist case failed: key exist, want not exist")
			}
		})

		// 获取一个正常存在的样例
		t.Run("HGetJson:Normal", func(t *testing.T) {
			key, field := "HGetJson:Normal", "Normal"
			caseStruct := testStruct{Key: key, Value: field}
			storedJson, marshalErr := json.Marshal(caseStruct)
			if marshalErr != nil {
				t.Errorf("HGetJson:Normal case failed when marshaling json: %v", marshalErr.Error())
			}

			addErr := impl.HSetValue(context.Background(), key, field, string(storedJson))
			if addErr != nil {
				t.Errorf("HGetJson:Normal case failed when adding json: %v", addErr.Error())
			}

			var receiver testStruct
			exist, getJsonErr := impl.HGetJson(context.Background(), key, field, &receiver)
			if getJsonErr != nil {
				t.Errorf("HGetJson:Normal case failed when getting json: %v", getJsonErr.Error())
			}
			if !exist {
				t.Errorf("HGetJson:Normal case failed: key not exist, want exist")
			}
			if receiver.Key != caseStruct.Key || receiver.Value != caseStruct.Value {
				t.Errorf("HGetJson:Normal case failed: incorrect json")
			}
		})
	}
}

func HGetAllFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("HGetAll:NotExist", func(t *testing.T) {
			key := "HGetAll:NotExist"
			result, getAllErr := impl.HGetAll(context.Background(), key)
			if getAllErr != nil {
				t.Errorf("HGetAll:NotExist case failed when getting all: %v", getAllErr.Error())
			}
			if result == nil {
				t.Errorf("HGetAll:NotExist case failed: result is nil")
			}
			if len(result) != 0 {
				t.Errorf("HGetAll:NotExist case failed: incorrect result")
			}
		})

		// 获取一个正常存在的样例
		t.Run("HGetAll:Normal", func(t *testing.T) {
			key, hashSet := "HGetAll:Normal", map[string]string{"Normal1": "Normal1", "Normal2": "Normal2"}
			addErr := impl.HSetValues(context.Background(), key, hashSet)
			if addErr != nil {
				t.Errorf("HGetAll:Normal case failed when adding values: %v", addErr.Error())
			}

			result, getAllErr := impl.HGetAll(context.Background(), key)
			if getAllErr != nil {
				t.Errorf("HGetAll:Normal case failed when getting all: %v", getAllErr.Error())
			}
			if result == nil {
				t.Errorf("HGetAll:Normal case failed: result is nil")
			}
			if len(result) != len(hashSet) {
				t.Errorf("HGetAll:Normal case failed: incorrect result")
			}
			for field, value := range hashSet {
				if result[field] != value {
					t.Errorf("HGetAll:Normal case failed: incorrect result")
				}
			}
		})
	}
}

func HGetAllJsonFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 获取一个不存在的样例
		t.Run("HGetAllJson:NotExist", func(t *testing.T) {
			key := "HGetAllJson:NotExist"
			var receiver map[string]testStruct
			getAllJsonErr := impl.HGetAllJson(context.Background(), key, &receiver)
			if getAllJsonErr != nil {
				t.Errorf("HGetAllJson:NotExist case failed when getting all json: %v", getAllJsonErr.Error())
			}
			if receiver == nil {
				t.Errorf("HGetAllJson:NotExist case failed: receiver is nil")
			}
			if len(receiver) != 0 {
				t.Errorf("HGetAllJson:NotExist case failed: incorrect receiver")
			}
		})

		// 获取一个正常存在的样例
		t.Run("HGetAllJson:Normal", func(t *testing.T) {
			key, hashSet := "HGetAllJson:Normal", map[string]string{"key": "HGetAllJson:Normal", "value": "Normal"}
			caseStruct := testStruct{Key: hashSet["key"], Value: hashSet["value"]}
			addErr := impl.HSetValues(context.Background(), key, hashSet)
			if addErr != nil {
				t.Errorf("HGetAllJson:Normal case failed when adding values: %v", addErr.Error())
			}

			var receiver testStruct
			getAllJsonErr := impl.HGetAllJson(context.Background(), key, &receiver)
			if getAllJsonErr != nil {
				t.Errorf("HGetAllJson:Normal case failed when getting all json: %v", getAllJsonErr.Error())
			}
			if receiver.Key != caseStruct.Key || receiver.Value != caseStruct.Value {
				t.Errorf("HGetAllJson:Normal case failed: incorrect receiver")
			}
		})
	}
}

func HSetValueFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 设置一个不存在的样例
		t.Run("HSetValue:NotExist", func(t *testing.T) {
			key, field, value := "HSetValue:NotExist", "NotExist", "NotExist"
			setValueErr := impl.HSetValue(context.Background(), key, field, value)
			if setValueErr != nil {
				t.Errorf("HSetValue:NotExist case failed when setting value: %v", setValueErr.Error())
			}

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HSetValue:NotExist case failed when getting value: %v", getValueErr.Error())
			}
			if !exist {
				t.Errorf("HSetValue:NotExist case failed: key not exist, want exist")
			}
			if getValue != value {
				t.Errorf("HSetValue:NotExist case failed: incorrect value")
			}
		})

		// 设置一个存在且未过期的样例
		t.Run("HSetValue:NotExpired", func(t *testing.T) {
			key, prevField, prevValue, field, value := "HSetValue:NotExpired", "HSetValue:NotExpired:Prev", "HSetValue:NotExpired:Prev", "HSetValue:NotExpired", "HSetValue:NotExpired"
			addErr := impl.HSetValue(context.Background(), key, prevField, prevValue)
			if addErr != nil {
				t.Errorf("HSetValue:NotExpired case failed when adding value: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("HSetValue:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			setValueErr := impl.HSetValue(context.Background(), key, field, value)
			if setValueErr != nil {
				t.Errorf("HSetValue:NotExpired case failed when setting value: %v", setValueErr.Error())
			}

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HSetValue:NotExpired case failed when getting value: %v", getValueErr.Error())
			}
			if !exist {
				t.Errorf("HSetValue:NotExpired case failed: key not exist, want exist")
			}
			if getValue != value {
				t.Errorf("HSetValue:NotExpired case failed: incorrect value")
			}
		})

		// 设置一个类型错误的样例
		t.Run("HSetValue:WrongType", func(t *testing.T) {
			key, field, value := "HSetValue:WrongType", "WrongType", "WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("HSetValue:WrongType case failed when storing key: %v", storeErr.Error())
			}

			setValueErr := impl.HSetValue(context.Background(), key, field, value)
			if setValueErr == nil {
				t.Errorf("HSetValue:WrongType case failed: no error")
			}
		})
	}
}

func HSetValuesFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 设置一个不存在的样例
		t.Run("HSetValues:NotExist", func(t *testing.T) {
			key, hashSet := "HSetValues:NotExist", map[string]string{"NotExist1": "NotExist1", "NotExist2": "NotExist2"}
			setValuesErr := impl.HSetValues(context.Background(), key, hashSet)
			if setValuesErr != nil {
				t.Errorf("HSetValues:NotExist case failed when setting values: %v", setValuesErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, "NotExist1", "NotExist2")
			if getValuesErr != nil {
				t.Errorf("HSetValues:NotExist case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HSetValues:NotExist case failed: values is nil")
			}
			if len(getValues) != len(hashSet) {
				t.Errorf("HSetValues:NotExist case failed: incorrect values")
			}
			for field, value := range hashSet {
				if getValues[field] != value {
					t.Errorf("HSetValues:NotExist case failed: incorrect values")
				}
			}
		})

		// 设置一个存在且未过期的样例
		t.Run("HSetValues:NotExpired", func(t *testing.T) {
			key, prevHashSet, hashSet := "HSetValues:NotExpired", map[string]string{"HSetValues:NotExpired:Prev1": "HSetValues:NotExpired:Prev1", "HSetValues:NotExpired:Prev2": "HSetValues:NotExpired:Prev2"}, map[string]string{"HSetValues:NotExpired1": "HSetValues:NotExpired1", "HSetValues:NotExpired2": "HSetValues:NotExpired2"}
			addErr := impl.HSetValues(context.Background(), key, prevHashSet)
			if addErr != nil {
				t.Errorf("HSetValues:NotExpired case failed when adding values: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("HSetValues:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			setValuesErr := impl.HSetValues(context.Background(), key, hashSet)
			if setValuesErr != nil {
				t.Errorf("HSetValues:NotExpired case failed when setting values: %v", setValuesErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, "HSetValues:NotExpired1", "HSetValues:NotExpired2")
			if getValuesErr != nil {
				t.Errorf("HSetValues:NotExpired case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HSetValues:NotExpired case failed: values is nil")
			}
			if len(getValues) != len(hashSet) {
				t.Errorf("HSetValues:NotExpired case failed: incorrect values")
			}
			for field, value := range hashSet {
				if getValues[field] != value {
					t.Errorf("HSetValues:NotExpired case failed: incorrect values")
				}
			}
		})

		// 设置一个类型错误的样例
		t.Run("HSetValues:WrongType", func(t *testing.T) {
			key, hashSet := "HSetValues:WrongType", map[string]string{"WrongType1": "WrongType1", "WrongType2": "WrongType2"}
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("HSetValues:WrongType case failed when storing key: %v", storeErr.Error())
			}

			setValuesErr := impl.HSetValues(context.Background(), key, hashSet)
			if setValuesErr == nil {
				t.Errorf("HSetValues:WrongType case failed: no error")
			}
		})
	}
}

func HRemoveValueFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 删除一个不存在的样例
		t.Run("HRemoveValue:NotExist", func(t *testing.T) {
			key, field := "HRemoveValue:NotExist", "NotExist"
			removeValueErr := impl.HRemoveValue(context.Background(), key, field)
			if removeValueErr != nil {
				t.Errorf("HRemoveValue:NotExist case failed when removing value: %v", removeValueErr.Error())
			}

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HRemoveValue:NotExist case failed when getting value: %v", getValueErr.Error())
			}
			if exist {
				t.Errorf("HRemoveValue:NotExist case failed: key exist, want not exist")
			}
			if getValue != "" {
				t.Errorf("HRemoveValue:NotExist case failed: incorrect value")
			}
		})

		// 删除一个存在且未过期的样例
		t.Run("HRemoveValue:NotExpired", func(t *testing.T) {
			key, field, value := "HRemoveValue:NotExpired", "HRemoveValue:NotExpired", "HRemoveValue:NotExpired"
			addErr := impl.HSetValue(context.Background(), key, field, value)
			if addErr != nil {
				t.Errorf("HRemoveValue:NotExpired case failed when adding value: %v", addErr.Error())
			}

			expireErr := impl.Expire(context.Background(), key, time.Hour)
			if expireErr != nil {
				t.Errorf("HRemoveValue:NotExpired case failed when setting expire: %v", expireErr.Error())
			}

			removeValueErr := impl.HRemoveValue(context.Background(), key, field)
			if removeValueErr != nil {
				t.Errorf("HRemoveValue:NotExpired case failed when removing value: %v", removeValueErr.Error())
			}

			exist, getValue, getValueErr := impl.HGetValue(context.Background(), key, field)
			if getValueErr != nil {
				t.Errorf("HRemoveValue:NotExpired case failed when getting value: %v", getValueErr.Error())
			}
			if exist {
				t.Errorf("HRemoveValue:NotExpired case failed: key exist, want not exist")
			}
			if getValue != "" {
				t.Errorf("HRemoveValue:NotExpired case failed: incorrect value")
			}
		})

		// 删除一个错误类型的样例
		t.Run("HRemoveValue:WrongType", func(t *testing.T) {
			key, field := "HRemoveValue:WrongType", "WrongType"
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("HRemoveValue:WrongType case failed when storing key: %v", storeErr.Error())
			}

			removeValueErr := impl.HRemoveValue(context.Background(), key, field)
			if removeValueErr == nil {
				t.Errorf("HRemoveValue:WrongType case failed: no error")
			}
		})
	}
}

func HRemoveValuesFunction(impl cache.Cache) func(t *testing.T) {
	return func(t *testing.T) {
		// 删除一个不存在的样例
		t.Run("HRemoveValues:NotExist", func(t *testing.T) {
			key, fields := "HRemoveValues:NotExist", []string{"NotExist1", "NotExist2"}
			removeValuesErr := impl.HRemoveValues(context.Background(), key, fields...)
			if removeValuesErr != nil {
				t.Errorf("HRemoveValues:NotExist case failed when removing values: %v", removeValuesErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, fields...)
			if getValuesErr != nil {
				t.Errorf("HRemoveValues:NotExist case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HRemoveValues:NotExist case failed: values is nil")
			}
			if len(getValues) != len(fields) {
				t.Errorf("HRemoveValues:NotExist case failed: incorrect values")
			}
			for _, field := range fields {
				if getValues[field] != "" {
					t.Errorf("HRemoveValues:NotExist case failed: incorrect values")
				}
			}
		})

		// 删除存在的样例
		t.Run("HRemoveValues:Normal", func(t *testing.T) {
			key, hashSet := "HRemoveValues:Normal", map[string]string{"Normal1": "Normal1", "Normal2": "Normal2", "Normal3": "Normal3"}
			deleteFields := []string{"Normal1", "Normal2"}
			addErr := impl.HSetValues(context.Background(), key, hashSet)
			if addErr != nil {
				t.Errorf("HRemoveValues:Normal case failed when adding values: %v", addErr.Error())
			}

			removeValuesErr := impl.HRemoveValues(context.Background(), key, deleteFields...)
			if removeValuesErr != nil {
				t.Errorf("HRemoveValues:Normal case failed when removing values: %v", removeValuesErr.Error())
			}

			getValues, getValuesErr := impl.HGetValues(context.Background(), key, deleteFields...)
			if getValuesErr != nil {
				t.Errorf("HRemoveValues:Normal case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HRemoveValues:Normal case failed: values is nil")
			}
			if len(getValues) != len(deleteFields) {
				t.Errorf("HRemoveValues:Normal case failed: incorrect values")
			}
			for _, field := range deleteFields {
				if getValues[field] != "" {
					t.Errorf("HRemoveValues:Normal case failed: incorrect values")
				}
			}

			getValues, getValuesErr = impl.HGetValues(context.Background(), key, "Normal3")
			if getValuesErr != nil {
				t.Errorf("HRemoveValues:Normal case failed when getting values: %v", getValuesErr.Error())
			}
			if getValues == nil {
				t.Errorf("HRemoveValues:Normal case failed: values is nil")
			}
			if len(getValues) != 1 {
				t.Errorf("HRemoveValues:Normal case failed: incorrect values")
			}
			if getValues["Normal3"] != hashSet["Normal3"] {
				t.Errorf("HRemoveValues:Normal case failed: incorrect values")
			}
		})

		// 删除一个错误类型的样例
		t.Run("HRemoveValues:WrongType", func(t *testing.T) {
			key, fields := "HRemoveValues:WrongType", []string{"WrongType1", "WrongType2"}
			storeErr := impl.Store(context.Background(), key, "WrongType")
			if storeErr != nil {
				t.Errorf("HRemoveValues:WrongType case failed when storing key: %v", storeErr.Error())
			}

			removeValuesErr := impl.HRemoveValues(context.Background(), key, fields...)
			if removeValuesErr == nil {
				t.Errorf("HRemoveValues:WrongType case failed: no error")
			}
		})
	}
}

func RunTestCase(t *testing.T, impl cache.Cache) {
	for _, i := range BaseCacheUnitTestCaseList {
		t.Run(i.CaseName, i.TestFunction(impl))
	}
}
