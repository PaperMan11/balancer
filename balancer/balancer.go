package balancer

import (
	"errors"
)

var (
	ErrorNoHost                = errors.New("no host")
	ErrorAlgorithmNotSupported = errors.New("algorithm not supported")
)

// Balancer 反向代理服务器
type Balancer interface {
	Add(string)                     // 添加主机
	Remove(string)                  // 删除主机
	Balance(string) (string, error) // 负载均衡
	Inc(string)                     // 代理主机连接数+1
	Done(string)                    // 代理主机连接数-1
}

// 工厂模式
type Factory func([]string) Balancer

var factories = make(map[string]Factory)

// Build 根据算法生成不同的 balancer
func Build(algorithm string, hosts []string) (Balancer, error) {
	factory, ok := factories[algorithm]
	if !ok {
		return nil, ErrorAlgorithmNotSupported
	}
	return factory(hosts), nil
}
