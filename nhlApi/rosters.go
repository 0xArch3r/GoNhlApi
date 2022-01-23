package nhlApi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Player struct {
	Person struct {
		ID       int    `json:"id"`
		FullName string `json:"fullName"`
		Link     string `json:"link"`
	} `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	Position     struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
}

type nhlRosterResponse struct {
	Roster []Player `json:"roster"`
}

func GetRoster(teamID int) ([]Player, error) {
	res, err := http.Get(fmt.Sprintf("%s/teams/%v/roster", baseUrl, teamID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response nhlRosterResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Roster, nil
}
