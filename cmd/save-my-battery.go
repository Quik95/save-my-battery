package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/distatus/battery"
	"github.com/gen2brain/beeep"
)

func main() {
	batteries, err := battery.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	const threshold int = 10
	// Flag to specify if we have already shown the notification
	// so we won't be bombarding user with notifications
	notificationShown := false

	for {
		for _, bat := range batteries {
			level := int(math.Floor(bat.Current / bat.Full * 100))
			if level > threshold && bat.State == battery.Charging && !notificationShown {
				message := fmt.Sprintf(
					"Your battery is charged in %d%%, which exceeds threshold of %d%%. Please consider disconnecting the charger to save battery life.",
					level, threshold)
				if err := beeep.Alert("Battery is overcharged", message, ""); err != nil {
					log.Fatal(err)
				}
				notificationShown = true
			} else if bat.State != battery.Charging {
				notificationShown = false
			}
		}
		time.Sleep(time.Second * 5)
	}
}
