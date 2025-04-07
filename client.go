package sdk

import (
	"io"
	"sync"
	"time"
)

type AuthType uint8

const (
	AuthTypePassword AuthType = iota
	AuthTypeAppSecret
)

type Config struct {
	AuthType       AuthType
	AcctID         string
	Username       string
	Password       string
	ServerUrl      string
	DownloadUrl    string
	AppID          string
	AppSecret      string
	Lcid           int
	OrgNum         int
	ConnectTimeout int
}

func InitConfig(at AuthType, acctID, username, password, serverUrl, downloadUrl, appid, appSecret string) *Config {
	ok := hasTrailingSlash(serverUrl)
	if !ok {
		serverUrl = serverUrl + "/"
	}

	return &Config{
		AcctID:         acctID,
		Username:       username,
		Password:       password,
		ServerUrl:      serverUrl,
		DownloadUrl:    downloadUrl,
		AuthType:       at,
		AppID:          appid,
		AppSecret:      appSecret,
		Lcid:           2052,
		OrgNum:         0,
		ConnectTimeout: 50,
	}
}

func (c *Config) IsValid() bool {
	if c.AcctID == "" || c.Username == "" || c.ServerUrl == "" {
		return false
	}
	if c.Password == "" && (c.AppID == "" || c.AppSecret == "") {
		return false
	}

	return true
}

type client struct {
	config       *Config
	cookiesStore *cookiesStore
	connect      *Cloud
}

func newClient(config *Config) *client {
	if config == nil {
		panic("config is nil")
	}
	if config.IsValid() == false {
		panic("config is invalid")
	}

	cloud := NewCloud(config.ConnectTimeout)
	loginCookies, err := cloud.login(config)
	if err != nil {
		panic(err)
	}

	cookies := newCookies()
	cookies.Set(loginCookies)

	return &client{
		config:       config,
		cookiesStore: cookies,
		connect:      cloud,
	}
}

func (c *client) getCookies() (string, error) {
	cookies, ok := c.cookiesStore.Get()
	if !ok {
		var err error
		cookies, err = c.connect.login(c.config)
		if err != nil {
			return "", err
		}

		c.cookiesStore.Set(cookies)
	}

	return cookies, nil
}

func (c *client) request(method string, path string, body []byte) ([]byte, error) {
	cookies, err := c.getCookies()
	if err != nil {
		return nil, err
	}
	return c.connect.request(method, c.config.ServerUrl+path, body, cookies)
}

func (c *client) requestFile(headerMap map[string]string) (io.ReadCloser, error) {
	cookies, err := c.getCookies()
	if err != nil {
		return nil, err
	}
	return c.connect.requestFile(c.config.DownloadUrl, c.config.ServerUrl, headerMap, cookies)
}

type cookiesStore struct {
	lock         sync.RWMutex
	isInitialize bool
	timeUnix     int64
	cookies      string
}

func newCookies() *cookiesStore {
	return &cookiesStore{
		isInitialize: false,
	}
}

func (c *cookiesStore) Set(cookies string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.isInitialize = true
	c.timeUnix = time.Now().Unix()
	c.cookies = cookies
}

func (c *cookiesStore) Get() (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if time.Now().Unix()-c.timeUnix > 120 {
		return "", false
	}
	return c.cookies, true
}
