package database

import (
	"math"
	"strings"
	"testing"
)

func testDB(t *testing.T) *DB {
	t.Helper()
	db, err := Initialize(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}

func TestZipcodeDatabaseQueries(t *testing.T) {
	db := testDB(t)
	if err := db.LoadFromJSON([]byte("{")); err == nil {
		t.Fatal("LoadFromJSON accepted malformed JSON")
	}
	data := []byte(`[
		{"state":"CA","city":"Beverly Hills","county":"Los Angeles","zip_code":90210,"latitude":34.0901,"longitude":-118.4065},
		{"state":"CA","city":"Los Angeles","county":"Los Angeles","zip_code":90001,"latitude":"33.9731","longitude":"-118.2479"},
		{"state":"NY","city":"New York","county":"New York","zip_code":10001,"latitude":40.7506,"longitude":-73.9972}
	]`)
	if err := db.LoadFromJSON(data); err != nil {
		t.Fatal(err)
	}
	if err := db.LoadFromJSON(data); err != nil {
		t.Fatalf("repeat load should skip unchanged data: %v", err)
	}

	zip, err := db.SearchByZipCode(90210)
	if err != nil || zip == nil || zip.City != "Beverly Hills" || zip.GetLatitudeString() != "34.090100" || zip.GetLongitudeString() != "-118.406500" {
		t.Fatalf("SearchByZipCode = %#v, %v", zip, err)
	}
	if zip, err := db.SearchByZipCode(12345); err != nil || zip != nil {
		t.Fatalf("missing zip = %#v, %v", zip, err)
	}
	if zip, err := db.SearchByCity("los angeles"); err != nil || len(zip) != 1 {
		t.Fatalf("SearchByCity = %#v, %v", zip, err)
	}
	if zip, err := db.SearchByState("ca"); err != nil || len(zip) != 2 {
		t.Fatalf("SearchByState = %#v, %v", zip, err)
	}
	if zip, err := db.SearchByStateAndCity("ca", "LOS ANGELES"); err != nil || len(zip) != 1 {
		t.Fatalf("SearchByStateAndCity = %#v, %v", zip, err)
	}
	if zip, err := db.SearchByPrefix("90"); err != nil || len(zip) != 2 {
		t.Fatalf("SearchByPrefix = %#v, %v", zip, err)
	}
	if suggestions, err := db.AutoComplete(" ", 0); err != nil || len(suggestions) != 0 {
		t.Fatalf("empty autocomplete = %#v, %v", suggestions, err)
	}
	if suggestions, err := db.AutoComplete("los", 0); err != nil || len(suggestions) != 1 || suggestions[0] != "Los Angeles, CA" {
		t.Fatalf("city autocomplete = %#v, %v", suggestions, err)
	}
	if suggestions, err := db.AutoComplete("NY", 10); err != nil || len(suggestions) != 1 || suggestions[0] != "New York, NY" {
		t.Fatalf("state autocomplete = %#v, %v", suggestions, err)
	}
	stats, err := db.GetStats()
	if err != nil || stats["total_zipcodes"] != 3 || stats["total_states"] != 2 || stats["total_cities"] != 3 {
		t.Fatalf("stats = %#v, %v", stats, err)
	}
	if zip, err := db.GetRandom(); err != nil || zip == nil {
		t.Fatalf("random zip = %#v, %v", zip, err)
	}
}

func TestZipcodeInsertRollbackAndCoordinates(t *testing.T) {
	db := testDB(t)
	duplicate := []byte(`[{"state":"CA","city":"A","zip_code":90210},{"state":"CA","city":"B","zip_code":90210}]`)
	if err := db.LoadFromJSON(duplicate); err == nil || !strings.Contains(err.Error(), "index 1") {
		t.Fatalf("duplicate insert error = %v", err)
	}
	stats, err := db.GetStats()
	if err != nil || stats["total_zipcodes"] != 0 {
		t.Fatalf("failed load was not rolled back: %#v, %v", stats, err)
	}

	if got := (&Zipcode{}).GetLatitudeString(); got != "" {
		t.Fatalf("nil latitude = %q", got)
	}
	if got := (&Zipcode{Latitude: "raw"}).GetLatitudeString(); got != "raw" {
		t.Fatalf("string latitude = %q", got)
	}
	if got := (&Zipcode{Latitude: 12.5}).GetLatitudeString(); got != "12.500000" {
		t.Fatalf("float latitude = %q", got)
	}
	if got := (&Zipcode{Latitude: 7}).GetLatitudeString(); got != "7" {
		t.Fatalf("integer latitude = %q", got)
	}
	if got := (&Zipcode{}).GetLongitudeString(); got != "" {
		t.Fatalf("nil longitude = %q", got)
	}
	if got := (&Zipcode{Longitude: "raw"}).GetLongitudeString(); got != "raw" {
		t.Fatalf("string longitude = %q", got)
	}
	if got := (&Zipcode{Longitude: 12.5}).GetLongitudeString(); got != "12.500000" {
		t.Fatalf("float longitude = %q", got)
	}
	if got := (&Zipcode{Longitude: 7}).GetLongitudeString(); got != "7" {
		t.Fatalf("integer longitude = %q", got)
	}
}

func TestSearchWithinRadiusValidatesAndSortsResults(t *testing.T) {
	db := testDB(t)
	data := []byte(`[{
		"state":"CA","city":"Beverly Hills","zip_code":90210,"latitude":34.0901,"longitude":-118.4065
	},{
		"state":"CA","city":"Los Angeles","zip_code":90001,"latitude":33.9731,"longitude":-118.2479
	},{
		"state":"AK","city":"West Point","zip_code":99501,"latitude":0,"longitude":179.9
	},{
		"state":"AK","city":"East Point","zip_code":99502,"latitude":0,"longitude":-179.9
	}]`)
	if err := db.LoadFromJSON(data); err != nil {
		t.Fatal(err)
	}
	results, err := db.SearchWithinRadius(34.09, -118.40, 10)
	if err != nil || len(results) != 1 || results[0].Zipcode.ZipCode != 90210 {
		t.Fatalf("nearby ZIP results = %#v, %v", results, err)
	}
	wrapped, err := db.SearchWithinRadius(0, 179.9, 20)
	if err != nil || len(wrapped) != 2 {
		t.Fatalf("dateline results = %#v, %v", wrapped, err)
	}
	for _, tc := range []struct{ lat, lon, radius float64 }{
		{91, 0, 1}, {0, 181, 1}, {0, 0, 0}, {math.NaN(), 0, 1}, {0, math.Inf(1), 1}, {0, 0, math.Inf(1)}, {0, 0, 101},
	} {
		if _, err := db.SearchWithinRadius(tc.lat, tc.lon, tc.radius); err == nil {
			t.Errorf("SearchWithinRadius(%v, %v, %v) accepted invalid input", tc.lat, tc.lon, tc.radius)
		}
	}
}
