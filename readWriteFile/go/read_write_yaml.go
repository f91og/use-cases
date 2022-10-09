// go mod init readwriteyaml
// go get gopkg.in/yaml.v3

package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func main() {
	buf, err := ioutil.ReadFile("./file.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	fmt.Printf("--- buf:\n%v\n\n", buf) // [65 58 10 45 32 65 49 10 45 32 65 50 10 10 66 58 10 32 66 49 58 32 34 66 49 34 10 32 66 50 58 32 34 66 50 34 10]

	// []map[string]string にマッピング, 只有在原yaml file是这个数据结构时才可以，https://qiita.com/daiching/items/a384a5ce229cf714a3b5
	// 否则是要用 interface{} 来包装复杂数据类型的
	data := make(map[string]interface{}) // or make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- data:\n%v\n\n", data)            // map[A:[A1 A2] B:map[B1:B1 B2:B2]]
	fmt.Printf("--- data['A']: \n%v\n\n", data["A"]) // [A1 A2], 这里必须是 "A"，不能是'A',否则打印不出来

	// []struct にマッピング
	t := T{}
	err = yaml.Unmarshal([]byte(data2), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t) // {Easy! {2 [3 4]}}

	// 根据这个https://qiita.com/daiching/items/a384a5ce229cf714a3b5，
	// 好像可以struct数组格式读出来，如果原yaml里是某个特定数据结构的数组的话
}

var data2 = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to correctly populate the data.
type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}
