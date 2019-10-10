package swlib

import (
	"github.com/mdlayher/netlink"
)

type PortMap struct {
	Number  int
	Segment string
	Virt    uint32
}

func (m *PortMap) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch AttributeType(ad.Type()) {
		case AttrPortMapSegment:
			m.Segment = ad.String()
		case AttrPortMapVirt:
			m.Virt = ad.Uint32()
		}
	}
	return ad.Err()
}
