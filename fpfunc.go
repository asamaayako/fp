/**
 * Author: asamaayako
 * Date: 2024/4/9
 * Time: 上午9:07
 * Describe: 配合slices库食用
 */
package fp

type Seq[T any] = seq[T]
type Seq2[K, V any] = seq2[K, V]
type Pair[T any, U any] struct {
	First  T
	Second U
}

// ==================== 公共门面函数（调用内部实现）====================

// Fold :: (a -> b -> a) -> a -> [b] -> a
//
//go:inline
func Fold[a, b any](acc a, f func(a, b) a) func(Seq[b]) a {
	return fold(acc, f)
}

// Zip :: [a] -> [b] -> [Pair<a,b>]
//
//go:inline
func Zip[T, U any](seq1 Seq[T], seq2 Seq[U]) Seq[Pair[T, U]] {
	return zip(seq1, seq2)
}

// UnZipSeq2 :: Seq2 a b -> ([a], [b])
//
//go:inline
func UnZipSeq2[Fir, Sec any](seq2 Seq2[Fir, Sec]) (Seq[Fir], Seq[Sec]) {
	return unZipSeq2(seq2)
}

// Collect :: [a] -> []a
//
//go:inline
func Collect[a any]() func(Seq[a]) []a {
	return collect[a]()
}

// FromSlice :: []a -> [a]
//
//go:inline
func FromSlice[a any](slice []a) Seq[a] {
	return fromSlice(slice)
}

// FromMap :: Map k v -> Seq2 k v
//
//go:inline
func FromMap[k comparable, v any](m map[k]v) Seq2[k, v] {
	return fromMap(m)
}

// MapKeys :: Map k v -> [k]
//
//go:inline
func MapKeys[k comparable, v any](m map[k]v) Seq[k] {
	return mapKeys(m)
}

// MapValues :: Map k v -> [v]
//
//go:inline
func MapValues[k comparable, v any](m map[k]v) Seq[v] {
	return mapValues(m)
}

// RepeatOne :: a -> [a]
//
//go:inline
func RepeatOne[a any](value a) Seq[a] {
	return func(yield func(a) bool) {
		for {
			if !yield(value) {
				return
			}
		}
	}
}

// ==================== 基础操作函数 ====================

// Head :: [a] -> [a]
//
//go:inline
func Head[a any]() func(Seq[a]) Seq[a] {
	return Take[a](1)
}

