I'll help you document the current lambda function and analyze the code. Let me start by exploring the repository structure and examining the logic folder.Now let me look at the logic folder more closely to understand the complete implementation:Now I have all the information needed to write a comprehensive specification. Let me create a detailed spec document.

## API Specification: GreenMoStalker Lambda Function

### 1. Overview

The GreenMoStalker Lambda function is a REST API endpoint deployed on AWS that serves as a proxy to query available GreenMobility cars and charging stations. When cars are placed on chargers while low on battery, users receive free driving minutes. This API facilitates automation of fetching car and charger information for a given geographic area and returns a static map visualization.

---

### 2. Current Endpoint Configuration

**Endpoint Path:** `/query`

**HTTP Method:** `GET`

**Authentication:** API Key (Header-based via API Gateway Usage Plan)

**CORS Support:** Enabled for `https://editor.swagger.io`

**Binary Media Support:** `image/png`

---

### 3. Query String Parameters

All parameters are received from the API Gateway event's `queryStringParameters`.

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `lat1` | float | Yes | N/A | Latitude of top-left corner of bounding box |
| `lon1` | float | Yes | N/A | Longitude of top-left corner of bounding box |
| `lat2` | float | Yes | N/A | Latitude of bottom-right corner of bounding box |
| `lon2` | float | Yes | N/A | Longitude of bottom-right corner of bounding box |
| `cars` | boolean | No | false | Whether to query GreenMobility cars (values: "true" or "false") |
| `chargers` | boolean | No | false | Whether to query Spirii chargers (values: "true" or "false") |
| `desiredFuelLevel` | integer | No | 40 | Battery level threshold (%) - only cars at or below this level are returned |

**Example Query:**
```
GET /query?lat1=55.794430&lon1=12.511368&lat2=55.779566&lon2=12.527933&cars=true&chargers=true&desiredFuelLevel=90
```

---

### 4. Function Output

#### 4.1 Success Response (200 OK) - With Results

**Content-Type:** `image/png`

