package swlib

import (
	"encoding"
	"fmt"
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

func (c *Conn) newRequest(cmd CommandType, data []byte) genetlink.Message {
	return genetlink.Message{
		Header: genetlink.Header{
			Command: uint8(cmd),
			Version: c.family.Version,
		},
		Data: data,
	}
}

func (c *Conn) query(cmd CommandType, flags netlink.HeaderFlags, data encoding.BinaryMarshaler) ([]genetlink.Message, error) {
	req := genetlink.Message{
		Header: genetlink.Header{
			Command: uint8(cmd),
			Version: c.family.Version,
		},
	}

	if data != nil {
		var err error
		if req.Data, err = data.MarshalBinary(); err != nil {
			return nil, fmt.Errorf("marshal data failed: %w", err)
		}
	}

	return c.conn.Execute(req, c.family.ID, netlink.Request|flags)
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

func marshalAttributes(g Group, msgs []genetlink.Message) ([]Attribute, error) {
	attributes := make([]Attribute, len(msgs)-1)
	for i, m := range msgs[:len(msgs)-1] {
		attribute := &attributes[i]
		attribute.AType = g
		if err := attribute.UnmarshalBinary(m.Data); err != nil {
			return nil, err
		}
	}
	return attributes, nil
}

func (c *Conn) ListGlobalAttributes(dev *Device) ([]Attribute, error) {
	msgs, err := c.query(CmdListGlobal, netlink.Acknowledge, &DeviceBase{ID: dev.ID})
	if err != nil {
		return nil, fmt.Errorf("netlink execute failed: %w", err)
	}

	return marshalAttributes(GroupGlobal, msgs)
}

func (c *Conn) ListPortAttributes(dev *Device) ([]Attribute, error) {
	msgs, err := c.query(CmdListPort, netlink.Acknowledge, &DeviceBase{ID: dev.ID})
	if err != nil {
		return nil, fmt.Errorf("netlink execute failed: %w", err)
	}

	return marshalAttributes(GroupPort, msgs)
}

func (c *Conn) ListVLANAttributes(dev *Device) ([]Attribute, error) {
	msgs, err := c.query(CmdListVLAN, netlink.Acknowledge, &DeviceBase{ID: dev.ID})
	if err != nil {
		return nil, fmt.Errorf("netlink execute failed: %w", err)
	}

	return marshalAttributes(GroupVLAN, msgs)
}
