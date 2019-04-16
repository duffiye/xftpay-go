package xftpay

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ApiBackend api相关的后端类型
type ApiBackend struct {
	Type       SupportedBackend
	URL        string
	HTTPClient *http.Client
}

// Call 后端处理请求方法
func (s ApiBackend) Call(method, path, key string, form *url.Values, params []byte, v interface{}) error {
	var body io.Reader
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" {
		body = bytes.NewBuffer(params)
	}

	if strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "DELETE" {
		if form != nil && len(*form) > 0 {
			data := form.Encode()
			path += "?" + data
		}
	}

	req, err := s.NewRequest(method, path, key, "application/json", body, params)

	if err != nil {
		return err
	}

	if err := s.Do(req, v); err != nil {
		return err
	}

	return nil
}

// NewRequest 建立http请求对象
func (s *ApiBackend) NewRequest(method, path, key, contentType string, body io.Reader, params []byte) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = s.URL + path
	req, err := http.NewRequest(method, path, body)
	if LogLevel > 2 {
		log.Printf("Request to pingpp is : \n %v\n", req)
	}

	if err != nil {
		if LogLevel > 0 {
			log.Printf("Cannot create pingpp request: %v\n", err)
		}
		return nil, err
	}
	var dataToBeSign string
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" {
		dataToBeSign = string(params)
	}
	requestTime := fmt.Sprintf("%d", time.Now().Unix())
	req.Header.Set("Xft-Request-Timestamp", requestTime)
	dataToBeSign = dataToBeSign + req.URL.RequestURI() + requestTime

	req.Header.Add("Xft-Version", apiVersion)
	req.Header.Add("User-Agent", "Xft go SDK version:"+Version())
	req.Header.Add("X-Xft-Client-User-Agent", OsInfo)
	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// Do 处理 http 请求
func (s *ApiBackend) Do(req *http.Request, v interface{}) error {
	if LogLevel > 1 {
		log.Printf("Requesting %v %v \n", req.Method, req.URL.String())
	}
	retryTimes := 1
	start := time.Now()

	var reqBody []byte
	var err error
	if req.Body != nil {
		reqBody, err = ioutil.ReadAll(req.Body)
		fmt.Printf("req.Body %s\n", reqBody)
		if err != nil {
			return err
		}
	}

	for i := 0; i <= retryTimes; i++ {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		res, err := s.HTTPClient.Do(req)
		fmt.Printf("Resp Body is %v\n", res)
		if LogLevel > 0 {
			log.Printf("Request to pingpp completed in %v\n", time.Since(start))
		}
		if err != nil {
			if LogLevel > 0 {
				log.Printf("Request to Pingpp failed: %v\n", err)
			}
			return err
		}
		defer res.Body.Close()
		if res.StatusCode == 502 {
			continue
		}

		resBody, err := ioutil.ReadAll(res.Body)
		fmt.Printf("req.Body %s\n", resBody)

		if err != nil {
			if LogLevel > 0 {
				log.Printf("Cannot parse Pingpp response: %v\n", err)
			}
			return err
		}

		if res.StatusCode >= 400 {
			var errMap map[string]interface{}
			JsonDecode(resBody, &errMap)

			if e, ok := errMap["error"]; !ok {
				log.Printf("%s", e)

				err := errors.New(string(resBody))
				if LogLevel > 0 {
					log.Printf("Unparsable error returned from Pingpp: %v\n", err)
				}
				return err
			}
		}

		if LogLevel > 2 {
			log.Printf("resBody from pingpp API: \n%v\n", string(resBody))
		}

		if v != nil {
			return JsonDecode(resBody, v)
		}
		return nil
	}
	return nil
}
