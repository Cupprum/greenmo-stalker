package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"function/geo"
	"function/greenmobility"
	"function/openstreetmaps"
	"function/spirii"
)

func coreLogic(params map[string]string) (int, string, []byte, error) {
	log.Println("Parsing parameters...")
	lat1, _ := strconv.ParseFloat(params["lat1"], 64)
	lon1, _ := strconv.ParseFloat(params["lon1"], 64)
	lat2, _ := strconv.ParseFloat(params["lat2"], 64)
	lon2, _ := strconv.ParseFloat(params["lon2"], 64)
	fuel, _ := strconv.Atoi(params["desiredFuelLevel"])
	if fuel == 0 {
		fuel = 40
	}
	qCars := params["cars"] == "true"
	qChargers := params["chargers"] == "true"
	log.Printf(
		"Lat1: %v, Lon1: %v, Lat2: %v, Lon2: %v, DesiredFuel: %v, Query Cars: %v, Query Chargers: %v\n",
		lat1, lon1, lat2, lon2, fuel, qCars, qChargers,
	)

	p1, p2 := geo.Position{Lat: lat1, Lon: lon1}, geo.Position{Lat: lat2, Lon: lon2}
	center := geo.Position{Lat: (p1.Lat + p2.Lat) / 2, Lon: (p1.Lon + p2.Lon) / 2}
	radius := geo.Distance(p1, p2) / 2

	var cars, chargers []geo.Position
	var err error

	if params["cars"] == "true" {
		log.Println("Querying cars...")
		url := "https://platform.api.gourban.services/v1/hb98ga69/front/vehicles"
		rCars, err := greenmobility.Query(url, center, radius, fuel)
		if err != nil {
			return 500, "", nil, fmt.Errorf("greenmo error: %w", err)
		}
		// Greenmo thinks in circles, we think in squares
		for _, c := range rCars {
			if p1.Lon < c.Lon && c.Lon < p2.Lon && p2.Lat < c.Lat && c.Lat < p1.Lat {
				cars = append(cars, c)
			}
		}
	}
	if params["chargers"] == "true" {
		log.Println("Querying chargers...")
		// Spirii uses NE/SW
		url := "https://app.spirii.dk/api/v2/clusters"
		if chargers, err = spirii.Query(url, geo.Position{Lat: lat1, Lon: lon2}, geo.Position{Lat: lat2, Lon: lon1}); err != nil {
			return 500, "", nil, fmt.Errorf("spirii error: %w", err)
		}
	}

	if len(cars) == 0 && len(chargers) == 0 {
		log.Println("No cars or charges found...")
		return 200, "application/json", []byte(`{"message":"No results found"}`), nil
	}

	log.Println("Generating map...")
	key := os.Getenv("GREENMO_OPEN_MAPS_API_TOKEN")
	url := "https://maps.geoapify.com/v1/staticmap"
	img, err := openstreetmaps.GenerateMap(url, center, cars, chargers, key)
	if err != nil {
		return 500, "", nil, fmt.Errorf("map error: %w", err)
	}

	return 200, "image/png", img, nil
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code, contentType, body, err := coreLogic(req.QueryStringParameters)
	if err != nil {
		log.Printf("Execution Error: %v", err)
		b, _ := json.Marshal(map[string]string{"error": "Internal Error"})
		return events.APIGatewayProxyResponse{StatusCode: code, Body: string(b)}, nil
	}

	isBase64 := contentType == "image/png"
	respBody := string(body)
	if isBase64 {
		respBody = base64.StdEncoding.EncodeToString(body)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      code,
		Headers:         map[string]string{"Content-Type": contentType, "Access-Control-Allow-Origin": "*"},
		Body:            respBody,
		IsBase64Encoded: isBase64,
	}, nil
}

func main() {
	log.Println("Starting execution...")
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Executing in lambda...")
		lambda.Start(handler)
	} else {
		log.Println("Executing locally...")
		res, _, body, err := coreLogic(map[string]string{
			"lat1": "55.739892", "lon1": "12.517685", "lat2": "55.734577", "lon2": "12.526059",
			"cars": "true", "chargers": "true",
		})
		log.Printf("Status: %d, Error: %v, Body Len: %d\n", res, err, len(body))
	}
	log.Println("Execution finished")
}
