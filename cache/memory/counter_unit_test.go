package memory

import (
	"sync"
	"testing"
	"time"

	"github.com/alioth-center/infrastructure/cache"
)

var BaseCounterUnitTestCaseList = []TestCase[cache.Counter]{
	{
		CaseName:     "Increase",
		TestFunction: IncreaseFunction,
	},
	{
		CaseName:     "IncreaseWithExpireWhenNotExist",
		TestFunction: IncreaseWithExpireWhenNOtExistFunction,
	},
	{
		CaseName:     "SetExpire",
		TestFunction: SetExpireFunction,
	},
	{
		CaseName:     "SetExpireWhenNotSet",
		TestFunction: SetExpireWhenNotSetFunction,
	},
	{
		CaseName:     "ExpireImmediately",
		TestFunction: ExpireImmediatelyFunction,
	},
}

func IncreaseFunction(impl cache.Counter) func(t *testing.T) {
	return func(t *testing.T) {
		// 增加一个不存在的计数器
		t.Run("Increase:NotExists", func(t *testing.T) {
			key := "Increase:NotExists"
			result := impl.Increase(nil, key, 1)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("Increase:NotExists failed, incorrect result code")
			}
			if result.GetValue() != 1 {
				t.Errorf("Increase:NotExists failed, incorrect result value")
			}
		})

		// 增加一个存在的计数器
		t.Run("Increase:Exists", func(t *testing.T) {
			key := "Increase:Exists"
			impl.Increase(nil, key, 1)
			result := impl.Increase(nil, key, 1)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("Increase:Exists failed, incorrect result code")
			}
			if result.GetValue() != 2 {
				t.Errorf("Increase:Exists failed, incorrect result value")
			}
		})

		// 增加一个非计数器的样例
		t.Run("Increase:NotCounter", func(t *testing.T) {
			key := "Increase:NotCounter"
			storeErr := impl.(*accessor).Store(nil, key, "NotCounter")
			if storeErr != nil {
				t.Errorf("Increase:NotCounter failed, store error: %v", storeErr)
			}

			result := impl.Increase(nil, key, 1)
			if result != cache.CounterResultEnumFailed {
				t.Errorf("Increase:NotCounter failed, incorrect result code")
			}
		})

		// 并发增加计数器
		t.Run("Increase:Concurrent", func(t *testing.T) {
			key, concurrentNum := "Increase:Concurrent", 1000

			impl.Increase(nil, key, 1)

			wg := sync.WaitGroup{}
			for i := 0; i < concurrentNum; i++ {
				wg.Add(1)
				go func() {
					impl.Increase(nil, key, 1)
					wg.Done()
				}()
			}

			wg.Wait()
			result := impl.Increase(nil, key, 1)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("Increase:Concurrent failed, incorrect result code")
			}
			if result.GetValue() != (concurrentNum + 2) {
				t.Errorf("Increase:Concurrent failed, incorrect result value: %d", result.GetValue())
			}
		})
	}
}

