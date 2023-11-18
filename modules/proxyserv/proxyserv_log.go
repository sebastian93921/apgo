package proxyserv

import (
	"time"
)

type ProxyLog struct {
	ID              int64
	Method          string
	URL             string
	RequestMessage  []byte
	RequestBody     []byte
	ResponseMessage []byte
	ResponseBody    []byte
	ContentType     string
	ContentLength   int64
	Target          string
	Host            string
	Time            time.Time
}
