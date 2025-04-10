package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client = http.Client{Timeout: 10 * time.Second}
)

type Body struct {
	Customer string `json:"customer"`
	Sign     string `json:"sign"`
	Param    *Param `json:"param"`
}

type Param struct {
	Com      string `json:"com"`
	Num      string `json:"num"`
	Phone    string `json:"phone"`
	From     string `json:"from"`
	To       string `json:"to"`
	Resultv2 string `json:"resultv2"`
	Show     string `json:"show"`
	Order    string `json:"order"`
}

type Expressage100Req struct {
	Message string `json:"message"`
	Nu      string `json:"nu"`
	Ischeck string `json:"ischeck"`
	Com     string `json:"com"`
	Status  string `json:"status"`
	Data    []struct {
		Time       string `json:"time"`
		Context    string `json:"context"`
		Ftime      string `json:"ftime"`
		AreaCode   string `json:"areaCode"`
		AreaName   string `json:"areaName"`
		Status     string `json:"status"`
		Location   string `json:"location"`
		AreaCenter string `json:"areaCenter"`
		AreaPinYin string `json:"areaPinYin"`
		StatusCode string `json:"statusCode"`
	} `json:"data"`
	State     string `json:"state"`
	Condition string `json:"condition"`
	RouteInfo struct {
		From struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"from"`
		Cur struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"cur"`
		To struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"to"`
	} `json:"routeInfo"`
	IsLoop bool `json:"isLoop"`
}

type Expressage100ReqErr struct {
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
	Message    string `json:"message"`
}

type Expressage struct {
	Url      string `json:"url" yaml:"url"`
	Key      string `json:"key" yaml:"key"`
	Customer string `json:"customer" yaml:"customer"`
}

func GetLogisticsExpressage100(urls, customer, key string, body *Body) (*Expressage100Req, error) {
	body.Customer = customer
	body.Param.Resultv2 = "4" // 完整数据
	body.Param.Show = "0"     // json返回
	b, err := json.Marshal(body.Param)
	if err != nil {
		return nil, err
	}

	sign := string(b) + key + customer

	hash := md5.Sum([]byte(sign))
	md5str := hex.EncodeToString(hash[:])

	// 将 MD5 值转换为大写形式
	md5str = strings.ToUpper(md5str)
	body.Sign = md5str

	paramByte, err := json.Marshal(body.Param)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("customer", body.Customer)
	data.Set("sign", body.Sign)
	data.Set("param", string(paramByte))

	req, err := http.NewRequest("POST", urls, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rsp.Body.Close(); err != nil {
			fmt.Println("Close body failed!", err.Error())
		}
	}()

	if rsp.StatusCode == 200 {
		res, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		expressage100Req := &Expressage100Req{}
		err = json.Unmarshal(res, expressage100Req)
		if err != nil {
			return nil, err
		}
		return expressage100Req, nil
	} else {
		resErr, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		expressage100ReqErr := &Expressage100ReqErr{}
		err = json.Unmarshal(resErr, expressage100ReqErr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(expressage100ReqErr.Message)
	}
}
