package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/distatus/battery"
	"github.com/gen2brain/beeep"
)

var (
	threshold int
)

func init() {
	flag.IntVar(&threshold, "threshold", 60, "Specify charging threshold after which notification will be shown")
}

func main() {
	flag.Parse()

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
