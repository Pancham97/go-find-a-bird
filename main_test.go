package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	//Here, we form a new HTTP request. This is the request that's going to be
	// passed to our handler.
	// The first argument is the method, the second argument is the route (which
	// we leave blank for now, and will get back to soon), and the third is the
	// request body, which we don't have in this case.
	request, error := http.NewRequest("GET", "", nil)

	if error != nil {
		t.Fatal(error)
	}

	// We use Go's httptest library to create an http recorder. This recorder
	// will act as the target of our http request
	// (you can think of it as a mini-browser, which will accept the result of
	// the http request that we make)
	recorder := httptest.NewRecorder()

	// Create an HTTP handler from our handler function. "handler" is the handler
	// function defined in our main.go file that we want to test
	hf := http.HandlerFunc(handler)

	// Serve the HTTP request to our recorder. This is the line that actually
	// executes our the handler that we want to test
	hf.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned incorrect status code. Expected %v, got %v", status, http.StatusOK)
	}

	expected := `Hello World`
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("handler returned unexpected body. Expected %v, got %v", expected, actual)
	}
}

func TestRouter(t *testing.T) {
	// Instantiate the router using the constructor function that
	// we defined previously
	r := newRouter()

	// Create a new server using the "httptest" libraries `NewServer` method
	// Documentation : https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	// The mock server we created runs a server and exposes its location in the
	// URL attribute
	// We make a GET request to the "hello" route we defined in the router
	response, error := http.Get(mockServer.URL + "/hello")

	if error != nil {
		t.Fatal(error)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Status should be okay. Got %d", response.StatusCode)
	}

	// In the next few lines, the response body is read, and converted to a string
	defer response.Body.Close()
	// read the body into a bunch of bytes (b)
	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		t.Fatal(error)
	}

	// Convert bytes to string
	responseString := string(body)
	expected := `Hello World`

	if responseString != expected {
		t.Errorf("Response is not correct. Expected %s, got %s", expected, responseString)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()

	mockServer := httptest.NewServer(r)

	// Most of the code is similar. The only difference is that now we make a
	// request to a route we know we didn't define, like the `POST /hello` route.
	response, error := http.Post(mockServer.URL+"/hello", "", nil)

	if error != nil {
		t.Fatal(error)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", response.StatusCode)
	}

	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		t.Fatal(error)
	}

	responseString := string(body)
	expected := ""

	if responseString != expected {
		t.Errorf("Expected body to be %s, got %s", expected, responseString)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// We want to hit the `GET /assets/` route to get the index.html file response
	response, error := http.Get(mockServer.URL + "/assets/")

	if error != nil {
		t.Fatal(error)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// It isn't wise to test the entire content of the HTML file.
	// Instead, we test that the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := response.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type. Expected %s, got %s", expectedContentType, contentType)
	}
}
