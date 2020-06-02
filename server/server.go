package main

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	travel_matepb "github.com/vbansal/travel_mate_service/proto_bufs"
	"github.com/vbansal/travel_mate_service/services"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) PlaceSearch(ctx context.Context, reqParams *travel_matepb.PlaceSearchRequest) (*travel_matepb.PlaceSearchResponse, error) {

	serviceChannel := make(chan services.ResponseData)

	for _, serviceType := range services.GetTypes() {
		go services.SearchPlaces(serviceType, *reqParams, serviceChannel)
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

	resCandidates := combinedCandidateList(googleBusinessList, yelpBusinessList)

	/*
		dataMarshalRes, dataMarshalError := json.Marshal(candidates)

		if dataMarshalError != nil {
			fmt.Println("Data Marshal error: ", dataMarshalError)
			http.Error(w, dataMarshalError.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(dataMarshalRes)
		}*/
	responseObj := &travel_matepb.PlaceSearchResponse{
		Candidates: resCandidates,
	}

	return responseObj, nil
}

func main() {
	fmt.Println("Hello World")

	//es, _ := elasticsearch.NewDefaultClient()
	//log.Println(es.Info())

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	travel_matepb.RegisterPlaceSearchServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
