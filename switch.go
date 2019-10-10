package swlib

import (
	"github.com/mdlayher/netlink"
)

type Switch struct {
	ID         uint32
	DeviceName string
	Alias      string
	Name       string
	VLANs      uint32
	Ports      uint32
	CPUPort    uint32
}

func (s *Switch) UnmarshalBinary(data []byte) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}

	for ad.Next() {
		switch AttributeType(ad.Type()) {
		case AttrID:
			s.ID = ad.Uint32()
		case AttrDevName:
			s.DeviceName = ad.String()
		case AttrName:
			s.Name = ad.String()
		case AttrAlias:
			s.Alias = ad.String()
		case AttrVLANs:
			s.VLANs = ad.Uint32()
		case AttrPorts:
			s.Ports = ad.Uint32()
		case AttrCPUPort:
			s.CPUPort = ad.Uint32()
		}
	}
	return nil
}
