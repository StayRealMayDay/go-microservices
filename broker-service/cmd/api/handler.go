package main

import (
	"borker/event"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"strings"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Application) Borker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.wirteJSON(w, http.StatusOK, payload)
}

func (app *Application) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.LogItemVisRPC(w, requestPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknow action"))
	}
}

func (app *Application) logEventViaRabbit(w http.ResponseWriter, payload LogPayload) {
	err := app.pushToQueue(payload.Name, payload.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Logged via RabbitMQ"
	app.wirteJSON(w, http.StatusAccepted, responsePayload)
}

func (app *Application) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	// payloadBytes := new(bytes.Buffer)
	// _ = json.NewEncoder(payloadBytes).Encode(payload)
	j, _ := json.MarshalIndent(&payload, "", "\n")

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest("POST", "http://logger-service/log", strings.NewReader(string(jsonData)))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling the service"))
		return
	}
	var jsonResponseFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonResponseFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonResponseFromService.Error {
		app.errorJSON(w, errors.New(jsonResponseFromService.Message), http.StatusBadRequest)
		return
	}

	var responseData = jsonResponse{
		Error:   false,
		Data:    jsonResponseFromService.Data,
		Message: "Logged succeed",
	}
	app.wirteJSON(w, http.StatusAccepted, responseData)
}

func (app *Application) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling the service"))
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload = jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	app.wirteJSON(w, http.StatusAccepted, payload)
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Application) LogItemVisRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	rpyPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	var result string
	err = client.Call("RPCServer.LogInfo", rpyPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	jsonResponse := jsonResponse{
		Error:   false,
		Message: result,
	}
	app.wirteJSON(w, http.StatusAccepted, jsonResponse)
}
