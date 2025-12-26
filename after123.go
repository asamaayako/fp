//go:build go1.23

package fp

import (
	"iter"
	"maps"
	"slices"
)

type seq[V any] = iter.Seq[V]
type seq2[K, V any] = iter.Seq2[K, V]

// ==================== 核心函数（有性能或API差异）====================

// fold 使用 for range 语法糖（代码更简洁）
//
//go:inline
func fold[a, b any](acc a, f func(a, b) a) func(Seq[b]) a {
	return func(bs Seq[b]) a {
		for bv := range bs {
			acc = f(acc, bv)
		}
		return acc
	}
}

// zip 使用 iter.Pull 实现真正的惰性求值
//
//go:inline
func zip[T, U any](seq1 Seq[T], seq2 Seq[U]) Seq[Pair[T, U]] {
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

// unZipSeq2 使用 slices.Values
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

	return slices.Values(firsts), slices.Values(seconds)
}

// ==================== 标准库优化函数 ====================

// collect 使用 slices.Collect
//
//go:inline
func collect[a any]() func(Seq[a]) []a {
	return func(seq Seq[a]) []a {
		return slices.Collect(seq)
	}
}

// fromSlice 使用 slices.Values
//
//go:inline
func fromSlice[a any](slice []a) Seq[a] {
	return slices.Values(slice)
}

// fromMap 使用 maps.All
//
//go:inline
func fromMap[k comparable, v any](m map[k]v) Seq2[k, v] {
	return maps.All(m)
}

// mapKeys 使用 maps.Keys
//
//go:inline
func mapKeys[k comparable, v any](m map[k]v) Seq[k] {
	return maps.Keys(m)
}

// mapValues 使用 maps.Values
//
//go:inline
func mapValues[k comparable, v any](m map[k]v) Seq[v] {
	return maps.Values(m)
}
