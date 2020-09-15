package api

var devices = []*Device{
	{
		Label: "esp32-01",
		Model: "ESP32",
		Ip:    "192.168.1.200",
		Port:  "esp32-01",
		Pins: []*Device_Pin{
			{
				Id:      2,
				Label:   "Interal LED",
				Purpose: "signal activity",
			},
			{
				Id:      5,
				Label:   "TX",
				Purpose: "Aux Serial",
			},
			{
				Id:      6,
				Label:   "RX",
				Purpose: "Aux Serial",
			},
		},
	},
	{
		Label: "nano-01",
		Model: "NANO",
		Port:  "nano-01",
		Pins: []*Device_Pin{
			{
				Id:      13,
				Label:   "Interal LED",
				Purpose: "signal activity",
			},
			{
				Id:      5,
				Label:   "TX",
				Purpose: "Aux Serial",
			},
			{
				Id:      6,
				Label:   "RX",
				Purpose: "Aux Serial",
			},
		},
	},
}
