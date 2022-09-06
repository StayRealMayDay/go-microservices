package main

import (
	"log-service/data"
	"net/http"
	"tools/jsonkit"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Application) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JsonPayload
	_ = jsonkit.ReadJSON(w, r, &requestPayload)

	event := data.LogoEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err := app.Models.LogoEntry.Insert(event)
	if err != nil {
		jsonkit.ErrorJSON(w, err)
	}
	resp := jsonkit.JsonResponse{
		Error:   false,
		Message: "logged",
	}
	jsonkit.WirteJSON(w, http.StatusAccepted, resp)
}
