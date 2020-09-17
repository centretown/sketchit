package api

// DeviceFax -
type DeviceFax struct {
	Sector string        `yaml:"sector,omitempty" json:"sector,omitempty"`
	Label  string        `yaml:"label,omitempty" json:"label,omitempty"`
	Model  string        `yaml:"model,omitempty" json:"model,omitempty"`
	IP     string        `yaml:"ip,omitempty" json:"ip,omitempty"`
	Port   string        `yaml:"port,omitempty" json:"port,omitempty"`
	Pins   []*Device_Pin `yaml:"pins,omitempty" json:"pins,omitempty"`
}

// MarshalYAML produces yaml output for schema
func (dev *Device) MarshalYAML() (out interface{}, err error) {
	yml := &DeviceFax{}
	yml.Sector = dev.Sector
	yml.Label = dev.Label
	yml.Model = dev.Model
	yml.IP = dev.Ip
	yml.Port = dev.Port
	yml.Pins = dev.Pins
	out = yml
	return
}

// // DeviceList is a wrapper that implements the Sprinter interface
// type DeviceList []*Device

// // Sprint -
// func (dl *DeviceList) Sprint(o interface{}, presentation *Presentation) (s string, err error) {

// 	return
// }
