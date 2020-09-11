package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// FakeSchema -
type FakeSchema struct {
	// schema title for presentation
	Title string `yaml:"Title,omitempty" json:"Title,omitempty" xml:"Title,omitempty"`
	// schema name used for keyed lookup
	Name string `yaml:"Name,omitempty" json:"Name,omitempty" xml:"Name,omitempty"`
	// type of data represented
	Type string `yaml:"Type,omitempty" json:"Type,omitempty" xml:"Type,omitempty"`
	// description of the schema
	Description string `yaml:"Description,omitempty" json:"Description,omitempty" xml:"Description,omitempty"`
	// uniqueItems constrains items to be unique
	UniqueItems bool `yaml:"UniqueItems,omitempty" json:"UniqueItems,omitempty" xml:"UniqueItems,omitempty"`
	// item list definition
	Items *Schema `yaml:"Items,omitempty" json:"Items,omitempty" xml:"Items,omitempty"`
	// required properties and order
	Required []string `yaml:"Required,omitempty" json:"Required,omitempty" xml:"Required,omitempty"`
	// options constrains schema to an array of choices
	Options []string `yaml:"Options,omitempty" json:"Options,omitempty" xml:"Options,omitempty"`
	// oneOf constrains schema to one of a selection of options
	OneOf []*Schema `yaml:"OneOf,omitempty" json:"OneOf,omitempty" xml:"OneOf,omitempty"`
	// properties defines an ordered list of children
	// order determined by required array
	Properties []*Schema `yaml:"Properties,omitempty" json:"Properties,omitempty" xml:"Properties,omitempty"`
}

// SchemaReducer -
type SchemaReducer int

const (
	// ReduceNone includes all detail
	ReduceNone SchemaReducer = iota
	// ReduceSummary includes less detail
	ReduceSummary
	// ReduceBrief includes name only
	ReduceBrief
)

var schemaReducer = ReduceNone

// SetReducer -
func (sch *Schema) SetReducer(r SchemaReducer) {
	schemaReducer = r
}

// Reduce schema within fake based on schemaReducer
func (sch *Schema) Reduce() (out interface{}) {
	fake := &FakeSchema{}
	out = fake

	// always include linkages
	fake.Items = sch.Items
	fake.OneOf = sch.OneOf
	fake.Properties = sch.Properties

	fake.Name = sch.Name
	if schemaReducer == ReduceBrief {
		return
	}

	fake.Type = sch.Type
	fake.Description = sch.Description
	if schemaReducer == ReduceSummary {
		return
	}

	fake.Title = sch.Title
	fake.UniqueItems = sch.UniqueItems
	fake.Required = sch.Required
	fake.Options = sch.Options
	return
}

// MarshalYAML produces yaml output for schema
func (sch *Schema) MarshalYAML() (out interface{}, err error) {
	out = sch.Reduce()
	return
}

// MarshalJSON produces json output for schema
func (sch *Schema) MarshalJSON() (b []byte, err error) {
	out := sch.Reduce()
	b, err = json.Marshal(out)
	return
}

// MarshalXML - xml Marshaler interface
func (sch *Schema) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	out := sch.Reduce()
	err = e.EncodeElement(out, start)
	return
}

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

// Dictionary is keyed map to data model
// keys are dot separated
type Dictionary map[string]*Schema

// DictionaryNew -
func DictionaryNew(collections []*Collection) (dict Dictionary) {
	dict = make(Dictionary)

	var f = func(s *Schema, t Traveler) {
		t.Push(s)
		dict[t.String()] = s
		t.Pop()
	}

	var km = KeyMaker{
		stack: make([]string, 0, maxStackDepth),
	}

	for _, c := range collections {
		c.Schema.Travel(&km, f)
	}
	return
}
