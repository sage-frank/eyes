package utility

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type F func(string, map[string]string) (any, error)

var ProxyDomain = viper.GetString("proxy.domain")

type IRRequest struct {
	domain string
	path   string
	params map[string]string
	logger *zap.SugaredLogger
}

func (r *IRRequest) Get() (any, error) {
	data := url.Values{}
	URL, err := url.Parse(strings.Trim(r.domain, "/") + "/" + strings.Trim(r.path, "/"))
	if err != nil {
		r.logger.Error("url.parse(r.domain): %+v", err)
		return nil, err
	}

	for k, v := range r.params {
		data.Set(k, v)
	}
	URL.RawQuery = data.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		r.logger.Error("http.Get(URL.String()) %+v", err)
		return nil, fmt.Errorf("http.Get(URL.String()): %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			r.logger.Error(" resp.Body.Close()", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("io.ReadAll ", zap.Error(err))
		return nil, fmt.Errorf(" io.ReadAll(resp.Body): %w", err)
	}

	var ret any
	err = json.Unmarshal(body, &ret)
	if err != nil {
		r.logger.Error("json.Unmarshal", zap.String("body", string(body)), zap.Error(err))
		return nil, fmt.Errorf("json.Unmarshal(body, &ret): %w", err)
	}
	return ret, nil
}

func (r *IRRequest) Post() (any, error) {
	data := url.Values{}
	for k, v := range r.params {
		data.Set(k, v)
	}
	reqBody := data.Encode()

	r.logger.Info("reqBody:  ", reqBody)
	URL := strings.Trim(r.domain, "/") + "/" + strings.Trim(r.path, "/")

	r.logger.Info("Post URL: " + URL)

	req, err := http.NewRequest("POST", URL, strings.NewReader(reqBody))
	if err != nil {
		r.logger.Error("Error creating HTTP request:", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.logger.Error("http.Post", zap.Error(err))
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			r.logger.Error("resp.Body.Close()", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error(" io.ReadAll", zap.Error(err))
		return nil, err
	}

	r.logger.Info("body:", body)

	var ret any
	err = json.Unmarshal(body, &ret)
	if err != nil {
		r.logger.Error("json.Unmarshal", zap.Error(err))
		return nil, err
	}

	return ret, nil
}

func IRGet(domain string, params map[string]string) (any, error) {
	data := url.Values{}
	URL, err := url.Parse(domain)
	if err != nil {
		zap.L().Error("IRGet url error: ", zap.String("url", domain), zap.Error(err))
		return nil, err
	}
	for k, v := range params {
		data.Set(k, v)
	}
	URL.RawQuery = data.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		zap.L().Error(" http.Get error ", zap.Error(err))
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			zap.L().Error(" resp.Body.Close ", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("io.ReadAll ", zap.Error(err))
		return nil, err
	}

	var ret any
	err = json.Unmarshal(body, &ret)
	if err != nil {
		zap.L().Error("json.Unmarshal", zap.String("body", string(body)), zap.Error(err))
		return nil, err
	}
	return ret, nil
}

func IRPost(domain string, params map[string]string) (any, error) {
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	reqBody := data.Encode()

	ProxyDomain = viper.GetString("proxy.domain")
	zap.L().Info("ProxyDomain:  " + ProxyDomain)

	uri := strings.Trim(ProxyDomain, "/") + "/" + strings.Trim(domain, "/")

	zap.L().Info("Post uri: " + uri)

	req, err := http.NewRequest("POST", uri, strings.NewReader(reqBody))
	if err != nil {
		zap.L().Error("Error creating HTTP request:", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		zap.L().Error("http.Post", zap.Error(err))
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			zap.L().Error("resp.body.close", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error(" io.ReadAll", zap.Error(err))
		return nil, err
	}

	var ret any
	err = json.Unmarshal(body, &ret)
	if err != nil {
		zap.L().Error("json.Unmarshal", zap.Error(err))
		return nil, err
	}
	return ret, nil
}

func DecorateRetry(times int, domain string, params map[string]string, f F) F {
	return func(s string, m map[string]string) (any, error) {
		resp, err := f(domain, params)
		if err == nil {
			return resp, err
		}
		for i := 1; i <= times; i++ {
			resp, err = f(domain, params)
			if err == nil {
				break
			} else {
				zap.L().Error(fmt.Sprintf("第%d次调用接口失败", i), zap.Error(err))
				fmt.Printf("第%d次调用接口失败\r\n", i)
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
		return resp, err
	}
}

func Retry(times int, domain string, params map[string]string, f F) (any, error) {
	resp, err := f(domain, params)
	if err == nil {
		return resp, err
	}

	for i := 1; i <= times; i++ {
		resp, err = f(domain, params)
		if err == nil {
			break
		} else {
			zap.L().Error(fmt.Sprintf("第%d次调用接口失败", i), zap.Error(err))
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
	return resp, err
}
