package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vbansal/travel_mate_service/helpers"
	"github.com/vbansal/travel_mate_service/services"
)

func createErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	errorResponse := services.ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    "Place id missing from request!",
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	dataMarshalRes, _ := json.Marshal(errorResponse)
	w.Write(dataMarshalRes)
}

func placesSearchHandler(w http.ResponseWriter, req *http.Request) {

	reqParams, err := helpers.ExtractParamsFromRequest(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	serviceChannel := make(chan services.ResponseData)

	for _, serviceType := range services.GetTypes() {
		go services.SearchPlaces(serviceType, reqParams, serviceChannel)
	}

	yelpBusinessList := services.YelpBusinessList{}
	googleBusinessList := services.GoogleBusinessList{}

	for i := 0; i < len(services.GetTypes()); i++ {
		serviceData := <-serviceChannel
		if serviceData.Type == services.Google && serviceData.Error == nil {
			googleBusinessList = services.GenerateGoogleDataList(serviceData.Data)
		} else if serviceData.Type == services.Yelp && serviceData.Error == nil {
			yelpBusinessList = services.GenerateYelpDataList(serviceData.Data)
		}
	}

	candidates := combinedCandidateList(googleBusinessList, yelpBusinessList)

	dataMarshalRes, dataMarshalError := json.Marshal(candidates)

	if dataMarshalError != nil {
		fmt.Println("Data Marshal error: ", dataMarshalError)
		http.Error(w, dataMarshalError.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(dataMarshalRes)
	}

}

func placesDetailsHandler(w http.ResponseWriter, req *http.Request) {

	reqParams, err := helpers.ExtractParamsFromRequest(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	serviceChannel := make(chan services.ResponseData)

	if len(reqParams.YelpID) > 0 {
		go services.GetPlaceDetails(services.Yelp, reqParams, serviceChannel)
	} else if len(reqParams.GoogleID) > 0 {
		go services.GetPlaceDetails(services.Google, reqParams, serviceChannel)
	} else {
		createErrorResponse(w, http.StatusBadRequest, "Place id missing from request!")
		return
	}

	serviceData := <-serviceChannel

	businessDetails := Candidate{}
	if serviceData.Type == services.Google && serviceData.Error == nil {
		googleBusinessData := services.GenerateGoogleData(serviceData.Data)
		businessDetails = convertFromGoogleModel(googleBusinessData.Candidate)
	} else {
		yelpBusinessData := services.GenerateYelpData(serviceData.Data)
		businessDetails = convertFromYelpModel(yelpBusinessData)
	}

	dataMarshalRes, dataMarshalError := json.Marshal(businessDetails)

	if dataMarshalError != nil {
		fmt.Println("Place Details data marshal error: ", dataMarshalError)
		http.Error(w, dataMarshalError.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(dataMarshalRes)
	}

}

func placesImageHandler(w http.ResponseWriter, req *http.Request) {

	reqParams, err := helpers.ExtractParamsFromRequest(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if len(reqParams.GooglePhotoReference) == 0 {
		createErrorResponse(w, http.StatusBadRequest, "Photo reference missing from request!")
		return
	}

	serviceChannel := make(chan services.ResponseData)
	go services.GetPhoto(services.Google, reqParams, serviceChannel)

	serviceData := <-serviceChannel

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(serviceData.Data)))
	w.Write(serviceData.Data)
}

func main() {
	//http.Handle("/events/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/v1/places/search", placesSearchHandler)
	http.HandleFunc("/v1/places/details", placesDetailsHandler)
	http.HandleFunc("/v1/places/photo", placesImageHandler)

	http.ListenAndServe(":8082", nil)
}
