package data

import "testing"

func TestReadEmbeddedZipcodeData(t *testing.T) {
	content, err := ReadFile("zipcodes.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal("embedded dataset is empty")
	}
	if _, err := ReadFile("missing.json"); err == nil {
		t.Fatal("ReadFile accepted missing data")
	}
}
