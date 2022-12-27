package balancer

import (
	"sync"

	fibHeap "github.com/starwander/GoFibonacciHeap"
)

func init() {
	factories[LeastLoadBalancer] = NewLeastLoad
}

// Tag .
func (h *host) Tag() interface{} { return h.name }

// Key .
func (h *host) Key() float64 { return float64(h.load) }

type LeastLoad struct {
	sync.RWMutex
	heap *fibHeap.FibHeap
}

func NewLeastLoad(hosts []string) Balancer {
	ll := &LeastLoad{heap: fibHeap.NewFibHeap()}
	for _, h := range hosts {
		ll.Add(h)
	}
	return ll
}

func (l *LeastLoad) Add(hostName string) {
	l.Lock()
	defer l.Unlock()
	if v := l.heap.GetValue(hostName); v != nil {
		return
	}
	_ = l.heap.InsertValue(&host{hostName, 0})
}

func (l *LeastLoad) Remove(hostName string) {
	l.Lock()
	defer l.Unlock()
	if v := l.heap.GetValue(hostName); v == nil {
		return
	}
	_ = l.heap.Delete(hostName)
}

func (l *LeastLoad) Balance(_ string) (string, error) {
	l.RLock()
	defer l.RUnlock()
	if l.heap.Num() == 0 {
		return "", ErrorNoHost
	}
	return l.heap.MinimumValue().Tag().(string), nil
}

func (l *LeastLoad) Inc(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	h := l.heap.GetValue(hostName)
	h.(*host).load++
	_ = l.heap.IncreaseKeyValue(h)
}

func (l *LeastLoad) Done(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	h := l.heap.GetValue(hostName)
	h.(*host).load--
	_ = l.heap.DecreaseKeyValue(h)
}
