package balancer

import (
	"hash/crc32"
	"sync"
)

func init() {
	factories[IPHashBalancer] = NewIPHash
}

type IPHash struct {
	sync.RWMutex
	hosts []string
}

var _ Balancer = (*IPHash)(nil)

func NewIPHash(hosts []string) Balancer {
	return &IPHash{
		hosts: hosts,
	}
}

// 添加主机
func (r *IPHash) Add(host string) {
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
func (r *IPHash) Remove(host string) {
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
func (r *IPHash) Balance(host string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", ErrorNoHost
	}
	i := crc32.ChecksumIEEE([]byte(host)) % uint32(len(r.hosts))
	return r.hosts[i], nil
}

// 代理主机连接数+1
func (r *IPHash) Inc(_ string) {}

// 代理主机连接数-1
func (r *IPHash) Done(_ string) {}
