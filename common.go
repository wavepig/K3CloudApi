package sdk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	ValidateUser      = "Kingdee.BOS.WebApi.ServicesStub.AuthService.ValidateUser.common.kdsvc"
	ValidateAppSecret = "Kingdee.BOS.WebApi.ServicesStub.AuthService.LoginByAppSecret.common.kdsvc"
	headerMap         = map[string]string{
		"Cache-Control":  "no-cache",
		"Content-Type":   "application/json",
		"Accept-Charset": "utf-8",
		"User-Agent":     "Kingdee/Golang WebApi SDK 7.3 (compatible; MSIE 6.0; Windows NT 5.1;SV1)",
	}
	SetCookie = "Set-Cookie"
)

type Cloud struct {
	httpClient         *http.Client
	kDServiceSessionId *LoginResult
	lock               sync.RWMutex
}

func NewCloud(t int) *Cloud {
	return &Cloud{
		httpClient: &http.Client{
			Timeout: time.Duration(t) * time.Second,
		},
		kDServiceSessionId: new(LoginResult),
	}
}

func (c *Cloud) login(config *Config) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	serverUrl := config.ServerUrl
	ReqBody := new(strings.Reader)
	if config.AuthType == AuthTypePassword {
		serverUrl = serverUrl + ValidateUser
		ReqBody = strings.NewReader(`{"parameters":["` + config.AcctID + `","` + config.Username + `","` + config.Password + `",2052]}`)
	} else if config.AuthType == AuthTypeAppSecret {
		serverUrl = serverUrl + ValidateAppSecret
		ReqBody = strings.NewReader(`{"parameters":["` + config.AcctID + `","` + config.Username + `","` + config.AppID + `","` + config.AppSecret + `",2052]}`)
	}
	request, err := http.NewRequest("POST", serverUrl, ReqBody)
	if err != nil {
		return "", err
	}

	for k, v := range headerMap {
		request.Header.Add(k, v)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("send to server failed: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request failed: %s", resp.Status)
	}

	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &c.kDServiceSessionId)
	if err != nil {
		return "", err
	}

	return strings.Join(resp.Header[SetCookie], ";"), nil
}

func (c *Cloud) request(method string, url string, body []byte, cookie string) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	request, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	for k, v := range headerMap {
		request.Header.Add(k, v)
	}

	request.Header.Add("Cookie", cookie)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send to cloud error: %s", err.Error())
	}
	defer resp.Body.Close()

	repBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		logBody := ""
		if body != nil {
			logBody = string(body)
		}

		return nil, fmt.Errorf("send to cloud error url: %s reqbody:%s status: %d body: %s", url, logBody, resp.StatusCode, repBody)
	}

	return repBody, nil
}

func (c *Cloud) requestFile(url, loginUrl string, headers map[string]string, cookie string) (io.ReadCloser, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	request, err := http.NewRequest("POST", url, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	request.Header.Add("Cookie", cookie)

	ctx := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("LoginUrl=%s&KDServiceSessionId=%s", loginUrl, c.kDServiceSessionId.KDSVCSessionId)))
	request.Header["CTX"] = []string{ctx}
	request.Header["PLM_ACCESS_TYPE"] = []string{"pure"}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send to cloud error: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		repBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("send to cloud error url: %s status: %d body: %s", url, resp.StatusCode, repBody)
	}

	return resp.Body, nil
}

func hasTrailingSlash(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return strings.HasSuffix(u.Path, "/")
}
