package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type AnswerPayload struct {
	Answer   string `json:"answer"`
	UserName string `json:"userId"`
	GameId   int64  `json:"gameId"`
}

func packPayload(userName string, answer string, gameId int64) AnswerPayload {
	return AnswerPayload{
		UserName: userName,
		Answer:   answer,
		GameId:   gameId,
	}
}

type AnswerResponse struct {
	ID       int64  `json:"id"`
	Question string `json:"question"`
	NextAt   int64  `json:"nextAt"`
	StartAt  int64  `json:"startAt"`

	Answer     *string `json:"answer,omitempty"`
	Winner     *string `json:"winner,omitempty"`
	NumOfTries *int64  `json:"numOfTries,omitempty"`
	Rewards    *int64  `json:"rewards,omitempty"`
	FinishedAt *int64  `json:"finishedAt,omitempty"`
}

type QuestionResponse struct {
	ID       int64  `json:"id"`
	Question string `json:"question"`
	Rewards  int64  `json:"rewards"`
}

func postAnswer(endpoint string, payload AnswerPayload, response *AnswerResponse) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, response)
}

func fetchQuestion(endpoint string, response *QuestionResponse) error {
	resp, err := http.Get(endpoint + "/state")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, response)
}

func playGame(userName string, endpoint string) {
	var (
		err          error
		ansPayload   AnswerPayload
		ansResponse  AnswerResponse
		quesResponse QuestionResponse
	)
	gameId := int64(0)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		answer := fmt.Sprintf("%02x", rand.Intn(0xFFF)) // generate a random hex string from 0 to FF
		ansPayload = packPayload(userName, answer, gameId)
		if err = postAnswer(endpoint, ansPayload, &ansResponse); err != nil {
			postError(endpoint, userName, gameId, "Error posting answer:", err)
			goto FAILURE_RETRY_LATER
		}

		if ansResponse.Answer == nil || *ansResponse.Answer == "" {
			continue
		}

		if ansResponse.Winner != nil && *ansResponse.Winner == userName {
			fmt.Printf("Bot %s won with answer %s for game %d, rewards: %d/%d\n", userName, *ansResponse.Answer, gameId, *ansResponse.Rewards, *ansResponse.NumOfTries)
		}

	FAILURE_RETRY_LATER:
		if err = fetchQuestion(endpoint, &quesResponse); err != nil {
			postError(endpoint, userName, gameId, "Error fetching question:", err)
			time.Sleep(10 * time.Second)
			ticker.Reset(1 * time.Second)
		}
		gameId = quesResponse.ID
		quesResponse = QuestionResponse{}
		ansResponse = AnswerResponse{}
	}
}

func main() {
	server := "http://localhost:5173"
	if len(os.Args) > 1 {
		server = os.Args[1]
	}
	botNames := []string{
		"Tra Pham",
		"John Doe",
		"Jane Smith",
		"Bob Johnson",
		"Alice Williams",
		"Charlie Brown",
		"Emily Davis",
		"David Miller",
		"Sarah Jones",
		"Tom Wilson",
		"Emma Moore",
		"Michael Taylor",
		"Anna Thomas",
		"James Anderson",
		"Jessica Jackson",
		"Robert White",
		"Patricia Harris",
		"William Martin",
		"Linda Thompson",
		"Richard Garcia",
		"Jennifer Martinez",
	}

	serverEndpoints := []string{
		server + "/api/game/1",
		server + "/api/game/2",
		server + "/api/game/3",
		server + "/api/game/4",
		server + "/api/game/5",
	}

	var wg sync.WaitGroup

	for _, botName := range botNames {
		for _, endpoint := range serverEndpoints {
			wg.Add(1)
			go func(name string, endpoint string) {
				defer wg.Done()
				playGame(name, endpoint)
			}(botName, endpoint)
		}
	}
	wg.Wait()
}

func postError(endpoint string, username string, gameid int64, note string, err error) {
	fmt.Printf("Server: %s, Bot: %s, Failed in game %d, at %s, with error %s\n", endpoint, username, gameid, note, err.Error())
}