func IncreaseWithExpireWhenNOtExistFunction(impl cache.Counter) func(t *testing.T) {
	return func(t *testing.T) {
		// 增加一个不存在的计数器
		t.Run("IncreaseWithExpireWhenNotExist:NotExists", func(t *testing.T) {
			key, expire := "IncreaseWithExpireWhenNotExist:NotExists", time.Second
			result := impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotExists failed, incorrect result code")
			}
			if result.GetValue() != 1 {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotExists failed, incorrect result value")
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotExists failed, get expired time error: %v", getErr)
			}
			if !exist {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotExists failed, incorrect exist value")
			}
			if !expiredTimeIsCorrect(time.Until(expired), expire) {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotExists failed, incorrect expired time")
			}
		})

		// 增加一个存在的计数器
		t.Run("IncreaseWithExpireWhenNotExist:Exists", func(t *testing.T) {
			key, expire := "IncreaseWithExpireWhenNotExist:Exists", time.Second
			impl.Increase(nil, key, 1)
			result := impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("IncreaseWithExpireWhenNotExist:Exists failed, incorrect result code")
			}

			if result.GetValue() != 2 {
				t.Errorf("IncreaseWithExpireWhenNotExist:Exists failed, incorrect result value")
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("IncreaseWithExpireWhenNotExist:Exists failed, get expired time error: %v", getErr)
			}
			if !exist {
				t.Errorf("IncreaseWithExpireWhenNotExist:Exists failed, incorrect exist value")
			}
			if !expired.IsZero() {
				t.Errorf("IncreaseWithExpireWhenNotExist:Exists failed, incorrect expired time")
			}
		})

		// 增加一个非计数器的样例
		t.Run("IncreaseWithExpireWhenNotExist:NotCounter", func(t *testing.T) {
			key, expire := "IncreaseWithExpireWhenNotExist:NotCounter", time.Second
			storeErr := impl.(*accessor).Store(nil, key, "NotCounter")
			if storeErr != nil {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotCounter failed, store error: %v", storeErr)
			}

			result := impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire)
			if result != cache.CounterResultEnumFailed {
				t.Errorf("IncreaseWithExpireWhenNotExist:NotCounter failed, incorrect result code")
			}
		})

		// 并发增加计数器
		t.Run("IncreaseWithExpireWhenNotExist:Concurrent", func(t *testing.T) {
			key, expire, concurrentNum := "IncreaseWithExpireWhenNotExist:Concurrent", time.Second, 1000

			impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire*10)

			wg := sync.WaitGroup{}
			for i := 0; i < concurrentNum; i++ {
				wg.Add(1)
				go func() {
					impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire)
					wg.Done()
				}()
			}

			wg.Wait()
			result := impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("IncreaseWithExpireWhenNotExist:Concurrent failed, incorrect result code")
			}
			if result.GetValue() != (concurrentNum + 2) {
				t.Errorf("IncreaseWithExpireWhenNotExist:Concurrent failed, incorrect result value: %d", result.GetValue())
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("IncreaseWithExpireWhenNotExist:Concurrent failed, get expired time error: %v", getErr)
			}
			if !exist {
				t.Errorf("IncreaseWithExpireWhenNotExist:Concurrent failed, incorrect exist value")
			}
			if !expiredTimeIsCorrect(time.Until(expired), expire*10) {
				t.Errorf("IncreaseWithExpireWhenNotExist:Concurrent failed, incorrect expired time")
			}
		})
	}
}

func SetExpireFunction(impl cache.Counter) func(t *testing.T) {
	return func(t *testing.T) {
		// 设置一个不存在的计数器
		t.Run("SetExpire:NotExists", func(t *testing.T) {
			key, expire := "SetExpire:NotExists", time.Second
			result := impl.SetExpire(nil, key, expire)
			if result != cache.CounterResultEnumNotEffective {
				t.Errorf("SetExpire:NotExists failed, incorrect result code")
			}
		})

		// 设置一个存在的计数器
		t.Run("SetExpire:Exists", func(t *testing.T) {
			key, expire := "SetExpire:Exists", time.Second
			impl.Increase(nil, key, 1)
			result := impl.SetExpire(nil, key, expire)
			if result != cache.CounterResultEnumSuccess {
				t.Errorf("SetExpire:Exists failed, incorrect result code")
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("SetExpire:Exists failed, get expired time error: %v", getErr)
			}
			if !exist {
				t.Errorf("SetExpire:Exists failed, incorrect exist value")
			}
			if !expiredTimeIsCorrect(time.Until(expired), expire) {
				t.Errorf("SetExpire:Exists failed, incorrect expired time")
			}
		})

		// 设置一个非计数器的样例
		t.Run("SetExpire:NotCounter", func(t *testing.T) {
			key, expire := "SetExpire:NotCounter", time.Second
			storeErr := impl.(*accessor).Store(nil, key, "NotCounter")
			if storeErr != nil {
				t.Errorf("SetExpire:NotCounter failed, store error: %v", storeErr)
			}

			result := impl.SetExpire(nil, key, expire)
			if result != cache.CounterResultEnumFailed {
				t.Errorf("SetExpire:NotCounter failed, incorrect result code: %v", result)
			}
		})
	}
}

