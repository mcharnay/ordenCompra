package models

import "html/template"

type Detail struct {
	Item        string
	Description string
	Quantity    string
	Price       string
	Tax         string
	Total       string
}

type PageData struct {
	Logo         template.URL
	OCNum        string
	Date         string
	Name         string
	Rut          string
	CommBusiness string
	Address      string
	Town         string
	Email        string
	Details      []Detail
	NetValue     string
	Tax          string
	Total        string
}
