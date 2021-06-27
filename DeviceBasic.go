package main

import "fmt"

type DeviceBasicResponse struct {
	Meta struct {
		Rc string `json:"rc"`
	} `json:"meta"`
	Data []DeviceBasic `json:"data"`
}

type DeviceBasic struct {
	Mac      string `json:"mac"`
	State    int    `json:"state"`
	Adopted  bool   `json:"adopted"`
	Disabled bool   `json:"disabled"`
	Type     string `json:"type"`
	Model    string `json:"model"`
	Name     string `json:"name"`
}

func (d *DeviceBasic) IsOK() bool {
	if d.Adopted && !d.Disabled && d.State != 1 {
		return false
	}
	return true
}

func (d *DeviceBasic) String() string {
	if d.State == 1 {
		return fmt.Sprintf("[%s] %s [%s] - OK", d.Model, d.Name, d.Mac)
	} else {
		return fmt.Sprintf("[%s] %s [%s] - failed", d.Model, d.Name, d.Mac)
	}
}
