package swlib

import (
	"encoding"
	"fmt"
	"io"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

func attributeAssertType(a *Attribute, t DataType) error {
	if a.Type != t {
		return fmt.Errorf("invalid data type: %s != %s", t, a.Type)
	}
	return nil
}

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

func (c *Conn) newRequest(cmd CommandType) genetlink.Message {
	return genetlink.Message{
		Header: genetlink.Header{
			Command: uint8(cmd),
			Version: c.family.Version,
		},
	}
}

func (c *Conn) query(cmd CommandType, flags netlink.HeaderFlags, data encoding.BinaryMarshaler) ([]genetlink.Message, error) {
	req := c.newRequest(cmd)

	if data != nil {
		var err error
		if req.Data, err = data.MarshalBinary(); err != nil {
			return nil, fmt.Errorf("marshal data failed: %w", err)
		}
	}

	return c.conn.Execute(req, c.family.ID, netlink.Request|flags)
}

func (c *Conn) getAttribute(a *Attribute, portOrVLAN uint32) ([]genetlink.Message, error) {
	var cmd CommandType
	switch a.Group {
	case GroupGlobal:
		cmd = CmdGetGlobal
	case GroupPort:
		cmd = CmdGetPort
	case GroupVLAN:
		cmd = CmdGetVLan
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidAttributeType, a.Group)
	}

	return c.query(cmd, 0, &sendAttribute{
		attribute:  a,
		portOrVLAN: portOrVLAN,
	})
}

func (c *Conn) ListSwitches() ([]Device, error) {
	msgs, err := c.query(CmdGetSwitch, netlink.Dump, nil)
	if err != nil {
		return nil, err
	}

	switches := make([]Device, len(msgs))
	for i, m := range msgs {
		if err := switches[i].UnmarshalBinary(m.Data); err != nil {
			return nil, err
		}
	}

	return switches, nil
}

func (c *Conn) ListGlobalAttributes(dev *Device) (a Attributes, err error) {
	msgs, err := c.query(CmdListGlobal, netlink.Acknowledge, &DeviceBase{ID: dev.ID})
	if err != nil {
		return nil, err
	}

	return AttributesFromMessages(dev, GroupGlobal, msgs[:len(msgs)-1])
}

func (c *Conn) ListPortAttributes(dev *Device) (a Attributes, err error) {
	msgs, err := c.query(CmdListPort, netlink.Acknowledge, &DeviceBase{ID: dev.ID})
	if err != nil {
		return nil, err
	}

	return AttributesFromMessages(dev, GroupPort, msgs[:len(msgs)-1])
}

func (c *Conn) ListVLANAttributes(dev *Device) (a Attributes, err error) {
	msgs, err := c.query(CmdListVLAN, netlink.Acknowledge, &DeviceBase{ID: dev.ID})
	if err != nil {
		return nil, err
	}

	return AttributesFromMessages(dev, GroupVLAN, msgs[:len(msgs)-1])
}

func (c *Conn) GetAttributeLink(a *Attribute, portOrVLAN uint32) (l *Link, err error) {
	if err := attributeAssertType(a, DataTypeLink); err != nil {
		return nil, err
	}

	msgs, err := c.getAttribute(a, portOrVLAN)
	if err != nil {
		return nil, err
	}

	for _, m := range msgs {
		ad, err := netlink.NewAttributeDecoder(m.Data)
		if err != nil {
			return nil, err
		}

		for ad.Next() {
			if AttributeType(ad.Type()) == AttrOPValueLink {
				l = &Link{}
				ad.Nested(l.UnmarshalAttributes)
			}
		}
	}

	return
}

func (c *Conn) GetAttributePorts(a *Attribute, portOrVLAN uint32) (p Ports, err error) {
	if err := attributeAssertType(a, DataTypePorts); err != nil {
		return nil, err
	}

	msgs, err := c.getAttribute(a, portOrVLAN)
	if err != nil {
		return nil, err
	}

	for _, m := range msgs {
		ad, err := netlink.NewAttributeDecoder(m.Data)
		if err != nil {
			return nil, err
		}

		for ad.Next() {
			if AttributeType(ad.Type()) == AttrOPValuePorts {
				ad.Nested(p.UnmarshalAttributes)
			}
		}
	}

	return
}
