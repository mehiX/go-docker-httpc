package data

type Port struct {
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}

type Config struct {
	NetworkMode string `json:"NetworkMode"`
}

type Network struct{}

type NSettings struct {
	Networks map[string]*Network `json:"Networks"`
}

type Mount struct {
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Driver      string `json:"Driver"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

type Container struct {
	ID              string            `json:"Id"`
	Names           []string          `json:"Names"`
	Image           string            `json:"Image"`
	ImageID         string            `json:"ImageID"`
	Command         string            `json:"Command"`
	Created         int               `json:"Created"`
	State           string            `json:"State"`
	Status          string            `json:"Status"`
	Ports           []Port            `json:"Ports"`
	Labels          map[string]string `json:"Labels"`
	SizeRw          int               `json:"SizeRw"`
	SizeRootFs      int               `json:"SizeRootFs"`
	HostConfig      Config            `json:"HostConfig"`
	NetworkSettings NSettings         `json:"NetworkSettings"`
	Mounts          []Mount           `json:"Mounts"`
}
