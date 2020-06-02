package services

import (
	"fmt"

	"github.com/vbansal/travel_mate_service/helpers"
	travel_matepb "github.com/vbansal/travel_mate_service/proto_bufs"
)

var googleAPIKey = "AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo"
var googleHost = "https://maps.googleapis.com/maps/api/place"

//GoogleGeometry represents google location model
type GoogleGeometry struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	Viewport struct {
		Northeast struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"northeast"`
		Southwest struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"southwest"`
	} `json:"viewport"`
}

//GoogleBusinessDetails represents data model for Google business
//This is what is returned when a place details request is made
type GoogleBusinessDetails struct {
	HTMLAttributions []interface{}  `json:"html_attributions"`
	Candidate        GoogleBusiness `json:"result"`
	Status           string         `json:"status"`
}

//GoogleBusiness represents a single Google Place data model
type GoogleBusiness struct {
	ID           string `json:"id"`
	PlaceID      string `json:"place_id"`
	Name         string `json:"name"`
	OpeningHours struct {
		OpenNow bool `json:"open_now"`
		Periods []struct {
			Close struct {
				Day  int    `json:"day"`
				Time string `json:"time"`
			} `json:"close"`
			Open struct {
				Day  int    `json:"day"`
				Time string `json:"time"`
			} `json:"open"`
		} `json:"periods"`
		WeekdayText []string `json:"weekday_text"`
	} `json:"opening_hours"`
	Geometry GoogleGeometry `json:"geometry"`
	Photos   []struct {
		Height           int      `json:"height"`
		HTMLAttributions []string `json:"html_attributions"`
		PhotoReference   string   `json:"photo_reference"`
		Width            int      `json:"width"`
	} `json:"photos"`
	Rating   float32 `json:"rating"`
	PlusCode struct {
		CompoundCode string `json:"compound_code"`
		GlobalCode   string `json:"global_code"`
	} `json:"plus_code"`
	PriceLevel        int      `json:"price_level"`
	Reference         string   `json:"reference"`
	Scope             string   `json:"scope"`
	Types             []string `json:"types"`
	Vicinity          string   `json:"vicinity"`
	AddressComponents []struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	} `json:"address_components"`
	AdrAddress               string `json:"adr_address"`
	FormattedAddress         string `json:"formatted_address"`
	FormattedPhoneNumber     string `json:"formatted_phone_number"`
	Icon                     string `json:"icon"`
	InternationalPhoneNumber string `json:"international_phone_number"`
	Reviews                  []struct {
		AuthorName              string `json:"author_name"`
		AuthorURL               string `json:"author_url"`
		Language                string `json:"language"`
		ProfilePhotoURL         string `json:"profile_photo_url"`
		Rating                  int32  `json:"rating"`
		RelativeTimeDescription string `json:"relative_time_description"`
		Text                    string `json:"text"`
		Time                    int    `json:"time"`
	} `json:"reviews"`
	URL       string `json:"url"`
	UtcOffset int    `json:"utc_offset"`
	Website   string `json:"website"`
}

//GoogleBusinessList represents list of Google Places data model
//This is returned when a search nearby request is made for given keyword
type GoogleBusinessList struct {
	HTMLAttributions []interface{}    `json:"html_attributions"`
	NextPageToken    string           `json:"next_page_token"`
	Candidates       []GoogleBusiness `json:"results"`
	DebugLog         struct {
		Line []interface{} `json:"line"`
	} `json:"debug_log"`
	Status string `json:"status"`
}

//searchPlacesOnGoogle is a helper function for searching places that are around a given lat-long position
//The passed params should describe lat-long position, radius and text that needs to be searched on Yelp
/*func searchPlacesOnGoogle(reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {

	path := "nearbysearch/json"

	//https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=indian%20food&inputtype=textquery&fields=photos,formatted_address,name,opening_hours,rating,price_level&locationbias=circle:9000@34.387616,-118.597237&key=AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo
	//https://maps.googleapis.com/maps/api/place/nearbysearch/json?keyword=thai%20food&location=34.387616,-118.597237&key=AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo&radius=9000
	uRL := fmt.Sprintf("%s/%s?keyword=%s&radius=%s&location=%s,%s&key=%s", googleHost, path,
		reqParams.Text, reqParams.Radius, reqParams.Latitude, reqParams.Longitude, googleAPIKey)

	reqHeaders := map[string]string{}
	makeNetworkRequest("GET", reqHeaders, uRL, Google, serviceChannel)
}*/

func searchPlacesOnGoogle(reqParams travel_matepb.PlaceSearchRequest, serviceChannel chan ResponseData) {

	path := "nearbysearch/json"

	//https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=indian%20food&inputtype=textquery&fields=photos,formatted_address,name,opening_hours,rating,price_level&locationbias=circle:9000@34.387616,-118.597237&key=AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo
	//https://maps.googleapis.com/maps/api/place/nearbysearch/json?keyword=thai%20food&location=34.387616,-118.597237&key=AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo&radius=9000
	uRL := fmt.Sprintf("%s/%s?keyword=%s&radius=%s&location=%s,%s&key=%s", googleHost, path,
		reqParams.GetText(), reqParams.GetRadius(), reqParams.GetLatitude(), reqParams.GetLongitude(), googleAPIKey)

	reqHeaders := map[string]string{}
	makeNetworkRequest("GET", reqHeaders, uRL, Google, serviceChannel)
}

func getPlaceDetailsOnGoogle(reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {
	path := "details/json"
	//https://maps.googleapis.com/maps/api/place/details/json?placeid=ChIJWW-OcJuHwoARYfu84ZSZNKY&key=AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo
	uRL := fmt.Sprintf("%s/%s?placeid=%s&key=%s", googleHost, path, reqParams.GoogleID, googleAPIKey)
	reqHeaders := map[string]string{}
	makeNetworkRequest("GET", reqHeaders, uRL, Google, serviceChannel)
}

func getPhoto(reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {
	path := "photo"
	width := "400"
	if len(reqParams.PhotoWidth) > 0 {
		width = reqParams.PhotoWidth
	}
	//https://maps.googleapis.com/maps/api/place/photo?photoreference=ChIJWW-OcJuHwoARYfu84ZSZNKY&key=AIzaSyA72CWmMBvvOhov3sbkcmyBHTC9yb4NCAo
	uRL := fmt.Sprintf("%s/%s?photoreference=%s&maxwidth=%s&key=%s", googleHost, path, reqParams.GooglePhotoReference, width, googleAPIKey)
	reqHeaders := map[string]string{}
	makeNetworkRequest("GET", reqHeaders, uRL, Google, serviceChannel)
}
