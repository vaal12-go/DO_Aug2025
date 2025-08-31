package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

type Cat struct {
	// gorm.Model
	ID                uint       `gorm:"primary_key"`
	CreatedAt         time.Time  `json:"-"`
	UpdatedAt         time.Time  `json:"-"`
	DeletedAt         *time.Time `json:"-"`
	Name              string     `json:"name"`
	YearsOfExperience uint32     `json:"years_of_experience"`
	Breed             string     `json:"breed" binding:"required,breedvalidator"`
	Salary            uint32     `json:"salary"`
	Mission           *Mission
}

const BREED_VALIDATION_URL = "https://api.thecatapi.com/v1/breeds"

func BreedValidator(fl validator.FieldLevel) bool {
	// fmt.Println("cat:17::", "hello from validator")
	// requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	res, err := http.Get(BREED_VALIDATION_URL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	// fmt.Printf("client: response body: %s\n", resBody)
	var breeds []any
	err = json.Unmarshal(resBody, &breeds)
	if err != nil {
		panic(err)
	}
	fmt.Println("Length of breeds:", len(breeds))
	breed_found := false
	breed_field_val := fl.Field().String()
	fmt.Println("cat:51 breed_field_val::", breed_field_val)
	for _, rec := range breeds {
		breed := rec.(map[string]any)
		// fmt.Println("cat:50 breed[name]::", breed["name"])
		if breed_field_val == breed["name"] {
			breed_found = true
		}
	}
	return breed_found
}
