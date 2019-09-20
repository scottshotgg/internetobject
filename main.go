package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

/*
name, age:{int, min:20}, address: {street, city, state}, active?:bool, tags?:[string]
---
Spiderman, 25, {Bond Street, New York, NY}, T, [agile, swift]
*/

const (
	schema = `age: { int, min: 20 }`
	// ,
	// address: { street, city, state },
	// active?: bool,
	// tags?: [ string ]`

	payload = `something_here, 23, T`
	// { Bond Street, New York, NY },
	// T,
	// [ agile, swift ]`
)

type IO_Address struct {
	Street string
	City   string
	State  string
}

type IO_Test struct {
	Name string
	Age  int
	// Address IO_Address
	Active *bool
	Tags   *[]string
}

var index int

func parseInt(f reflect.Value) {
	fmt.Println("its me lil int boi")

	var number []rune

	// TODO: make something that will test the current token index vs the length
	for _, char := range payload[index:] {
		// fmt.Println("char - int", string(char))

		if char == ',' {
			index++
			break
		}

		if unicode.IsSpace(char) {
			index++
			continue
		}

		if unicode.IsDigit(char) {
			number = append(number, char)
		} else {
			log.Fatal("wtf not a digit", string(char))
		}

		index++
	}

	i, err := strconv.ParseInt(string(number), 10, 64)
	if err != nil {
		log.Fatal("fucking err", err)
	}

	f.SetInt(i)
}

func parseString(f reflect.Value) {
	fmt.Println("lil stringy gurl")

	var text string

	for _, char := range payload[index:] {
		// fmt.Println("char - string", string(char))

		if char == ',' {
			index++
			break
		}

		if unicode.IsSpace(char) {
			index++
			continue
		}

		text += string(char)
		index++
	}

	f.SetString(text)
}

// TODO: nest nest ze parser to piecewise parse ze data
func parseStruct() {}

func parseBool(f reflect.Value) {
	fmt.Println("booley p gang")

	var booley bool

	for _, char := range payload[index:] {
		// fmt.Println("char - bool", string(char))

		if char == ',' {
			index++
			break
		}

		if unicode.IsSpace(char) {
			index++
			continue
		}

		if char == 'T' {
			booley = true
		} else if char == 'F' {
			booley = false
		} else {
			log.Fatal("WTFF not a bool char", string(char))
		}

		index++
	}

	f.Set(reflect.New(f.Type().Elem()))
	f.Elem().SetBool(booley)
}

func listStruct(test interface{}) {
	fmt.Println(reflect.TypeOf(test))

	var elem = reflect.ValueOf(test).Elem()

	for i := 0; i < elem.NumField(); i++ {
		var field = elem.Type().Field(i)
		// just assume the tags are lower case for now
		var tag = strings.ToLower(field.Name)

		// elem.Field(i).Set()

		switch field.Type.Kind() {
		case reflect.Int:
			parseInt(elem.FieldByName(field.Name))

		case reflect.String:
			parseString(elem.FieldByName(field.Name))

		case reflect.Ptr:
			// fmt.Println("me", elem.FieldByName(field.Name).Type().Elem().Kind())

			switch elem.FieldByName(field.Name).Type().Elem().Kind() {
			case reflect.Bool:
				parseBool(elem.FieldByName(field.Name))
			}

			// parseBoolP()
		}

		fmt.Println("tag:", tag)
	}
}

func main() {
	fmt.Println("hi")
	var test IO_Test

	listStruct(&test)

	fmt.Println("test:", *test.Active)
	bytes, err := json.Marshal(test)
	if err != nil {
		log.Fatal("somthing happene", bytes)
	}

	fmt.Println(string(bytes))
}
