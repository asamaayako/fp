/**
 * Author: asamaayako
 * Date: 2024/4/9
 * Describe: Go 语言中的 Monad 实现
 *
 * Monad 是一种设计模式，用于处理带有上下文的计算（如可能失败、异步、可选值等）。
 * Monad 需要满足三个定律：
 * 1. 左单位元: unit(a).flatMap(f) == f(a)
 * 2. 右单位元: m.flatMap(unit) == m
 * 3. 结合律: m.flatMap(f).flatMap(g) == m.flatMap(x => f(x).flatMap(g))
 */
package fp

import (
	"fmt"
	"runtime"
	"strings"
)

// ==================== Maybe Monad ====================
// Maybe 用于处理可能不存在的值，避免 nil 指针问题

type Maybe[T any] struct {
	value   T
	isValid bool
}

// Just 创建一个包含值的 Maybe
func Just[T any](value T) Maybe[T] {
	return Maybe[T]{value: value, isValid: true}
}

// Nothing 创建一个空的 Maybe
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{isValid: false}
}

// IsJust 检查 Maybe 是否包含值
func (m Maybe[T]) IsJust() bool {
	return m.isValid
}

// IsNothing 检查 Maybe 是否为空
func (m Maybe[T]) IsNothing() bool {
	return !m.isValid
}

// Unwrap 获取 Maybe 中的值，如果为空则返回默认值
func (m Maybe[T]) Unwrap() (T, bool) {
	return m.value, m.isValid
}

// UnwrapOr 获取 Maybe 中的值，如果为空则返回默认值
func (m Maybe[T]) UnwrapOr(defaultValue T) T {
	if m.isValid {
		return m.value
	}
	return defaultValue
}

// Map :: Maybe a -> (a -> b) -> Maybe b
func MaybeMap[A, B any](m Maybe[A], f func(A) B) Maybe[B] {
	if m.isValid {
		return Just(f(m.value))
	}
	return Nothing[B]()
}

// FlatMap :: Maybe a -> (a -> Maybe b) -> Maybe b (这是 Monad 的核心操作 bind/>>=/flatMap)
func MaybeFlatMap[A, B any](m Maybe[A], f func(A) Maybe[B]) Maybe[B] {
	if m.isValid {
		return f(m.value)
	}
	return Nothing[B]()
}

// ==================== Result Monad ====================
// Result 用于处理可能失败的计算，类似于 Rust 的 Result 或 Haskell 的 Either

type Result[T any] struct {
	value T
	err   error
	isOk  bool
}

// Ok 创建一个成功的 Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, isOk: true}
}

// Err 创建一个失败的 Result
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, isOk: false}
}

// IsOk 检查 Result 是否成功
func (r Result[T]) IsOk() bool {
	return r.isOk
}

// IsErr 检查 Result 是否失败
func (r Result[T]) IsErr() bool {
	return !r.isOk
}

// Unwrap 获取 Result 中的值和错误
func (r Result[T]) Unwrap() T {
	if r.isOk {
		return r.value
	}
	panic(r.panicWithStackTrace())
}

