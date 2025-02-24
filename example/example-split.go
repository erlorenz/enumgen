package example

//go:generate go run ../ -split

var shippingMethodValues = map[string]string{
	"":        "Unknown",
	"ups_GND": "UPS Ground",
	"ups_1DA": "UPS Next Day Air",
	"ups_2DA": "UPS 2 Day Air",
}

var paymentMethodValues = map[string]string{
	"":            "Unknown",
	"card":        "Credit Card",
	"netterms":    "Net Terms",
	"ach_payment": "ACH",
	"free":        "Free",
}
