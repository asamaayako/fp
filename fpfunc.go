/**
 * Author: asamaayako
 * Date: 2024/4/9
 * Time: 上午9:07
 * Describe:
 */

package fp

import (
	"math/bits"
	"math/rand/v2"
	"strings"
)

// 右折叠
// Foldr:: ((a -> b -> a) -> a) -> [b] -> a
func Foldr[a, b any](f func(a, b) a, acc a) func([]b) a {
	return func(bs []b) a {
		if len(bs) == 0 {
			return acc
		}
		return f(Foldr(f, acc)(bs[1:]), bs[0])
	}
}

// 左折叠
// FoldL:: ((a -> b -> a) -> a) -> [b] -> b
func FoldL[a, b any](f func(a, b) a, acc a) func([]b) a {
	return func(bs []b) a {
		if len(bs) == 0 {
			return acc
		}
		return FoldL(f, f(acc, bs[0]))(bs[1:])
	}
}

// Map :: (a -> b) -> [a] -> [b]
func Map[a, b any](f func(a) b) func([]a) []b {
	return func(as []a) []b {
		return FoldL(func(b1 []b, a1 a) []b { return append(b1, f(a1)) }, make([]b, 0, len(as)))(as)
	}
}

// Filter :: (a -> Bool) -> [a] -> [a]
func Filter[a any](p func(a) bool) func([]a) []a {
	return func(as []a) []a {
		return FoldL(func(a2 []a, a1 a) []a {
			if p(a1) {
				return append(a2, a1)
			} else {
				return a2
			}
		}, make([]a, 0, len(as)))(as)
	}
}

// Compose :: f (function [t]->[t]) t =>[([t]->[t])]->[t]->[t]
func Compose[T any](fnList ...func(...T) []T) func(...T) []T {
	return func(s ...T) []T {
		f := fnList[0]
		nextFnList := fnList[1:]

		if len(fnList) == 1 {
			return f(s...)
		}

		return f(Compose(nextFnList...)(s...)...)
	}
}

// ComposeInterface Compose the functions from right to left (Math: f(g(x)) Compose: Compose(f, g)(x))
func ComposeInterface(fnList ...func(...interface{}) []interface{}) func(...interface{}) []interface{} {
	return Compose(fnList...)
}

// Pipe Pipe the functions from left to right
func Pipe[T any](fnList ...func(...T) []T) func(...T) []T {
	return func(s ...T) []T {
		lastIndex := len(fnList) - 1
		f := fnList[lastIndex]
		nextFnList := fnList[:lastIndex]

		if len(fnList) == 1 {
			return f(s...)
		}

		return f(Pipe(nextFnList...)(s...)...)
	}
}

// PipeInterface Pipe the functions from left to right
func PipeInterface(fnList ...func(...interface{}) []interface{}) func(...interface{}) []interface{} {
	return Pipe(fnList...)
}

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
