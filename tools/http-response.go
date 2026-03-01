package tools

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type HttpResponse struct {
	Status int
	Message string
	Headers map[string]string
	Body []byte
}

func Send(conn net.Conn, status int, message string, body string) {
	res := HttpResponse{
		Status: status,
		Message: message,
		Headers: map[string]string {
			"Content-Type": "text/plain",
		},
		Body: []byte(body),
	}
	res.Send(conn)
}

func (r *HttpResponse) Serialize() string {
	const SEP = "\r\n"
	var content strings.Builder;
	fmt.Fprintf(&content, "HTTP/1.1 %v %v%v", r.Status, r.Message, SEP)
	for key, value := range r.Headers {
		fmt.Fprintf(&content, "%v: %v%v", key, value, SEP)
	}
	content.WriteString(SEP)
	content.WriteString(string(r.Body))
	return content.String()
}

func ParseResponse(data []byte) HttpResponse {
	res := HttpResponse{}
	res.Parse(data)
	return res
}

func (r *HttpResponse) Parse(data []byte) {
	req := ParseRequest(data)
	status, err := strconv.Atoi(req.Url)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	r.Status = status
	r.Message = req.Version
	r.Headers = req.Headers
	copy(r.Body, req.Body)
}

func (r *HttpResponse) Send(conn net.Conn) {
	content := r.Serialize()
	conn.Write([]byte(content))
}
