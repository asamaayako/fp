package fp

import (
	"slices"
	"testing"
)

func TestFold(t *testing.T) {
	// 测试 Fold 函数，计算数组元素之和
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	sum := Fold(0, func(acc, val int) int {
		return acc + val
	})(seq)

	if sum != 15 {
		t.Errorf("Expected sum to be 15, got %d", sum)
	}
}

func TestHead(t *testing.T) {
	// 测试 Head 函数，获取序列的第一个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	headSeq := Head[int]()(seq)

	result := slices.Collect(headSeq)

	if len(result) != 1 || result[0] != 1 {
		t.Errorf("Expected [1], got %v", result)
	}
}

func TestTake(t *testing.T) {
	// 测试 Take 函数，获取序列的前n个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	takeSeq := Take[int](3)(seq)

	result := slices.Collect(takeSeq)

	expected := []int{1, 2, 3}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestLen(t *testing.T) {
	// 测试 Len 函数，计算序列长度
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	length := Len(seq)

	if length != 5 {
		t.Errorf("Expected length to be 5, got %d", length)
	}
}

func TestTail(t *testing.T) {
	// 测试 Tail 函数，去掉序列的第一个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	tailSeq := Tail[int]()(seq)

	result := slices.Collect(tailSeq)

	expected := []int{2, 3, 4, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDrop(t *testing.T) {
	// 测试 Drop 函数，去掉序列的前n个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	dropSeq := Drop[int](2)(seq)

	result := slices.Collect(dropSeq)

	expected := []int{3, 4, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReduce(t *testing.T) {
	// 测试 Reduce 函数，计算数组元素之积
	numbers := []int{1, 2, 3, 4}
	seq := slices.Values(numbers)
	product := Reduce(func(acc, val int) int {
		return acc * val
	})(seq)

	if product != 24 {
		t.Errorf("Expected product to be 24, got %d", product)
	}
}

func TestZip(t *testing.T) {
	// 测试 Zip 函数，将两个序列合并成键值对
	keys := []string{"a", "b", "c"}
	values := []int{1, 2, 3}

	keySeq := slices.Values(keys)
	valSeq := slices.Values(values)
	zippedSeq := Zip(keySeq, valSeq)

	result := slices.Collect(zippedSeq)

	expected := []Pair[string, int]{
		{First: "a", Second: 1},
		{First: "b", Second: 2},
		{First: "c", Second: 3},
	}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
		return
	}

	for i, pair := range result {
		if pair.First != expected[i].First || pair.Second != expected[i].Second {
			t.Errorf("At index %d: expected %v, got %v", i, expected[i], pair)
		}
	}
}

func TestUnZip(t *testing.T) {
	// 测试 UnZip 函数，将键值对序列拆分为两个序列
	pairs := []Pair[string, int]{
		{First: "a", Second: 1},
		{First: "b", Second: 2},
		{First: "c", Second: 3},
	}

	seq := slices.Values(pairs)
	unzippedSeq := UnZip(seq)

	var keys []string
	var values []int

	unzippedSeq(func(key string, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true
	})

	expectedKeys := []string{"a", "b", "c"}
	expectedValues := []int{1, 2, 3}

	if !slices.Equal(keys, expectedKeys) {
		t.Errorf("Expected keys %v, got %v", expectedKeys, keys)
	}

	if !slices.Equal(values, expectedValues) {
		t.Errorf("Expected values %v, got %v", expectedValues, values)
	}
}

func TestSplit(t *testing.T) {
	// 测试 Split 函数，将序列分割为两个相同的序列
	numbers := []int{1, 2, 3}
	seq := slices.Values(numbers)
	seq1, seq2 := Split(seq)

	result1 := slices.Collect(seq1)
	result2 := slices.Collect(seq2)

	expected := []int{1, 2, 3}

	if !slices.Equal(result1, expected) {
		t.Errorf("Expected first sequence %v, got %v", expected, result1)
	}

	if !slices.Equal(result2, expected) {
		t.Errorf("Expected second sequence %v, got %v", expected, result2)
	}
}

func TestChunk(t *testing.T) {
	// 测试 Chunk 函数，将序列分块
	numbers := []int{1, 2, 3, 4, 5, 6, 7}
	seq := slices.Values(numbers)
	chunkSeq := Chunk[int](3)(seq)

	result := slices.Collect(chunkSeq)

	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}

	if len(result) != len(expected) {
		t.Errorf("Expected %d chunks, got %d", len(expected), len(result))
		return
	}

	for i, chunk := range result {
		if !slices.Equal(chunk, expected[i]) {
			t.Errorf("Chunk %d: expected %v, got %v", i, expected[i], chunk)
		}
	}
}

func TestSlidingWindow(t *testing.T) {
	// 测试 SlidingWindow 函数，创建滑动窗口
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	windowSeq := SlidingWindow[int](3)(seq)

	result := slices.Collect(windowSeq)

	expected := [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}

	if len(result) != len(expected) {
		t.Errorf("Expected %d windows, got %d", len(expected), len(result))
		return
	}

	for i, window := range result {
		if !slices.Equal(window, expected[i]) {
			t.Errorf("Window %d: expected %v, got %v", i, expected[i], window)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	// 测试 TakeWhile 函数，当满足条件时继续迭代
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	seq := slices.Values(numbers)
	takeWhileSeq := TakeWhile(func(n int) bool {
		return n < 5
	})(seq)

	result := slices.Collect(takeWhileSeq)

	expected := []int{1, 2, 3, 4}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAdjacentPairs(t *testing.T) {
	// 测试 AdjacentPairs 函数，处理相邻元素对
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	pairsSeq := AdjacentPairs(seq)

	result := slices.Collect(pairsSeq)

	expected := [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}

	if len(result) != len(expected) {
		t.Errorf("Expected %d pairs, got %d", len(expected), len(result))
		return
	}

	for i, pair := range result {
		if pair[0] != expected[i][0] || pair[1] != expected[i][1] {
			t.Errorf("Pair %d: expected %v, got %v", i, expected[i], pair)
		}
	}
}

func TestPairs(t *testing.T) {
	// 测试 Pairs 函数，将序列元素两两分组
	numbers := []int{1, 2, 3, 4, 5, 6}
	seq := slices.Values(numbers)
	pairsSeq := Pairs(seq)

	result := slices.Collect(pairsSeq)

	expected := [][2]int{{1, 2}, {3, 4}, {5, 6}}

	if len(result) != len(expected) {
		t.Errorf("Expected %d pairs, got %d", len(expected), len(result))
		return
	}

	for i, pair := range result {
		if pair[0] != expected[i][0] || pair[1] != expected[i][1] {
			t.Errorf("Pair %d: expected %v, got %v", i, expected[i], pair)
		}
	}
}

func TestMap(t *testing.T) {
	// 测试 Map 函数，对序列中每个元素应用函数
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	mapSeq := Map(func(n int) int {
		return n * 2
	})(seq)

	result := slices.Collect(mapSeq)

	expected := []int{2, 4, 6, 8, 10}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFilter(t *testing.T) {
	// 测试 Filter 函数，过滤序列中满足条件的元素
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	seq := slices.Values(numbers)
	filterSeq := Filter(func(n int) bool {
		return n%2 == 0
	})(seq)

	result := slices.Collect(filterSeq)

	expected := []int{2, 4, 6, 8, 10}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCompose(t *testing.T) {
	// 测试 Compose 函数，组合两个函数
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	// 先过滤偶数，再将每个元素乘以2
	filterEven := Filter(func(n int) bool { return n%2 == 0 })
	double := Map(func(n int) int { return n * 2 })

	composed := Compose(double, filterEven)(seq)
	result := slices.Collect(composed)

	expected := []int{4, 8} // 偶数2,4 -> 4,8
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPipe(t *testing.T) {
	// 测试 Pipe 函数，从左到右管道函数
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	// 先将每个元素乘以2，再过滤偶数结果
	double := Map(func(n int) int { return n * 2 })
	filterEven := Filter(func(n int) bool { return n%2 == 0 })

	piped := Pipe(double, filterEven)(seq)
	result := slices.Collect(piped)

	expected := []int{2, 4, 6, 8, 10} // 1,2,3,4,5 -> 2,4,6,8,10 -> 2,4,6,8,10
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
