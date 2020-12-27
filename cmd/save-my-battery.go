package main

import (
	"fmt"
	"log"

	"github.com/distatus/battery"
)

func main() {
	batteries, err := battery.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, battery := range batteries {
		fmt.Printf("Battery number: %d, charge level: %v", i, battery.Current/battery.Full*100)
	}
}
