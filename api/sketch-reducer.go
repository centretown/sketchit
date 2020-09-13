package api

// Reduce -
func (x *Sketch) Reduce(reduction ...Reduction) (out interface{}) {
	if len(reduction) < 1 {
		return
	}
	level := reduction[0]
	f := &FakeSketch{}
	out = f
	f.Model = x.Model
	f.Label = x.Label
	if level == Brief {
		return
	}

	f.Device = x.Device
	f.Purpose = x.Purpose
	if level == Summary {
		return
	}

	f.Setup = x.Setup
	f.Loop = x.Loop
	return
}

// FakeSketch -
type FakeSketch struct {
	Model   string    `yaml:"model,omitempty" json:"model,omitempty"`
	Label   string    `yaml:"label,omitempty" json:"label,omitempty"`
	Device  string    `yaml:"device,omitempty" json:"device,omitempty"`
	Purpose string    `yaml:"purpose,omitempty" json:"purpose,omitempty"`
	Setup   []*Action `yaml:"setup,omitempty" json:"setup,omitempty"`
	Loop    []*Action `yaml:"loop,omitempty" json:"loop,omitempty"`
}