// panicWithStackTrace 获取调用堆栈信息并构造错误消息
func (r Result[T]) panicWithStackTrace() string {
	var buf strings.Builder
	buf.WriteString("called Result.Unwrap() on an error result\n")
	buf.WriteString("error: " + r.err.Error() + "\n")
	buf.WriteString("stack trace:\n")

	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		buf.WriteString(fmt.Sprintf("  %s\n", frame.Function))
		buf.WriteString(fmt.Sprintf("    %s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}

	return buf.String()
}
func (r Result[T]) SafeUnwarp() (T, error) {
	return r.value, r.err
}

// UnwrapOr 获取 Result 中的值，如果失败则返回默认值
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// Error 获取 Result 中的错误
func (r Result[T]) Error() error {
	return r.err
}

// Map :: Result a -> (a -> b) -> Result b
func ResultMap[A, B any](r Result[A], f func(A) B) Result[B] {
	if r.isOk {
		return Ok(f(r.value))
	}
	return Err[B](r.err)
}

// FlatMap :: Result a -> (a -> Result b) -> Result b
func ResultFlatMap[A, B any](r Result[A], f func(A) Result[B]) Result[B] {
	if r.isOk {
		return f(r.value)
	}
	return Err[B](r.err)
}

// MapErr :: Result a -> (error -> error) -> Result a
func ResultMapErr[T any](r Result[T], f func(error) error) Result[T] {
	if r.isOk {
		return r
	}
	return Err[T](f(r.err))
}

// ==================== IO Monad ====================
// IO Monad 用于封装副作用操作，使其变得纯粹和可组合

type IO[T any] struct {
	run func() T
}

// NewIO 创建一个 IO Monad
func NewIO[T any](f func() T) IO[T] {
	return IO[T]{run: f}
}

// Run 执行 IO 操作并返回结果
func (io IO[T]) Run() T {
	return io.run()
}

// Map :: IO a -> (a -> b) -> IO b
func IOMap[A, B any](io IO[A], f func(A) B) IO[B] {
	return NewIO(func() B {
		return f(io.Run())
	})
}

// FlatMap :: IO a -> (a -> IO b) -> IO b
func IOFlatMap[A, B any](io IO[A], f func(A) IO[B]) IO[B] {
	return NewIO(func() B {
		return f(io.Run()).Run()
	})
}

// ==================== List Monad ====================
// List Monad 用于处理非确定性计算（多个可能的结果）

type List[T any] struct {
	values []T
}

// NewList 从切片创建 List
func NewList[T any](values ...T) List[T] {
	return List[T]{values: values}
}

// ToSlice 将 List 转换为切片
func (l List[T]) ToSlice() []T {
	return l.values
}

// Map :: List a -> (a -> b) -> List b
func ListMap[A, B any](l List[A], f func(A) B) List[B] {
	result := make([]B, len(l.values))
	for i, v := range l.values {
		result[i] = f(v)
	}
	return List[B]{values: result}
}

// FlatMap :: List a -> (a -> List b) -> List b
func ListFlatMap[A, B any](l List[A], f func(A) List[B]) List[B] {
	var result []B
	for _, v := range l.values {
		result = append(result, f(v).values...)
	}
	return List[B]{values: result}
}

// ==================== State Monad ====================
// State Monad 用于处理带状态的计算

type State[S, A any] struct {
	run func(S) (A, S)
}

// NewState 创建一个 State Monad
func NewState[S, A any](f func(S) (A, S)) State[S, A] {
	return State[S, A]{run: f}
}

// Run 执行 State 操作
func (s State[S, A]) Run(initialState S) (A, S) {
	return s.run(initialState)
}

// Get 获取当前状态
func StateGet[S any]() State[S, S] {
	return NewState(func(s S) (S, S) {
		return s, s
	})
}

// Put 设置新状态
func StatePut[S any](newState S) State[S, struct{}] {
	return NewState(func(_ S) (struct{}, S) {
		return struct{}{}, newState
	})
}

// Modify 修改状态
func StateModify[S any](f func(S) S) State[S, struct{}] {
	return NewState(func(s S) (struct{}, S) {
		return struct{}{}, f(s)
	})
}

// Map :: State s a -> (a -> b) -> State s b
func StateMap[S, A, B any](st State[S, A], f func(A) B) State[S, B] {
	return NewState(func(s S) (B, S) {
		a, newS := st.Run(s)
		return f(a), newS
	})
}

// FlatMap :: State s a -> (a -> State s b) -> State s b
func StateFlatMap[S, A, B any](st State[S, A], f func(A) State[S, B]) State[S, B] {
	return NewState(func(s S) (B, S) {
		a, newS := st.Run(s)
		return f(a).Run(newS)
	})
}

// ==================== Reader Monad ====================
// Reader Monad 用于依赖注入和共享只读环境

type Reader[E, A any] struct {
	run func(E) A
}

// NewReader 创建一个 Reader Monad
func NewReader[E, A any](f func(E) A) Reader[E, A] {
	return Reader[E, A]{run: f}
}

// Run 执行 Reader 操作
func (r Reader[E, A]) Run(env E) A {
	return r.run(env)
}

// Ask 获取环境
func ReaderAsk[E any]() Reader[E, E] {
	return NewReader(func(e E) E {
		return e
	})
}

// Asks 从环境中提取信息
func ReaderAsks[E, A any](f func(E) A) Reader[E, A] {
	return NewReader(f)
}

// Map :: Reader e a -> (a -> b) -> Reader e b
func ReaderMap[E, A, B any](r Reader[E, A], f func(A) B) Reader[E, B] {
	return NewReader(func(e E) B {
		return f(r.Run(e))
	})
}

// FlatMap :: Reader e a -> (a -> Reader e b) -> Reader e b
func ReaderFlatMap[E, A, B any](r Reader[E, A], f func(A) Reader[E, B]) Reader[E, B] {
	return NewReader(func(e E) B {
		return f(r.Run(e)).Run(e)
	})
}

// ==================== Writer Monad ====================
// Writer Monad 用于累积日志或其他附加信息

type Writer[W, A any] struct {
	value A
	log   W
}

// NewWriter 创建一个 Writer Monad
func NewWriter[W, A any](value A, log W) Writer[W, A] {
	return Writer[W, A]{value: value, log: log}
}

// Run 获取 Writer 的值和日志
func (w Writer[W, A]) Run() (A, W) {
	return w.value, w.log
}

// Tell 添加日志
func WriterTell[W any](log W) Writer[W, struct{}] {
	return Writer[W, struct{}]{value: struct{}{}, log: log}
}

// Map :: Writer w a -> (a -> b) -> Writer w b
func WriterMap[W, A, B any](w Writer[W, A], f func(A) B) Writer[W, B] {
	return Writer[W, B]{value: f(w.value), log: w.log}
}

// FlatMap :: Writer w a -> (a -> Writer w b) -> Writer w b (需要 W 是 monoid)
func WriterFlatMapSlice[W any, A, B any](w Writer[[]W, A], f func(A) Writer[[]W, B]) Writer[[]W, B] {
	next := f(w.value)
	return Writer[[]W, B]{
		value: next.value,
		log:   append(w.log, next.log...),
	}
}

// FlatMap for string logs
func WriterFlatMapString[A, B any](w Writer[string, A], f func(A) Writer[string, B]) Writer[string, B] {
	next := f(w.value)
	return Writer[string, B]{
		value: next.value,
		log:   w.log + next.log,
	}
}

// ==================== 辅助函数 ====================

// MaybeFromPtr 从指针创建 Maybe
func MaybeFromPtr[T any](ptr *T) Maybe[T] {
	if ptr == nil {
		return Nothing[T]()
	}
	return Just(*ptr)
}

// ResultFromFunc 从可能返回错误的函数创建 Result
func ResultFromFunc[T any](f func() (T, error)) Result[T] {
	value, err := f()
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Sequence 将 []Maybe[T] 转换为 Maybe[[]T]
func MaybeSequence[T any](maybes []Maybe[T]) Maybe[[]T] {
	result := make([]T, 0, len(maybes))
	for _, m := range maybes {
		if m.IsNothing() {
			return Nothing[[]T]()
		}
		v, _ := m.Unwrap()
		result = append(result, v)
	}
	return Just(result)
}

// ResultSequence 将 []Result[T] 转换为 Result[[]T]
func ResultSequence[T any](results []Result[T]) Result[[]T] {
	values := make([]T, 0, len(results))
	for _, r := range results {
		if r.IsErr() {
			return Err[[]T](r.err)
		}
		values = append(values, r.value)
	}
	return Ok(values)
}

// ==================== 实用示例函数 ====================

// SafeDiv 安全除法，避免除零错误
func SafeDiv(a, b int) Maybe[int] {
	if b == 0 {
		return Nothing[int]()
	}
	return Just(a / b)
}

// SafeDivResult 安全除法，返回 Result
func SafeDivResult(a, b int) Result[int] {
	if b == 0 {
		return Err[int](fmt.Errorf("division by zero: %d / %d", a, b))
	}
	return Ok(a / b)
}

// ParseInt 解析整数，返回 Maybe
func ParseIntMaybe(s string) Maybe[int] {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return Nothing[int]()
	}
	return Just(result)
}
