package service

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	UserToken string
}

func NewService(token string) *Service {
	return &Service{UserToken: token}
}

func (s *Service) GetPlacesInsideTimePeriod(isTimePeriodFilter func(time.Time) bool) Places {
	locationIDs := getPlaceIDs()

	places := make(Places)
	for _, locationID := range locationIDs {
		placeInfo := *GetPlaceInfo(locationID, s.UserToken)
		places[locationID] = TimeSlotsPlace{
			TimeSlots:    placeInfo.TimeSlots,
			LocationName: placeInfo.LocationName,
			MaxCount:     placeInfo.MaxCount,
		}
	}

	placesInTimePeriod := make(Places)
	for placeId, timeSlots := range places {
		var timeSlotsInTimePeriod []TimeSlot
		for _, timeSlot := range timeSlots.TimeSlots {
			if isTimePeriodFilter(timeSlot.Time) {
				timeSlotsInTimePeriod = append(timeSlotsInTimePeriod, timeSlot)
			}
		}
		if len(timeSlotsInTimePeriod) > 0 {
			placesInTimePeriod[placeId] = TimeSlotsPlace{
				TimeSlots:    timeSlotsInTimePeriod,
				LocationName: places[placeId].LocationName,
				MaxCount:     places[placeId].MaxCount,
			}
		}
	}
	return placesInTimePeriod
}

func getPlaceIDs() []string {
	const urlWithPlaces = "https://moscowseasons.com/event/katki-moskvy-2022-2023-adresa-vrema-raboty-registracia-pravila-posesenia/"
	resp, err := http.Get(urlWithPlaces)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	placeIDs := make([]string, 0)
	doc.Find("div.raw-html").First().Find("li>a").Each(func(i int, selection *goquery.Selection) {
		place, _ := selection.Attr("href")
		u, _ := url.Parse(place)
		const ticketsHostName = "mftickets.technolab.com.ru"
		if u.Host == ticketsHostName {
			placeID := strings.Split(u.Path, "/")[2]
			placeIDs = append(placeIDs, placeID)
		}
	})
	return placeIDs
}
