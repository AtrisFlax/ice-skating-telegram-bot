package commander

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ice-skating-telegram-bot/internal/app/service"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func createPlaceTimeSlotsText(timeSlotsPlace service.TimeSlotsPlace) string {

	msgText := strings.Builder{}

	msgText.WriteString(fmt.Sprintf(
		"%s %s\n",
		timeSlotsPlace.LocationName,
		getRandomWinterEmoji()))

	for _, slot := range timeSlotsPlace.TimeSlots {
		slotLocalTime := getMoscowLocalTime(slot.Time)

		msgText.WriteString(fmt.Sprintf(
			"%02d:%02d мест - %d\n",
			slotLocalTime.Hour(),
			slotLocalTime.Minute(),
			slot.FreePlacesCount))
	}
	return msgText.String()
}

func createKeyboardWithTimeSlots(timeSlotsPlace service.TimeSlotsPlace, placeId string) tgbotapi.InlineKeyboardMarkup {
	buttonRow := make([]tgbotapi.InlineKeyboardButton, 0, len(timeSlotsPlace.TimeSlots))

	for _, slot := range timeSlotsPlace.TimeSlots {
		slotLocalTime := getMoscowLocalTime(slot.Time)

		replyURL := createReplyURL(placeId, slot.Time, slot.PlaceDateId)

		buttonText := fmt.Sprintf(
			"%02d:%02d\n",
			slotLocalTime.Hour(),
			slotLocalTime.Minute())

		replyButton := createReplyButton(buttonText, replyURL)
		buttonRow = append(buttonRow, replyButton)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttonRow)
	return keyboard
}

func createReplyButton(text string, url string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonURL(text, url)
}

func createReplyURL(placeID string, time time.Time, placeDateID string) string {
	path := "mc/" + placeID + "/registration"

	const timeLayout = "2006-01-02T15:04:05.000Z"
	params := fmt.Sprintf("timeSlot=%s&placeDateId=%s", time.Format(timeLayout), placeDateID)

	replyUrl := url.URL{
		Scheme:   "https",
		Host:     "mftickets.technolab.com.ru",
		Path:     path,
		RawQuery: params,
	}

	return replyUrl.String()
}

func getMoscowLocalTime(slot time.Time) time.Time {
	location, _ := time.LoadLocation("Europe/Moscow")
	slotLocalTime := slot.In(location)
	return slotLocalTime
}

var emojiSlice = []string{
	"🏂",
	"🌲",
	"🏔️",
	"🌁",
	"🌬️",
	"☔",
	"❄",
	"☃",
	"⛄",
	"🎿",
	"🧣",
	"🧤",
	"🧥",
	"☃",
	"🧊",
	"🌨️",
	"🎄",
	"🥂",
	"⛷",
	"🏂",
	"☕",
	"🌃",
	"⚪",
	"🛷",
	"🎣",
	"🐻‍",
	"🐧‍",
	"🎁‍",
	"🧦‍",
}

func getRandomWinterEmoji() string {
	return emojiSlice[(rand.Int() % len(emojiSlice))]
}
