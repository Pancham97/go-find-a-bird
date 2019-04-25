package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestCreateBirdHandler(t *testing.T) {
	birds = []Bird{
		{"sparrow", "A small harmless bird"},
	}

	form := newCreateBirdForm()
	request, error := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if error != nil {
		t.Fatal(error)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(createBirdHandler)
	hf.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code. Expected %v, got %v", http.StatusOK, status)
	}

	expected := Bird{"eagle", "A bird of prey"}

	if error != nil {
		t.Fatal(error)
	}

	actual := birds[1]

	if actual != expected {
		t.Errorf("handler returned unexpected body. Expected %s, got %s", expected, actual)
	}
}

func TestGetBirdHandler(t *testing.T) {
	birds = []Bird{
		{"sparrow", "A small harmless bird"},
	}

	request, error := http.NewRequest("GET", "", nil)

	if error != nil {
		t.Fatal(error)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(getBirdHandler)
	hf.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned incorrect status code. Expected %d, got %d", http.StatusOK, status)
	}

	expected := Bird{"sparrow", "A small harmless bird"}
	body := []Bird{}
	error = json.NewDecoder(recorder.Body).Decode(&body)

	if error != nil {
		t.Fatal(error)
	}

	actual := body[0]

	if actual != expected {
		t.Errorf("handler returned incorrect body. Expected %v, got %v", expected, actual)
	}
}

func newCreateBirdForm() *url.Values {
	form := url.Values{}
	form.Set("species", "eagle")
	form.Set("description", "A bird of prey")

	return &form
}
