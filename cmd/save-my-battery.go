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
	const threshold int = 10
	// Flag to specify if we have already shown the notification
	// so we won't be bombarding user with notifications
	notificationShown := false

	for {
		batteries, err := battery.GetAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, bat := range batteries {
			level := int(math.Floor(bat.Current / bat.Full * 100))
			if checkNotification(bat, level, notificationShown, threshold) {
				message := fmt.Sprintf(
					"Your battery is charged in %d%%, which exceeds threshold of %d%%. Please consider disconnecting the charger to save battery life.",
					level, threshold)
				if err := beeep.Alert("Battery is overcharged", message, ""); err != nil {
					log.Fatal(err)
				}
				notificationShown = true
			} else if bat.State != battery.Charging {
				// when battery state changes reset notification
				notificationShown = false
			}
		}
		time.Sleep(time.Second * 5)
	}
}

// checkNotification abstracts logic for showing notification
func checkNotification(bat *battery.Battery, level int, shown bool, threshold int) bool {
	if bat.State == battery.Charging && level > threshold && !shown {
		return true
	}
	return false
}
