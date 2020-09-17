package api

// ModelFax resquences Model and add tags
type ModelFax struct {
	// brief
	Name string `yaml:"Name,omitempty" json:"Name,omitempty" xml:"Name,omitempty"`
	Type string `yaml:"Type,omitempty" json:"Type,omitempty" xml:"Type,omitempty"`
	// summary
	Description string   `yaml:"Description,omitempty" json:"Description,omitempty" xml:"Description,omitempty"`
	Options     []string `yaml:"Options,omitempty" json:"Options,omitempty" xml:"Options,omitempty"`
	// full
	Title       string   `yaml:"Title,omitempty" json:"Title,omitempty" xml:"Title,omitempty"`
	UniqueItems bool     `yaml:"UniqueItems,omitempty" json:"UniqueItems,omitempty" xml:"UniqueItems,omitempty"`
	Required    []string `yaml:"Required,omitempty" json:"Required,omitempty" xml:"Required,omitempty"`
	// links
	Items      *Model   `yaml:"Items,omitempty" json:"Items,omitempty" xml:"Items,omitempty"`
	OneOf      []*Model `yaml:"OneOf,omitempty" json:"OneOf,omitempty" xml:"OneOf,omitempty"`
	Properties []*Model `yaml:"Properties,omitempty" json:"Properties,omitempty" xml:"Properties,omitempty"`
}

// Full projection
func (sch *Model) Full() (out interface{}) {
	fax := sch.Summary().(*ModelFax)
	fax.Title = sch.Title
	fax.UniqueItems = sch.UniqueItems
	fax.Required = sch.Required
	out = fax
	return
}

// Summary -
func (sch *Model) Summary() (out interface{}) {
	fax := sch.Brief().(*ModelFax)
	fax.Description = sch.Description
	fax.Options = sch.Options
	out = fax
	return
}

// Brief -
func (sch *Model) Brief() (out interface{}) {
	fax := &ModelFax{}
	// always include links
	fax.Items = sch.Items
	fax.OneOf = sch.OneOf
	fax.Properties = sch.Properties

	fax.Name = sch.Label
	fax.Type = sch.Type
	out = fax
	return
}

// Present Presenter interface contract
func (sch *Model) Present(presentation *Presentation) (s string) {

	var f = func(tsch *Model, t Traveler) {
		t.Push(tsch)
		Marshal(tsch, presentation)
		t.Pop()
	}

	var km = KeyMaker{
		stack: make([]string, 0, KeyStackDepth),
	}

	sch.Travel(&km, f)
	return
}
