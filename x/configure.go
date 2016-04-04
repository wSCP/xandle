package x

import (
	"log"
	"sort"
)

type ConfigFn func(Handle) error

type Config interface {
	Order() int
	Configure(Handle) error
}

type config struct {
	order int
	fn    ConfigFn
}

func DefaultConfig(fn ConfigFn) Config {
	return config{50, fn}
}

func NewConfig(order int, fn ConfigFn) Config {
	return config{order, fn}
}

func (c config) Order() int {
	return c.order
}

func (c config) Configure(h Handle) error {
	return c.fn(h)
}

type List []Config

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	return l[i].Order() < l[j].Order()
}

type Configuration interface {
	Add(...Config)
	AddFn(...ConfigFn)
	Configure(Handle) error
	Configured() bool
}

type configuration struct {
	configured bool
	list       List
}

func newConfiguration(conf ...Config) *configuration {
	c := &configuration{
		list: make([]Config, 0),
	}
	c.Add(conf...)
	return c
}

func (c *configuration) Add(conf ...Config) {
	c.list = append(c.list, conf...)
}

func (c *configuration) AddFn(fns ...ConfigFn) {
	for _, fn := range fns {
		c.list = append(c.list, DefaultConfig(fn))
	}
}

func configure(h Handle, conf ...Config) error {
	for _, c := range conf {
		err := c.Configure(h)
		if err != nil {
			log.Printf("x:handle:configuration: %s", err.Error())
			return err
		}
	}
	return nil
}

func (c *configuration) Configure(h Handle) error {
	sort.Sort(c.list)

	err := configure(h, c.list...)
	if err == nil {
		c.configured = true
	}

	return err
}

func (c *configuration) Configured() bool {
	return c.configured
}
