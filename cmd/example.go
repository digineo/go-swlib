package main

import (
	"fmt"

	"github.com/digineo/swlib"
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

	attributes, err := c.ListGlobalAttributes(&switches[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attributes)

	attributes, err = c.ListPortAttributes(&switches[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attributes)

	attributes, err = c.ListVLANAttributes(&switches[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attributes)
}
