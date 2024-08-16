package fp

import "iter"

type Pair[first, second any] struct {
	First  first
	Second second
}

// Chan2Sep 必须在外部写入端主动关闭 否则会造成内存泄露
func Chan2Sep[CT <-chan T, T any](c CT) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range c {
			if !yield(v) {
				for range c { //清空剩余的数据
				}
				return
			}
		}
	}
}
