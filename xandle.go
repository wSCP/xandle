package xandle

import (
	"github.com/wSCP/xandle/x"
)

type Xandle interface {
	x.Handle
	Swap(x.Handle)
}

type xandle struct {
	x.Handle
}

func New(h x.Handle) Xandle {
	return &xandle{h}
}

func (xh *xandle) Swap(h x.Handle) {
	xh.Handle = h
}
