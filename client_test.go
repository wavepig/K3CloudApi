package sdk

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	newClient(InitConfig(0, "67525a835e7b0f", "administrator", "kingdee@123", "http://192.168.32.32/k3cloud/", "http://192.168.32.32/CloudPLMWarehouse/Download", "", ""))
}

func TestNewK3CloudApiSdk(t *testing.T) {
	sdk, _ := NewK3CloudApiSdk(
		AuthTypeAppSecret,
		"67cfe0c793cfb7",
		"administrator",
		//"kingdee@123",
		"",
		"http://192.168.32.32/k3cloud/",
		"http://192.168.32.32/CloudPLMWarehouse/Download",
		"303090_X68p59EJTuo9Tf/p5/Qo7ZVE0N7a1KLH",
		"ecb7d282810748d5bc436fa1da09f0d5")
	//str := `{"parameters": ["SUB_SUBREQORDER","{\"Number\":\"SUB00000014\",\"Id\":\"\"}"]}`
	postData := map[string]string{
		"Number": "SUB00000018",
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

func TestRequestAny(t *testing.T) {
	sdk, err := NewK3CloudApiSdk(
		AuthTypeAppSecret,
		"6806b4475169ac",
		"administrator",
		//"kingdee@123",
		"",
		"http://192.168.32.32/k3cloud/",
		"http://192.168.32.34/CloudPLMWarehouse/Download",
		"303090_X68p59EJTuo9Tf/p5/Qo7ZVE0N7a1KLH",
		"ecb7d282810748d5bc436fa1da09f0d5")
	postData := map[string]string{
		//"FileID": "MDAwYzI5ZmQtMTliNi05YjVlLTExZjAtMDk0ODFhOWIwNmQ1",
		//"FileID": "MDAwYzI5ZmQtMTliNi05YjVlLTExZjAtMDk0NzZhZjNhNTY1",
		"FileID": "MDAwYzI5ZWUtNTAzYS1iZWU2LTExZjAtMjVkM2NlMjZkOWI2",
		"token":  "876A1CD9-B944-C2EF-E824-07BC51D3E2AB",
	}
	body, err := sdk.RequestFile(postData)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer body.Close()

	all, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	//fmt.Println(all)
	//fmt.Println(string(all))

	f, err := os.OpenFile("file1.bin", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(all)

	//bufWriter := bufio.NewWriter(f)
	//
	//_, err = io.Copy(bufWriter, body)
	//if err != nil {
	//	return
	//}

	// 这里不要忘记最后把缓冲区中剩余的数据写入磁盘，默认情况下，4096byte后会自动进行一次磁盘写入
	// 比如文件为5000byte, flush会将缓冲区中的904byte写入磁盘中
	//bufWriter.Flush()
}
