package event

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type ChallengeRequest struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}

// decrypt 解密飞书的加密请求
func decrypt(content string, secretKey string) (decrypted []byte, err error) {
	// 解码加密的字符串
	buffer, decodeErr := base64.StdEncoding.DecodeString(content)
	if decodeErr != nil {
		return []byte{}, fmt.Errorf("failed to decode challenge request encrypt: %w", decodeErr)
	}
	if len(buffer) < aes.BlockSize {
		return []byte{}, fmt.Errorf("cipher too short")
	}

	// 使用 SHA256 生成密钥
	keyBuffer := sha256.Sum256([]byte(secretKey))

	// 创建 AES 解密器
	blockData, createCipherErr := aes.NewCipher(keyBuffer[:])
	if createCipherErr != nil {
		return []byte{}, fmt.Errorf("failed to calculate aes cipher: %w", createCipherErr)
	}

	// 检查密文长度是否合法
	if len(buffer)%aes.BlockSize != 0 {
		return []byte{}, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// 解密数据
	iv, buf := buffer[:aes.BlockSize], buffer[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(blockData, iv)
	mode.CryptBlocks(buf, buf)

	// 查找 JSON 对象的起始和终止位置
	n := bytes.Index(buf, []byte("{"))
	m := bytes.LastIndex(buf, []byte("}"))
	if n == -1 || m == -1 {
		return []byte{}, fmt.Errorf("decrypted data does not contain a valid JSON object")
	}

	return buf[n : m+1], nil
}

// ParseChallenge 解析飞书的 challenge 请求，无需解密
// reference: https://open.feishu.cn/document/server-docs/event-subscription-guide/event-subscription-configure-/request-url-configuration-case
func ParseChallenge(requestBody []byte, verificationToken string) (challengeBody ChallengeResponse, err error) {
	var challengeRequest ChallengeRequest
	if err := json.Unmarshal(requestBody, &challengeRequest); err != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to unmarshal challenge request: %w", err)
	}
	if challengeRequest.Type != "url_verification" {
		return ChallengeResponse{}, fmt.Errorf("invalid challenge request type: %s", challengeRequest.Type)
	}
	if challengeRequest.Token != verificationToken {
		return ChallengeResponse{}, fmt.Errorf("invalid challenge request token: %s", challengeRequest.Token)
	}

	return ChallengeResponse{Challenge: challengeRequest.Challenge}, nil
}

type EncryptedRequest struct {
	Encrypt string `json:"encrypt"`
}

// ParseChallengeWithEncryption 解析飞书的 challenge 请求，并进行解密
// reference: https://open.feishu.cn/document/server-docs/event-subscription-guide/event-subscription-configure-/encrypt-key-encryption-configuration-case
func ParseChallengeWithEncryption(requestBody []byte, verificationToken, secretKey string) (challengeBody ChallengeResponse, err error) {
	var encryptedRequest EncryptedRequest

	// 解析加密的请求
	if err := json.Unmarshal(requestBody, &encryptedRequest); err != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to unmarshal challenge request: %w", err)
	} else if encryptedRequest.Encrypt == "" {
		return ChallengeResponse{}, fmt.Errorf("invalid challenge request: empty encrypt")
	}

	// 解密请求
	challenge, decryptErr := decrypt(encryptedRequest.Encrypt, secretKey)
	if decryptErr != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to decrypt challenge request: %w", decryptErr)
	}

	// 解析请求
	var challengeData ChallengeRequest
	if err := json.Unmarshal(challenge, &challengeData); err != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to unmarshal decrypted challenge data: %w", err)
	}
	if challengeData.Token != verificationToken {
		return ChallengeResponse{}, fmt.Errorf("invalid verification token in decrypted challenge data: %s", challengeData.Token)
	}

	// 封装 challenge 响应
	return ChallengeResponse{Challenge: challengeData.Challenge}, nil
}

// ParseEncryptedRequest 解析飞书的回调请求，并进行解密
// reference: https://open.feishu.cn/document/server-docs/event-subscription-guide/event-subscription-configure-/encrypt-key-encryption-configuration-case
func ParseEncryptedRequest(requestBody []byte, secretKey string) (challengeBody CallbackRequest, err error) {
	var encryptedRequest EncryptedRequest
	if err := json.Unmarshal(requestBody, &encryptedRequest); err != nil {
		return CallbackRequest{}, fmt.Errorf("failed to unmarshal challenge request: %w", err)
	} else if encryptedRequest.Encrypt == "" {
		return CallbackRequest{}, fmt.Errorf("invalid challenge request: empty encrypt")
	}

	return ParseEncryptedRequestString(encryptedRequest.Encrypt, secretKey)
}

// ParseEncryptedRequestString 解析飞书的回调请求的encrypted字段，并进行解密
// reference: https://open.feishu.cn/document/server-docs/event-subscription-guide/event-subscription-configure-/encrypt-key-encryption-configuration-case
func ParseEncryptedRequestString(encrypted string, secretKey string) (challengeBody CallbackRequest, err error) {
	challenge, decryptErr := decrypt(encrypted, secretKey)
	if decryptErr != nil {
		return CallbackRequest{}, fmt.Errorf("failed to decrypt challenge request: %w", decryptErr)
	}

	var challengeData CallbackRequest
	if err := json.Unmarshal(challenge, &challengeData); err != nil {
		return CallbackRequest{}, fmt.Errorf("failed to unmarshal decrypted challenge data: %w", err)
	}
	return challengeData, nil
}
