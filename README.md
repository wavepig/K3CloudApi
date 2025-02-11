# K3CloudApi
## About

- 金蝶Golang WebApi调用库

## Requirement

- [Golang](https://go.dev/dl/) 1.23.4

## Dependencies

## Usage

- Install

  ```shell
  $ go get -u github.com/wavepig/K3CloudApi
  ```

- test

  ```go
  func TestNewK3CloudApiSdk(t *testing.T) {
  	sdk := NewK3CloudApiSdk(
  		AuthTypeAppSecret,
  		"xxx",
  		"xxx",
  		"xxx",
  		"http://xxx.xxx.xxx.xxx/k3cloud/",
  		"xxxx",
  		"xxx")
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
  ```

  
