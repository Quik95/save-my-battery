package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	notify "github.com/TheCreeper/go-notify"
	"github.com/distatus/battery"
)

var (
	threshold    int
	irritateMode bool
	updateRate   int
	urgencyLevel int
)

func init() {
	flag.IntVar(&threshold, "threshold", 60, "Specify charging threshold after which notification will be shown")
	flag.BoolVar(&irritateMode, "irritate", false, "Irritate mode sends notification on every update")
	flag.IntVar(&updateRate, "rate", 30, "Specify how often battery level should be checked. Value should be in seconds")
	flag.IntVar(&urgencyLevel, "urgency", 2, "Specify notification urgency level. 0 for Low, 1 for Normal and 2 for Urgent")
}

func main() {
	flag.Parse()

	// match urgencyLevel to notify enums
	var urgency byte
	switch urgencyLevel {
	case int(notify.UrgencyLow):
		urgency = notify.UrgencyLow
	case int(notify.UrgencyNormal):
		urgency = notify.UrgencyNormal
	case int(notify.UrgencyCritical):
		urgency = notify.UrgencyCritical
	default:
		log.Fatalf("%d is not a valid urgency setting.", urgencyLevel)
	}

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
				showNotification(threshold, level, urgency)

				notificationShown = true
			} else if bat.State != battery.Charging {
				// when battery state changes reset notification
				notificationShown = false
			}
		}
		time.Sleep(time.Second * time.Duration(updateRate))
	}
}

func showNotification(threshold, currentLevel int, urgencyLevel byte) error {
	message := fmt.Sprintf(
		"Your battery is charged in %d%%, which exceeds threshold of %d%%. Please consider disconnecting the charger to save battery life.",
		currentLevel, threshold)
	ntf := notify.NewNotification("The battery is overcharged", message)
	ntf.Timeout = notify.ExpiresNever
	ntf.AppIcon = "battery"
	ntf.Hints = make(map[string]interface{})
	ntf.Hints[notify.HintUrgency] = notify.UrgencyCritical

	if _, err := ntf.Show(); err != nil {
		return err
	}
	return nil
}

// checkNotification abstracts logic for showing notification
func checkNotification(bat *battery.Battery, level int, shown bool, threshold int) bool {
	baseCase := bat.State == battery.Charging && level > threshold
	if baseCase && irritateMode {
		return true
	} else if baseCase && !shown {
		return true
	}
	return false
}
