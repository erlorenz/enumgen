package example

//go:generate go run ../

var shippingMethodValues = map[string]string{
	"":        "Unknown",
	"ups_GND": "UPS Ground",
	"ups_1DA": "UPS Next Day Air",
	"ups_2DA": "UPS 2 Day Air",
}
