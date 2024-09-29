package common_mistake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// 将 JSON 中的数字解码为 interface 类型
// (1) 在 encode/decode JSON 数据时，Go 默认会将数值当做 float64 处理; 若 decode 的 JSON 字段是整型解决办法如下：
//     a. 将 int 值转为 float 统一使用
//     b. 将 decode 后需要的 float 值转为 int 使用
//     c. 使用 Decoder 类型来 decode JSON 数据，明确表示字段的值类型
//     d. 使用 struct 类型将你需要的数据映射为数值型
//     e. 可以使用 struct 将数值类型映射为 json.RawMessage 原生数据类型,
//    	  适用于如果 JSON 数据不着急 decode 或 JSON 某个字段的值类型不固定等情况：

func TestMistake_038(t *testing.T) {
	wrong038()
	right038_1()
	right038_2()
	right038_3()
	right038_4()
	right038_5()
}

// 在 encode/decode JSON 数据时，Go 默认会将数值当做 float64 处理
// 下面代码会造成panic: panic: interface conversion: interface {} is float64, not int
func wrong038() {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T\n", result["status"]) // float64
	var status = result["status"].(int)  // 类型断言错误
	fmt.Println("Status value: ", status)
}

func right038_1() {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalln(err)
	}

	// b. 将 decode 后需要的 float 值转为 int 使用
	var status = uint64(result["status"].(float64))
	fmt.Println("Status value: ", status)
}

func right038_2() {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}

	var decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&result); err != nil {
		log.Fatalln(err)
	}

	// c. 使用 Decoder 类型来 decode JSON 数据，明确表示字段的值类型
	var status, _ = result["status"].(json.Number).Int64()
	fmt.Println("Status value: ", status)
}

func right038_3() {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}
	var decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		log.Fatalln(err)
	}
	var status uint64
	// 使用 string 来存储数值数据，在 decode 时再决定按 int 还是 float 使用
	err := json.Unmarshal([]byte(result["status"].(json.Number).String()), &status)
	checkError(err)
	fmt.Println("Status value: ", status)
}

func right038_4() {
	var data = []byte(`{"status": 200}`)
	// d. 使用 struct 类型将你需要的数据映射为数值型
	var result struct {
		Status uint64 `json:"status"`
	}

	err := json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	checkError(err)
	fmt.Printf("Result: %+v", result)
}

func right038_5() {
	records := [][]byte{
		[]byte(`{"status":200, "tag":"one"}`),
		[]byte(`{"status":"ok", "tag":"two"}`),
	}

	// e. 可以使用 struct 将数值类型映射为 json.RawMessage 原生数据类型,
	//    适用于如果 JSON 数据不着急 decode 或 JSON 某个字段的值类型不固定等情况：
	for idx, record := range records {
		var result struct {
			StatusCode uint64
			StatusName string
			Status     json.RawMessage `json:"status"`
			Tag        string          `json:"tag"`
		}

		err := json.NewDecoder(bytes.NewReader(record)).Decode(&result)
		checkError(err)

		var name string
		err = json.Unmarshal(result.Status, &name)
		if err == nil {
			result.StatusName = name
		}

		var code uint64
		err = json.Unmarshal(result.Status, &code)
		if err == nil {
			result.StatusCode = code
		}

		fmt.Printf("[%v] result => %+v\n", idx, result)
	}
}
