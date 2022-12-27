package proxy

import (
	"balancer/balancer"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	// eg. accept-encoding 转换为 Accept-Encoding
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")       // 记录真实 client ip
	XProxy        = http.CanonicalHeaderKey("X-Proxy")         // Balancer-Reverse-Proxy
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For") // 请求链路是 client -> proxy1 -> proxy2 -> webapp
)

var ReverseProxy = "Balancer-Reverse-Proxy"

type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer

	sync.RWMutex // protect alive
	alive        map[string]bool
}

func NewHTTPProxy(targetHosts []string, algorithm string) (*HTTPProxy, error) {
	hosts := make([]string, 0)
	hostMap := make(map[string]*httputil.ReverseProxy)
	alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		url, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(url)

		originDirector := proxy.Director
		proxy.Director = func(req *http.Request) { // 修改 request，新增标签
			originDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
		}

		host := GetHost(url)
		alive[host] = true
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}

	lb, err := balancer.Build(algorithm, hosts) // 根据算法构建负载均衡器
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{
		hostMap: hostMap,
		lb:      lb,
		alive:   alive,
	}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy cause panic: %s", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()

	host, err := h.lb.Balance(GetIP(r))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}
	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, r) // 请求转发
}
