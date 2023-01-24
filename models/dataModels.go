package models

import (
	"time"
)

type DetailData struct {
	ItemData        int
	DescriptionData int
	QuantityData    int
	PriceData       int
	TaxData         int
	TotalData       int
}

type OCPageData struct {
	OCNumData        int
	DateData         time.Time
	NameData         string
	RutData          int
	RutDVData        byte
	CommBusinessData string
	AddressData      string
	TownData         string
	EmailData        string
	DetailsData      []DetailData
	NetValueData     int
	TaxData          int
	TotalData        int
}
