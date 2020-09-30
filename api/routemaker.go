package api

import (
	"github.com/golang/glog"
)

// RouteMaker supports the creation of dot notated keys
type RouteMaker struct {
	stack []string
}

// KeyStackDepth defines the maximum levels of recursion
const KeyStackDepth = 32

// Push adds to tail
func (rm *RouteMaker) Push(sch *Model) {
	if len(rm.stack) >= KeyStackDepth {
		glog.Error("RouteMaker.Push maxDepth exceeded")
	}
	rm.stack = append(rm.stack, sch.Label)
}

// Pop removes tail
func (rm *RouteMaker) Pop() {
	if len(rm.stack) < 1 {
		glog.Error("RouteMaker.Pop empty stack")
		return
	}
	rm.stack = rm.stack[:len(rm.stack)-1]
}

func (rm *RouteMaker) String() (s string) {
	sep := "."
	for _, k := range rm.stack {
		s += sep + k
	}
	return
}
