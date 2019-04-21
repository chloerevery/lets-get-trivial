package main

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "math/rand"
    "net/http"
    "time"
)

const InputErrorMessage = "Please send \"<prompt>\" \"<answer>\" \"<attribution>\""
const NoQuestionsMessage = "This channel doesn't have any trivia questions yet! Add one now!"
const DEFAULT_CHANNEL = "allTrivia"

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/trivium", GetTrivium).Methods("GET")
    router.HandleFunc("/trivium", CreateTrivium).Methods("POST")

    allTrivia.idx = 0

    // Seed some initial trivia questions.
    allTrivia.trivia = append(allTrivia.trivia, Trivium{Prompt: "What is Dr. Seuss's real name?", Answer: "Theodore Geisel"})
    allTrivia.trivia = append(allTrivia.trivia, Trivium{Prompt: "In what country is the region of Andalusia located?", Answer: "Spain"})
    allTrivia.trivia = append(allTrivia.trivia, Trivium{Prompt: "In what theater was Lincoln killed?", Answer: "Ford"})

    Channels[DEFAULT_CHANNEL] = allTrivia

    log.Fatal(http.ListenAndServe(":8080", router))
}

func GetTrivium(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    channel := DEFAULT_CHANNEL
    k, ok := params["channel_name"]
    if ok {
        channel = k[0]
    }
    group, ok := Channels[channel]
    if !ok {
        // Channel name not recognized.
        e := TriviaError{message: NoQuestionsMessage}
        json.NewEncoder(w).Encode(e.message)
        return
    }

    if len(group.trivia) == 0 {
        e := TriviaError{message: NoQuestionsMessage}
        json.NewEncoder(w).Encode(e.message)
        return
    }
    t := group.trivia[group.idx]

    group.idx += 1 // Increment pointer to next question for group.
    if group.idx >= len(group.trivia) {
        group.idx = 0 // Wrap around.
    }

    // Save new idx.
    Channels[channel] = group

    json.NewEncoder(w).Encode(t.Prompt)
    return
}

func CreateTrivium(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var tr CreateTriviumRequest
    err := decoder.Decode(&tr)
    if err != nil {
        e := TriviaError{message: InputErrorMessage}
        json.NewEncoder(w).Encode(e.message)
        return
    }

    t := Trivium{
        Prompt:        tr.Prompt,
        Answer:        tr.Answer,
        AnswerDetails: tr.AnswerDetails,
        Attribution:   tr.Attribution,
    }

    _, ok := Channels[tr.ChannelName]
    if !ok {
        // Channel name not recognized. Save channel as new trivia group.
        Channels[tr.ChannelName] = TriviaGroup{groupName: tr.ChannelName, idx: 0}
    }

    // Add new trivium to group.
    group := Channels[tr.ChannelName]
    group.trivia = append(group.trivia, t)
    Channels[tr.ChannelName] = group

    json.NewEncoder(w).Encode(t)
    return
}

type CreateTriviumRequest struct {
    ChannelName   string `json:"channel_name"`
    Prompt        string `json:"prompt"`
    Answer        string `json:"answer"`
    AnswerDetails string `json:"answer_details"`
    Attribution   string `json:"attribution"`
}

type Trivium struct {
    Prompt        string `json:"prompt"`
    Answer        string `json:"answer"`
    AnswerDetails string `json:"answer_details"`
    Attribution   string `json:"attribution"`
}

type TriviaGroup struct {
    groupName string
    trivia    []Trivium
    idx       int
}

type TriviaError struct {
    message string
}

type Response struct {
    responseType string `json:"response_type"`
    text         string `json:"text"`
}

var Channels = map[string]TriviaGroup{}

var allTrivia TriviaGroup
