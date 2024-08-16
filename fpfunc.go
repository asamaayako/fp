/**
 * Author: asamaayako
 * Date: 2024/4/9
 * Time: 上午9:07
 * Describe:
 */

package fp

import "iter"

// Fold:: ((a -> b -> a) -> a) -> [b] -> a
func Fold[a, b any](f func(a, b) a, acc a) func(iter.Seq[b]) a {
	return func(bs iter.Seq[b]) a {
		for bv := range bs {
			acc = f(acc, bv)
		}
		return acc
	}
}

func Take[a any](n int) func(seq iter.Seq[a]) iter.Seq[a] {
	return func(seq iter.Seq[a]) iter.Seq[a] {
		var count int
		return func(yield func(a) bool) {
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

func Reduce[a, b any](f func(a, b) a) func(iter.Seq[b]) a {
	var av a
	return Fold(f, av)
}

func Zip[Fir, Sec any](seq2 iter.Seq2[Fir, Sec]) iter.Seq[Pair[Fir, Sec]] {
	return func(yield func(Pair[Fir, Sec]) bool) {
		seq2(func(fir Fir, sec Sec) bool {
			return yield(Pair[Fir, Sec]{fir, sec})
		})
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

// Split :: [a] -> ([a],[a])
func Split[a any](seq iter.Seq[a]) (iter.Seq[a], iter.Seq[a]) {
	seq1 := seq // iter.Seq 是只读的因此简单复制不会出现问题
	return seq1, seq
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
