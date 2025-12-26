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

func TestFlatMap(t *testing.T) {
	// 测试 FlatMap 函数，扁平化映射
	numbers := []int{1, 2, 3}
	seq := slices.Values(numbers)

	flatMapSeq := FlatMap(func(n int) Seq[int] {
		return slices.Values([]int{n, n * 2})
	})(seq)

	result := slices.Collect(flatMapSeq)

	expected := []int{1, 2, 2, 4, 3, 6}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAny(t *testing.T) {
	// 测试 Any 函数，检查是否有任何元素满足条件
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	hasEven := Any(func(n int) bool { return n%2 == 0 })(seq)
	if !hasEven {
		t.Errorf("Expected Any to find even number, got false")
	}

	seq2 := slices.Values([]int{1, 3, 5})
	hasEven2 := Any(func(n int) bool { return n%2 == 0 })(seq2)
	if hasEven2 {
		t.Errorf("Expected Any to return false for all odd numbers, got true")
	}
}

func TestAll(t *testing.T) {
	// 测试 All 函数，检查所有元素是否满足条件
	numbers := []int{2, 4, 6, 8}
	seq := slices.Values(numbers)

	allEven := All(func(n int) bool { return n%2 == 0 })(seq)
	if !allEven {
		t.Errorf("Expected All to return true for even numbers, got false")
	}

	numbers2 := []int{1, 2, 3, 4}
	seq2 := slices.Values(numbers2)
	allEven2 := All(func(n int) bool { return n%2 == 0 })(seq2)
	if allEven2 {
		t.Errorf("Expected All to return false, got true")
	}
}

func TestNone(t *testing.T) {
	// 测试 None 函数，检查没有任何元素满足条件
	numbers := []int{1, 3, 5, 7}
	seq := slices.Values(numbers)

	noEven := None(func(n int) bool { return n%2 == 0 })(seq)
	if !noEven {
		t.Errorf("Expected None to return true for odd numbers, got false")
	}

	numbers2 := []int{1, 2, 3, 4}
	seq2 := slices.Values(numbers2)
	noEven2 := None(func(n int) bool { return n%2 == 0 })(seq2)
	if noEven2 {
		t.Errorf("Expected None to return false, got true")
	}
}

func TestContains(t *testing.T) {
	// 测试 Contains 函数，检查序列是否包含某个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	contains3 := Contains(3)(seq)
	if !contains3 {
		t.Errorf("Expected Contains to find 3, got false")
	}

	seq2 := slices.Values(numbers)
	contains10 := Contains(10)(seq2)
	if contains10 {
		t.Errorf("Expected Contains to not find 10, got true")
	}
}

func TestFind(t *testing.T) {
	// 测试 Find 函数，查找满足条件的第一个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	found, ok := Find(func(n int) bool { return n > 2 })(seq)
	if !ok || found != 3 {
		t.Errorf("Expected to find 3, got %d (ok=%v)", found, ok)
	}

	seq2 := slices.Values([]int{1, 3, 5})
	_, ok2 := Find(func(n int) bool { return n > 10 })(seq2)
	if ok2 {
		t.Errorf("Expected Find to return false, got true")
	}
}

func TestFirst(t *testing.T) {
	// 测试 First 函数，获取序列的第一个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	first, ok := First[int]()(seq)
	if !ok || first != 1 {
		t.Errorf("Expected 1, got %d (ok=%v)", first, ok)
	}

	seq2 := slices.Values([]int{})
	_, ok2 := First[int]()(seq2)
	if ok2 {
		t.Errorf("Expected First to return false for empty sequence, got true")
	}
}

func TestLast(t *testing.T) {
	// 测试 Last 函数，获取序列的最后一个元素
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	last, ok := Last[int]()(seq)
	if !ok || last != 5 {
		t.Errorf("Expected 5, got %d (ok=%v)", last, ok)
	}

	seq2 := slices.Values([]int{})
	_, ok2 := Last[int]()(seq2)
	if ok2 {
		t.Errorf("Expected Last to return false for empty sequence, got true")
	}
}

func TestUnique(t *testing.T) {
	// 测试 Unique 函数，去除重复元素
	numbers := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
	seq := slices.Values(numbers)
	uniqueSeq := Unique[int]()(seq)

	result := slices.Collect(uniqueSeq)

	expected := []int{1, 2, 3, 4, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestUniqueBy(t *testing.T) {
	// 测试 UniqueBy 函数，根据键函数去除重复元素
	type User struct {
		ID   int
		Name string
	}

	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "Alice2"},
		{ID: 3, Name: "Charlie"},
	}

	seq := slices.Values(users)
	uniqueSeq := UniqueBy(func(u User) int { return u.ID })(seq)

	result := slices.Collect(uniqueSeq)

	if len(result) != 3 {
		t.Errorf("Expected 3 unique users, got %d", len(result))
	}

	if result[0].ID != 1 || result[1].ID != 2 || result[2].ID != 3 {
		t.Errorf("Expected IDs [1, 2, 3], got %v", []int{result[0].ID, result[1].ID, result[2].ID})
	}
}

func TestDropWhile(t *testing.T) {
	// 测试 DropWhile 函数，当条件不满足时停止丢弃
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	dropWhileSeq := DropWhile(func(n int) bool { return n < 4 })(seq)

	result := slices.Collect(dropWhileSeq)

	expected := []int{4, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestConcat(t *testing.T) {
	// 测试 Concat 函数，连接多个序列
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5})
	seq3 := slices.Values([]int{6, 7, 8})

	concatSeq := Concat(seq1, seq2, seq3)
	result := slices.Collect(concatSeq)

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFlatten(t *testing.T) {
	// 测试 Flatten 函数，扁平化序列
	seq1 := slices.Values([]int{1, 2})
	seq2 := slices.Values([]int{3, 4})
	seq3 := slices.Values([]int{5})

	seqOfSeqs := slices.Values([]Seq[int]{seq1, seq2, seq3})
	flatSeq := Flatten[int]()(seqOfSeqs)

	result := slices.Collect(flatSeq)

	expected := []int{1, 2, 3, 4, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReverse(t *testing.T) {
	// 测试 Reverse 函数，反转序列
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)
	reverseSeq := Reverse[int]()(seq)

	result := slices.Collect(reverseSeq)

	expected := []int{5, 4, 3, 2, 1}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGroupBy(t *testing.T) {
	// 测试 GroupBy 函数，按键函数分组
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	seq := slices.Values(numbers)

	grouped := GroupBy(func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})(seq)

	if len(grouped) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(grouped))
	}

	if !slices.Equal(grouped["even"], []int{2, 4, 6, 8}) {
		t.Errorf("Expected even group [2, 4, 6, 8], got %v", grouped["even"])
	}

	if !slices.Equal(grouped["odd"], []int{1, 3, 5, 7}) {
		t.Errorf("Expected odd group [1, 3, 5, 7], got %v", grouped["odd"])
	}
}

func TestForEach(t *testing.T) {
	// 测试 ForEach 函数，遍历序列
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	var collected []int
	ForEach(func(n int) {
		collected = append(collected, n*2)
	})(seq)

	expected := []int{2, 4, 6, 8, 10}
	if !slices.Equal(collected, expected) {
		t.Errorf("Expected %v, got %v", expected, collected)
	}
}

func TestFromSlice(t *testing.T) {
	// 测试 FromSlice 函数，从切片创建序列
	numbers := []int{1, 2, 3, 4, 5}
	seq := FromSlice(numbers)

	result := slices.Collect(seq)

	if !slices.Equal(result, numbers) {
		t.Errorf("Expected %v, got %v", numbers, result)
	}
}

func TestFromMap(t *testing.T) {
	// 测试 FromMap 函数，从映射创建Seq2
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	seq2 := FromMap(m)

	var keys []string
	var values []int
	seq2(func(k string, v int) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// 检查所有键都存在
	expectedKeys := []string{"a", "b", "c"}
	slices.Sort(keys)
	slices.Sort(expectedKeys)
	if !slices.Equal(keys, expectedKeys) {
		t.Errorf("Expected keys %v, got %v", expectedKeys, keys)
	}
}

func TestMapKeys(t *testing.T) {
	// 测试 MapKeys 函数，获取映射的所有键
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keySeq := MapKeys(m)

	keys := slices.Collect(keySeq)

	slices.Sort(keys)
	expected := []string{"a", "b", "c"}
	if !slices.Equal(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestMapValues(t *testing.T) {
	// 测试 MapValues 函数，获取映射的所有值
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	valueSeq := MapValues(m)

	values := slices.Collect(valueSeq)
	slices.Sort(values)

	expected := []int{1, 2, 3}
	if !slices.Equal(values, expected) {
		t.Errorf("Expected %v, got %v", expected, values)
	}
}

func TestSeqKeys(t *testing.T) {
	// 测试 SeqKeys 函数，从Seq2获取键序列
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	seq2 := FromMap(m)
	keySeq := SeqKeys[string, int]()(seq2)

	keys := slices.Collect(keySeq)
	slices.Sort(keys)

	expected := []string{"a", "b", "c"}
	if !slices.Equal(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestSeqValues(t *testing.T) {
	// 测试 SeqValues 函数，从Seq2获取值序列
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	seq2 := FromMap(m)
	valueSeq := SeqValues[string, int]()(seq2)

	values := slices.Collect(valueSeq)
	slices.Sort(values)

	expected := []int{1, 2, 3}
	if !slices.Equal(values, expected) {
		t.Errorf("Expected %v, got %v", expected, values)
	}
}

func TestSeq2Filter(t *testing.T) {
	// 测试 Seq2Filter 函数，过滤Seq2
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	seq2 := FromMap(m)

	filteredSeq2 := Seq2Filter(func(k string, v int) bool {
		return v > 2
	})(seq2)

	var keys []string
	var values []int
	filteredSeq2(func(k string, v int) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})

	if len(keys) != 2 {
		t.Errorf("Expected 2 items, got %d", len(keys))
	}
}

func TestEnumerate(t *testing.T) {
	// 测试 Enumerate 函数，添加索引
	items := []string{"a", "b", "c", "d"}
	seq := slices.Values(items)

	enumSeq := Enumerate[string]()(seq)

	var indices []int
	var values []string
	enumSeq(func(i int, v string) bool {
		indices = append(indices, i)
		values = append(values, v)
		return true
	})

	expectedIndices := []int{0, 1, 2, 3}
	if !slices.Equal(indices, expectedIndices) {
		t.Errorf("Expected indices %v, got %v", expectedIndices, indices)
	}

	if !slices.Equal(values, items) {
		t.Errorf("Expected values %v, got %v", items, values)
	}
}

func TestRepeatOne(t *testing.T) {
	// 测试 RepeatOne 函数，无限重复单个元素
	repeatSeq := RepeatOne(5)

	// 使用Take限制次数
	taken := Take[int](5)(repeatSeq)
	result := slices.Collect(taken)

	expected := []int{5, 5, 5, 5, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCollect(t *testing.T) {
	// 测试 Collect 函数，将序列转换为切片
	numbers := []int{1, 2, 3, 4, 5}
	seq := slices.Values(numbers)

	result := Collect[int]()(seq)

	if !slices.Equal(result, numbers) {
		t.Errorf("Expected %v, got %v", numbers, result)
	}
}

func TestUnZipSeq2(t *testing.T) {
	// 测试 UnZipSeq2 函数，将Seq2分解为两个Seq
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	seq2 := FromMap(m)

	keySeq, valSeq := UnZipSeq2(seq2)

	keys := slices.Collect(keySeq)
	vals := slices.Collect(valSeq)

	slices.Sort(keys)
	slices.Sort(vals)

	expectedKeys := []string{"a", "b", "c"}
	expectedVals := []int{1, 2, 3}

	if !slices.Equal(keys, expectedKeys) {
		t.Errorf("Expected keys %v, got %v", expectedKeys, keys)
	}

	if !slices.Equal(vals, expectedVals) {
		t.Errorf("Expected values %v, got %v", expectedVals, vals)
	}
}
