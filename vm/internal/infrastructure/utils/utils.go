package utils

import "regexp"

// VerifyEmail
// @Description: 验证邮箱格式
// @param        email string
// @return       bool
// @Author tianjiajie 2025-03-20 16:19:24
func VerifyEmail(email string) bool {
	// 正则表达式验证邮箱格式
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	re, err := regexp.Compile(emailRegex)
	if err != nil {
		// 如果正则表达式有语法错误，返回 false
		return false
	}

	// 使用正则表达式匹配邮箱
	return re.MatchString(email)
}
