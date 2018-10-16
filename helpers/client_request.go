package helpers

import (
	"log"
	"net/http"
	"strings"
)

//ClientRequestParams represents parameters in Client;s request
type ClientRequestParams struct {
	GoogleID             string
	YelpID               string
	Latitude             string
	Longitude            string
	Radius               string
	Text                 string
	GooglePhotoReference string
	PhotoWidth           string
}

//ExtractParamsFromRequest is used to extract parameters from Client's request and put it inside ClientRequestParams object
func ExtractParamsFromRequest(req *http.Request) (ClientRequestParams, error) {

	var reqParams ClientRequestParams
	if err := req.ParseForm(); err != nil {

		log.Println("Error parsing form: ", err)
		return reqParams, err
	}

	reqParams = ClientRequestParams{
		Latitude:             req.Form.Get("latitude"),
		Longitude:            req.Form.Get("longitude"),
		Radius:               req.Form.Get("radius"),
		Text:                 strings.Replace(req.Form.Get("text"), " ", "+", -1),
		GoogleID:             req.Form.Get("google_id"),
		YelpID:               req.Form.Get("yelp_id"),
		GooglePhotoReference: req.Form.Get("google_photoreference"),
		PhotoWidth:           req.Form.Get("photo_width"),
	}
	return reqParams, nil
}
