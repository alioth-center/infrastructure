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

type ChallengeWithEncryptionRequest struct {
	Encrypt string `json:"encrypt"`
}

// ParseChallengeWithEncryption 解析飞书的 challenge 请求，并进行解密
// reference: https://open.feishu.cn/document/server-docs/event-subscription-guide/event-subscription-configure-/encrypt-key-encryption-configuration-case
func ParseChallengeWithEncryption(requestBody []byte, verificationToken, secretKey string) (challengeBody ChallengeResponse, err error) {
	var challengeRequest ChallengeWithEncryptionRequest
	if unmarshalErr := json.Unmarshal(requestBody, &challengeRequest); unmarshalErr != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to unmarshal challenge request: %w", unmarshalErr)
	}
	if challengeRequest.Encrypt == "" {
		return ChallengeResponse{}, fmt.Errorf("invalid challenge request encrypt: %s", challengeRequest.Encrypt)
	}

	// 解码加密的字符串
	buffer, decodeErr := base64.StdEncoding.DecodeString(challengeRequest.Encrypt)
	if decodeErr != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to decode challenge request encrypt: %w", decodeErr)
	}
	if len(buffer) < aes.BlockSize {
		return ChallengeResponse{}, fmt.Errorf("cipher too short")
	}

	// 使用 SHA256 生成密钥
	keyBuffer := sha256.Sum256([]byte(secretKey))

	// 创建 AES 解密器
	blockData, createCipherErr := aes.NewCipher(keyBuffer[:])
	if createCipherErr != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to calculate aes cipher: %w", createCipherErr)
	}

	// 检查密文长度是否合法
	if len(buffer)%aes.BlockSize != 0 {
		return ChallengeResponse{}, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// 解密数据
	iv, buf := buffer[:aes.BlockSize], buffer[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(blockData, iv)
	mode.CryptBlocks(buf, buf)

	// 查找 JSON 对象的起始和终止位置
	n := bytes.Index(buf, []byte("{"))
	m := bytes.LastIndex(buf, []byte("}"))
	if n == -1 || m == -1 {
		return ChallengeResponse{}, fmt.Errorf("decrypted data does not contain a valid JSON object")
	}

	// 提取 JSON 对象并验证 token
	challenge := buf[n : m+1]
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

type ChallengeAdaptionRequest struct {
	Encrypt   string `json:"encrypt,omitempty"`
	Challenge string `json:"challenge,omitempty"`
	Token     string `json:"token,omitempty"`
	Type      string `json:"type,omitempty"`
}

// ParseChallengeAdaption 解析飞书的 challenge 请求，自适应解密
func ParseChallengeAdaption(requestBody []byte, verificationToken, secretKey string) (challengeBody ChallengeResponse, err error) {
	fullRequest := ChallengeAdaptionRequest{}
	if unmarshalErr := json.Unmarshal(requestBody, &fullRequest); unmarshalErr != nil {
		return ChallengeResponse{}, fmt.Errorf("failed to unmarshal challenge request: %w", unmarshalErr)
	}

	if fullRequest.Encrypt != "" {
		// 使用加密方式解析
		return ParseChallengeWithEncryption(requestBody, verificationToken, secretKey)
	} else {
		// 使用非加密方式解析
		return ParseChallenge(requestBody, verificationToken)
	}
}
