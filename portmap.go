package swlib

import (
	"github.com/mdlayher/netlink"
)

type PortMapAttribute int

const (
	PortMapAttrPorts PortMapAttribute = iota
	PortMapAttrSegment
	PortMapAttrVirt
)

type PortMap struct {
	Number  uint32
	Segment string
	Virt    uint32
}

func (m *PortMap) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch PortMapAttribute(ad.Type()) {
		case PortMapAttrSegment:
			m.Segment = ad.String()
		case PortMapAttrVirt:
			m.Virt = ad.Uint32()
		}
	}
	return ad.Err()
}
