package etcd

import clientv3 "go.etcd.io/etcd/client/v3"

type IEtcdService interface {
	GetEtcd(key string) (string, error)
}

type EtcdService struct {
	cli *clientv3.Client
}

func (e EtcdService) GetEtcd(key string) (string, error) {
	resp, err := e.cli.Get(e.cli.Ctx(), key)
	if err != nil {
		return "", err
	}
	return string(resp.Kvs[0].Value), nil
}

func NewEtcdService(cli *clientv3.Client) *EtcdService {
	return &EtcdService{cli: cli}
}
