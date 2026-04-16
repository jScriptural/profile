package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrNoRows = errors.New("No rows found")
	Err502 = errors.New("Invalid Data")
)

type Profile struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Gender             string    `json:"gender"`
	GenderProbability  float64   `json:"gender_probability,omitempty"`
	SampleSize         int       `json:"sample_size,omitempty"`
	Age                int       `json:"age"`
	AgeGroup           string    `json:"age_group"`
	CountryID          string    `json:"country_id"`
	CountryProbability float64   `json:"country_probability,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
}

type PostData struct {
	Name string `json:"name"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Count   *int   `json:"count,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type GenderizeResponse struct {
	Name        string  `json:"name"`
	Gender      *string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
	CountryID   *string  `json:"country_id,omitempty"`
}

type AgifyResponse struct {
	Count     int    `json:"count"`
	Name      string `json:"name"`
	Age       *int    `json:"age"`
	CountryID string `json:"country_id,omitempty"`
}

type NationalizeResponse struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func (s AgifyResponse) GetAgeGroup() string {
	age := *s.Age
	switch {
	case age <= 12:
		return "child"
	case age <= 19:
		return "teenager"
	case age <= 59:
		return "adult"
	default:
		return "senior"
	}
}

func (s NationalizeResponse) GetMostProbableNationality() Country {
	a := s.Country
	mostProb := Country{}
	for _, v := range a {
		if mostProb.Probability < v.Probability {
			mostProb = v
		}
	}

	return mostProb
}
