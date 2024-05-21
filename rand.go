package fp

import (
	"math/bits"
	"math/rand/v2"
	"strings"
)

func RoundStringBase64(n int) string {
	const Base64Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-"
	return RoundStringInString(n, Base64Chars)
}
func RoundStringEn(n int) string {
	const EnChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return RoundStringInString(n, EnChars)
}
func RoundStringEnAndNum(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RoundStringInString(n, letterBytes)
}

// 最差的情况 O(2n) 最好的情况O(n)
func RoundStringInString(n int, in string) string {
	var (
		letterIdxBits = bits.Len64(uint64(len(in)))  // 可以表示len(in)的最小位数
		letterIdxMask = uint64(1<<letterIdxBits - 1) // 获取最小位数的掩码 111111
		letterIdxMax  = 64 / letterIdxBits           // 64位数 可以获取10次 6位
	)
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Uint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Uint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(in) {
			sb.WriteByte(in[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
