package commander

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (c *Commander) Tomorrow(inputMessage *tgbotapi.Message) {
	isTomorrow := isInTomorrowTimePeriod()

	tomorrowPlaces := c.service.GetPlacesInsideTimePeriod(isTomorrow)

	if len(tomorrowPlaces) == 0 {
		c.SendMessage(inputMessage.Chat.ID,
			fmt.Sprintf("Свободных места на завтра нет 🥶\n"))
	} else {
		c.SendMessage(inputMessage.Chat.ID, fmt.Sprintf("Свободные места на завтра ⛸\n"))

		for placeId, timeSlotsPlace := range tomorrowPlaces {
			msgText := createPlaceTimeSlotsText(timeSlotsPlace)
			keyboard := createKeyboardWithTimeSlots(timeSlotsPlace, placeId)
			c.SendMessageWithKeyboard(inputMessage.Chat.ID, msgText, keyboard)
		}
	}
}

func isInTomorrowTimePeriod() func(t time.Time) bool {
	return func(t time.Time) bool {
		todayDayStart := time.Now().UTC().Truncate(time.Hour * 24)
		startOfTomorrowDay := todayDayStart.AddDate(0, 0, 1)
		endOfTomorrowDay := startOfTomorrowDay.AddDate(0, 0, 1).Add(-time.Second)
		return t.After(startOfTomorrowDay) && t.Before(endOfTomorrowDay)
	}
}
