package api

// SchemaYAML -
type SchemaYAML struct {
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

// Reduce schema within fake based on schemaReducer
func (sch *Schema) Reduce(projection ...Projection) (out interface{}) {
	out = sch
	if len(projection) < 1 {
		return
	}
	level := projection[0]
	if level == Projection_FULL {
		return
	}

	yaml := &SchemaYAML{}
	out = yaml

	// always include links
	yaml.Items = sch.Items
	yaml.OneOf = sch.OneOf
	yaml.Properties = sch.Properties

	yaml.Name = sch.Name
	yaml.Type = sch.Type
	if level == Projection_BRIEF {
		return
	}

	yaml.Description = sch.Description
	yaml.Options = sch.Options
	return
}
