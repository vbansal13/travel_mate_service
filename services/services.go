package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vbansal/travel_mate_service/helpers"
	travel_matepb "github.com/vbansal/travel_mate_service/proto_bufs"
)

//Types respresents all supported services types available
type Types int

const (
	//Yelp service type
	Yelp Types = iota
	//Google service type
	Google
)

func (t Types) String() string {
	return [...]string{"yelp", "google"}[t]
}

//ErrorResponse is data model for replying error data to client
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"error_message"`
}

//GetTypes is a helper function that will return a list of supported service types
func GetTypes() []Types {
	types := []Types{Yelp, Google}
	return types
}

//ResponseData describes data that service returns
type ResponseData struct {
	Type       Types
	Data       []byte
	Error      error
	RequestURL string
}

//SearchPlaces is a helper function for searching places that are around a given lat-long position
//The passed params should describe lat-long position, radius and text that needs to be searched on given serviceType
/*
func SearchPlaces(serviceType Types, reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {
	switch serviceType {
	case Yelp:
		searchPlacesOnYelp(reqParams, serviceChannel)
	case Google:
		searchPlacesOnGoogle(reqParams, serviceChannel)
	}
}
*/

//SearchPlaces is a helper function for searching places that are around a given lat-long position
//The passed params should describe lat-long position, radius and text that needs to be searched on given serviceType
func SearchPlaces(serviceType Types, reqParams travel_matepb.PlaceSearchRequest, serviceChannel chan ResponseData) {
	switch serviceType {
	case Yelp:
		searchPlacesOnYelp(reqParams, serviceChannel)
	case Google:
		searchPlacesOnGoogle(reqParams, serviceChannel)
	}
}

//GetPlaceDetails is a helper function for fetching place details for a given service type
func GetPlaceDetails(serviceType Types, reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {

	switch serviceType {
	case Yelp:
		yelpServiceChannel := make(chan ResponseData)
		go getPlaceDetailsOnYelp(reqParams, yelpServiceChannel)
		go getPlaceReviewsOnYelp(reqParams, yelpServiceChannel)
		yelpData := YelpBusiness{}
		yelpReviews := YelpBusinessReviews{}
		for i := 0; i < 2; i++ {
			serviceData := <-yelpServiceChannel
			if strings.Contains(serviceData.RequestURL, "reviews") {
				yelpReviews = generateYelpReviewsData(serviceData.Data)
			} else {
				yelpData = GenerateYelpData(serviceData.Data)
			}
		}

		for _, review := range yelpReviews.Reviews {
			yelpData.Reviews = append(yelpData.Reviews, review)
		}
		dataMarshalRes, dataMarshalError := json.Marshal(yelpData)
		if dataMarshalError != nil {
			fmt.Println("Yelp Places details marshal error :", dataMarshalError)
		}
		serviceChannel <- ResponseData{
			Type:       Yelp,
			RequestURL: "yelpPlaceDetails",
			Data:       dataMarshalRes,
			Error:      nil,
		}
	case Google:
		getPlaceDetailsOnGoogle(reqParams, serviceChannel)
	}
}

//GetPhoto is a helper function for fetching image data with given photo reference.
//Currently this method is only needed for Google.
func GetPhoto(serviceType Types, reqParams helpers.ClientRequestParams, serviceChannel chan ResponseData) {
	switch serviceType {
	case Google:
		getPhoto(reqParams, serviceChannel)
	}
}

//GenerateYelpDataList is a helper function for generating Yelp business from []byte
func GenerateYelpDataList(data []byte) YelpBusinessList {
	businessList := YelpBusinessList{}
	unmarshalError := json.Unmarshal([]byte(data), &businessList)
	if unmarshalError != nil {
		fmt.Println("Yelp Places data unmarshal error :", unmarshalError)
	}
	return businessList
}

//GenerateGoogleDataList is a helper function for generating Google business list from []byte
func GenerateGoogleDataList(data []byte) GoogleBusinessList {
	businessList := GoogleBusinessList{}
	unmarshalError := json.Unmarshal([]byte(data), &businessList)
	if unmarshalError != nil {
		fmt.Println("Google Places data unmarshal error :", unmarshalError)
	}
	return businessList
}

//GenerateYelpData is a helper function for generating Yelp business from []byte
func GenerateYelpData(data []byte) YelpBusiness {
	businessData := YelpBusiness{}
	unmarshalError := json.Unmarshal([]byte(data), &businessData)
	if unmarshalError != nil {
		fmt.Println("Yelp Places details unmarshal error :", unmarshalError)
	}
	return businessData
}

//GenerateGoogleData is a helper function for generating Google business list from []byte
func GenerateGoogleData(data []byte) GoogleBusinessDetails {
	businessData := GoogleBusinessDetails{}
	unmarshalError := json.Unmarshal([]byte(data), &businessData)
	if unmarshalError != nil {
		fmt.Println("Google Places details unmarshal error :", unmarshalError)
	}
	return businessData
}

func generateYelpReviewsData(data []byte) YelpBusinessReviews {
	businessReviews := YelpBusinessReviews{}
	unmarshalError := json.Unmarshal([]byte(data), &businessReviews)
	if unmarshalError != nil {
		fmt.Println("Yelp Reviews data unmarshal error :", unmarshalError)
	}
	return businessReviews
}
func makeNetworkRequest(requestType string, requestHeaders map[string]string,
	urlString string, serviceType Types, serviceChannel chan ResponseData) {
	client := &http.Client{}
	req, err := http.NewRequest(requestType, urlString, nil)

	for key, value := range requestHeaders {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(" Error: ", requestType, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	serviceChannel <- ResponseData{
		Type:       serviceType,
		Data:       body,
		Error:      err,
		RequestURL: urlString,
	}
}
