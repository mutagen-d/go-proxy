package tools

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HttpRequest struct {
	Method  string
	Url     string
	Version string
	Headers map[string]string
	Body    []byte
}

func ParseRequest(data []byte) HttpRequest {
	req := HttpRequest{}
	req.Parse(data)
	return req
}

func IsSupportedMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodConnect:
		return true
	case http.MethodDelete:
		return true
	case http.MethodGet:
		return true
	case http.MethodHead:
		return true
	case http.MethodOptions:
		return true
	case http.MethodPatch:
		return true
	case http.MethodPost:
		return true
	case http.MethodPut:
		return true
	case http.MethodTrace:
		return true
	}
	return false
}

type FixHttp struct {
	Reader io.Reader
	buf    []byte
}

func (r *FixHttp) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if n > 0 {
		// fmt.Printf("fixed:\n%v\n", string(p[:n]))
		transformed, ferr := FixHttpRequest(p[:n])
		if ferr != nil {
			return n, ferr
		}
		copy(p, transformed)
		n = len(transformed)
	}
	return n, err
}
func FixFast(data []byte, origin string) []byte {
	s := string(data)
	r := fmt.Sprintf("http://%v", origin)
	return []byte(strings.Replace(s, r, "", 1))
}
func FixHttpRequest(data []byte) ([]byte, error) {
	req := HttpRequest{}
	req.Parse(data)
	if !strings.HasPrefix(req.Url, "http:") &&
		!strings.HasPrefix(req.Url, "ws:") &&
		!strings.HasPrefix(req.Url, "wss:") &&
		!strings.HasPrefix(req.Url, "https:") {
		return data, nil
	}
	u, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}
	xHost := req.Headers["host"]
	if xHost == "" {
		req.Headers["host"] = u.Host
	}
	req.Url = u.Path
	return []byte(req.Serialize()), nil
}

func (r *HttpRequest) Log() string {
	const SEP = "\r\n"
	var content strings.Builder
	fmt.Fprintf(&content, "%v %v HTTP/1.1%v", r.Method, r.Url, SEP)
	for key, value := range r.Headers {
		fmt.Fprintf(&content, "%v: %v%v", key, value, SEP)
	}
	fmt.Fprintf(&content, "%v: %v%v", "Body-Length", len(r.Body), SEP)
	content.WriteString(SEP)
	return content.String()
}

func (r *HttpRequest) Serialize() string {
	const SEP = "\r\n"
	var content strings.Builder
	fmt.Fprintf(&content, "%v %v HTTP/1.1%v", r.Method, r.Url, SEP)
	for key, value := range r.Headers {
		fmt.Fprintf(&content, "%v: %v%v", key, value, SEP)
	}
	content.WriteString(SEP)
	content.WriteString(string(r.Body))
	return content.String()
}

func (r *HttpRequest) Parse(data []byte) {
	const LF = 0x0A
	const CR = 0x0D
	var headers []string
	for i, offset := 0, 0; i < len(data); i += 1 {
		if data[i] == LF && i < len(data)-1 && data[i+1] == LF {
			copy(r.Body, data[i+2:])
			break
		}
		if data[i] == LF {
			header := strings.Trim(string(data[offset:i]), " \t")
			if header != "" {
				headers = append(headers, header)
			}
			offset = i + 1
		}
	}
	for i := 1; i < len(headers); i += 1 {
		header := headers[i]
		found := strings.Contains(header, ":")
		if found {
			params := strings.Split(header, ":")
			name := strings.ToLower(params[0])
			if r.Headers == nil {
				r.Headers = make(map[string]string)
			}
			r.Headers[name] = strings.Trim(strings.Join(params[1:], ":"), " \r\n\t")
		}
	}
	params := strings.Fields(headers[0])
	len := len(params)
	r.Method = params[0]
	if len > 1 {
		r.Url = params[1]
	}
	if len > 2 {
		r.Version = strings.Join(params[2:], " ")
	}
}
