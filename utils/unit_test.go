package utils

import (
	"math/rand"
	"testing"
	"time"
)

func TestEncoding(t *testing.T) {
	secret := "1234567890123456"
	encrypted, ee := EncryptMessageWithAES("i love u", secret)
	if ee != nil {
		t.Error(ee)
	}

	decrypted, de := DecryptMessageWithAES(encrypted, secret)
	if de != nil {
		t.Error(de)
	}

	if decrypted != "i love u" {
		t.Error("decrypted message not match")
	}

	t.Log("encrypted message:", encrypted, "decrypted message:", decrypted)
}

func TestPasswd(t *testing.T) {
	encoded, ee := BcryptEncode("i love u")
	if ee != nil {
		t.Error(ee)
	}

	if !BcryptCheck("i love u", encoded) {
		t.Error("bcrypt check failed")
	}

	t.Log("encoded message:", encoded)
}

func TestRandom(t *testing.T) {
	// for random check, must set seed
	rand.Seed(0)
	six := GenerateRandomSixDigitNumberCode()
	if six != "795303" {
		t.Error("six digit number code not match")
	}
	t.Log("random six digit number code:", six)

	four := GenerateRandomFourDigitNumberCode()
	if four != "1125" {
		t.Error("four digit number code not match")
	}
	t.Log("random four digit number code:", four)

	prefixSix := GenerateRandomSixDigitNumberCodeWithPrefix("T")
	if prefixSix != "T656761" {
		t.Error("random six digit number code with prefix not match")
	}
	t.Log("random six digit number code with prefix:", prefixSix)

	prefixFour := GenerateRandomFourDigitNumberCodeWithPrefix("T")
	if prefixFour != "T4177" {
		t.Error("random four digit number code with prefix not match")
	}
	t.Log("random four digit number code with prefix:", prefixFour)

	base62 := GenerateRandomBase62(10)
	if base62 != "i2Wa2Tnen9" {
		t.Error("random base62 not match")
	}
	t.Log("random base62:", base62)

	base64 := GenerateRandomBase64(10)
	if base64 != "VI1PEcESnu" {
		t.Error("random base64 not match")
	}
	t.Log("random base64:", base64)

	base62WithPrefix := GenerateRandomBase62WithPrefix("T", 10)
	if base62WithPrefix != "T6yuuf9a2H" {
		t.Error("random base62 with prefix not match")
	}
	t.Log("random base62 with prefix:", base62WithPrefix)

	base64WithPrefix := GenerateRandomBase64WithPrefix("T", 10)
	if base64WithPrefix != "T5LOlIxi8H" {
		t.Error("random base64 with prefix not match")
	}
	t.Log("random base64 with prefix:", base64WithPrefix)
}

type entry struct {
	Name string
	Age  int
}

func TestHashMD5(t *testing.T) {
	t.Run("HashMD5", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Log("hash md5:", HashMD5("i love u"))
		}
	})

	t.Run("HashEntryMD5", func(t *testing.T) {
		entry := &entry{
			Name: "alice",
			Age:  114,
		}

		for i := 0; i < 10; i++ {
			t.Log("hash md5:", HashEntryMD5(entry))
		}
	})
}

func TestTimeZone(t *testing.T) {
	_ = SetLocal(LocationShanghai)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now()
	t.Log("now:", now)
	utc := now.In(time.UTC).Unix()
	t.Log("utc:", time.Unix(utc, 0))
	tky := now.In(loc).Unix()
	t.Log("tokyo:", time.Unix(tky, 0))
	fixednow := TimestampInLocal(utc)
	t.Log("fixed now:", fixednow)
	convertnow := FixedTimestampInLocal(utc, LocationTokyo)
	t.Log("convert now:", convertnow)
}
