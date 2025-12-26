package fp

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

// ==================== Maybe Monad 测试 ====================

func TestMaybeJust(t *testing.T) {
	// 测试 Just 创建包含值的 Maybe
	m := Just(42)

	if !m.IsJust() {
		t.Error("Expected IsJust to be true")
	}

	if m.IsNothing() {
		t.Error("Expected IsNothing to be false")
	}

	v, ok := m.Unwrap()
	if !ok || v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
}

func TestMaybeNothing(t *testing.T) {
	// 测试 Nothing 创建空的 Maybe
	m := Nothing[int]()

	if m.IsJust() {
		t.Error("Expected IsJust to be false")
	}

	if !m.IsNothing() {
		t.Error("Expected IsNothing to be true")
	}

	_, ok := m.Unwrap()
	if ok {
		t.Error("Expected Unwrap to return false")
	}
}

func TestMaybeUnwrapOr(t *testing.T) {
	// 测试 UnwrapOr 默认值
	just := Just(42)
	nothing := Nothing[int]()

	if just.UnwrapOr(0) != 42 {
		t.Error("Expected 42 for Just")
	}

	if nothing.UnwrapOr(100) != 100 {
		t.Error("Expected 100 for Nothing")
	}
}

func TestMaybeMap(t *testing.T) {
	// 测试 Maybe 的 Map 操作
	just := Just(21)
	nothing := Nothing[int]()

	doubled := MaybeMap(just, func(x int) int { return x * 2 })
	if v, ok := doubled.Unwrap(); !ok || v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}

	nothingDoubled := MaybeMap(nothing, func(x int) int { return x * 2 })
	if nothingDoubled.IsJust() {
		t.Error("Expected Nothing after mapping Nothing")
	}
}

func TestMaybeFlatMap(t *testing.T) {
	// 测试 Maybe 的 FlatMap 操作（Monad 核心）
	// FlatMap 允许链式调用返回 Maybe 的函数

	// 示例：安全除法链
	result := MaybeFlatMap(Just(100), func(a int) Maybe[int] {
		return MaybeFlatMap(SafeDiv(a, 5), func(b int) Maybe[int] {
			return SafeDiv(b, 2)
		})
	})

	if v, ok := result.Unwrap(); !ok || v != 10 {
		t.Errorf("Expected 10, got %d", v)
	}

	// 测试除零失败
	failResult := MaybeFlatMap(Just(100), func(a int) Maybe[int] {
		return MaybeFlatMap(SafeDiv(a, 0), func(b int) Maybe[int] {
			return SafeDiv(b, 2)
		})
	})

	if failResult.IsJust() {
		t.Error("Expected Nothing for division by zero")
	}
}

func TestMaybeMonadLaws(t *testing.T) {
	// 测试 Monad 三定律

	f := func(x int) Maybe[int] { return Just(x * 2) }
	g := func(x int) Maybe[int] { return Just(x + 10) }

	// 1. 左单位元: unit(a).flatMap(f) == f(a)
	a := 5
	left1 := MaybeFlatMap(Just(a), f)
	right1 := f(a)
	if l, _ := left1.Unwrap(); l != func() int { v, _ := right1.Unwrap(); return v }() {
		t.Error("Left identity law violated")
	}

	// 2. 右单位元: m.flatMap(unit) == m
	m := Just(42)
	left2 := MaybeFlatMap(m, func(x int) Maybe[int] { return Just(x) })
	if l, _ := left2.Unwrap(); l != 42 {
		t.Error("Right identity law violated")
	}

	// 3. 结合律: m.flatMap(f).flatMap(g) == m.flatMap(x => f(x).flatMap(g))
	left3 := MaybeFlatMap(MaybeFlatMap(m, f), g)
	right3 := MaybeFlatMap(m, func(x int) Maybe[int] {
		return MaybeFlatMap(f(x), g)
	})
	l3, _ := left3.Unwrap()
	r3, _ := right3.Unwrap()
	if l3 != r3 {
		t.Errorf("Associativity law violated: %d != %d", l3, r3)
	}
}

// ==================== Result Monad 测试 ====================

func TestResultOk(t *testing.T) {
	// 测试 Ok 创建成功的 Result
	r := Ok(42)

	if !r.IsOk() {
		t.Error("Expected IsOk to be true")
	}

	if r.IsErr() {
		t.Error("Expected IsErr to be false")
	}

	v, err := r.Unwrap()
	if err != nil || v != 42 {
		t.Errorf("Expected 42 and nil error, got %d and %v", v, err)
	}
}

