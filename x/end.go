package x

type Ender interface {
	End()
	Ending() bool
}

type ender struct {
	end bool
}

func newEnder() *ender {
	return &ender{}
}

func (e *ender) End() {
	e.end = true
}

func (e *ender) Ending() bool {
	return e.end
}
