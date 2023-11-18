package proxyserv

import (
	"apgo/modules/intercept"
	"apgo/system"
	"apgo/util/decoder"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
)

type ProxyImpl struct {
	sync.RWMutex
	Module   *Proxy
	Sess     *system.Session
	LogIDs   []int64
	LogEntry map[int64]*ProxyLog

	interceptor *intercept.Interceptor
}

// NewProxyImpl creates a new controller for the core intercetp proxy
func NewProxyImpl(proxy *Proxy, s *system.Session, interceptor *intercept.Interceptor) *ProxyImpl {
	c := &ProxyImpl{
		Module:      proxy,
		Sess:        s,
		LogEntry:    make(map[int64]*ProxyLog),
		interceptor: interceptor,
	}

	c.Module.OnReq = c.onReq
	c.Module.OnResp = c.onResp
	c.Module.Proxyh.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(c.connectionHandle))

	return c
}

// Executed when a response arrives
func (c *ProxyImpl) onResp(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

	if r != nil {
		bodyBytes, err := decoder.DecodeBodyResponse(r)
		if err != nil {
			log.Panicln("Error decoding response:", err)
		}

		r.ContentLength = int64(len(bodyBytes))

		c.RLock()
		logEntry, ok := c.LogEntry[ctx.Session]
		if !ok {
			// Handle the case where the log entry is not found
			return r
		}
		c.RUnlock()

		dump, err := httputil.DumpResponse(r, false)
		if err != nil {
			log.Panicln("Error dumping response:", err)
			return r
		}

		c.Lock()
		logEntry.ContentLength = r.ContentLength
		logEntry.Time = time.Now()
		logEntry.ResponseMessage = dump
		logEntry.ResponseBody = bodyBytes
		c.Unlock()

		c.interceptor.OnReceive(dump, bodyBytes, r, ctx)
	}

	return r
}

// Executed when a request arrives
func (c *ProxyImpl) onReq(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	var resp *http.Response
	if r != nil {
		bodyBytes, err := decoder.DecodeBodyRequest(r)
		if err != nil {
			log.Panicln("Error decoding request:", err)
		}

		r.ContentLength = int64(len(bodyBytes))

		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			log.Panicln("Error dumping request:", err)
			return r, resp
		}

		logEntry := &ProxyLog{
			ID:             ctx.Session,
			Method:         r.Method,
			URL:            r.URL.String(),
			ContentLength:  r.ContentLength,
			Time:           time.Now(),
			Target:         r.URL.Scheme + "://" + r.URL.Host,
			Host:           r.Host,
			RequestMessage: dump,
			RequestBody:    bodyBytes,
		}

		c.Lock()
		c.LogEntry[ctx.Session] = logEntry // Add the log entry to the log slice
		c.LogIDs = append(c.LogIDs, ctx.Session)
		c.Unlock()

		// Intercept the request
		r, resp = c.interceptorRequest(dump, bodyBytes, r, ctx)
	}

	return r, resp
}

func (c *ProxyImpl) connectionHandle(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.MitmConnect, host
}

func (c *ProxyImpl) interceptorRequest(messageByte []byte, bodyBytes []byte, req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	c.Lock()
	newRequest, newResponse := c.interceptor.OnInterception(messageByte, bodyBytes, req, ctx)
	c.Unlock()
	return newRequest, newResponse
}
