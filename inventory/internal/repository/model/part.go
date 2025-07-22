package model

import "time"

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value interface {
	isValue()
}

type StringValue struct {
	StringValue string
}
type Int64Value struct {
	Int64Value int64
}
type DoubleValue struct {
	DoubleValue float64
}
type BoolValue struct {
	BoolValue bool
}

func (v StringValue) isValue() {}
func (v Int64Value) isValue()  {}
func (v DoubleValue) isValue() {}
func (v BoolValue) isValue()   {}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
