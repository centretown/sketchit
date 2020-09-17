package api

import (
	"fmt"
	"strings"
)

type indent int

// Traveler - interface
type Traveler interface {
	Push(*Model)
	Pop()
	String() string
}

// Travel the schema after performing supplied function
func (sch *Model) Travel(traveler Traveler,
	f func(sch *Model, traveler Traveler)) {

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

func (i *indent) Push(*Model) {
	*i++
}

func (sch *Model) showSchema() {

	var f = func(s *Model, l Traveler) {
		fmt.Printf("%s%s\n", l, s.Title)
		fmt.Printf("%s Name: %v\n", l, s.Label)
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
