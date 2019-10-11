package swlib

import (
	"github.com/mdlayher/netlink"
)

type DeviceBase struct {
	ID uint32
}

func (d *DeviceBase) MarshalBinary() (data []byte, err error) {
	ae := netlink.NewAttributeEncoder()
	ae.Uint32(uint16(AttrID), d.ID)
	return ae.Encode()
}

type Device struct {
	ID         uint32
	DeviceName string
	Alias      string
	Name       string
	VLANs      uint32
	Ports      uint32
	CPUPort    uint32
	PortMap    []*PortMap
}

func (d *Device) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch AttributeType(ad.Type()) {
		case AttrID:
			d.ID = ad.Uint32()
		case AttrDevName:
			d.DeviceName = ad.String()
		case AttrName:
			d.Name = ad.String()
		case AttrAlias:
			d.Alias = ad.String()
		case AttrVLANs:
			d.VLANs = ad.Uint32()
		case AttrPorts:
			d.Ports = ad.Uint32()
		case AttrCPUPort:
			d.CPUPort = ad.Uint32()
		case AttrPortMap:
			ad.Nested(func(nad *netlink.AttributeDecoder) error {
				for nad.Next() {
					if AttributeType(nad.Type()) == AttrPorts {
						pm := &PortMap{Number: uint32(len(d.PortMap) + 1)}
						nad.Nested(pm.UnmarshalAttributes)
						d.PortMap = append(d.PortMap, pm)
					}
				}
				return nad.Err()
			})
		}
	}
	return ad.Err()
}

func (d *Device) UnmarshalBinary(data []byte) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	return d.UnmarshalAttributes(ad)
}
