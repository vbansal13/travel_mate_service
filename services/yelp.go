package services

import (
	"fmt"

	"github.com/vbansal/travel_mate_service/helpers"
	travel_matepb "github.com/vbansal/travel_mate_service/proto_bufs"
)

var yelpAPIKey = "QukW_DfpYicraXEUIxX7OW8AbsDELMd8xu0YXZHKVZKJctc9x4CABR0uJTCw19vnpRQogbzTbedX-A5YW4zGfwSvqubQy9q3-XH4Q7iDTXnR16Xjr9th9eiaBxd8W3Yx"
var yelpHost = "https://api.yelp.com/v3/businesses"

//YelpBusinessReview represents data model for yelp business review
type YelpBusinessReview struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Text        string `json:"text"`
	Rating      int32  `json:"rating"`
	TimeCreated string `json:"time_created"`
	User        struct {
		ID         string `json:"id"`
		ProfileURL string `json:"profile_url"`
		ImageURL   string `json:"image_url"`
		Name       string `json:"name"`
	} `json:"user"`
}

//YelpBusinessReviews represents data model containing yelp business reviews
type YelpBusinessReviews struct {
	Reviews           []YelpBusinessReview `json:"reviews"`
	Total             int                  `json:"total"`
	PossibleLanguages []string             `json:"possible_languages"`
}

//YelpBusiness represents Yelp Business data model
type YelpBusiness struct {
	ID          string `json:"id"`
	Alias       string `json:"alias"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	IsClosed    bool   `json:"is_closed"`
	URL         string `json:"url"`
	ReviewCount int32  `json:"review_count"`
	Categories  []struct {
		Alias string `json:"alias"`
		Title string `json:"title"`
	} `json:"categories"`
	Rating      float32 `json:"rating"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	Transactions []string `json:"transactions"`
	Price        string   `json:"price"`
	Location     struct {
		Address1       string   `json:"address1"`
		Address2       string   `json:"address2"`
		Address3       string   `json:"address3"`
		City           string   `json:"city"`
		ZipCode        string   `json:"zip_code"`
		Country        string   `json:"country"`
		State          string   `json:"state"`
		DisplayAddress []string `json:"display_address"`
		CrossStreets   string   `json:"cross_streets"`
	} `json:"location"`
	Phone        string   `json:"phone"`
	DisplayPhone string   `json:"display_phone"`
	Distance     float64  `json:"distance"`
	Photos       []string `json:"photos"`
	Hours        []struct {
		Open []struct {
			IsOvernight bool   `json:"is_overnight"`
			Start       string `json:"start"`
			End         string `json:"end"`
			Day         int    `json:"day"`
		} `json:"open"`
		HoursType string `json:"hours_type"`
		IsOpenNow bool   `json:"is_open_now"`
	} `json:"hours"`
	Reviews []YelpBusinessReview `json:"reviews"`
}

//YelpBusinessList represents Yelp Business data model
type YelpBusinessList struct {
	Businesses []YelpBusiness `json:"businesses"`
	Total      int            `json:"total"`
	Region     struct {
		Center struct {
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		} `json:"center"`
	} `json:"region"`
}

//searchPlacesOnYelp is a helper function for searching places that are around a given lat-long position
//The passed params should describe lat-long position, radius and text that needs to be searched on Yelp
/*
func searchPlacesOnYelp(reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {

	reqPath := "search"

	uRL := fmt.Sprintf("%s/%s?term=%s&latitude=%s&longitude=%s&radius=%s", yelpHost, reqPath, reqParams.Text, reqParams.Latitude, reqParams.Longitude, reqParams.Radius)

	reqHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", yelpAPIKey),
	}
	makeNetworkRequest("GET", reqHeaders, uRL, Yelp, serviceChannel)
}
*/

func searchPlacesOnYelp(reqParams travel_matepb.PlaceSearchRequest, serviceChannel chan ResponseData) {

	reqPath := "search"

	uRL := fmt.Sprintf("%s/%s?term=%s&latitude=%s&longitude=%s&radius=%s", yelpHost, reqPath, reqParams.GetText(),
		reqParams.GetLatitude(), reqParams.GetLongitude(), reqParams.GetRadius())

	reqHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", yelpAPIKey),
	}
	makeNetworkRequest("GET", reqHeaders, uRL, Yelp, serviceChannel)
}

//getPlaceDetailsOnYelp is a helper function for searching detailed information for a given business id.
func getPlaceDetailsOnYelp(reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {
	uRL := fmt.Sprintf("%s/%s", yelpHost, reqParams.YelpID)

	reqHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", yelpAPIKey),
	}
	makeNetworkRequest("GET", reqHeaders, uRL, Yelp, serviceChannel)
}

//getPlaceReviewsOnYelp is a helper function for searching reviews  for a given business id.
func getPlaceReviewsOnYelp(reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {
	uRL := fmt.Sprintf("%s/%s/reviews", yelpHost, reqParams.YelpID)

	reqHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", yelpAPIKey),
	}
	makeNetworkRequest("GET", reqHeaders, uRL, Yelp, serviceChannel)
}
