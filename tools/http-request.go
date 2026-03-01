package tools

import (
	"net/http"
	"strings"
)

type HttpRequest struct {
	Method string
	Url string
	Version string
	Headers map[string]string
	Body []byte
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

func (r *HttpRequest) Parse(data []byte) {
	const LF = 0x0A
	const CR = 0x0D
	var headers []string
	for i, offset := 0, 0; i < len(data); i += 1 {
		if data[i] == LF &&  i < len(data) - 1 && data[i + 1] == LF {
			copy(r.Body, data[i + 2:])
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
