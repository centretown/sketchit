package api

// DeviceYAML -
type DeviceYAML struct {
	Sector string        `yaml:"sector,omitempty" json:"sector,omitempty"`
	Label  string        `yaml:"label,omitempty" json:"label,omitempty"`
	Model  string        `yaml:"model,omitempty" json:"model,omitempty"`
	IP     string        `yaml:"ip,omitempty" json:"ip,omitempty"`
	Port   string        `yaml:"port,omitempty" json:"port,omitempty"`
	Pins   []*Device_Pin `yaml:"pins,omitempty" json:"pins,omitempty"`
}

// MarshalYAML produces yaml output for schema
func (dev *Device) MarshalYAML() (out interface{}, err error) {
	yaml := &DeviceYAML{}
	yaml.Sector = dev.Sector
	yaml.Label = dev.Label
	yaml.Model = dev.Model
	yaml.IP = dev.Ip
	yaml.Port = dev.Port
	yaml.Pins = dev.Pins
	out = yaml
	return
}
