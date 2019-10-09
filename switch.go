package swlib

type Switch struct {
	ID         uint32
	DeviceName string
	Alias      string
	Name       string
	VLANs      uint32
	Ports      uint32
	CPUPort    uint32
}
