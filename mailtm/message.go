package mailtm

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type Message struct {
	ID          string       `json:"id"`
	From        Adressee     `json:"from"`
	To          []Adressee   `json:"to"`
	Subject     string       `json:"subject"`
	Intro       string       `json:"intro"`
	Seen        bool         `json:"seen"`
	IsDeleted   bool         `json:"isDeleted"`
	Size        int          `json:"size"`
	Text        string       `json:"text"`
	HTML        []string     `json:"html"`
	Attachments []Attachment `json:"attachments"`
}

type Adressee struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Attachment struct {
	Id               string `json:"id"`
	Filename         string `json:"filename"`
	ContentType      string `json:"contentType"`
	Disposition      string `json:"disposition"`
	TransferEncoding string `json:"transferEncoding"`
	Related          bool   `json:"related"`
	Size             int    `json:"size"`
	DownloadURL      string `json:"downloadUrl"`
}

func (account *Account) MessagesAt(page int) ([]Message, error) {
	var data map[string][]Message
	URI := URI_MESSAGES + "?page=" + strconv.Itoa(page)
	request := requestData{
		uri:    URI,
		method: "GET",
		bearer: account.bearer,
	}
	response, err := makeRequest(request)
	if err != nil {
		return nil, err
	}
	if response.code != 200 {
		return nil, errors.New("failed to get messages")
	}
	json.Unmarshal(response.body, &data)
	messages := data["hydra:member"]
	for i, _ := range messages {
		msg, err := account.MessageById(messages[i].ID)
		if err != nil {
			return nil, err
		}
		messages[i] = msg
	}
	return messages, nil
}

func (account *Account) MessageById(id string) (Message, error) {
	var msg Message
	URI := URI_MESSAGES + "/" + id
	request := requestData{
		uri:    URI,
		method: "GET",
		bearer: account.bearer,
	}
	response, err := makeRequest(request)
	if err != nil {
		return Message{}, err
	}
	if response.code != 200 {
		return Message{}, errors.New("failed to get message")
	}
	json.Unmarshal(response.body, &msg)
	return msg, nil
}

func (account *Account) MessagesChan(ctx context.Context) <-chan Message {
	msgChan := make(chan Message)
	go func() {
		lastMsg, _ := account.LastMessage()
		lastMsgId := lastMsg.ID
	loop:
		for {
			select {
			case <-ctx.Done():
				close(msgChan)
				break loop
			default:
			}
			msg, err := account.LastMessage()
			if err != nil {
				continue
			}
			if msg.ID != lastMsgId {
				msgChan <- msg
				lastMsgId = msg.ID
			}
		}
	}()
	return msgChan
}

func (account *Account) LastMessage() (Message, error) {
	msgs, err := account.MessagesAt(1)
	if err != nil {
		return Message{}, err
	}
	if len(msgs) == 0 {
		return Message{}, errors.New("no messages")
	}
	return msgs[0], nil
}

func (account *Account) DeleteMessage(id string) error {
	URI := URI_MESSAGES + "/" + id
	request := requestData{
		uri:    URI,
		method: "DELETE",
		bearer: account.bearer,
	}
	response, err := makeRequest(request)
	if err != nil {
		return err
	}
	if response.code == 404 {
		return errors.New("message with id " + id + " was not found")
	}
	if response.code != 204 {
		return errors.New("failed to delete message")
	}
	return nil
}

func (account *Account) MarkMessage(id string) error {
	URI := URI_MESSAGES + "/" + id
	request := requestData{
		uri:    URI,
		method: "PATCH",
		bearer: account.bearer,
	}
	response, err := makeRequest(request)
	if err != nil {
		return err
	}
	if response.code == 404 {
		return errors.New("message with id " + id + " was not found")
	}
	if response.code != 200 {
		return errors.New("failed to mark message")
	}
	return nil
}
