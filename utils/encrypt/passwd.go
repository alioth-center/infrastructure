package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// BcryptEncode 使用 bcrypt 加密算法加密密码
//   - password: 明文密码
//   - encoded: 加密后的密码
//   - err: 错误信息
//
// 没有出错重试机制，出错直接返回 err
func BcryptEncode(password string) (encoded string, err error) {
	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), hashErr
}

// BcryptCheck 使用 bcrypt 加密算法校验密码
//   - password: 明文密码
//   - hashed: 加密后的密码
//   - bool: 校验结果
func BcryptCheck(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

// RsaKeyGenerate generates RSA private key and public key.
//   - length: length of the key, minimum 386, recommended 2048 and above
//   - privateKey: RSA private key, pem format
//   - publicKey: RSA public key, pem format
//   - err: error message
//
// example:
//
//	privateKey, publicKey, err := RsaKeyGenerate(2048)
func RsaKeyGenerate(length int) (privateKey, publicKey string, err error) {
	privKey, errGen := rsa.GenerateKey(rand.Reader, length)
	if errGen != nil {
		return "", "", errGen
	}

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	pubKeyBytes, errMarshal := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if errMarshal != nil {
		return "", "", errMarshal
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(privKeyPEM), string(pubKeyPEM), nil
}

// RsaDecrypt decrypts encrypted data using the provided RSA private key.
//   - privateKey: RSA private key, pem format
//   - encrypted: encrypted data, pem format
//   - message: decrypted data
//   - err: error message
//
// example:
//
//	privateKey, _, generateKeyErr := RsaKeyGenerate(2048)
//	if generateKeyErr != nil {
//	    t.Error(generateKeyErr)
//	}
//
//	encrypted, encryptErr := RsaEncrypt(publicKey, "i love u")
//	if encryptErr != nil {
//	    t.Error(encryptErr)
//	}
//
//	decrypted, decryptErr := RsaDecrypt(privateKey, encrypted)
//	if decryptErr != nil {
//	    t.Error(decryptErr)
//	}
func RsaDecrypt(privateKey, encrypted string) (message string, err error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the key")
	}

	privKey, errParse := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errParse != nil {
		return "", errParse
	}

	encryptedBytes, _ := pem.Decode([]byte(encrypted))
	if encryptedBytes == nil {
		return "", errors.New("failed to parse PEM block containing encrypted data")
	}

	decryptedBytes, errDecrypt := rsa.DecryptPKCS1v15(rand.Reader, privKey, encryptedBytes.Bytes)
	if errDecrypt != nil {
		return "", errDecrypt
	}

	return string(decryptedBytes), nil
}

// RsaEncrypt encrypts the given message using the RSA public key.
//   - publicKey: RSA public key, pem format
//   - message: message to be encrypted
//   - encrypted: encrypted data, pem format
//   - err: error message
//
// example:
//
//	privateKey, publicKey, generateKeyErr := RsaKeyGenerate(2048)
//	if generateKeyErr != nil {
//	    t.Error(generateKeyErr)
//	}
//
//	encrypted, encryptErr := RsaEncrypt(publicKey, "i love u")
//	if encryptErr != nil {
//	    t.Error(encryptErr)
//	}
func RsaEncrypt(publicKey, message string) (encrypted string, err error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the key")
	}

	pubKey, errParse := x509.ParsePKIXPublicKey(block.Bytes)
	if errParse != nil {
		return "", errParse
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("key type is not RSA")
	}

	encryptedBytes, errEncrypt := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(message))
	if errEncrypt != nil {
		return "", errEncrypt
	}

	encryptedPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA ENCRYPTED",
		Bytes: encryptedBytes,
	})

	return string(encryptedPEM), nil
}
