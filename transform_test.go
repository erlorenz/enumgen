package main

import (
	"maps"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPascalCase(t *testing.T) {
	table := map[string]string{
		"UPS Ground Shipping": "UPSGroundShipping",
		"Free Shipping":       "FreeShipping",
		"freeShipping":        "FreeShipping",
		"Free_Shipping":       "FreeShipping",
		"Free-Shipping":       "FreeShipping",
	}

	for name, want := range table {
		t.Run(name, func(t *testing.T) {
			got := toPascalCase(name)
			if want != got {
				t.Fatalf("wanted %s, got %s", want, got)
			}
		})
	}
}

func TestMapToEnumData(t *testing.T) {
	name := "shippingMethodValues"
	vals := map[string]string{
		"ups_GND": "UPS Ground",
		"ups_1DA": "UPS Next Day Air",
		"ups_2DA": "UPS 2 Day Air",
	}

	ed := mapToEnumData(MapParseResult{name, vals})
	want := EnumData{Type: "ShippingMethod", Value: "shippingmethodvalue", Items: map[string]string{
		"ups_GND": "UPSGround",
		"ups_1DA": "UPSNextDayAir",
		"ups_2DA": "UPS2DayAir",
	}}

	if want, got := want.Type, ed.Type; want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
	if want, got := want.Value, ed.Value; want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
	if want, got := want.Items, ed.Items; !maps.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}
