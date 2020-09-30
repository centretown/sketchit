package api

// DeviceFax -
type DeviceFax struct {
	Sector  string        `yaml:"sector,omitempty" json:"sector,omitempty"`
	Label   string        `yaml:"label,omitempty" json:"label,omitempty"`
	Toolkit string        `yaml:"toolkit,omitempty" json:"toolkit,omitempty"`
	IP      string        `yaml:"ip,omitempty" json:"ip,omitempty"`
	Port    string        `yaml:"port,omitempty" json:"port,omitempty"`
	Pins    []*Device_Pin `yaml:"pins,omitempty" json:"pins,omitempty"`
}

// MarshalYAML produces yaml output for schema
func (dev *Device) MarshalYAML() (out interface{}, err error) {
	yml := &DeviceFax{}
	yml.Sector = dev.Sector
	yml.Label = dev.Label
	yml.Toolkit = dev.Toolkit
	yml.IP = dev.Ip
	yml.Port = dev.Port
	yml.Pins = dev.Pins
	out = yml
	return
}
