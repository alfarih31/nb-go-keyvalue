package keyvalue_test

import (
	"fmt"
	"github.com/alfarih31/nb-go-keyvalue"
	"log"
	"testing"
)

type Account struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func TestFromStruct(t *testing.T) {
	p1 := Account{
		Name:        "John Doe",
		Email:       "johndoe@mail.com",
		PhoneNumber: "+1234567890",
	}
	p1kv, err := keyvalue.FromStruct(p1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("(current) p1kv:\n\t%v\n", p1kv.JSON())

	if p1kv["name"] != "John Doe" {
		log.Fatalf(`keyvalue.FromStruct(p1), want convert p1 to map[string]interface{}`)
	}
}

func TestKeyValue_AssignTo(t *testing.T) {
	p1 := Account{
		Name:        "John Doe",
		Email:       "johndoe@mail.com",
		PhoneNumber: "+1234567890",
	}

	p2 := Account{
		Name:        "Foo Bar",
		Email:       "foobar@mail.com",
		PhoneNumber: "",
	}

	p1kv, err := keyvalue.FromStruct(p1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("(current) p1kv:\n\t%v\n", p1kv.JSON())

	p2kv, err := keyvalue.FromStruct(p2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("(current) p2kv:\n\t%v\n", p2kv.JSON())

	// Assigning p1 data to p2
	p1kv.AssignTo(p2kv)
	if p2kv["phone_number"] != "+1234567890" {
		log.Fatalf(`p1kv.AssignTo(p2kv), want match for p2kv["phone_number""]=p1kv["phone_number"]`)
	}
	fmt.Printf("(now) p2kv:\n\t%v\n", p2kv.JSON())
}

func TestKeyValue_Assign(t *testing.T) {
	p1 := Account{
		Name:        "John Doe",
		Email:       "johndoe@mail.com",
		PhoneNumber: "+1234567890",
	}

	p2 := Account{
		Name:        "Foo Bar",
		Email:       "foobar@mail.com",
		PhoneNumber: "",
	}

	p1kv, err := keyvalue.FromStruct(p1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("(current) p1kv:\n\t%v\n", p1kv.JSON())

	p2kv, err := keyvalue.FromStruct(p2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("(current) p2kv:\n\t%v\n", p2kv.JSON())

	// Assigning p2 data with p1 data
	p2kv.Assign(p1kv)
	if p2kv["phone_number"] != "+1234567890" {
		log.Fatalf(`p2kv.Assign(p1kv), want match for p1kv["phone_number""]=p2kv["phone_number"]`)
	}
	fmt.Printf("(now) p2kv:\n\t%v\n", p2kv.JSON())
}

func TestKeyValue_Unmarshal(t *testing.T) {
	p1 := Account{
		Name:        "John Doe",
		Email:       "johndoe@mail.com",
		PhoneNumber: "+1234567890",
	}
	fmt.Printf("(current) p1:\n\t%v\n", p1)

	p1kv, err := keyvalue.FromStruct(p1)
	if err != nil {
		log.Fatal(err)
	}

	p1kv["name"] = "Foo Bar"

	// Unmarshal p1kv back to p1
	err = p1kv.Unmarshal(&p1)
	if err != nil {
		log.Fatal(err)
	}

	if p1.Name != "Foo Bar" {
		log.Fatalf(`p1kv.Unmarshal(), want p1kv value to p1`)
	}
	fmt.Printf("(now) p1:\n\t%v\n", p1)
}
