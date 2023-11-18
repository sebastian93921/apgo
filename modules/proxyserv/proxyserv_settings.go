package proxyserv

type ProxySettings struct {
	IP            string
	Port          int
	Interceptor   bool
	ReqIntercept  bool
	RespIntercept bool
}

var Stg = &ProxySettings{
	IP:            "127.0.0.1",
	Port:          8080,
	Interceptor:   false,
	ReqIntercept:  true,
	RespIntercept: false,
}
