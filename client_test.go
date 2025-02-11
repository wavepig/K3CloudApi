package main

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	newClient(InitConfig(0, "67525a835e7b0f", "administrator", "kingdee@123", "http://192.168.32.32/k3cloud/", "", ""))
}

func TestNewK3CloudApiSdk(t *testing.T) {
	sdk := NewK3CloudApiSdk(
		AuthTypeAppSecret,
		"xxx",
		"xxx",
		"xxx",
		"http://xxx.xxx.xxx.xxx/k3cloud/",
		"xxxx",
		"xxx")
	//str := `{"parameters": ["SUB_SUBREQORDER","{\"Number\":\"SUB00000014\",\"Id\":\"\"}"]}`
	postData := map[string]string{
		"Number": "SUB00000014",
	}
	view, err := sdk.Request(View, "SUB_SUBREQORDER", postData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(view))
	view, err = sdk.Request(View, "SUB_SUBREQORDER", postData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(view))
}