func SetExpireWhenNotSetFunction(impl cache.Counter) func(t *testing.T) {
	return func(t *testing.T) {
		// 设置一个不存在的计数器
		t.Run("SetExpireWhenNotSet:NotExists", func(t *testing.T) {
			key, expire := "SetExpireWhenNotSet:NotExists", time.Second
			result := impl.SetExpireWhenNotSet(nil, key, expire)
			if result != cache.CounterResultEnumNotEffective {
				t.Errorf("SetExpireWhenNotSet:NotExists failed, incorrect result code")
			}
		})

		// 设置一个存在的计数器
		t.Run("SetExpireWhenNotSet:Exists", func(t *testing.T) {
			key, expire := "SetExpireWhenNotSet:Exists", time.Second
			impl.Increase(nil, key, 1)
			result := impl.SetExpireWhenNotSet(nil, key, expire)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("SetExpireWhenNotSet:Exists failed, incorrect result code")
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("SetExpireWhenNotSet:Exists failed, get expired time error: %v", getErr)
			}
			if !exist {
				t.Errorf("SetExpireWhenNotSet:Exists failed, incorrect exist value")
			}
			if !expiredTimeIsCorrect(time.Until(expired), expire) {
				t.Errorf("SetExpireWhenNotSet:Exists failed, incorrect expired time: %v", time.Until(expired))
			}
		})

		// 设置一个已经存在且有过期时间的计数器
		t.Run("SetExpireWhenNotSet:ExistsWithExpired", func(t *testing.T) {
			key, expire1, expire2 := "SetExpireWhenNotSet:ExistsWithExpired", time.Minute, time.Second
			impl.IncreaseWithExpireWhenNotExist(nil, key, 1, expire1)
			result := impl.SetExpireWhenNotSet(nil, key, expire2)
			if result != cache.CounterResultEnumNotEffective {
				t.Errorf("SetExpireWhenNotSet:ExistsWithExpired failed, incorrect result code: %v", result)
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("SetExpireWhenNotSet:ExistsWithExpired failed, get expired time error: %v", getErr)
			}
			if !exist {
				t.Errorf("SetExpireWhenNotSet:ExistsWithExpired failed, incorrect exist value")
			}
			if !expiredTimeIsCorrect(time.Until(expired), expire1) {
				t.Errorf("SetExpireWhenNotSet:ExistsWithExpired failed, incorrect expired time")
			}
		})

		// 设置一个非计数器的样例
		t.Run("SetExpireWhenNotSet:NotCounter", func(t *testing.T) {
			key, expire := "SetExpireWhenNotSet:NotCounter", time.Second
			storeErr := impl.(*accessor).Store(nil, key, "NotCounter")
			if storeErr != nil {
				t.Errorf("SetExpireWhenNotSet:NotCounter failed, store error: %v", storeErr)
			}

			result := impl.SetExpireWhenNotSet(nil, key, expire)
			if result != cache.CounterResultEnumFailed {
				t.Errorf("SetExpireWhenNotSet:NotCounter failed, incorrect result code: %v", result)
			}
		})
	}
}

func ExpireImmediatelyFunction(impl cache.Counter) func(t *testing.T) {
	return func(t *testing.T) {
		// 设置一个不存在的计数器
		t.Run("ExpireImmediately:NotExists", func(t *testing.T) {
			key := "ExpireImmediately:NotExists"
			result := impl.ExpireImmediately(nil, key)
			if result != cache.CounterResultEnumNotEffective {
				t.Errorf("ExpireImmediately:NotExists failed, incorrect result code")
			}
		})

		// 设置一个存在的计数器
		t.Run("ExpireImmediately:Exists", func(t *testing.T) {
			key := "ExpireImmediately:Exists"
			impl.Increase(nil, key, 1)
			result := impl.ExpireImmediately(nil, key)
			if result == cache.CounterResultEnumFailed || result == cache.CounterResultEnumNotEffective {
				t.Errorf("ExpireImmediately:Exists failed, incorrect result code")
			}

			exist, expired, getErr := impl.(*accessor).GetExpiredTime(nil, key)
			if getErr != nil {
				t.Errorf("ExpireImmediately:Exists failed, get expired time error: %v", getErr)
			}
			if exist {
				t.Errorf("ExpireImmediately:Exists failed, incorrect exist value")
			}
			if !expired.IsZero() {
				t.Errorf("ExpireImmediately:Exists failed, incorrect expired time")
			}
		})

		// 设置一个非计数器的样例
		t.Run("ExpireImmediately:NotCounter", func(t *testing.T) {
			key := "ExpireImmediately:NotCounter"
			storeErr := impl.(*accessor).Store(nil, key, "NotCounter")
			if storeErr != nil {
				t.Errorf("ExpireImmediately:NotCounter failed, store error: %v", storeErr)
			}

			result := impl.ExpireImmediately(nil, key)
			if result != cache.CounterResultEnumSuccess {
				t.Errorf("ExpireImmediately:NotCounter failed, incorrect result code: %v", result)
			}
		})
	}
}

func RunCounterTestCases(t *testing.T, impl cache.Counter) {
	for _, v := range BaseCounterUnitTestCaseList {
		t.Run(v.CaseName, v.TestFunction(impl))
	}
}
