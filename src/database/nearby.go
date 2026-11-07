package database

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// NearbyZipcode pairs a postal record with its distance from the search origin.
type NearbyZipcode struct {
	Zipcode       Zipcode `json:"zipcode"`
	DistanceMiles float64 `json:"distance_miles"`
}

// SearchWithinRadius returns records no farther than radiusMiles from a point.
func (db *DB) SearchWithinRadius(latitude, longitude, radiusMiles float64) ([]NearbyZipcode, error) {
	if math.IsNaN(latitude) || math.IsInf(latitude, 0) || latitude < -90 || latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if math.IsNaN(longitude) || math.IsInf(longitude, 0) || longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}
	if math.IsNaN(radiusMiles) || math.IsInf(radiusMiles, 0) || radiusMiles <= 0 || radiusMiles > 100 {
		return nil, fmt.Errorf("radius must be greater than 0 and no more than 100 miles")
	}

	latitudeDelta := radiusMiles / 69.0
	cosLatitude := math.Abs(math.Cos(latitude * math.Pi / 180))
	longitudeDelta := radiusMiles / (69.0 * math.Max(cosLatitude, 0.01))
	minLatitude := math.Max(-90, latitude-latitudeDelta)
	maxLatitude := math.Min(90, latitude+latitudeDelta)
	minLongitude := longitude - longitudeDelta
	maxLongitude := longitude + longitudeDelta
	longitudeCondition := "CAST(longitude AS REAL) BETWEEN ? AND ?"
	arguments := []any{minLatitude, maxLatitude}
	switch {
	case minLongitude < -180:
		longitudeCondition = "(CAST(longitude AS REAL) >= ? OR CAST(longitude AS REAL) <= ?)"
		arguments = append(arguments, minLongitude+360, maxLongitude)
	case maxLongitude > 180:
		longitudeCondition = "(CAST(longitude AS REAL) >= ? OR CAST(longitude AS REAL) <= ?)"
		arguments = append(arguments, minLongitude, maxLongitude-360)
	default:
		arguments = append(arguments, minLongitude, maxLongitude)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes
		WHERE CAST(latitude AS REAL) BETWEEN ? AND ? AND ` + longitudeCondition
	rows, err := db.conn.QueryContext(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	results := make([]NearbyZipcode, 0)
	for rows.Next() {
		var zipcode Zipcode
		var latitudeText, longitudeText, areaCodes string
		if err := rows.Scan(&zipcode.State, &zipcode.City, &zipcode.County, &zipcode.ZipCode, &latitudeText, &longitudeText, &zipcode.StateFull, &zipcode.Timezone, &areaCodes); err != nil {
			return nil, err
		}
		zipcode.Latitude = latitudeText
		zipcode.Longitude = longitudeText
		zipcode.Zip = fmt.Sprintf("%05d", zipcode.ZipCode)
		_ = json.Unmarshal([]byte(areaCodes), &zipcode.AreaCodes)
		zipLatitude, latitudeErr := strconv.ParseFloat(latitudeText, 64)
		zipLongitude, longitudeErr := strconv.ParseFloat(longitudeText, 64)
		if latitudeErr != nil || longitudeErr != nil {
			continue
		}
		distance := haversineMiles(latitude, longitude, zipLatitude, zipLongitude)
		if distance <= radiusMiles {
			results = append(results, NearbyZipcode{
				Zipcode:       zipcode,
				DistanceMiles: math.Round(distance*100) / 100,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].DistanceMiles == results[j].DistanceMiles {
			return results[i].Zipcode.ZipCode < results[j].Zipcode.ZipCode
		}
		return results[i].DistanceMiles < results[j].DistanceMiles
	})
	return results, nil
}

func haversineMiles(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	const earthRadiusMiles = 3958.7613
	lat1 := latitude1 * math.Pi / 180
	lat2 := latitude2 * math.Pi / 180
	deltaLatitude := (latitude2 - latitude1) * math.Pi / 180
	deltaLongitude := (longitude2 - longitude1) * math.Pi / 180
	a := math.Sin(deltaLatitude/2)*math.Sin(deltaLatitude/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLongitude/2)*math.Sin(deltaLongitude/2)
	return earthRadiusMiles * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
