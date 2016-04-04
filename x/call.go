package x

import "sync"

type Call func(Event) error

type Caller interface {
	SetCaller(Xevent, Call)
	Call(Event) error
}

type caller struct {
	calls map[Xevent][]Call
	clk   *sync.RWMutex
}

func newCaller() *caller {
	return &caller{
		calls: make(map[Xevent][]Call),
		clk:   &sync.RWMutex{},
	}
}

func (c *caller) SetCaller(tag Xevent, fn Call) {
	c.clk.Lock()
	defer c.clk.Unlock()
	c.calls[tag] = append(c.calls[tag], fn)
}

func (c *caller) Call(evt Event) error {
	c.clk.Lock()
	defer c.clk.Unlock()
	if fns, ok := c.calls[evt.Tag]; ok {
		for _, fn := range fns {
			err := fn(evt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
