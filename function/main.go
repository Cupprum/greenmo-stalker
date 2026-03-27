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
	lat1, _ := strconv.ParseFloat(params["lat1"], 64)
	lon1, _ := strconv.ParseFloat(params["lon1"], 64)
	lat2, _ := strconv.ParseFloat(params["lat2"], 64)
	lon2, _ := strconv.ParseFloat(params["lon2"], 64)
	fuel, _ := strconv.Atoi(params["desiredFuelLevel"])
	if fuel == 0 {
		fuel = 40
	}

	p1, p2 := geo.Position{Lat: lat1, Lon: lon1}, geo.Position{Lat: lat2, Lon: lon2}
	center := geo.Position{Lat: (p1.Lat + p2.Lat) / 2, Lon: (p1.Lon + p2.Lon) / 2}

	radius := geo.Distance(p1, p2) / 2

	var cars, chargers []geo.Position
	var err error

	if params["cars"] == "true" {
		if cars, err = greenmobility.Query("https://platform.api.gourban.services/v1/hb98ga69/front/vehicles", center, radius, fuel); err != nil {
			return 500, "", nil, fmt.Errorf("greenmo error: %w", err)
		}
	}
	if params["chargers"] == "true" {
		// Spirii uses NE/SW
		if chargers, err = spirii.Query("https://app.spirii.dk/api/v2/clusters", geo.Position{Lat: lat1, Lon: lon2}, geo.Position{Lat: lat2, Lon: lon1}); err != nil {
			return 500, "", nil, fmt.Errorf("spirii error: %w", err)
		}
	}

	if len(cars) == 0 && len(chargers) == 0 {
		return 200, "application/json", []byte(`{"message":"No results found"}`), nil
	}

	key := os.Getenv("MAP_API_KEY")
	img, err := openstreetmaps.GenerateMap("https://maps.geoapify.com/v1/staticmap", center, cars, chargers, key)
	if err != nil {
		return 500, "", nil, fmt.Errorf("map error: %w", err)
	}

	return 200, "image/png", img, nil
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code, contentType, body, err := coreLogic(req.QueryStringParameters)
	if err != nil {
		log.Printf("Execution Error: %v", err)
		b, _ := json.Marshal(map[string]string{"error": err.Error()})
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
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(handler)
	} else {
		res, _, body, err := coreLogic(map[string]string{"lat1": "55.79", "lon1": "12.51", "lat2": "55.77", "lon2": "12.52", "cars": "true"})
		fmt.Printf("Status: %d, Error: %v, Body Len: %d\n", res, err, len(body))
	}
}
