package fp

import (
	"fmt"
	"slices"
	"testing"
)

func TestIterrator(t *testing.T) {

	si := slices.Values(slices.Repeat([]int{1, 2, 3, 4}, 3))

	s1, s2 := Split(si)

	fmt.Println(slices.Collect(Take[int](2)(s1)))
	fmt.Println(slices.Collect(Drop[int](2)(s2)))
	fmt.Println(slices.Collect(Drop[int](0)(Compose[int, int, int](Map(func(i int) int { return i / 2 }), Filter(func(i int) bool { return i%2 == 0 }))(s1))))
}
