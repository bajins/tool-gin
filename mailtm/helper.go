package mailtm

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

type requestData struct {
	uri    string
	method string
	body   map[string]string
	bearer string
}

type responseData struct {
	code int
	body []byte
}

func generateString(length int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := 0; i < len(b); i++ {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func makeRequest(data requestData) (responseData, error) {
	var body []byte
	if data.body != nil {
		var err error
		body, err = json.Marshal(data.body)
		if err != nil {
			return responseData{}, err
		}
	}
	request, err := http.NewRequest(data.method, data.uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return responseData{}, err
	}
	if data.bearer != "" {
		request.Header.Add("Authorization", "Bearer "+data.bearer)
	}
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return responseData{}, err
	}
	if response.StatusCode == 429 {
		time.Sleep(1 * time.Second)
		return makeRequest(data)
	}
	resBody, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return responseData{}, err
	}
	return responseData{
		code: response.StatusCode,
		body: resBody,
	}, nil
}
