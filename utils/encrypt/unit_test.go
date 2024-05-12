package encrypt

import (
	"math/rand"
	"testing"

	"github.com/alioth-center/infrastructure/utils/generate"
)

func TestEncoding(t *testing.T) {
	secret := "1234567890123456"
	encrypted, ee := AesEncrypt("i love u", secret)
	if ee != nil {
		t.Error(ee)
	}

	decrypted, de := AesDecrypt(encrypted, secret)
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
	six := generate.RandomSixDigitNumberCode()
	if six != "795303" {
		t.Error("six digit number code not match")
	}
	t.Log("random six digit number code:", six)

	four := generate.RandomFourDigitNumberCode()
	if four != "1125" {
		t.Error("four digit number code not match")
	}
	t.Log("random four digit number code:", four)

	prefixSix := generate.RandomSixDigitNumberCodeWithPrefix("T")
	if prefixSix != "T656761" {
		t.Error("random six digit number code with prefix not match")
	}
	t.Log("random six digit number code with prefix:", prefixSix)

	prefixFour := generate.RandomFourDigitNumberCodeWithPrefix("T")
	if prefixFour != "T4177" {
		t.Error("random four digit number code with prefix not match")
	}
	t.Log("random four digit number code with prefix:", prefixFour)

	base62 := generate.RandomBase62(10)
	if base62 != "i2Wa2Tnen9" {
		t.Error("random base62 not match")
	}
	t.Log("random base62:", base62)

	base64 := generate.RandomBase64(10)
	if base64 != "VI1PEcESnu" {
		t.Error("random base64 not match")
	}
	t.Log("random base64:", base64)

	base62WithPrefix := generate.RandomBase62WithPrefix("T", 10)
	if base62WithPrefix != "T6yuuf9a2H" {
		t.Error("random base62 with prefix not match")
	}
	t.Log("random base62 with prefix:", base62WithPrefix)

	base64WithPrefix := generate.RandomBase64WithPrefix("T", 10)
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

func TestRsa(t *testing.T) {
	t.Run("RsaEncrypt:Success", func(t *testing.T) {
		pri, pub, err := RsaKeyGenerate(256)
		if err != nil {
			t.Error(err)
		}

		t.Log("private key:", pri)
		t.Log("public key:", pub)

		encrypted, ee := RsaEncrypt(pub, "i love u")
		if ee != nil {
			t.Error(ee)
		}

		decrypted, de := RsaDecrypt(pri, encrypted)
		if de != nil {
			t.Error(de)
		}

		if decrypted != "i love u" {
			t.Error("decrypted message not match")
		}
	})

	t.Run("RsaEncrypt:NoKey", func(t *testing.T) {
		pri, pub, err := RsaKeyGenerate(256)
		if err != nil {
			t.Error(err)
		}

		t.Log("private key:", pri)
		t.Log("public key:", pub)

		encrypted, ee := RsaEncrypt("", "i love u")
		if ee == nil {
			t.Error("want error, but got nil")
		}

		_, de := RsaDecrypt("", encrypted)
		if de == nil {
			t.Error("want error, but got nil")
		}
	})
}