**Body:** Base64-encoded PNG image showing:
- **Green markers** (color: #3ea635): GreenMobility cars with battery <= desiredFuelLevel
- **Red markers** (color: #f30e0e): Available charging stations from Spirii
- **Center:** Geographic center of the bounding box
- **Zoom Level:** 14
- **Style:** Maptiler 3D
- **Size:** 600x600 pixels

```json
{
  "statusCode": 200,
  "headers": {
    "Content-Type": "image/png",
    "Access-Control-Allow-Origin": "https://editor.swagger.io"
  },
  "body": "<base64-encoded-image-data>",
  "isBase64Encoded": true
}
```

#### 4.2 Success Response (200 OK) - No Results

**Content-Type:** `application/json`

```json
{
  "statusCode": 200,
  "headers": {
    "Content-Type": "application/json",
    "Access-Control-Allow-Origin": "https://editor.swagger.io"
  },
  "body": "{\"message\": \"No available cars and chargers were found.\"}"
}
```

#### 4.3 Error Response (400 Bad Request)

**Content-Type:** `application/json`

**Triggers:**
- Missing query string parameters
- Invalid position format (lat/lon not parseable as floats)

```json
{
  "statusCode": 400,
  "headers": {
    "Content-Type": "application/json",
    "Access-Control-Allow-Origin": "https://editor.swagger.io"
  },
  "body": "{\"message\": \"The positions are not in a valid format.\"}"
}
```

#### 4.4 Error Response (403 Forbidden)

**Content-Type:** `application/json`

**Triggers:**
- Invalid response code from external APIs (GreenMobility, Spirii, or Geoapify)
- Network error with external services

```json
{
  "statusCode": 403,
  "headers": {
    "Content-Type": "application/json",
    "Access-Control-Allow-Origin": "https://editor.swagger.io"
  },
  "body": "{\"message\": \"Invalid response code - GreenMo. Got 401, expected 200\"}"
}
```

#### 4.5 Error Response (500 Internal Server Error)

**Content-Type:** `application/json`

**Triggers:**
- Unexpected exceptions during processing

```json
{
  "statusCode": 500,
  "headers": {
    "Content-Type": "application/json",
    "Access-Control-Allow-Origin": "https://editor.swagger.io"
  },
  "body": "{\"message\": \"unknown exception\"}"
}
```

---

### 5. External API Calls

#### 5.1 GreenMobility API - Query Cars

**Purpose:** Fetch GreenMobility vehicles in a specified radius from a center point.

**API Details:**
- **Protocol:** HTTPS
- **Hostname:** `platform.api.gourban.services`
- **Endpoint:** `v1/hb98ga69/front/vehicles`
- **HTTP Method:** GET
- **Query Parameters:**

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `lat` | float | Yes | Center latitude | 55.787 |
| `lng` | float | Yes | Center longitude | 12.519 |
| `rad` | integer | Yes | Search radius in kilometers | 2 |
| `excludeStationedVehicles` | boolean | Yes | Exclude vehicles at stations | true |

**Request URL:**
```
https://platform.api.gourban.services/v1/hb98ga69/front/vehicles?lat=55.787&lng=12.519&rad=2&excludeStationedVehicles=true
```

**Response Format (JSON Array):**
```json
[
  {
    "id": 12345,
    "stateOfCharge": 35,
    "position": {
      "coordinates": [12.515, 55.790]
    }
  },
  {
    "id": 12346,
    "stateOfCharge": 50,
    "position": {
      "coordinates": [12.525, 55.785]
    }
  }
]
```

**Filtering Logic:**
- Only vehicles with `stateOfCharge <= desiredFuelLevel` are retained
- Coordinates are converted from `[longitude, latitude]` to `{lat, lon}` format

**Expected Status Code:** 200

---

#### 5.2 Spirii API - Query Chargers

**Purpose:** Fetch available charging stations within a geographic bounding box.

**API Details:**
- **Protocol:** HTTPS
- **Hostname:** `app.spirii.dk`
- **Endpoint:** `api/v2/clusters`
- **HTTP Method:** GET
- **Query Parameters:**

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `includeOccupied` | boolean | Yes | Include occupied chargers | true |
| `includeOutOfService` | boolean | Yes | Include out-of-service chargers | true |
| `includeRoaming` | boolean | Yes | Include roaming chargers | true |
| `onlyIncludeFavourite` | boolean | Yes | Filter to favorites only | false |
| `neCoordinates` | string | Yes | Northeast corner (lat, lon format) | "55.794430, 12.527933" |
| `swCoordinates` | string | Yes | Southwest corner (lat, lon format) | "55.779566, 12.511368" |
| `zoomLevel` | integer | Yes | Zoom level (for clustering) | 22 |

**Request URL:**
```
https://app.spirii.dk/api/v2/clusters?includeOccupied=true&includeOutOfService=true&includeRoaming=true&onlyIncludeFavourite=false&neCoordinates=55.794430,%2012.527933&swCoordinates=55.779566,%2012.511368&zoomLevel=22
```

**Response Format (GeoJSON FeatureCollection):**
```json
{
  "features": [
    {
      "properties": {
        "id": "charger_001",
        "availableConnectors": 2
      },
      "geometry": {
        "coordinates": [12.520, 55.787]
      }
    },
    {
      "properties": {
        "id": "charger_002",
        "availableConnectors": 0
      },
      "geometry": {
        "coordinates": [12.525, 55.785]
      }
    }
  ]
}
```

**Filtering Logic:**
- Only chargers with `availableConnectors > 0` are retained
- Coordinates are converted from `[longitude, latitude]` to `{lat, lon}` format

**Expected Status Code:** 200

---

#### 5.3 Geoapify Static Maps API - Generate Map Image

**Purpose:** Generate a static PNG map with markers for cars and chargers.

**API Details:**
- **Protocol:** HTTPS
- **Hostname:** `maps.geoapify.com`
- **Endpoint:** `v1/staticmap`
- **HTTP Method:** GET
- **Query Parameters:**

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `style` | string | Yes | Map style | `maptiler-3d` |
| `width` | integer | Yes | Image width in pixels | 600 |
| `height` | integer | Yes | Image height in pixels | 600 |
| `center` | string | Yes | Map center (format: lonlat:lng,lat) | `lonlat:12.519,55.787` |
| `zoom` | integer | Yes | Zoom level | 14 |
| `marker` | string | Yes | Marker pins (pipe-separated) | See below |
| `apiKey` | string | Yes | API authentication token | Retrieved from SSM Parameter Store |

**Marker Format:**
Each marker follows the pattern: `lonlat:longitude,latitude;color:%23hexcolor;size:size`

- **Cars (Green):** `lonlat:12.515,55.790;color:%233ea635;size:medium`
- **Chargers (Red):** `lonlat:12.520,55.787;color:%23f30e0e;size:medium`

**Multiple Markers Example:**
```
marker=lonlat:12.515,55.790;color:%233ea635;size:medium|lonlat:12.525,55.785;color:%233ea635;size:medium|lonlat:12.520,55.787;color:%23f30e0e;size:medium
```

**Complete Request URL:**
```
https://maps.geoapify.com/v1/staticmap?style=maptiler-3d&width=600&height=600&center=lonlat:12.519,55.787&zoom=14&marker=lonlat:12.515,55.790;color:%233ea635;size:medium|lonlat:12.520,55.787;color:%23f30e0e;size:medium&apiKey=YOUR_API_KEY
```

**Response Format:**
- **Content-Type:** `image/png`
- **Body:** Binary PNG image data (Uint8Array)

**API Key Storage:**
- Retrieved from AWS SSM Parameter Store at path: `/greenmo/mapsApiToken`
- Uses AWS Lambda Powertools for secure parameter retrieval

**Expected Status Code:** 200

---

### 6. Internal Data Flow

#### 6.1 Processing Steps

```
1. Parse Input Parameters
   ├─ Extract lat1, lon1, lat2, lon2
   └─ Validate as floats

2. Calculate Center Position
   ├─ Center lat = (lat1 + lat2) / 2
   └─ Center lon = (lon1 + lon2) / 2

3. Parallel Execution (if enabled)
   ├─ Query GreenMobility API (if cars=true)
   │  ├─ Filter by desiredFuelLevel
   │  └─ Extract coordinates
   └─ Query Spirii API (if chargers=true)
      ├─ Filter by availableConnectors > 0
      └─ Extract coordinates

4. Generate Map (if results found)
   ├─ Build marker strings for each position
   └─ Call Geoapify Static Maps API

5. Transform and Return
   ├─ Encode map image to base64
   └─ Return with image/png content type
```

#### 6.2 Error Handling

- **Parse Errors:** Return 400 Bad Request
- **Network Errors:** Return 403 Forbidden
- **Unexpected Errors:** Return 500 Internal Server Error
- **No Results:** Return 200 OK with JSON message
- **Partial Results:** Continue if one external API fails gracefully

---

### 7. Type Definitions

```typescript
// Core types
type Position = {
  lat: number;
  lon: number;
};

type Params = {
  queryCars: boolean;
  queryChargers: boolean;
  desiredFuelLevel: number;
};

// GreenMobility API response item
type Car = Position & {
  id: number;
  stateOfCharge: number;
  position: {
    coordinates: [number, number]; // [lon, lat]
  };
};

// Spirii API response item
type Charger = {
  properties: {
    id: string;
    availableConnectors: number;
  };
  geometry: {
    coordinates: [number, number]; // [lon, lat]
  };
};

// Lambda response
type APIGatewayProxyResult = {
  statusCode: number;
  headers: Record<string, string>;
  body: string;
  isBase64Encoded?: boolean;
};
```

---

### 8. Environment & Configuration

**SSM Parameter Store Paths:**
- `/greenmo/mapsApiToken` - Geoapify API token

**AWS Lambda Configuration:**
- Runtime: Node.js (TypeScript)
- Timeout: Recommended 30+ seconds (due to external API calls)
- Memory: Recommended 512+ MB

---

### 9. Test Cases Guidance

Based on this specification, the following test scenarios should be covered:

1. **Input Validation**
   - Valid coordinates
   - Missing parameters
   - Invalid coordinate formats (non-numeric)
   - Out-of-range coordinates

2. **GreenMobility API Integration**
   - Successful car query with results
   - Successful car query with no results
   - API returns 401/403 errors
   - Filtering by desiredFuelLevel
   - Network timeout/connection error

3. **Spirii API Integration**
   - Successful charger query with results
   - Successful charger query with no results
   - Filtering by availableConnectors
   - API returns errors
   - Network timeout/connection error

4. **Map Generation**
   - Successful map with cars and chargers
   - Map with only cars
   - Map with only chargers
   - Geoapify API errors
   - Base64 encoding correctness

5. **Response Handling**
   - Correct HTTP status codes
   - Correct Content-Type headers
   - CORS headers presence
   - Image data integrity

6. **Edge Cases**
   - All parameters optional combinations
   - Empty result sets
   - Single result vs. multiple results
   - Boundary coordinates
   - Simultaneous API failures

---

This specification provides comprehensive details for generating tests for the next generation of the Lambda function, covering all external API interactions, parameter validation, and response formats.
