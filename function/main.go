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
	osm "function/openstreetmaps"
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
		fuel = 100
	}
	qCars := params["cars"] == "true"
	qChargers := params["chargers"] == "true"
	log.Printf(
		"Lat1: %v, Lon1: %v, Lat2: %v, Lon2: %v, DesiredFuel: %v, Query Cars: %v, Query Chargers: %v\n",
		lat1, lon1, lat2, lon2, fuel, qCars, qChargers,
	)

	nw, se := geo.Position{Lat: lat1, Lon: lon1}, geo.Position{Lat: lat2, Lon: lon2}
	var cars []greenmobility.Car
	var chargers []geo.Position
	var err error

	if params["cars"] == "true" {
		log.Println("Querying cars...")
		url := "https://platform.api.gourban.services/v1/hb98ga69/front/vehicles"
		if cars, err = greenmobility.Query(url, nw, se, fuel); err != nil {
			return 500, "", nil, fmt.Errorf("greenmo error: %w", err)
		}
	}
	if params["chargers"] == "true" {
		log.Println("Querying chargers...")
		url := "https://app.spirii.dk/api/v2/clusters"
		if chargers, err = spirii.Query(url, nw, se); err != nil {
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
	center := geo.Position{Lat: (nw.Lat + se.Lat) / 2, Lon: (nw.Lon + se.Lon) / 2}
	markers := []osm.Marker{}
	for _, c := range cars {
		color := "#3ea635"
		if c.Discounted {
			color = "#e3e31f"
		}
		m := osm.Marker{Pos: c.Pos, Color: color, Text: strconv.Itoa(c.Charge)}
		markers = append(markers, m)
	}
	for _, c := range chargers {
		markers = append(markers, osm.Marker{Pos: c, Color: "#5588d0", Icon: "ev_station"})
	}
	img, err := osm.GenerateMap(url, center, markers, key)
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
			"lat1": "55.740000", "lon1": "12.515000", "lat2": "55.730000", "lon2": "12.530000",
			"cars": "true", "chargers": "true",
		})
		log.Printf("Status: %d, Error: %v, Body Len: %d\n", res, err, len(body))
		fmt.Println(string(base64.StdEncoding.EncodeToString(body)))
	}
	log.Println("Execution finished")
}
