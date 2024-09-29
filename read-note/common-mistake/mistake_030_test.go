package common_mistake

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 不导出的 struct 字段无法被 encode
// (1) 以小写字母开头的字段成员是无法被外部直接访问的，
//
//	所以 struct 在进行 json、xml、gob 等格式的 encode 操作时，这些私有字段会被忽略，decode时得到零值：
func TestMistake_030(t *testing.T) {
	wrong030()
	right030()
}

func wrong030() {
}

func right030() {
	type MyData struct {
		One int
		two string
	}
	in := MyData{1, "two"}
	fmt.Printf("%#v\n", in) // common_mistake.MyData{One:1, two:"two"}

	encoded, _ := json.Marshal(in)
	fmt.Println(string(encoded)) // {"One":1}    // 私有字段 two 被忽略了

	var out MyData
	json.Unmarshal(encoded, &out)
	fmt.Printf("%#v\n", out) // common_mistake.MyData{One:1, two:""}
}
