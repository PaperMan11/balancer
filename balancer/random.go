package balancer

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	factories[RandomBalancer] = NewRandom
}

type Random struct {
	sync.RWMutex
	hosts []string
	rnd   *rand.Rand
}

var _ Balancer = (*Random)(nil)

func NewRandom(hosts []string) Balancer {
	return &Random{
		hosts: hosts,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// 添加主机
func (r *Random) Add(host string) {
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
func (r *Random) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
		}
	}
}

// 负载均衡
func (r *Random) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", ErrorNoHost
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}

// 代理主机连接数+1
func (r *Random) Inc(_ string) {}

// 代理主机连接数-1
func (r *Random) Done(_ string) {}
