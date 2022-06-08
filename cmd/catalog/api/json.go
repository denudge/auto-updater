package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Validatable interface {
	Validate() error
}

// parseAndValidateJsonRequest unmarshalls the request body into v. Errors are written to the HTTP response
func (api *Api) parseAndValidateJsonRequest(w http.ResponseWriter, r *http.Request, v Validatable) error {
	var body []byte
	_, err := r.Body.Read(body)
	if err != nil {
		err := fmt.Errorf("error reading request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		err = fmt.Errorf("error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err = v.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return err
	}

	return nil
}

func (api *Api) writeJsonResponse(w http.ResponseWriter, v any) error {
	responseBody, err := json.Marshal(v)
	if err != nil {
		log.Println(fmt.Printf("Error marshalling JSON response of: %#v\n", v))
		http.Error(w, "error serializing JSON response", http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(responseBody)
	if err != nil {
		log.Println(fmt.Printf("Error writing JSON response of: %#v\n", v))
		http.Error(w, "error serializing JSON response", http.StatusInternalServerError)
		return err
	}

	return nil
}

// parseAndValidateJsonPostRequest unmarshalls the POST request body into v. Errors are written to the HTTP response
func (api *Api) parseAndValidateJsonPostRequest(w http.ResponseWriter, r *http.Request, v Validatable) error {
	if err := api.validateMethodIs(w, r, http.MethodPost); err != nil {
		// Errors are already written to the response
		return err
	}

	return api.parseAndValidateJsonRequest(w, r, v)
}
