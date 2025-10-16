package mailtm

import (
	"encoding/json"
	"errors"
)

type Properties map[string]any

type Account struct {
	address    string
	password   string
	bearer     string
	properties Properties
}

func (account *Account) Address() string {
	return account.address
}

func (account *Account) Password() string {
	return account.password
}

func (account *Account) Bearer() string {
	return account.bearer
}

func (account *Account) Property(name string) (p any, ok bool) {
	property, exists := account.properties[name]
	return property, exists
}

func (account *Account) Delete() error {
	URI := URI_ACCOUNTS + "/" + account.properties["id"].(string)
	request := requestData{
		uri:    URI,
		method: "DELETE",
		bearer: account.bearer,
	}
	response, err := makeRequest(request)
	if err != nil {
		return err
	}
	if response.code != 204 {
		return errors.New("failed to delete account")
	}
	return nil
}

type Options struct {
	Domain   string
	Username string
	Password string
}

func NewAccount() (*Account, error) {
	domains, err := AvailableDomains()
	if err != nil {
		return nil, err
	}
	if len(domains) == 0 {
		return nil, errors.New("no domains available")
	}
	return NewAccountWithOptions(Options{
		Domain:   domains[0].Domain,
		Username: generateString(16),
		Password: generateString(16),
	})
}

func Login(address string, password string) (*Account, error) {
	id, token, err := GetIdAndToken(address, password)
	if err != nil {
		return nil, err
	}
	account, err := LoginWithIdAndToken(id, token)
	if err != nil {
		return nil, err
	}
	account.password = password
	if violations, ok := account.Property("violations"); ok {
		return nil, errors.New(violations.([]any)[0].(map[string]any)["message"].(string))
	}
	return account, nil
}

func LoginWithToken(token string) (*Account, error) {
	account := new(Account)
	request := requestData{
		uri:    URI_ME,
		method: "GET",
		bearer: token,
	}
	response, err := makeRequest(request)
	if err != nil {
		return nil, err
	}
	if response.code != 200 {
		return nil, errors.New("failed to get account")
	}
	json.Unmarshal(response.body, &account.properties)
	account.address = account.properties["address"].(string)
	account.bearer = token
	return account, nil
}

func GetIdAndToken(address string, password string) (string, string, error) {
	data := map[string]string{
		"address":  address,
		"password": password,
	}
	body := make(map[string]any)
	request := requestData{
		uri:    URI_TOKEN,
		method: "POST",
		body:   data,
	}
	response, err := makeRequest(request)
	if err != nil {
		return "", "", err
	}
	if response.code != 200 {
		return "", "", errors.New("failed to get id and token")
	}
	err = json.Unmarshal(response.body, &body)
	if err != nil {
		return "", "", err
	}
	return body["id"].(string), body["token"].(string), nil
}

func LoginWithIdAndToken(id string, token string) (*Account, error) {
	account := new(Account)
	uri := URI_ACCOUNTS + "/" + id
	request := requestData{
		uri:    uri,
		method: "GET",
		bearer: token,
	}
	response, err := makeRequest(request)
	if err != nil {
		return nil, err
	}
	if response.code != 200 {
		return nil, errors.New("failed to get account")
	}
	json.Unmarshal(response.body, &account.properties)
	account.address = account.properties["address"].(string)
	account.bearer = token
	return account, nil
}

func NewAccountWithOptions(options Options) (*Account, error) {
	address := options.Username + "@" + options.Domain
	password := options.Password
	data := map[string]string{
		"address":  address,
		"password": password,
	}
	request := requestData{
		uri:    URI_ACCOUNTS,
		method: "POST",
		body:   data,
	}
	response, err := makeRequest(request)
	if err != nil {
		return nil, err
	}
	if response.code != 201 {
		return nil, errors.New("failed to create an account")
	}
	account, err := Login(address, password)
	return account, err
}
