package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (app *Application) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	log.Println(user.Email, user.Password, valid)

	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}
	err = app.LogRequest("authentication", fmt.Sprintf("%s logged in", requestPayload.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged as user %s", user.Email),
		Data:    user,
	}
	app.wirteJSON(w, http.StatusAccepted, payload)

}

func (app *Application) LogRequest(name, data string) error {
	type logData struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry := logData{
		Name: name,
		Data: data,
	}
	byteData, _ := json.MarshalIndent(entry, "", "\t")
	log.Println(string(byteData))
	request, err := http.NewRequest("POST", "http://logger-service/log", strings.NewReader(string(byteData)))

	if err != nil {
		log.Println(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
