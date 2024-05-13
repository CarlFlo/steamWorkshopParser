package steamWorkshop

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WorkshopData struct {
	AppID      string
	WorkshopID string
}

var linkToGameSelector = "#ig_bottom > div.breadcrumbs > a:nth-child(1)"

var (
	errFailedToFetchAppID    = status.Errorf(codes.NotFound, "unable to fetch the games appID for the provided workshop item")
	errNoWorkshopURLProvided = status.Errorf(codes.InvalidArgument, "no workshop url provided")
	errInvalidURL            = status.Errorf(codes.InvalidArgument, "invalid workshop URL")
)

func checkInput(url string) error {

	if len(url) == 0 {
		return errNoWorkshopURLProvided
	}

	r := regexp.MustCompile(`^https:\/\/steamcommunity\.com\/sharedfiles\/filedetails\/\?id=\d+$`)
	if !r.MatchString(url) {
		return errInvalidURL
	}

	return nil
}

func RunProgram(url string) (string, error) {

	if err := checkInput(url); err != nil {
		return "", err
	}

	item := WorkshopData{}

	item.WorkshopID = parseWorkshopID(url)

	var err error
	if item.AppID, err = visitPage(url); err != nil {
		return "", err
	}

	output := item.toString()
	return output, nil
}

func (item WorkshopData) toString() string {

	return fmt.Sprintf("workshop_download_item %s %s", item.AppID, item.WorkshopID)
}

func visitPage(url string) (string, error) {
	c := colly.NewCollector()

	var appID string
	c.OnHTML(linkToGameSelector, func(e *colly.HTMLElement) {
		linkToGame := e.Attr("href")
		appID = parseAppID(linkToGame)
	})

	// https://go-colly.org/docs/introduction/start/
	err := c.Visit(url)

	if len(appID) == 0 {
		return "", errFailedToFetchAppID
	}

	return appID, err
}

func parseAppID(linkToGame string) string {
	split := strings.Split(linkToGame, "/")
	appID := split[len(split)-1]
	return appID
}

func parseWorkshopID(linkToWorkshop string) string {
	split := strings.Split(linkToWorkshop, "=")
	workshopID := split[len(split)-1]
	return workshopID
}
