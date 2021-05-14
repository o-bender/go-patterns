package main

import "fmt"

type DocumentProperties map[string]interface{}

type Document interface {
	Put(key string, value interface{})
	Get(key string) interface{}
}

type HasType interface {
	GetType() string
}

func GetType(document Document) string {
	return document.Get("type").(string)
}

type HasPrice interface {
	GetPrice() float64
}

func GetPrice(document Document) float64 {
	return document.Get("price").(float64)
}

type AbstractDocument struct {
	Document
	properties DocumentProperties
}

func (ad *AbstractDocument) Put(key string, value interface{}) {
	ad.properties[key] = value
}

func (ad AbstractDocument) Get(key string) interface{} {
	return ad.properties[key]
}

func NewAbstractDocument(properties DocumentProperties) AbstractDocument {
	return AbstractDocument{
		properties: properties,
	}
}

type Car struct {
	AbstractDocument
	HasType
	HasPrice
}

func (c Car) GetPrice() float64 {
	return GetPrice(&c.AbstractDocument)
}

func (c Car) GetType() string {
	return GetType(&c.AbstractDocument)
}

func NewCar(properties DocumentProperties) *Car {
	return &Car{
		AbstractDocument: AbstractDocument{
			properties: properties,
		},
	}
}

func main() {
	doc := NewAbstractDocument(DocumentProperties{
		"test": "best",
	})
	fmt.Println(doc.Get("test"))

	car := NewCar(DocumentProperties{
		"type":  "Москвич",
		"price": 3.0,
	})
	car.Put("price", 1000.0)
	fmt.Println(car.Get("price"))
	fmt.Println(car.GetPrice())
	fmt.Println(car.GetType())
}
