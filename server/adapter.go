package main

import (
	"fmt"
	"strings"

	travel_matepb "github.com/vbansal/travel_mate_service/proto_bufs"
	"github.com/vbansal/travel_mate_service/services"
)

//CandidateSource describes sources of this candidate, could be Google, Yelp or some other source.
type CandidateSource struct {
	Name        string  `json:"name"`
	ID          string  `json:"id"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
}

//ReviewAuthor described data model for author who adds a review
type ReviewAuthor struct {
	ProfileURL string `json:"profile_url"`
	ImageURL   string `json:"image_url"`
	Name       string `json:"name"`
}

//Review desribes data model for user review
type Review struct {
	Author                  ReviewAuthor `json:"author"`
	Rating                  int          `json:"rating"`
	RelativeTimeDescription string       `json:"relative_time_description"`
	Text                    string       `json:"text"`
	Time                    string       `json:"time"`
	Source                  string       `json:"source"`
}

//Candidate represents TravelMate Business data model
type Candidate struct {
	Sources  []CandidateSource `json:"sources"`
	Name     string            `json:"name"`
	IsClosed bool              `json:"is_closed"`
	Website  string            `json:"website"`

	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	FormattedAddress string   `json:"formatted_address"`
	Price            string   `json:"price"`
	Phone            string   `json:"phone"`
	DisplayPhone     string   `json:"display_phone"`
	Distance         float64  `json:"distance"`
	Photos           []string `json:"photos"`
	Reviews          []Review `json:"reviews"`
}

//CandidateList represents TravelMate Business data model
type CandidateList struct {
	Candidates []Candidate `json:"candidates"`
	Total      int         `json:"total"`
}

//AppendSource adds another candidate source to received candidate
func appendSource(sourceCandidate *travel_matepb.Candidate, anotherCandidate travel_matepb.Candidate) {
	if sourceExists(sourceCandidate, anotherCandidate.Sources[0].Name) {
		return
	}
	sourceCandidate.Sources = append(sourceCandidate.Sources, anotherCandidate.Sources[0])
}

//SourceExists will find given source in the Candidate
func sourceExists(sourceCandidate *travel_matepb.Candidate, sourceName string) bool {
	for _, source := range sourceCandidate.Sources {
		if source.Name == sourceName {
			return true
		}
	}
	return false
}

//convertFromYelpModel converts YelpBusiness model to TravelMate Candidate model
func convertFromYelpModel(yelpModel services.YelpBusiness) travel_matepb.Candidate {
	var candidate travel_matepb.Candidate
	candidate.Name = yelpModel.Name
	candidate.Photos = append(candidate.Photos, yelpModel.ImageURL)
	candidate.IsClosed = yelpModel.IsClosed
	candidate.Website = yelpModel.URL
	candidate.Coordinates = &travel_matepb.Coordinates{
		Latitude:  yelpModel.Coordinates.Latitude,
		Longitude: yelpModel.Coordinates.Longitude,
	}
	candidate.FormattedAddress = fmt.Sprintf("%s, %s", strings.Join(yelpModel.Location.DisplayAddress, ", "), yelpModel.Location.Country)
	candidate.Price = yelpModel.Price
	candidate.Phone = yelpModel.Phone
	candidate.DisplayPhone = yelpModel.DisplayPhone
	candidate.Distance = yelpModel.Distance

	for _, photoURL := range yelpModel.Photos {
		candidate.Photos = append(candidate.Photos, photoURL)
	}

	for _, review := range yelpModel.Reviews {
		candidateReview := &travel_matepb.Review{
			Author: &travel_matepb.ReviewAuthor{
				Name:       review.User.Name,
				ProfileUrl: review.User.ProfileURL,
				ImageUrl:   review.User.ImageURL,
			},
			Text: review.Text,
			Time: review.TimeCreated,
			/*RelativeTimeDescription: ?*/ //TBD
			Rating:                        review.Rating,
			Source:                        services.Yelp.String(),
		}
		candidate.Reviews = append(candidate.Reviews, candidateReview)
	}

	source := &travel_matepb.CandidateSource{
		Id:          yelpModel.ID,
		Rating:      yelpModel.Rating,
		ReviewCount: yelpModel.ReviewCount,
		Name:        services.Yelp.String(),
	}

	//generatedDeck = append(generatedDeck, value+" of "+suit)
	candidate.Sources = append(candidate.Sources, source)

	return candidate
}

//convertFromGoogleModel converts GoogleBusiness model to TravelMate Candidate model
func convertFromGoogleModel(googleModel services.GoogleBusiness) travel_matepb.Candidate {

	var candidate travel_matepb.Candidate

	candidate.Name = googleModel.Name
	//candidate.ImageURL = googleModel.
	candidate.IsClosed = !googleModel.OpeningHours.OpenNow
	//candidate.URL = yelpModel.URL
	//candidate.Coordinates = yelpModel.Coordinates
	candidate.FormattedAddress = googleModel.Vicinity
	candidate.Price = strings.Repeat("$", googleModel.PriceLevel)
	//candidate.Phone = googleModel.
	//candidate.DisplayPhone = yelpModel.DisplayPhone
	//candidate.Distance = yelpModel.Distance
	candidate.Website = googleModel.Website

	for _, photoStruct := range googleModel.Photos {
		candidate.Photos = append(candidate.Photos, photoStruct.PhotoReference)
	}

	for _, review := range googleModel.Reviews {
		candidateReview := &travel_matepb.Review{
			Author: &travel_matepb.ReviewAuthor{
				Name:       review.AuthorName,
				ProfileUrl: review.AuthorURL,
				ImageUrl:   review.ProfilePhotoURL,
			},
			Text: review.Text,
			/*Time: review.Time,*/   //TBD This needs to be converted to string time
			RelativeTimeDescription: review.RelativeTimeDescription,
			Rating:                  review.Rating,
			Source:                  services.Google.String(),
		}
		candidate.Reviews = append(candidate.Reviews, candidateReview)
	}

	source := &travel_matepb.CandidateSource{
		Id:          googleModel.PlaceID,
		Rating:      googleModel.Rating,
		ReviewCount: 0,
		Name:        services.Google.String(),
	}

	//generatedDeck = append(generatedDeck, value+" of "+suit)
	candidate.Sources = append(candidate.Sources, source)
	return candidate
}

func combinedCandidateList(googleBusinessList services.GoogleBusinessList,
	yelpBusinessList services.YelpBusinessList) []*travel_matepb.Candidate {
	candidates := []*travel_matepb.Candidate{}
	googleCandidates := []travel_matepb.Candidate{}
	googleProcessedIndex := map[int]bool{}

	for index, googleBusiness := range googleBusinessList.Candidates {
		googleCandidate := convertFromGoogleModel(googleBusiness)
		googleCandidates = append(googleCandidates, googleCandidate)
		googleProcessedIndex[index] = false
	}

	for _, yelpBusiness := range yelpBusinessList.Businesses {
		yelpCandidate := convertFromYelpModel(yelpBusiness)
		for index, googleCandidate := range googleCandidates {
			if strings.Contains(yelpCandidate.FormattedAddress, googleCandidate.FormattedAddress) {
				appendSource(&yelpCandidate, googleCandidate)
				googleProcessedIndex[index] = true
			}
		}
		candidates = append(candidates, &yelpCandidate)
	}

	for index, googleCandidate := range googleCandidates {
		if googleProcessedIndex[index] == false {
			candidates = append(candidates, &googleCandidate)
		}
	}
	return candidates
}