// Take :: Int -> [a] -> [a]
//
//go:inline
func Take[a any](n int) func(Seq[a]) Seq[a] {
	return func(seq Seq[a]) Seq[a] {
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

//go:inline
func Len[T any](seq Seq[T]) int {
	return Fold(0, func(sum int, _ T) int {
		return sum + 1
	})(seq)
}

// Tail :: [a] -> [a]
//
//go:inline
func Tail[a any]() func(Seq[a]) Seq[a] {
	return Drop[a](1)
}

// Drop :: Int -> [a] -> [a]
//
//go:inline
func Drop[a any](n int) func(seq Seq[a]) Seq[a] {
	return func(seq Seq[a]) Seq[a] {
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
//
//go:inline
func Reduce[a any](f func(a, a) a) func(Seq[a]) a {
	var result a
	var initialized bool
	return Fold(result, func(acc a, bv a) a {
		if !initialized {
			initialized = true
			return bv
		}
		return f(acc, bv)
	})
}

// UnZip :: [Pair<a,b>] -> [a] -> [b]
//
//go:inline
func UnZip[Fir, Sec any](seq Seq[Pair[Fir, Sec]]) Seq2[Fir, Sec] {
	return func(yield func(Fir, Sec) bool) {
		seq(func(p Pair[Fir, Sec]) bool {
			return yield(p.First, p.Second)
		})
	}
}

// Split :: [a] -> ([a],[a])
//
//go:inline
func Split[a any](seq Seq[a]) (Seq[a], Seq[a]) {
	seq1 := seq //流可以重复使用 注意作用域即可
	return seq1, seq
}

// 将序列分块处理
// Chunk:: int->[a]->[[a]]
//
//go:inline
func Chunk[T any](size int) func(seq Seq[T]) Seq[[]T] {
	return func(seq Seq[T]) Seq[[]T] {
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
//
//go:inline
func SlidingWindow[T any](n int) func(Seq[T]) Seq[[]T] {
	return func(seq Seq[T]) Seq[[]T] {
		return func(yield func([]T) bool) {
			var window = make([]T, 0, n)

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
//
//go:inline
func TakeWhile[T any](predicate func(T) bool) func(Seq[T]) Seq[T] {
	return func(seq Seq[T]) Seq[T] {
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
//
//go:inline
func AdjacentPairs[T any](seq Seq[T]) Seq[[2]T] {
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

//go:inline
func Pairs[T any](seq Seq[T]) Seq[[2]T] {
	return func(yield func([2]T) bool) {
		chunks := Chunk[T](2)(seq)
		chunks(func(item []T) bool {
			return yield([2]T{item[0], item[1]})
		})
	}
}

// Map :: (a -> b) -> [a] -> [b]
//
//go:inline
func Map[a, b any](f func(a) b) func(Seq[a]) Seq[b] {
	return func(as Seq[a]) Seq[b] {
		return func(yield func(b) bool) {
			as(func(a a) bool {
				return yield(f(a))
			})
		}
	}
}

// Filter :: (a -> Bool) -> [a] -> [a]
//
//go:inline
func Filter[a any](p func(a) bool) func(Seq[a]) Seq[a] {
	return func(as Seq[a]) Seq[a] {
		return func(yield func(a) bool) {
			as(func(av a) bool {
				if p(av) {
					yield(av)
				}
				return true
			})
		}
	}
}

// FlatMap :: (a -> [b]) -> [a] -> [b]
//
//go:inline
func FlatMap[a, b any](f func(a) Seq[b]) func(Seq[a]) Seq[b] {
	return func(as Seq[a]) Seq[b] {
		return func(yield func(b) bool) {
			as(func(a a) bool {
				cont := true
				f(a)(func(b b) bool {
					if !yield(b) {
						cont = false
						return false
					}
					return true
				})
				return cont
			})
		}
	}
}

// ==================== 判断类函数 ====================

// Any :: (a -> Bool) -> [a] -> Bool
//
//go:inline
func Any[a any](predicate func(a) bool) func(Seq[a]) bool {
	return func(seq Seq[a]) bool {
		found := false
		seq(func(item a) bool {
			if predicate(item) {
				found = true
				return false
			}
			return true
		})
		return found
	}
}

// All :: (a -> Bool) -> [a] -> Bool
//
//go:inline
func All[a any](predicate func(a) bool) func(Seq[a]) bool {
	return func(seq Seq[a]) bool {
		allMatch := true
		seq(func(item a) bool {
			if !predicate(item) {
				allMatch = false
				return false
			}
			return true
		})
		return allMatch
	}
}

// None :: (a -> Bool) -> [a] -> Bool
//
//go:inline
func None[a any](predicate func(a) bool) func(Seq[a]) bool {
	return func(seq Seq[a]) bool {
		return !Any(predicate)(seq)
	}
}

// Contains :: a -> [a] -> Bool
//
//go:inline
func Contains[a comparable](target a) func(Seq[a]) bool {
	return func(seq Seq[a]) bool {
		return Any(func(item a) bool { return item == target })(seq)
	}
}

// ==================== 查找类函数 ====================

// Find :: (a -> Bool) -> [a] -> Maybe a
//
//go:inline
func Find[a any](predicate func(a) bool) func(Seq[a]) (a, bool) {
	return func(seq Seq[a]) (a, bool) {
		var result a
		found := false
		seq(func(item a) bool {
			if predicate(item) {
				result = item
				found = true
				return false
			}
			return true
		})
		return result, found
	}
}

// First :: [a] -> Maybe a
//
//go:inline
func First[a any]() func(Seq[a]) (a, bool) {
	return func(seq Seq[a]) (a, bool) {
		var result a
		found := false
		seq(func(item a) bool {
			result = item
			found = true
			return false
		})
		return result, found
	}
}

// Last :: [a] -> Maybe a
//
//go:inline
func Last[a any]() func(Seq[a]) (a, bool) {
	return func(seq Seq[a]) (a, bool) {
		var result a
		found := false
		seq(func(item a) bool {
			result = item
			found = true
			return true
		})
		return result, found
	}
}

// ==================== 去重类函数 ====================

// Unique :: [a] -> [a]
//
//go:inline
func Unique[a comparable]() func(Seq[a]) Seq[a] {
	return func(seq Seq[a]) Seq[a] {
		return func(yield func(a) bool) {
			seen := make(map[a]struct{})
			seq(func(item a) bool {
				if _, exists := seen[item]; !exists {
					seen[item] = struct{}{}
					return yield(item)
				}
				return true
			})
		}
	}
}

// UniqueBy :: (a -> k) -> [a] -> [a]
//
//go:inline
func UniqueBy[a any, k comparable](keyFunc func(a) k) func(Seq[a]) Seq[a] {
	return func(seq Seq[a]) Seq[a] {
		return func(yield func(a) bool) {
			seen := make(map[k]struct{})
			seq(func(item a) bool {
				key := keyFunc(item)
				if _, exists := seen[key]; !exists {
					seen[key] = struct{}{}
					return yield(item)
				}
				return true
			})
		}
	}
}

// ==================== 序列操作函数 ====================

// DropWhile :: (a -> Bool) -> [a] -> [a]
//
//go:inline
func DropWhile[a any](predicate func(a) bool) func(Seq[a]) Seq[a] {
	return func(seq Seq[a]) Seq[a] {
		return func(yield func(a) bool) {
			dropping := true
			seq(func(item a) bool {
				if dropping && predicate(item) {
					return true
				}
				dropping = false
				return yield(item)
			})
		}
	}
}

// Concat :: [[a]] -> [a]
//
//go:inline
func Concat[a any](seqs ...Seq[a]) Seq[a] {
	return func(yield func(a) bool) {
		for _, seq := range seqs {
			cont := true
			seq(func(item a) bool {
				if !yield(item) {
					cont = false
					return false
				}
				return true
			})
			if !cont {
				return
			}
		}
	}
}

// Flatten :: [[a]] -> [a]
//
//go:inline
func Flatten[a any]() func(Seq[Seq[a]]) Seq[a] {
	return func(seqs Seq[Seq[a]]) Seq[a] {
		return func(yield func(a) bool) {
			seqs(func(seq Seq[a]) bool {
				cont := true
				seq(func(item a) bool {
					if !yield(item) {
						cont = false
						return false
					}
					return true
				})
				return cont
			})
		}
	}
}

// Reverse :: [a] -> [a]
//
//go:inline
func Reverse[a any]() func(Seq[a]) Seq[a] {
	return func(seq Seq[a]) Seq[a] {
		return func(yield func(a) bool) {
			var items []a
			seq(func(item a) bool {
				items = append(items, item)
				return true
			})
			for i := len(items) - 1; i >= 0; i-- {
				if !yield(items[i]) {
					return
				}
			}
		}
	}
}

// ==================== 分组和遍历 ====================

// GroupBy :: (a -> k) -> [a] -> Map k [a]
//
//go:inline
func GroupBy[a any, k comparable](keyFunc func(a) k) func(Seq[a]) map[k][]a {
	return func(seq Seq[a]) map[k][]a {
		result := make(map[k][]a)
		seq(func(item a) bool {
			key := keyFunc(item)
			result[key] = append(result[key], item)
			return true
		})
		return result
	}
}

// ForEach :: (a -> ()) -> [a] -> ()
//
//go:inline
func ForEach[a any](action func(a)) func(Seq[a]) {
	return func(seq Seq[a]) {
		seq(func(item a) bool {
			action(item)
			return true
		})
	}
}

// ==================== Seq2 操作函数 ====================

// SeqKeys :: Seq2 k v -> Seq k
//
//go:inline
func SeqKeys[k, v any]() func(Seq2[k, v]) Seq[k] {
	return func(seq2 Seq2[k, v]) Seq[k] {
		return func(yield func(k) bool) {
			seq2(func(key k, _ v) bool {
				return yield(key)
			})
		}
	}
}

// SeqValues :: Seq2 k v -> Seq v
//
//go:inline
func SeqValues[k, v any]() func(Seq2[k, v]) Seq[v] {
	return func(seq2 Seq2[k, v]) Seq[v] {
		return func(yield func(v) bool) {
			seq2(func(_ k, val v) bool {
				return yield(val)
			})
		}
	}
}

// Seq2Filter :: (k -> v -> Bool) -> Seq2 k v -> Seq2 k v
//
//go:inline
func Seq2Filter[k, v any](predicate func(k, v) bool) func(Seq2[k, v]) Seq2[k, v] {
	return func(seq2 Seq2[k, v]) Seq2[k, v] {
		return func(yield func(k, v) bool) {
			seq2(func(key k, val v) bool {
				if predicate(key, val) {
					return yield(key, val)
				}
				return true
			})
		}
	}
}

// Enumerate :: [a] -> Seq2 Int a
//
//go:inline
func Enumerate[a any]() func(Seq[a]) Seq2[int, a] {
	return func(seq Seq[a]) Seq2[int, a] {
		return func(yield func(int, a) bool) {
			i := 0
			seq(func(item a) bool {
				result := yield(i, item)
				i++
				return result
			})
		}
	}
}

// ==================== 函数组合 ====================

// Compose :: ([b]->[c])->([a]->[b])->[t]->[t]
//
//go:inline
func Compose[a, b, c any](f func(Seq[b]) Seq[c], g func(Seq[a]) Seq[b]) func(Seq[a]) Seq[c] {
	return func(as Seq[a]) Seq[c] {
		return f(g(as))
	}
}

// Pipe Pipe the functions from left to right
//
//go:inline
func Pipe[a, b, c any](f func(Seq[a]) Seq[b], g func(Seq[b]) Seq[c]) func(Seq[a]) Seq[c] {
	return func(as Seq[a]) Seq[c] {
		return g(f(as))
	}
}
