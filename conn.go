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

	msgs, err := c.conn.Execute(req, c.family.ID, netlink.Request|netlink.Dump)
	if err != nil {
		return nil, err
	}

	switches := make([]Switch, len(msgs))
	for i, m := range msgs {
		if err := switches[i].UnmarshalBinary(m.Data); err != nil {
			return nil, err
		}
	}

	return switches, nil
}
