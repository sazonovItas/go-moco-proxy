package connpool

import (
	"github.com/jackc/puddle/v2"
)

type Conn struct {
	p   *Pool
	res *puddle.Resource[*connResource]
}

// Conn method returns underlying network connection.
func (c *Conn) Conn() PoolConn {
	return c.connResource().PoolConn
}

// Release method returns connection to the pool.
func (c *Conn) Release() {
	if c.res == nil {
		return
	}

	res := c.res
	c.res = nil
	if c.p.isExpired(res) {
		res.Destroy()
		c.p.triggerHealthCheck()
		return
	}

	res.Release()
}

// Hijack mthod hijack connection from the pool or panics
// if connection already released or hijacked.
func (c *Conn) Hijack() PoolConn {
	if c.res == nil {
		panic("cannot hijack already released or hijacked connection")
	}

	conn := c.Conn()
	res := c.res
	c.res = nil

	res.Hijack()

	return conn
}

// connResource method returns undelying connection resource.
func (c *Conn) connResource() *connResource {
	if c.res == nil {
		return nil
	}

	return c.res.Value()
}
