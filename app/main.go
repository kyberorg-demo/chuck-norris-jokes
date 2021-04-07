package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const chuckNorrisApi = "http://api.icndb.com/jokes/random?limitTo=[nerdy]"

const (
	errorApiTimeout       = "ERR_API_TIMEOUT"
	errorApiMalformedJson = "ERR_API_BAD_JSON"
)

type Joke struct {
	ID   uint32 `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

func getJoke() (string, error) {
	c := http.Client{}

	resp, err := c.Get(chuckNorrisApi)
	if err != nil {
		return "jokes API not responding", errors.New(errorApiTimeout)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	joke := JokeResponse{}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "joke error", errors.New(errorApiMalformedJson)
	}
	return joke.Value.Joke, nil
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	respMessage, jokeError := getJoke()
	if jokeError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(respMessage))
		log.Println(writeErr)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write([]byte(respMessage))
	if writeErr != nil {
		log.Println(writeErr)
		return
	}
}

func main() {
	http.HandleFunc("/", HandleRequest)
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln(err)
}
