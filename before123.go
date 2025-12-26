//go:build go1.18 && !go1.23

package fp

type Seq[V any] func(yield func(V) bool)
type Seq2[K, V any] func(yield func(K, V) bool)

// ==================== 核心函数（兼容实现）====================

// fold 手动调用序列函数
//
//go:inline
func fold[a, b any](acc a, f func(a, b) a) func(Seq[b]) a {
	return func(bs Seq[b]) a {
		bs(func(bv b) bool {
			acc = f(acc, bv)
			return true
		})
		return acc
	}
}

// zip 通过物化序列实现（无 iter.Pull）
//
//go:inline
func zip[T, U any](seq1 Seq[T], seq2 Seq[U]) Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		var items1 []T
		var items2 []U

		seq1(func(v T) bool {
			items1 = append(items1, v)
			return true
		})

		seq2(func(v U) bool {
			items2 = append(items2, v)
			return true
		})

		minLen := len(items1)
		if len(items2) < minLen {
			minLen = len(items2)
		}

		for i := 0; i < minLen; i++ {
			if !yield(Pair[T, U]{items1[i], items2[i]}) {
				return
			}
		}
	}
}

// unZipSeq2 手动创建 Seq（无 slices.Values）
//
//go:inline
func unZipSeq2[Fir, Sec any](seq2 Seq2[Fir, Sec]) (Seq[Fir], Seq[Sec]) {
	var firsts []Fir
	var seconds []Sec

	seq2(func(f Fir, s Sec) bool {
		firsts = append(firsts, f)
		seconds = append(seconds, s)
		return true
	})

	seqFir := func(yield func(Fir) bool) {
		for _, v := range firsts {
			if !yield(v) {
				return
			}
		}
	}

	seqSec := func(yield func(Sec) bool) {
		for _, v := range seconds {
			if !yield(v) {
				return
			}
		}
	}

	return seqFir, seqSec
}

// ==================== 兼容的辅助函数 ====================

// collect 手动收集
//
//go:inline
func collect[a any]() func(Seq[a]) []a {
	return func(seq Seq[a]) []a {
		var result []a
		seq(func(item a) bool {
			result = append(result, item)
			return true
		})
		return result
	}
}

// fromSlice 手动创建 Seq
//
//go:inline
func fromSlice[a any](slice []a) Seq[a] {
	return func(yield func(a) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

// fromMap 手动遍历 map
//
//go:inline
func fromMap[k comparable, v any](m map[k]v) Seq2[k, v] {
	return func(yield func(k, v) bool) {
		for key, val := range m {
			if !yield(key, val) {
				return
			}
		}
	}
}

// mapKeys 手动获取 keys
//
//go:inline
func mapKeys[k comparable, v any](m map[k]v) Seq[k] {
	return func(yield func(k) bool) {
		for key := range m {
			if !yield(key) {
				return
			}
		}
	}
}

// mapValues 手动获取 values
//
//go:inline
func mapValues[k comparable, v any](m map[k]v) Seq[v] {
	return func(yield func(v) bool) {
		for _, val := range m {
			if !yield(val) {
				return
			}
		}
	}
}
