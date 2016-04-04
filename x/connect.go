package x

import "github.com/BurntSushi/xgb"

type Connection interface {
	Conn() *xgb.Conn
}

type connection struct {
	conn *xgb.Conn
}

func newConnection(c *xgb.Conn) *connection {
	return &connection{c}
}

func (c *connection) Conn() *xgb.Conn {
	return c.conn
}
