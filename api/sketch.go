package api

// SketchFax -
type SketchFax struct {
	Toolkit string    `yaml:"toolkit,omitempty" json:"toolkit,omitempty"`
	Label   string    `yaml:"label,omitempty" json:"label,omitempty"`
	Device  string    `yaml:"device,omitempty" json:"device,omitempty"`
	Purpose string    `yaml:"purpose,omitempty" json:"purpose,omitempty"`
	Setup   []*Action `yaml:"setup,omitempty" json:"setup,omitempty"`
	Loop    []*Action `yaml:"loop,omitempty" json:"loop,omitempty"`
}

// MarshalYAML yaml Marshaler interface contract
func (skch *Sketch) MarshalYAML() (out interface{}, err error) {
	yml := &SketchFax{}
	yml.Toolkit = skch.Toolkit
	yml.Label = skch.Label
	yml.Device = skch.Device
	yml.Purpose = skch.Purpose
	yml.Setup = skch.Setup
	yml.Loop = skch.Loop
	out = yml
	return
}
