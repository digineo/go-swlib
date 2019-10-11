package swlib

import (
	"github.com/mdlayher/netlink"
)

type sendAttribute struct {
	attribute  *Attribute
	portOrVLAN uint32
}

func (s *sendAttribute) MarshalBinary() (data []byte, err error) {
	ae := netlink.NewAttributeEncoder()
	ae.Uint32(uint16(AttrID), s.attribute.Device.ID)
	ae.Uint32(uint16(AttrOPID), s.attribute.ID)
	switch s.attribute.Group {
	case GroupPort:
		ae.Uint32(uint16(AttrOPPort), s.portOrVLAN)
	case GroupVLAN:
		ae.Uint32(uint16(AttrOPVLAN), s.portOrVLAN)
	}
	return ae.Encode()
}
