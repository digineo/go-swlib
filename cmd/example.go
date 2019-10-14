package main

import (
	"fmt"

	"github.com/digineo/go-swlib"
)

func main() {
	c, err := swlib.Dial(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	switches, err := c.ListSwitches()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(switches)

	fmt.Println()
	attributes, err := c.ListGlobalAttributes(&switches[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attributes)

	fmt.Println()
	attributes, err = c.ListPortAttributes(&switches[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attributes)

	for i := uint32(0); i < switches[0].Ports; i++ {
		l, err := c.GetAttributeLink(attributes["link"], i)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(i, l)

		pvid, err := c.GetAttributeInt(attributes["pvid"], i)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(i, pvid)
	}

	fmt.Println()
	attributes, err = c.ListVLANAttributes(&switches[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attributes)

	for i := uint32(0); i < 5; i++ {
		p, err := c.GetAttributePorts(attributes["ports"], i)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(i, p)
	}
}
