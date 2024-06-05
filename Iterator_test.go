package fp

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestIterrator(t *testing.T) {
	// 使用 NewIntRange 函数生成迭代器
	s := "dagsdaga"
	person := &Person{Name: "Bob", Age: 25}
	err := printPersonWithoutEscape(person)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error printing person:", err)
	}

	println(&s)
	println(&person)
}

type Person struct {
	Name string
	Age  int
}

func (m Person) String() string {
	return fmt.Sprintf("MyStruct{Name: %s, Age: %d}", m.Name, m.Age)
}
func (m Person) GoString() string {
	return fmt.Sprintf("MyStruct{Name: %#v, Age: %#v}", m.Name, m.Age)
}
func printPersonWithoutEscape(p *Person) error {
	var output strings.Builder
	output.WriteString(p.String())
	output.Write([]byte("\n"))
	output.WriteString(p.GoString())
	output.Write([]byte("\n"))
	// 将内容写入标准输出，避免直接使用fmt.Println可能引起的逃逸。
	_, err := os.Stdout.WriteString(output.String())
	return err
}
