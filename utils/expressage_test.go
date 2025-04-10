package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExpressage(t *testing.T) {
	expressage := Expressage{

		Url:      "https://poll.kuaidi100.com/poll/query.do",
		Customer: "",
		Key:      "",
	}
	body := &Body{
		Param: &Param{
			Num: "",
		},
	}

	expressage100, err := GetLogisticsExpressage100(expressage.Url, expressage.Customer, expressage.Key, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	indent, err := json.MarshalIndent(expressage100, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(indent))
}
