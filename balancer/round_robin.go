package balancer

import (
	"sync"
)

type RoundRobin struct {
	sync.RWMutex
	i     uint64
	hosts []string
}

var _ Balancer = (*RoundRobin)(nil)

func init() {
	factories[R2Balancer] = NewRoundRobin
}

func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		hosts: hosts,
		i:     0,
	}
}

// 添加主机
func (r *RoundRobin) Add(host string) {
	r.Lock()
	defer r.Unlock()
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

// 删除主机
func (r *RoundRobin) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return
		}
	}
}

// 负载均衡
func (r *RoundRobin) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", ErrorNoHost
	}
	host := r.hosts[r.i]
	r.i = (r.i + 1) % uint64(len(r.hosts))
	return host, nil
}

// 代理主机连接数+1
func (r *RoundRobin) Inc(_ string) {}

// 代理主机连接数-1
func (r *RoundRobin) Done(_ string) {}
