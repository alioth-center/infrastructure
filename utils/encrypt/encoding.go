package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// AesEncrypt 使用 AES 加密算法加密消息
//   - message: 明文消息
//   - secret: 密钥，长度必须是 16、24 或 32 位
//   - encrypted: 加密后的消息
//   - err: 错误信息
//
// 没有出错重试机制，出错直接返回 err
func AesEncrypt(message string, secret string) (encrypted string, err error) {
	if block, buildAesBlockErr := aes.NewCipher([]byte(secret)); buildAesBlockErr != nil {
		return "", buildAesBlockErr
	} else {
		plainTextBytes := []byte(message)
		cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
		iv := cipherText[:aes.BlockSize]
		if _, buildIvErr := rand.Read(iv); buildIvErr != nil {
			return "", buildIvErr
		}

		cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)
		encryptedString := base64.StdEncoding.EncodeToString(cipherText)

		return encryptedString, nil
	}
}

// AesDecrypt 使用 AES 加密算法解密消息
//   - encrypted: 加密后的消息
//   - secret: 密钥，长度必须是 16、24 或 32 位
//   - decrypted: 解密后的消息
//   - err: 错误信息
//
// 没有出错重试机制，出错直接返回 err
func AesDecrypt(encrypted string, secret string) (decrypted string, err error) {
	if block, buildAesBlockErr := aes.NewCipher([]byte(secret)); buildAesBlockErr != nil {
		return "", buildAesBlockErr
	} else {
		encryptedBytes, decodeErr := base64.StdEncoding.DecodeString(encrypted)
		if decodeErr != nil {
			return "", decodeErr
		}

		if len(encryptedBytes) < aes.BlockSize {
			return "", fmt.Errorf("encrypted message too short")
		}

		iv := encryptedBytes[:aes.BlockSize]
		encryptedBytes = encryptedBytes[aes.BlockSize:]

		cipher.NewCFBDecrypter(block, iv).XORKeyStream(encryptedBytes, encryptedBytes)

		return fmt.Sprintf("%s", encryptedBytes), nil
	}
}

// HashMD5 使用 MD5 算法计算消息的哈希值
func HashMD5(message string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(message)))
}

// HashEntryMD5 使用 MD5 算法计算结构体的哈希值
func HashEntryMD5[T any](entry T) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%#v", entry))))
}
