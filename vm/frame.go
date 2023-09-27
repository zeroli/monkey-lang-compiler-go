package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	cl *object.Closure
	ip int
	bp int
}

func NewFrame(cl *object.Closure, sp int) *Frame {
	return &Frame{cl: cl, ip: 0, bp: sp}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
