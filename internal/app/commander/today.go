package commander

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (c *Commander) Today(inputMessage *tgbotapi.Message) {
	isToday := isInTodayTimePeriod()

	todayPlaces := c.service.GetPlacesInsideTimePeriod(isToday)

	if len(todayPlaces) == 0 {
		c.SendMessage(inputMessage.Chat.ID,
			fmt.Sprintf("Свободных места на сегодня нет 🥶\n"))
	} else {
		c.SendMessage(inputMessage.Chat.ID, fmt.Sprintf("Свободные места на сегодня ⛸\n"))

		for placeId, timeSlotsPlace := range todayPlaces {
			msgText := createPlaceTimeSlotsText(timeSlotsPlace)
			keyboard := createKeyboardWithTimeSlots(timeSlotsPlace, placeId)
			c.SendMessageWithKeyboard(inputMessage.Chat.ID, msgText, keyboard)
		}
	}
}

func isInTodayTimePeriod() func(t time.Time) bool {
	return func(t time.Time) bool {
		timeStartOfTodayDay := time.Now().UTC().Truncate(time.Hour * 24)
		timeEndOfTodayDay := timeStartOfTodayDay.AddDate(0, 0, 1).Add(-time.Second)
		return t.After(timeStartOfTodayDay) && t.Before(timeEndOfTodayDay)
	}
}
