package registry

import (
	"net"

	consul "github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
)

type Client struct {
	inner *consul.Client
}

func IPAddr() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			if ipnet.IP.To4() != nil || ipnet.IP.To16() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, nil
}

func NewWithConfig(conf *consul.Config) (*Client, error) {
	inner, err := consul.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &Client{
		inner: inner,
	}, nil
}

func New(addr string) (*Client, error) {
	conf := consul.DefaultConfig()
	conf.Address = addr

	return NewWithConfig(conf)
}

// RegisterService takes a service definition and adds it to consul
// It disregards whatever ID was given and uses a uuid internal so as to prevent clashes
func (c *Client) RegisterService(svc *consul.AgentServiceRegistration) (string, error) {

	id := uuid.New()
	svc.ID = id

	return id, c.inner.Agent().ServiceRegister(svc)
}

func (c *Client) Register(name string, port int) (string, error) {

	ip, err := IPAddr()
	if err != nil {
		return "", err
	}

	reg := &consul.AgentServiceRegistration{
		Name:    name,
		Address: ip.String(),
		Port:    port,
	}

	return c.RegisterService(reg)
}

func (c *Client) DeRegister(id string) error {
	return c.inner.Agent().ServiceDeregister(id)
}
