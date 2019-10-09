package swlib

import (
	"io"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

// connector represents a
type connector interface {
	io.Closer

	Execute(m genetlink.Message, family uint16, flags netlink.HeaderFlags) ([]genetlink.Message, error)
}

// Conn represents a genetlink connection to the OpenWRT swconfig
// subsystem and implements their available actions.
type Conn struct {
	conn   connector
	family genetlink.Family
}

// Dial opens a new genetlink connection and returns a Conn structure
// instance that implements the switch API.
func Dial(config *netlink.Config) (*Conn, error) {
	c, err := genetlink.Dial(config)
	if err != nil {
		return nil, err
	}

	family, err := c.GetFamily(FamilyName)
	if err != nil {
		return nil, err
	}

	return &Conn{conn: c, family: family}, nil
}

func (c *Conn) ListSwitches() ([]Switch, error) {
	req := genetlink.Message{
		Header: genetlink.Header{
			Command: uint8(CmdGetSwitch),
			Version: c.family.Version,
		},
	}

	flags := netlink.Request | netlink.Dump
	msgs, err := c.conn.Execute(req, c.family.ID, flags)
	if err != nil {
		return nil, err
	}

	switches := make([]Switch, len(msgs))
	for i, m := range msgs {
		ad, err := netlink.NewAttributeDecoder(m.Data)
		if err != nil {
			return nil, err
		}

		for ad.Next() {
			sw := &switches[i]
			switch AttributeType(ad.Type()) {
			case AttrId:
				sw.ID = ad.Uint32()
			case AttrDevName:
				sw.DeviceName = ad.String()
			case AttrName:
				sw.Name = ad.String()
			case AttrAlias:
				sw.Alias = ad.String()
			case AttrVLANs:
				sw.VLANs = ad.Uint32()
			case AttrPorts:
				sw.Ports = ad.Uint32()
			case AttrCPUPort:
				sw.CPUPort = ad.Uint32()
			}
		}
	}

	return switches, nil
}
