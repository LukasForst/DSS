package main

import (
	"encoding/json"
	"fmt"
)

type Hello struct {
	Bytes []byte
}

func main() {
	data := Hello{[]byte("hello world")}

	marshaled, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(marshaled))

	var receivedHello Hello
	if err := json.Unmarshal(marshaled, &receivedHello); err != nil {
		panic(err)
	}

	fmt.Println(string(receivedHello.Bytes))
}
