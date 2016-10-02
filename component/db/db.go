package db

import "github.com/garyburd/redigo/redis"

type Connector interface {
	Connect(network, address string, options ...redis.DialOption) (redis.Conn, error)
}

type simpleCon struct {
}

func (simpleCon) Connect(network, address string, options ...redis.DialOption) (redis.Conn, error)  {
	return redis.Dial(network, address, options...)
}

func NewConnector() Connector {
	return new (simpleCon)
}
