package etcdInit

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

var EtcdReg registry.Registry

func init() {
	EtcdReg = etcd.NewRegistry(
		registry.Addrs("172.19.0.10:2379"),
	)
}
