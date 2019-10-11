package swlib

import (
	"github.com/mdlayher/genetlink"
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
	Device *Device
	Group  Group

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

type Attributes map[string]*Attribute

func AttributesFromMessages(d *Device, g Group, msgs []genetlink.Message) (Attributes, error) {
	a := make(Attributes, len(msgs))
	for _, m := range msgs {
		attr := &Attribute{
			Device: d,
			Group:  g,
		}
		if err := attr.UnmarshalBinary(m.Data); err != nil {
			return nil, err
		}
		a[attr.Name] = attr
	}
	return a, nil
}
