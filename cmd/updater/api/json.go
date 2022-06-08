package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (client *CatalogClient) doJson(req any, method string, url string, resp any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, client.makeUrl(url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	response, err := client.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("got error status code %d", response.StatusCode)
	}

	if resp != nil {
		return json.Unmarshal(respBody, resp)
	}

	return nil
}

func (client *CatalogClient) makeUrl(url string) string {
	return strings.TrimRight(client.baseUrl, "/") +
		"/" +
		strings.TrimLeft(url, "/")
}
