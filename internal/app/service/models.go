package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

type Places map[string]TimeSlotsPlace

type TimeSlotsPlace struct {
	LocationName string
	MaxCount     int
	TimeSlots    []TimeSlot
}

type TimeSlot struct {
	PlaceDateId     string
	Time            time.Time
	FreePlacesCount int
}

type DtoTimeSlotsByPlaceID struct {
	Rows []struct {
		PlaceId string `json:"place_id"`
		Count   int    `json:"count"`
		Dates   []struct {
			PlaceDateId string         `json:"place_date_id"`
			Count       int            `json:"count"`
			Title       string         `json:"title"`
			TimeSlots   map[string]int `json:"time_slots"`
		} `json:"dates"`
	} `json:"rows"`
}

func GetPlaceInfo(placeID string, userToken string) *TimeSlotsPlace {
	timeSlotsByPlaceID := getTimeSlots(placeID, userToken)

	timeSlots := convertTimeSlots(timeSlotsByPlaceID, timeSlotsByPlaceID.Rows[0].Dates[0].Count)

	return &TimeSlotsPlace{
		LocationName: timeSlotsByPlaceID.Rows[0].Dates[0].Title,
		MaxCount:     timeSlotsByPlaceID.Rows[0].Dates[0].Count,
		TimeSlots:    timeSlots,
	}
}

func getTimeSlots(placeID string, userToken string) DtoTimeSlotsByPlaceID {
	const timeSlotURL = "https://mf.technolab.com.ru/v1/place/dates-list"
	urlGet := fmt.Sprintf("%s?place_id=%s", timeSlotURL, placeID)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100)*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, urlGet, nil)
	req.Header.Set("x-app-token", userToken)
	httpClient := http.Client{Timeout: time.Duration(100) * time.Second}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	timeSlotsByPlaceID := DtoTimeSlotsByPlaceID{}
	jsonErr := json.Unmarshal(body, &timeSlotsByPlaceID)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return timeSlotsByPlaceID
}

func convertTimeSlots(timeSlotsByPlaceID DtoTimeSlotsByPlaceID, maxCount int) []TimeSlot {
	timeSlots := make([]TimeSlot, 0, len(timeSlotsByPlaceID.Rows[0].Dates[0].TimeSlots))

	for slotTime, reservedCount := range timeSlotsByPlaceID.Rows[0].Dates[0].TimeSlots {
		goTime, _ := parseTime(slotTime)
		timeSlot := TimeSlot{PlaceDateId: timeSlotsByPlaceID.Rows[0].Dates[0].PlaceDateId, Time: goTime, FreePlacesCount: maxCount - reservedCount}

		if timeSlot.FreePlacesCount > 0 {
			timeSlots = append(timeSlots, timeSlot)
		}
	}
	sort.Slice(timeSlots, func(i, j int) bool {
		return timeSlots[i].Time.Before(timeSlots[j].Time)
	})
	return timeSlots
}

func parseTime(slotTime string) (time.Time, error) {
	const timeLayout = "2006-01-02T15:04:05.000Z"
	return time.Parse(timeLayout, slotTime)
}