func TestResultErr(t *testing.T) {
	// 测试 Err 创建失败的 Result
	expectedErr := errors.New("something went wrong")
	r := Err[int](expectedErr)

	if r.IsOk() {
		t.Error("Expected IsOk to be false")
	}

	if !r.IsErr() {
		t.Error("Expected IsErr to be true")
	}

	_, err := r.Unwrap()
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestResultMap(t *testing.T) {
	// 测试 Result 的 Map 操作
	ok := Ok(21)
	err := Err[int](errors.New("error"))

	doubled := ResultMap(ok, func(x int) int { return x * 2 })
	if v, _ := doubled.Unwrap(); v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}

	errDoubled := ResultMap(err, func(x int) int { return x * 2 })
	if errDoubled.IsOk() {
		t.Error("Expected error after mapping error")
	}
}

func TestResultFlatMap(t *testing.T) {
	// 测试 Result 的 FlatMap 操作
	result := ResultFlatMap(Ok(100), func(a int) Result[int] {
		return ResultFlatMap(SafeDivResult(a, 5), func(b int) Result[int] {
			return SafeDivResult(b, 2)
		})
	})

	if v, err := result.Unwrap(); err != nil || v != 10 {
		t.Errorf("Expected 10, got %d with error %v", v, err)
	}

	// 测试除零失败
	failResult := ResultFlatMap(Ok(100), func(a int) Result[int] {
		return ResultFlatMap(SafeDivResult(a, 0), func(b int) Result[int] {
			return SafeDivResult(b, 2)
		})
	})

	if failResult.IsOk() {
		t.Error("Expected error for division by zero")
	}
}

func TestResultMapErr(t *testing.T) {
	// 测试 ResultMapErr 转换错误
	err := Err[int](errors.New("original"))
	mapped := ResultMapErr(err, func(e error) error {
		return fmt.Errorf("wrapped: %w", e)
	})

	if mapped.IsOk() {
		t.Error("Expected error")
	}

	if !strings.Contains(mapped.Error().Error(), "wrapped") {
		t.Error("Expected wrapped error message")
	}
}

// ==================== IO Monad 测试 ====================

func TestIOMonad(t *testing.T) {
	// 测试 IO Monad 基本操作
	counter := 0

	// IO 操作被包装，不会立即执行
	io := NewIO(func() int {
		counter++
		return 42
	})

	// 验证操作尚未执行
	if counter != 0 {
		t.Error("IO should not execute until Run is called")
	}

	// 执行 IO
	result := io.Run()
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	if counter != 1 {
		t.Errorf("Expected counter to be 1, got %d", counter)
	}
}

func TestIOMap(t *testing.T) {
	// 测试 IO 的 Map 操作
	io := NewIO(func() int { return 21 })
	doubled := IOMap(io, func(x int) int { return x * 2 })

	if doubled.Run() != 42 {
		t.Error("Expected 42")
	}
}

func TestIOFlatMap(t *testing.T) {
	// 测试 IO 的 FlatMap 操作
	io := NewIO(func() int { return 10 })

	result := IOFlatMap(io, func(x int) IO[int] {
		return NewIO(func() int { return x * 2 })
	})

	if result.Run() != 20 {
		t.Error("Expected 20")
	}
}

// ==================== List Monad 测试 ====================

