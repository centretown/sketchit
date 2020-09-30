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
func (mod *Model) Full() (out interface{}) {
	fax := mod.Summary().(*ModelFax)
	fax.Title = mod.Title
	fax.UniqueItems = mod.UniqueItems
	fax.Required = mod.Required
	out = fax
	return
}

// Summary -
func (mod *Model) Summary() (out interface{}) {
	fax := mod.Brief().(*ModelFax)
	fax.Description = mod.Description
	fax.Options = mod.Options
	out = fax
	return
}

// Brief -
func (mod *Model) Brief() (out interface{}) {
	fax := &ModelFax{}
	// always include links
	fax.Items = mod.Items
	fax.OneOf = mod.OneOf
	fax.Properties = mod.Properties

	fax.Name = mod.Label
	fax.Type = mod.Type
	out = fax
	return
}

// Present Presenter interface contract
func (mod *Model) Present(presentation *Presentation) (s string) {

	var marshal = func(model *Model, traveler Traveler) {
		traveler.Push(model)
		b, err := Marshal(model, presentation)
		if err == nil {
			s += string(b)
		}
		traveler.Pop()
	}

	var routeMaker = &RouteMaker{
		stack: make([]string, 0, KeyStackDepth),
	}

	mod.Travel(routeMaker, marshal)
	return
}
