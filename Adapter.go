package fp

//这个文件用于将各种集合转换未Iterator接口

type Iterator[T any] interface {
	Iter() <-chan T
}

type SliceIterator[T any] []T

func NewSliceIterator[T any](p []T) SliceIterator[T] {
	return SliceIterator[T](p)
}
func (t SliceIterator[T]) Iter() <-chan T {
	r := make(chan T)
	go func() {
		for _, v := range t {
			r <- v
		}
		close(r)
	}()
	return r
}

type resultIterator[R any] func() <-chan R

func (r *resultIterator[R]) Iter() <-chan R {
	return (*r)()
}

func (r *resultIterator[R]) name() {

}
func NewIterator[R any](fn func() <-chan R) Iterator[R] {
	r := new(resultIterator[R])
	*r = fn
	return r
}

func MapIter[T, R any](fn func(T) R, iter Iterator[T]) Iterator[R] {
	return NewIterator[R](func() <-chan R {
		r := make(chan R)
		go func() {
			for V := range iter.Iter() {
				r <- fn(V)
			}
		}()
		return r
	})
}

func SliceToIterator[T any, TS interface{ ~[]T }](s TS) Iterator[T] {
	return SliceIterator[T](s)
}
