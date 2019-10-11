package swlib

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

type PortAttribute int

const (
	PortAttrUnspec PortAttribute = iota
	PortAttrID
	PortAttrFlagTagged
)

type PortFlags int

const (
	PortFlagTagged PortFlags = 1 << iota
)

type Port struct {
	ID    uint32
	Flags PortFlags
}

func (p *Port) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	*p = Port{}
	for ad.Next() {
		switch a := PortAttribute(ad.Type()); a {
		case PortAttrID:
			p.ID = ad.Uint32()
		case PortAttrFlagTagged:
			if ad.Flag() {
				p.Flags |= PortFlagTagged
			}
		default:
			return fmt.Errorf("%w: %s", ErrIncompatibleNetlink, a)
		}
	}
	return ad.Err()
}

type Ports []Port

func (s *Ports) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		if AttributeType(ad.Type()) == AttrPort {
			p := Port{}
			ad.Nested(p.UnmarshalAttributes)
			*s = append(*s, p)
		}
	}
	return ad.Err()
}
