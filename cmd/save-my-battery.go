package main

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/distatus/battery"
	"github.com/gen2brain/beeep"
)

func main() {
	batteries, err := battery.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	for {
		for _, battery := range batteries {
			levelFloat := math.Floor(battery.Current / battery.Full * 100)
			levelString := strconv.Itoa(int(levelFloat))
			if err := beeep.Alert("Battery Level", levelString, ""); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
