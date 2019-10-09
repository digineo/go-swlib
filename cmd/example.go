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
}
