package api

import (
	"github.com/golang/glog"
)

// KeyMaker supports the creation of dot notated keys
type KeyMaker struct {
	stack []string
}

// KeyStackDepth defines the maximum levels of recursion
const KeyStackDepth = 32

// Push adds to tail
func (km *KeyMaker) Push(sch *Model) {
	if len(km.stack) >= KeyStackDepth {
		glog.Error("keyMaker.Push maxDepth exceeded")
	}
	km.stack = append(km.stack, sch.Label)
}

// Pop removes tail
func (km *KeyMaker) Pop() {
	if len(km.stack) < 1 {
		glog.Error("keymaker.Pop empty stack")
		return
	}
	km.stack = km.stack[:len(km.stack)-1]
}

func (km *KeyMaker) String() (s string) {
	sep := "."
	for _, k := range km.stack {
		s += sep + k
	}
	return
}