func TestListMonad(t *testing.T) {
	// 测试 List Monad 基本操作
	list := NewList(1, 2, 3)

	result := list.ToSlice()
	expected := []int{1, 2, 3}

	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestListMap(t *testing.T) {
	// 测试 List 的 Map 操作
	list := NewList(1, 2, 3)
	doubled := ListMap(list, func(x int) int { return x * 2 })

	expected := []int{2, 4, 6}
	if !slices.Equal(doubled.ToSlice(), expected) {
		t.Errorf("Expected %v, got %v", expected, doubled.ToSlice())
	}
}

func TestListFlatMap(t *testing.T) {
	// 测试 List 的 FlatMap 操作（用于非确定性计算）
	// 示例：生成所有可能的组合
	list := NewList(1, 2, 3)

	// 每个元素生成 [x, x*10]
	result := ListFlatMap(list, func(x int) List[int] {
		return NewList(x, x*10)
	})

	expected := []int{1, 10, 2, 20, 3, 30}
	if !slices.Equal(result.ToSlice(), expected) {
		t.Errorf("Expected %v, got %v", expected, result.ToSlice())
	}
}

func TestListFlatMapCombinations(t *testing.T) {
	// 测试使用 List Monad 生成笛卡尔积
	xs := NewList(1, 2)
	ys := NewList("a", "b")

	type Pair struct {
		X int
		Y string
	}

	result := ListFlatMap(xs, func(x int) List[Pair] {
		return ListMap(ys, func(y string) Pair {
			return Pair{X: x, Y: y}
		})
	})

	expected := []Pair{
		{1, "a"}, {1, "b"},
		{2, "a"}, {2, "b"},
	}

	resultSlice := result.ToSlice()
	if len(resultSlice) != len(expected) {
		t.Errorf("Expected %d pairs, got %d", len(expected), len(resultSlice))
	}

	for i, p := range resultSlice {
		if p != expected[i] {
			t.Errorf("At index %d: expected %v, got %v", i, expected[i], p)
		}
	}
}

// ==================== State Monad 测试 ====================

func TestStateMonad(t *testing.T) {
	// 测试 State Monad 基本操作
	// State Monad 用于线程化状态

	// 创建一个增加计数器的 State
	increment := NewState(func(s int) (int, int) {
		return s, s + 1
	})

	value, newState := increment.Run(0)
	if value != 0 || newState != 1 {
		t.Errorf("Expected (0, 1), got (%d, %d)", value, newState)
	}
}

func TestStateGet(t *testing.T) {
	// 测试 StateGet 获取当前状态
	get := StateGet[int]()

	value, state := get.Run(42)
	if value != 42 || state != 42 {
		t.Errorf("Expected (42, 42), got (%d, %d)", value, state)
	}
}

func TestStatePut(t *testing.T) {
	// 测试 StatePut 设置新状态
	put := StatePut(100)

	_, state := put.Run(42)
	if state != 100 {
		t.Errorf("Expected state 100, got %d", state)
	}
}

func TestStateModify(t *testing.T) {
	// 测试 StateModify 修改状态
	double := StateModify(func(s int) int { return s * 2 })

	_, state := double.Run(21)
	if state != 42 {
		t.Errorf("Expected state 42, got %d", state)
	}
}

func TestStateFlatMap(t *testing.T) {
	// 测试 State 的 FlatMap 操作
	// 模拟一个简单的栈操作

	push := func(x int) State[[]int, struct{}] {
		return NewState(func(stack []int) (struct{}, []int) {
			return struct{}{}, append(stack, x)
		})
	}

	pop := NewState(func(stack []int) (int, []int) {
		if len(stack) == 0 {
			return 0, stack
		}
		return stack[len(stack)-1], stack[:len(stack)-1]
	})

	// push(1) >> push(2) >> pop
	program := StateFlatMap(push(1), func(_ struct{}) State[[]int, int] {
		return StateFlatMap(push(2), func(_ struct{}) State[[]int, int] {
			return pop
		})
	})

	value, finalStack := program.Run(nil)
	if value != 2 {
		t.Errorf("Expected popped value 2, got %d", value)
	}
	if !slices.Equal(finalStack, []int{1}) {
		t.Errorf("Expected stack [1], got %v", finalStack)
	}
}

// ==================== Reader Monad 测试 ====================

func TestReaderMonad(t *testing.T) {
	// 测试 Reader Monad 基本操作
	// Reader 用于依赖注入

	type Config struct {
		Host string
		Port int
	}

	getHost := NewReader(func(c Config) string {
		return c.Host
	})

	config := Config{Host: "localhost", Port: 8080}
	host := getHost.Run(config)

	if host != "localhost" {
		t.Errorf("Expected localhost, got %s", host)
	}
}

func TestReaderAsk(t *testing.T) {
	// 测试 ReaderAsk 获取整个环境
	ask := ReaderAsk[int]()

	result := ask.Run(42)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestReaderAsks(t *testing.T) {
	// 测试 ReaderAsks 从环境中提取信息
	type Config struct {
		Value int
	}

	getValue := ReaderAsks(func(c Config) int {
		return c.Value * 2
	})

	result := getValue.Run(Config{Value: 21})
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestReaderFlatMap(t *testing.T) {
	// 测试 Reader 的 FlatMap 操作
	type Config struct {
		Multiplier int
		Addend     int
	}

	multiply := NewReader(func(c Config) int {
		return 10 * c.Multiplier
	})

	add := func(x int) Reader[Config, int] {
		return NewReader(func(c Config) int {
			return x + c.Addend
		})
	}

	result := ReaderFlatMap(multiply, add)

	config := Config{Multiplier: 2, Addend: 5}
	if result.Run(config) != 25 { // 10*2 + 5 = 25
		t.Errorf("Expected 25, got %d", result.Run(config))
	}
}

// ==================== Writer Monad 测试 ====================

func TestWriterMonad(t *testing.T) {
	// 测试 Writer Monad 基本操作
	w := NewWriter(42, "initialized")

	value, log := w.Run()
	if value != 42 {
		t.Errorf("Expected value 42, got %d", value)
	}
	if log != "initialized" {
		t.Errorf("Expected log 'initialized', got %s", log)
	}
}

func TestWriterMap(t *testing.T) {
	// 测试 Writer 的 Map 操作
	w := NewWriter(21, "start")
	doubled := WriterMap(w, func(x int) int { return x * 2 })

	value, log := doubled.Run()
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
	if log != "start" {
		t.Errorf("Expected log 'start', got %s", log)
	}
}

func TestWriterFlatMapString(t *testing.T) {
	// 测试 Writer 的 FlatMap 操作（字符串日志）
	w := NewWriter(10, "step1 ")

	result := WriterFlatMapString(w, func(x int) Writer[string, int] {
		return NewWriter(x*2, "step2 ")
	})

	value, log := result.Run()
	if value != 20 {
		t.Errorf("Expected 20, got %d", value)
	}
	if log != "step1 step2 " {
		t.Errorf("Expected 'step1 step2 ', got '%s'", log)
	}
}

func TestWriterFlatMapSlice(t *testing.T) {
	// 测试 Writer 的 FlatMap 操作（切片日志）
	w := NewWriter(10, []string{"log1"})

	result := WriterFlatMapSlice(w, func(x int) Writer[[]string, int] {
		return NewWriter(x*2, []string{"log2", "log3"})
	})

	value, log := result.Run()
	if value != 20 {
		t.Errorf("Expected 20, got %d", value)
	}

	expectedLog := []string{"log1", "log2", "log3"}
	if !slices.Equal(log, expectedLog) {
		t.Errorf("Expected %v, got %v", expectedLog, log)
	}
}

// ==================== 辅助函数测试 ====================

func TestMaybeFromPtr(t *testing.T) {
	// 测试从指针创建 Maybe
	value := 42
	ptr := &value

	just := MaybeFromPtr(ptr)
	if v, ok := just.Unwrap(); !ok || v != 42 {
		t.Error("Expected Just(42)")
	}

	nothing := MaybeFromPtr[int](nil)
	if nothing.IsJust() {
		t.Error("Expected Nothing for nil pointer")
	}
}

func TestResultFromFunc(t *testing.T) {
	// 测试从函数创建 Result
	okResult := ResultFromFunc(func() (int, error) {
		return 42, nil
	})
	if v, err := okResult.Unwrap(); err != nil || v != 42 {
		t.Error("Expected Ok(42)")
	}

	errResult := ResultFromFunc(func() (int, error) {
		return 0, errors.New("error")
	})
	if errResult.IsOk() {
		t.Error("Expected Err")
	}
}

func TestMaybeSequence(t *testing.T) {
	// 测试 MaybeSequence
	allJust := []Maybe[int]{Just(1), Just(2), Just(3)}
	result := MaybeSequence(allJust)
	if values, ok := result.Unwrap(); !ok || !slices.Equal(values, []int{1, 2, 3}) {
		t.Error("Expected Just([1,2,3])")
	}

	withNothing := []Maybe[int]{Just(1), Nothing[int](), Just(3)}
	nothingResult := MaybeSequence(withNothing)
	if nothingResult.IsJust() {
		t.Error("Expected Nothing when any element is Nothing")
	}
}

func TestResultSequence(t *testing.T) {
	// 测试 ResultSequence
	allOk := []Result[int]{Ok(1), Ok(2), Ok(3)}
	result := ResultSequence(allOk)
	if values, err := result.Unwrap(); err != nil || !slices.Equal(values, []int{1, 2, 3}) {
		t.Error("Expected Ok([1,2,3])")
	}

	withErr := []Result[int]{Ok(1), Err[int](errors.New("error")), Ok(3)}
	errResult := ResultSequence(withErr)
	if errResult.IsOk() {
		t.Error("Expected Err when any element is Err")
	}
}

func TestSafeDiv(t *testing.T) {
	// 测试安全除法
	result := SafeDiv(10, 2)
	if v, ok := result.Unwrap(); !ok || v != 5 {
		t.Error("Expected Just(5)")
	}

	zero := SafeDiv(10, 0)
	if zero.IsJust() {
		t.Error("Expected Nothing for division by zero")
	}
}

func TestParseIntMaybe(t *testing.T) {
	// 测试整数解析
	valid := ParseIntMaybe("42")
	if v, ok := valid.Unwrap(); !ok || v != 42 {
		t.Error("Expected Just(42)")
	}

	invalid := ParseIntMaybe("not a number")
	if invalid.IsJust() {
		t.Error("Expected Nothing for invalid input")
	}
}

// ==================== 实际使用示例 ====================

func ExampleMaybe_chainedOperations() {
	// 示例：链式操作处理可能失败的计算
	// 计算 100 / x / y，其中任何除零都会导致 Nothing

	calculate := func(x, y int) Maybe[int] {
		return MaybeFlatMap(SafeDiv(100, x), func(a int) Maybe[int] {
			return SafeDiv(a, y)
		})
	}

	result1 := calculate(5, 2) // 100 / 5 / 2 = 10
	result2 := calculate(0, 2) // 除零，返回 Nothing
	result3 := calculate(5, 0) // 除零，返回 Nothing

	fmt.Println(result1.UnwrapOr(-1)) // 10
	fmt.Println(result2.UnwrapOr(-1)) // -1
	fmt.Println(result3.UnwrapOr(-1)) // -1

	// Output:
	// 10
	// -1
	// -1
}

func ExampleResult_errorHandling() {
	// 示例：使用 Result 进行错误处理
	divide := func(a, b int) Result[int] {
		if b == 0 {
			return Err[int](fmt.Errorf("cannot divide %d by zero", a))
		}
		return Ok(a / b)
	}

	// 链式计算
	result := ResultFlatMap(divide(100, 5), func(a int) Result[int] {
		return divide(a, 2)
	})

	if v, err := result.Unwrap(); err == nil {
		fmt.Printf("Result: %d\n", v)
	}

	// Output:
	// Result: 10
}

func ExampleList_nonDeterministic() {
	// 示例：使用 List Monad 进行非确定性计算
	// 生成所有可能的点坐标

	xs := NewList(1, 2)
	ys := NewList(10, 20)

	points := ListFlatMap(xs, func(x int) List[string] {
		return ListMap(ys, func(y int) string {
			return fmt.Sprintf("(%d,%d)", x, y)
		})
	})

	for _, p := range points.ToSlice() {
		fmt.Println(p)
	}

	// Output:
	// (1,10)
	// (1,20)
	// (2,10)
	// (2,20)
}

func ExampleState_counter() {
	// 示例：使用 State Monad 实现计数器
	increment := NewState(func(count int) (int, int) {
		return count, count + 1
	})

	// 连续增加三次
	program := StateFlatMap(increment, func(a int) State[int, int] {
		return StateFlatMap(increment, func(b int) State[int, int] {
			return StateMap(increment, func(c int) int {
				return a + b + c // 返回三次增加前的值之和
			})
		})
	})

	result, finalState := program.Run(0)
	fmt.Printf("Sum of values: %d, Final state: %d\n", result, finalState)

	// Output:
	// Sum of values: 3, Final state: 3
}

func ExampleReader_dependencyInjection() {
	// 示例：使用 Reader Monad 进行依赖注入
	type AppConfig struct {
		DatabaseURL string
		APIKey      string
	}

	getDBUrl := NewReader(func(c AppConfig) string {
		return c.DatabaseURL
	})

	formatConnection := ReaderFlatMap(getDBUrl, func(url string) Reader[AppConfig, string] {
		return ReaderAsks(func(c AppConfig) string {
			return fmt.Sprintf("Connecting to %s with key %s", url, c.APIKey)
		})
	})

	config := AppConfig{
		DatabaseURL: "postgres://localhost:5432/mydb",
		APIKey:      "secret123",
	}

	fmt.Println(formatConnection.Run(config))

	// Output:
	// Connecting to postgres://localhost:5432/mydb with key secret123
}

func ExampleWriter_logging() {
	// 示例：使用 Writer Monad 进行日志记录
	type Log = []string

	withLog := func(value int, msg string) Writer[Log, int] {
		return NewWriter(value, []string{msg})
	}

	// 带日志的计算
	result := WriterFlatMapSlice(
		withLog(10, "Starting with 10"),
		func(x int) Writer[Log, int] {
			return WriterFlatMapSlice(
				withLog(x*2, "Doubled to 20"),
				func(y int) Writer[Log, int] {
					return withLog(y+5, "Added 5 to get 25")
				},
			)
		},
	)

	value, logs := result.Run()
	fmt.Printf("Final value: %d\n", value)
	fmt.Printf("Logs: %v\n", logs)

	// Output:
	// Final value: 25
	// Logs: [Starting with 10 Doubled to 20 Added 5 to get 25]
}
