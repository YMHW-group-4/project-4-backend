package database

import (
	"io"
	"net/http"
)

const supabaseUrl = "https://hvttiajsasltwhgfhdnc.supabase.co"
const supabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Imh2dHRpYWpzYXNsdHdoZ2ZoZG5jIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NjM1NTMsImV4cCI6MTk4NTAzOTU1M30.kbumzE2DTKWYCKbmfz1KSPdnqSskPC0K9Je77Ql80qE"
const supabaseBearer = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Imh2dHRpYWpzYXNsdHdoZ2ZoZG5jIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NjM1NTMsImV4cCI6MTk4NTAzOTU1M30.kbumzE2DTKWYCKbmfz1KSPdnqSskPC0K9Je77Ql80qE"

func Get(url string, filters map[string]string) []byte {
	client := &http.Client{}
	req, _ := createRequest(supabaseUrl+url, "GET", nil)

	if filters != nil {
		for filter, value := range filters {
			req.Header.Set(filter, value)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return errorHandler(err)
	}
	return readBody(response.Body)
}

func errorHandler(err error) []byte {
	return nil
}

func readBody(responseBody io.ReadCloser) []byte {
	body, resErr := io.ReadAll(responseBody)
	bodyErr := responseBody.Close()
	if bodyErr != nil || resErr != nil {
		return nil
	}
	return body
}

func POST(url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, _ := createRequest(supabaseUrl+url, "POST", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	return client.Do(req)
}

func createRequest(url string, method string, body io.Reader) (*http.Request, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", supabaseBearer)
	return req, nil
}
