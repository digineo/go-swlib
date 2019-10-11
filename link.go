package swlib

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

type LinkAttr int

const (
	LinkAttrUnspec LinkAttr = iota
	LinkAttrFlagLink
	LinkAttrFlagDuplex
	LinkAttrFlagANeg
	LinkAttrFlagTXFlow
	LinkAttrFlagRXFlow
	LinkAttrSpeed
	LinkAttrFlagEEE100BaseT
	LinkAttrFlagEEE1000BaseT
)

type LinkEEE int

const (
	LinkEEE100BaseT LinkEEE = 1 << iota
	LinkEEE1000BaseT
)

type Link struct {
	Link   bool
	Duplex bool
	ANeg   bool
	TXFlow bool
	RXFlow bool
	Speed  uint32
	EEE    LinkEEE
}

func (l *Link) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	*l = Link{}
	for ad.Next() {
		switch attrType := LinkAttr(ad.Type()); attrType {
		case LinkAttrFlagLink:
			l.Link = ad.Flag()
		case LinkAttrFlagDuplex:
			l.Duplex = ad.Flag()
		case LinkAttrFlagANeg:
			l.ANeg = ad.Flag()
		case LinkAttrFlagTXFlow:
			l.TXFlow = ad.Flag()
		case LinkAttrFlagRXFlow:
			l.RXFlow = ad.Flag()
		case LinkAttrSpeed:
			l.Speed = ad.Uint32()
		case LinkAttrFlagEEE100BaseT:
			if ad.Flag() {
				l.EEE |= LinkEEE100BaseT
			}
		case LinkAttrFlagEEE1000BaseT:
			if ad.Flag() {
				l.EEE |= LinkEEE1000BaseT
			}
		default:
			return fmt.Errorf("%w: %s", ErrIncompatibleNetlink, attrType)
		}
	}
	return ad.Err()
}
