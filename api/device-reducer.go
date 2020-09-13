package api

// Reduce -
func (dev *Device) Reduce(reduction ...Reduction) (out interface{}) {
	if len(reduction) < 1 {
		return
	}
	level := reduction[0]
	fake := &FakeDevice{}
	out = fake
	fake.Sector = dev.Sector
	fake.Label = dev.Label
	if level == Brief {
		return
	}
	fake.Model = dev.Model
	fake.IP = dev.Ip
	fake.Port = dev.Port
	if level == Summary {
		return
	}
	fake.Pins = dev.Pins
	return
}

// FakeDevice -
type FakeDevice struct {
	Sector string        `yaml:"sector,omitempty" json:"sector,omitempty"`
	Label  string        `yaml:"label,omitempty" json:"label,omitempty"`
	Model  string        `yaml:"model,omitempty" json:"model,omitempty"`
	IP     string        `yaml:"ip,omitempty" json:"ip,omitempty"`
	Port   string        `yaml:"port,omitempty" json:"port,omitempty"`
	Pins   []*Device_Pin `yaml:"pins,omitempty" json:"pins,omitempty"`
}
