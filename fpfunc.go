//go:build go1.23

/**
 * Author: asamaayako
 * Date: 2024/4/9
 * Time: 上午9:07
 * Describe: 配合slices库食用
 */
package fp

import (
	"iter"
	"slices"
)

// Fold:: ((a -> b -> a) -> a) -> [b] -> a
func Fold[a, b any](acc a, f func(a, b) a) func(iter.Seq[b]) a {
	return func(bs iter.Seq[b]) a {
		for bv := range bs {
			acc = f(acc, bv)
		}
		return acc
	}
}

// Head :: [a] -> [a]
func Head[a any]() func(iter.Seq[a]) iter.Seq[a] {
	return Take[a](1)
}

// Take :: Int -> [a] -> [a]
func Take[a any](n int) func(iter.Seq[a]) iter.Seq[a] {
	return func(seq iter.Seq[a]) iter.Seq[a] {
		return func(yield func(a) bool) {
			var count int
			seq(func(v a) bool {
				count++
				if count <= n {
					return yield(v)
				}
				return false
			})
		}
	}
}
func Len[T any](seq iter.Seq[T]) int {
	return Fold(0, func(sum int, _ T) int {
		return sum + 1
	})(seq)
}

// Tail :: [a] -> [a]
func Tail[a any]() func(iter.Seq[a]) iter.Seq[a] {
	return Drop[a](1)
}

// Drop :: Int -> [a] -> [a]
func Drop[a any](n int) func(seq iter.Seq[a]) iter.Seq[a] {
	return func(seq iter.Seq[a]) iter.Seq[a] {
		return func(yield func(a) bool) {
			var count int
			seq(func(v a) bool {
				count++
				if count <= n {
					return true
				}
				return yield(v)
			})
		}
	}
}

// Reduce :: (a -> b -> a) -> a -> [b] -> a
func Reduce[a, b any](f func(a, b) a) func(iter.Seq[b]) a {
	var av a
	return Fold(av, f)
}

// 需要同时从两个序列中获取元素
// Zip :: [a] -> [b] -> [Pair<a,b>]
func Zip[T, U any](seq1 iter.Seq[T], seq2 iter.Seq[U]) iter.Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		next1, stop1 := iter.Pull(seq1)
		defer stop1()
		next2, stop2 := iter.Pull(seq2)
		defer stop2()

		for {
			v1, ok1 := next1()
			v2, ok2 := next2()
			if !ok1 || !ok2 {
				return
			}
			if !yield(Pair[T, U]{v1, v2}) {
				return
			}
		}
	}
}

// UnZip :: [Pair<a,b>] -> [a] -> [b]
func UnZip[Fir, Sec any](seq iter.Seq[Pair[Fir, Sec]]) iter.Seq2[Fir, Sec] {
	return func(yield func(Fir, Sec) bool) {
		seq(func(p Pair[Fir, Sec]) bool {
			return yield(p.First, p.Second)
		})
	}
}

func UnZipSeq2[Fir, Sec any](seq2 iter.Seq2[Fir, Sec]) (iter.Seq[Fir], iter.Seq[Sec]) {
	var firsts []Fir
	var seconds []Sec

	// 提取所有元素到缓存中
	seq2(func(f Fir, s Sec) bool {
		firsts = append(firsts, f)
		seconds = append(seconds, s)
		return true
	})

	return slices.Values(firsts), slices.Values(seconds)
}

// Split :: [a] -> ([a],[a])
func Split[a any](seq iter.Seq[a]) (iter.Seq[a], iter.Seq[a]) {
	seq1 := seq //流可以重复使用 注意作用域即可
	return seq1, seq
}

// 将序列分块处理
// Chunk:: int->[a]->[[a]]
func Chunk[T any](size int) func(seq iter.Seq[T]) iter.Seq[[]T] {
	return func(seq iter.Seq[T]) iter.Seq[[]T] {
		return func(yield func([]T) bool) {
			var chunk = make([]T, 0, size)
			seq(func(i T) bool {
				chunk = append(chunk, i)
				if len(chunk) == size {
					res := yield(chunk)
					chunk = make([]T, 0)
					return res
				}
				return true
			})
			if len(chunk) > 0 {
				yield(chunk)
			}
		}
	}

}

// 创建滑动窗口
// SlidingWindow:: int -> [T]->[[T]]
// used read only
func SlidingWindow[T any](n int) func(iter.Seq[T]) iter.Seq[[]T] {
	return func(seq iter.Seq[T]) iter.Seq[[]T] {
		return func(yield func([]T) bool) {
			var window []T = make([]T, 0, n)

			seq(func(item T) bool {
				window = append(window, item)
				if len(window) > n {
					window = window[1:]
				}
				//积累窗口到达n时 产出一个窗口
				if len(window) == n {
					windowCopy := make([]T, len(window))
					copy(windowCopy, window)
					if !yield(windowCopy) {
						return false
					}
				}
				return true
			})
		}
	}
}

// 当满足某个条件时才继续迭代
// TakeWhile:: (T->bool)->[T]->[T]
func TakeWhile[T any](predicate func(T) bool) func(iter.Seq[T]) iter.Seq[T] {
	return func(seq iter.Seq[T]) iter.Seq[T] {
		return func(yield func(T) bool) {
			seq(func(item T) bool {
				if !predicate(item) {
					return false
				}
				return yield(item)
			})
		}
	}
}

// 处理相邻元素对
// AdjacentPairs:: [T]->[Pair(T,T)]
func AdjacentPairs[T any](seq iter.Seq[T]) iter.Seq[[2]T] {
	return func(yield func([2]T) bool) {
		var prev T
		var hasPrev bool
		seq(func(item T) bool {
			if !hasPrev {
				prev = item
				hasPrev = true
				return true
			}
			if !yield([2]T{prev, item}) {
				return false
			}
			prev = item
			return true
		})
	}
}

func Pairs[T any](seq iter.Seq[T]) iter.Seq[[2]T] {
	return func(yield func([2]T) bool) {
		chunks := Chunk[T](2)(seq)
		chunks(func(item []T) bool {
			return yield([2]T{item[0], item[1]})
		})
	}
}

// Map :: (a -> b) -> [a] -> [b]
func Map[a, b any](f func(a) b) func(iter.Seq[a]) iter.Seq[b] {
	return func(as iter.Seq[a]) iter.Seq[b] {
		return func(yield func(b) bool) {
			as(func(a a) bool {
				return yield(f(a))
			})
		}
	}
}

// Filter :: (a -> Bool) -> [a] -> [a]
func Filter[a any](p func(a) bool) func(iter.Seq[a]) iter.Seq[a] {
	return func(as iter.Seq[a]) iter.Seq[a] {
		return func(yield func(a) bool) {
			as(func(av a) bool {
				if p(av) {
					yield(av)
				}
				return true
			})
			return
		}
	}
}

// Compose :: ([b]->[c])->([a]->[b])->[t]->[t]
func Compose[a, b, c any](f func(iter.Seq[b]) iter.Seq[c], g func(iter.Seq[a]) iter.Seq[b]) func(iter.Seq[a]) iter.Seq[c] {
	return func(as iter.Seq[a]) iter.Seq[c] {
		return f(g(as))
	}
}

// Pipe Pipe the functions from left to right
func Pipe[a, b, c any](f func(iter.Seq[a]) iter.Seq[b], g func(iter.Seq[b]) iter.Seq[c]) func(iter.Seq[a]) iter.Seq[c] {
	return func(as iter.Seq[a]) iter.Seq[c] {
		return g(f(as))
	}
}
