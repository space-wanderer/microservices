package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UUID          string             `bson:"uuid" json:"uuid"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	Price         float64            `bson:"price" json:"price"`
	StockQuantity int64              `bson:"stock_quantity" json:"stock_quantity"`
	Category      Category           `bson:"category" json:"category"`
	Dimensions    *Dimensions        `bson:"dimensions" json:"dimensions"`
	Manufacturer  *Manufacturer      `bson:"manufacturer" json:"manufacturer"`
	Tags          []string           `bson:"tags" json:"tags"`
	Metadata      map[string]*Value  `bson:"metadata" json:"metadata"`
	CreatedAt     primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt     primitive.DateTime `bson:"updated_at" json:"updated_at"`
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
