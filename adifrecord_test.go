package adifparser

import (
	"testing"
)

func TestGetNextField(t *testing.T) {
	buf := []byte("<blah:2>AB<FOO:3>XYZ <bar:4:s>1234")

	expected := []struct {
		n string
		v string
	}{
		{"blah", "AB"},
		{"foo", "XYZ"},
		{"bar", "1234"},
	}

	var err error
	var data *fieldData

	for _, el := range expected {
		data, buf, err = getNextField(buf)
		if err != nil {
			t.Fatal(err)
		}
		if data.name != el.n || data.value != el.v {
			t.Fatalf("Got %q=%q, expected %q=%q.", data.name, data.value, el.n, el.v)
		}
	}
}

func TestParseADIFRecord(t *testing.T) {
	testData := "<call:4>W1AW<STATION_CALL:6>KF4MDV"
	record, err := ParseADIFRecord([]byte(testData))
	if err != nil {
		t.Fatal(err)
	}
	if n, ok := record.values["call"]; ok {
		if n != "W1AW" {
			t.Fatalf("CALL: %q != %q", n, "W1AW")
		}
	} else {
		t.Fatal("No 'call' value.")
	}
	if n, ok := record.values["station_call"]; ok {
		if n != "KF4MDV" {
			t.Fatalf("STATION_CALL: %q != %q", n, "KF4MDV")
		}
	} else {
		t.Fatal("No 'station_call' value.")
	}
}

func TestGetFields(t *testing.T) {
	testData := "<call:4>W1AW<STATION_CALL:6>KF4MDV"
	expected := [2]string{"call", "station_call"}

	record, err := ParseADIFRecord([]byte(testData))
	if err != nil {
		t.Fatal(err)
	}
	fieldNames := record.GetFields()

	if len(fieldNames) != len(expected) {
		t.Fatalf("Expected %d fields but got %d", len(expected), len(fieldNames))
	}

OUTER:
	for _, exp := range expected {
		for _, field := range fieldNames {
			if exp == field {
				continue OUTER
			}
		}
		t.Fatalf("Expected field %v wasn't in the actual fields", exp)
	}
}

func TestSetValue(t *testing.T) {
	testData := map[string]string{"call": "W1AW", "STATION_CALL": "KF4MDV"}
	expected := [2]string{"call", "station_call"}

	record := NewADIFRecord()
	for k, v := range testData {
		record.SetValue(k, v)
	}
	fieldNames := record.GetFields()

	if len(fieldNames) != len(expected) {
		t.Fatalf("Expected %d fields but got %d", len(expected), len(fieldNames))
	}

OUTER:
	for _, exp := range expected {
		for _, field := range fieldNames {
			if exp == field {
				continue OUTER
			}
		}
		t.Fatalf("Expected field %v wasn't in the actual fields", exp)
	}
}
