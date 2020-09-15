package api

// SketchYAML -
type SketchYAML struct {
	Model   string    `yaml:"model,omitempty" json:"model,omitempty"`
	Label   string    `yaml:"label,omitempty" json:"label,omitempty"`
	Device  string    `yaml:"device,omitempty" json:"device,omitempty"`
	Purpose string    `yaml:"purpose,omitempty" json:"purpose,omitempty"`
	Setup   []*Action `yaml:"setup,omitempty" json:"setup,omitempty"`
	Loop    []*Action `yaml:"loop,omitempty" json:"loop,omitempty"`
}

// MarshalYAML produces yaml output for schema
func (skch *Sketch) MarshalYAML() (out interface{}, err error) {
	yaml := &SketchYAML{}
	yaml.Model = skch.Model
	yaml.Label = skch.Label
	yaml.Device = skch.Device
	yaml.Purpose = skch.Purpose
	yaml.Setup = skch.Setup
	yaml.Loop = skch.Loop
	out = yaml
	return
}
