package mailtm

import (
	"encoding/json"
	"time"
)

type Domain struct {
	Id        string    `json:"id"`
	Domain    string    `json:"domain"`
	IsActive  bool      `json:"isActive"`
	IsPrivate bool      `json:"isPrivate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func AvailableDomains() ([]Domain, error) {
	request := requestData{
		uri:    URI_DOMAINS,
		method: "GET",
	}
	response, err := makeRequest(request)
	if err != nil {
		return nil, err
	}
	if response.code != 200 {
		return nil, err
	}
	data := map[string][]Domain{}
	json.Unmarshal(response.body, &data)
	return data["hydra:member"], nil
}
