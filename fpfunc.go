/**
 * Author: asamaayako
 * Date: 2024/4/9
 * Time: 上午9:07
 * Describe:
 */

package fp

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
