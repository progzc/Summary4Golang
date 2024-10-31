package __

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack"
)

// Protobuf编码原理及优化技巧探讨
// https://cloud.tencent.com/developer/article/2368061
func TestProtobuf(t *testing.T) {
	person := &Person{
		Name: "John Doe",
		Age:  30,
	}

	// 序列化
	data, err := proto.Marshal(person)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	// data: [10 8 74 111 104 110 32 68 111 101 16 30], len(data) = 12
	fmt.Printf("data: %v, len(data) = %d\n", data, len(data))

	// 反序列化
	newPerson := &Person{}
	if err := proto.Unmarshal(data, newPerson); err != nil {
		log.Fatalf("failed to unmarshal: %v", err)
	}
	// person: name:"John Doe"  age:30
	fmt.Printf("person: %+v\n", newPerson)
}

func TestMsgPack(t *testing.T) {
	person := &Person{
		Name: "John Doe",
		Age:  30,
	}

	// 序列化
	data, err := msgpack.Marshal(person)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	// data: [130 164 78 97 109 101 168 74 111 104 110 32 68 111 101 163 65 103 101 210 0 0 0 30], len(data) = 24
	fmt.Printf("data: %v, len(data) = %d\n", data, len(data))

	// 反序列化
	newPerson := &Person{}
	if err := msgpack.Unmarshal(data, newPerson); err != nil {
		log.Fatalf("failed to unmarshal: %v", err)
	}
	// person: name:"John Doe" age:30
	fmt.Printf("person: %+v\n", newPerson)
}

func TestJson(t *testing.T) {
	person := &Person{
		Name: "John Doe",
		Age:  30,
	}

	// 序列化

	data, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	// data: [123 34 110 97 109 101 34 58 34 74 111 104 110 32 68 111 101 34 44 34 97 103 101 34 58 51 48 125], len(data) = 28
	fmt.Printf("data: %v, len(data) = %d\n", data, len(data))

	// 反序列化
	newPerson := &Person{}
	if err := json.Unmarshal(data, newPerson); err != nil {
		log.Fatalf("failed to unmarshal: %v", err)
	}
	// person: name:"John Doe"  age:30
	fmt.Printf("person: %+v\n", newPerson)
}
