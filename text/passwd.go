package text

import "golang.org/x/crypto/bcrypt"

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
