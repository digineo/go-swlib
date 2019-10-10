package swlib

import (
	"github.com/mdlayher/netlink"
)

type Group int

const (
	GroupGlobal Group = iota
	GroupVLAN
	GroupPort
)

type DataType uint32

const (
	DataTypeUnspec DataType = iota
	DataTypeInt
	DataTypeString
	DataTypePorts
	DataTypeLink
	DataTypeNoVal
)

type Attribute struct {
	AType       Group
	ID          uint32
	Type        DataType
	Name        string
	Description string
}

func (a *Attribute) UnmarshalAttributes(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch AttributeType(ad.Type()) {
		case AttrOPID:
			a.ID = ad.Uint32()
		case AttrOPType:
			a.Type = DataType(ad.Uint32())
		case AttrOPName:
			a.Name = ad.String()
		case AttrOPDescription:
			a.Description = ad.String()
		}
	}
	return ad.Err()
}

func (a *Attribute) UnmarshalBinary(data []byte) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	return a.UnmarshalAttributes(ad)
}
