package api

// Reduce schema within fake based on schemaReducer
func (sch *Schema) Reduce(reduction ...Reduction) (out interface{}) {
	if len(reduction) < 1 {
		return
	}
	level := reduction[0]
	fake := &FakeSchema{}
	out = fake

	// always include links
	fake.Items = sch.Items
	fake.OneOf = sch.OneOf
	fake.Properties = sch.Properties

	fake.Name = sch.Name
	fake.Type = sch.Type
	if level == Brief {
		return
	}

	fake.Description = sch.Description
	fake.Options = sch.Options
	if level == Summary {
		return
	}

	fake.Title = sch.Title
	fake.UniqueItems = sch.UniqueItems
	fake.Required = sch.Required
	return
}

// FakeSchema -
type FakeSchema struct {
	Title       string    `yaml:"Title,omitempty" json:"Title,omitempty" xml:"Title,omitempty"`
	Name        string    `yaml:"Name,omitempty" json:"Name,omitempty" xml:"Name,omitempty"`
	Type        string    `yaml:"Type,omitempty" json:"Type,omitempty" xml:"Type,omitempty"`
	Description string    `yaml:"Description,omitempty" json:"Description,omitempty" xml:"Description,omitempty"`
	UniqueItems bool      `yaml:"UniqueItems,omitempty" json:"UniqueItems,omitempty" xml:"UniqueItems,omitempty"`
	Items       *Schema   `yaml:"Items,omitempty" json:"Items,omitempty" xml:"Items,omitempty"`
	Required    []string  `yaml:"Required,omitempty" json:"Required,omitempty" xml:"Required,omitempty"`
	Options     []string  `yaml:"Options,omitempty" json:"Options,omitempty" xml:"Options,omitempty"`
	OneOf       []*Schema `yaml:"OneOf,omitempty" json:"OneOf,omitempty" xml:"OneOf,omitempty"`
	Properties  []*Schema `yaml:"Properties,omitempty" json:"Properties,omitempty" xml:"Properties,omitempty"`
}
