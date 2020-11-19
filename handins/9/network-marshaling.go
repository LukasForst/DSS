package main

import (
	"encoding/json"
	"log"
)

func UnmarshalString(m json.RawMessage) string {
	var str string
	if err := json.Unmarshal(m, &str); err != nil {
		log.Fatal("It was not possible to decode string.")
	}
	return str
}

func MessageTypeAndRest(objmap map[string]json.RawMessage) (string, json.RawMessage) {
	typeJ, _ := objmap["Type"]
	messageType := UnmarshalString(typeJ)
	dataJ, _ := objmap["Data"]
	return messageType, dataJ
}
