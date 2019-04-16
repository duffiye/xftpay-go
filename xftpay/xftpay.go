package xftpay

import (
	"bytes"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	apiBase            = "https://xlink.91xft.cn"
	defaultHTTPTimeout = 60 * time.Second
	apiVersion         = "2017-06-30"

	TotalBackends                  = 1
	APIBackend    SupportedBackend = "api"
)

var (
	Key        string
	httpClient = &http.Client{
		Timeout: defaultHTTPTimeout,
	}
	backends Backends
	OsInfo   string
	LogLevel = 2
)

type SupportedBackend string

type Backend interface {
	Call(method, path, key string, body *url.Values, params []byte, v interface{}) error
}

type Backends struct {
	API Backend
}

func Version() string {
	return apiVersion
}

// 通过不同的参数获取不同的后端对象
func GetBackend(backend SupportedBackend) Backend {
	var ret Backend
	switch backend {
	case APIBackend:
		if backends.API == nil {
			backends.API = ApiBackend{backend, apiBase, httpClient}
		}

		ret = backends.API
	}
	return ret
}

//设定后端处理对象
func SetBackend(backend SupportedBackend, b Backend) {
	switch backend {
	case APIBackend:
		backends.API = b
	}
}

func init() {
	var uname string
	switch runtime.GOOS {
	case "windows":
		uname = "windows"
	default:
		cmd := exec.Command("uname", "-a")
		cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()
		uname = out.String()
	}
	m := map[string]interface{}{
		"lang":             "golang",
		"lang_version":     runtime.Version(),
		"bindings_version": Version(),
		"publisher":        "pingpp",
		"uname":            uname,
	}
	content, _ := JsonEncode(m)
	OsInfo = string(content)
}
