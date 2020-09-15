package api

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

type indent int

// Traveler - interface
type Traveler interface {
	Push(*Schema)
	Pop()
	String() string
}

// Travel the schema after performing supplied function
func (sch *Schema) Travel(traveler Traveler,
	f func(sch *Schema, traveler Traveler)) {

	f(sch, traveler)

	if len(sch.OneOf) > 1 {
		traveler.Push(sch)
		for _, s := range sch.OneOf {
			s.Travel(traveler, f)
		}
		traveler.Pop()
	}

	traveler.Push(sch)
	if sch.Items != nil {
		sch.Items.Travel(traveler, f)
	}
	if len(sch.Properties) > 0 {
		for _, s := range sch.Properties {
			s.Travel(traveler, f)
		}
	}
	traveler.Pop()
}

func (i *indent) String() string {
	return strings.Repeat("  ", int(*i))
}

func (i *indent) Pop() {
	*i--
}

func (i *indent) Push(*Schema) {
	*i++
}

func (sch *Schema) showSchema() {

	var f = func(s *Schema, l Traveler) {
		fmt.Printf("%s%s\n", l, s.Title)
		fmt.Printf("%s Name: %v\n", l, s.Name)
		fmt.Printf("%s Type: %v\n", l, s.Type)
		fmt.Printf("%s Description: %v\n", l, s.Description)
		if len(s.Required) > 1 {
			fmt.Printf("%s Required: %v\n", l, s.Required)
		}
		if len(s.Options) > 0 {
			fmt.Printf("%s Options: %v\n", l, s.Options)
		}
	}

	var level = indent(0)
	sch.Travel(&level, f)
}

// KeyMaker supports the creation of dot notated keys
type KeyMaker struct {
	stack []string
}

const maxStackDepth = 32

// Push adds to tail
func (km *KeyMaker) Push(sch *Schema) {
	if len(km.stack) >= maxStackDepth {
		glog.Error("keyMaker.Push maxDepth exceeded")
	}
	km.stack = append(km.stack, sch.Name)
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
