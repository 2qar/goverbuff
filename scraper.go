package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	cloudfront = "https://dtmwra1jsgyb0.cloudfront.net/"
)

// Get information on the enemy team in a round of Open Division
func GetOtherTeam(tournamentLink, teamID string, round int) (teamInfo, error) {
	cutIndex := strings.LastIndex(tournamentLink, "/") + 1
	stageID := tournamentLink[cutIndex:]

	m, err := GetMatch(stageID, teamID, round)
	if err != nil {
		return teamInfo{}, err
	}

	var info teamInfo
	info, err = m.Team()
	if err != nil {
		return info, err
	}

	info.MatchLink = tournamentLink + "/match/" + m.ID
	return info, nil
}

func getPlayerInfo(activeIDs []string, p player, captain bool) playerInfo {
	active := captain
	if !captain {
		for i, id := range activeIDs {
			if id == p.ID {
				active = true
				activeIDs = append(activeIDs[:i], activeIDs[i+1:]...)
				break
			}
		}
	}

	stats, _ := GetPlayer(p.Btag())
	return playerInfo{active, p.User.Name, stats}
}

func (m *match) Team() (teamInfo, error) {
	var t team
	if m.Pos == "top" {
		t = m.Top
	} else {
		t = m.Bottom
	}

	resp, err := http.Get(cloudfront + "persistent-teams/" + t.Info.PID)
	if err != nil {
		return teamInfo{}, err
	}
	defer resp.Body.Close()

	var pts [1]persistentTeam
	err = json.NewDecoder(resp.Body).Decode(&pts)
	if err != nil {
		return teamInfo{}, err
	}
	pt := pts[0]

	var activeIDs []string
	ids := t.Info.ActiveIDS
	for _, p := range t.Info.Players {
		for i, id := range ids {
			if id == p.ID {
				activeIDs = append(activeIDs, p.PID)
				ids = append(ids[:i], ids[i+1:]...)
				break
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(pt.Players) + 1)
	pCh := make(chan playerInfo, len(pt.Players)+1)

	getPlayer := func(p player, captain bool) {
		defer wg.Done()
		info := getPlayerInfo(activeIDs, p, captain)
		pCh <- info
	}

	go getPlayer(pt.Captain, true)
	for _, p := range pt.Players {
		go getPlayer(p, false)
	}

	wg.Wait()
	close(pCh)

	var players []playerInfo
	for i := 0; i < len(pt.Players)+1; i++ {
		players = append(players, <-pCh)
	}

	return teamInfo{"", pt.Name, pt.Logo, players}, nil
}

type match struct {
	ID     string `json:"_id"`
	Pos    string
	Top    team `json:"top"`
	Bottom team `json:"bottom"`
}

type team struct {
	ID   string `json:"teamID"`
	Info struct {
		ActiveIDS []string `json:"playerIDs"`
		PID       string   `json:"persistentTeamID"`
		Players   []struct {
			ID  string `json:"_id"`
			PID string `json:"persistentPlayerID"`
		} `json:"players"`
	} `json:"team"`
}

type persistentTeam struct {
	Name    string   `json:"name"`
	Logo    string   `json:"logoUrl"`
	Captain player   `json:"persistentCaptain"`
	Players []player `json:"persistentPlayers"`
}

type player struct {
	ID   string `json:"_id"`
	PID  string `json:"persistentPlayerID"`
	IGN  string `json:"inGameName"`
	User struct {
		Name  string `json:"username"`
		Accts struct {
			Bnet struct {
				Btag string `json:"battletag"`
			} `'json:"battlenet"`
		} `json:"accounts"`
	} `json:"user"`
}

func (p *player) Btag() string {
	if p.IGN != "" {
		return p.IGN
	} else if p.User.Accts.Bnet.Btag != "" {
		return p.User.Accts.Bnet.Btag
	}
	return ""
}

type teamInfo struct {
	MatchLink string
	Name      string
	Logo      string
	Players   []playerInfo
}

type playerInfo struct {
	Active bool
	Name   string
	Stats  Player
}

// Find a match in the given round where a team with the given id is playing
func GetMatch(stageID, teamID string, round int) (match, error) {
	matchesLink := fmt.Sprintf(cloudfront+"stages/%s/rounds/%d/matches", stageID, round)

	resp, err := http.Get(matchesLink)
	if err != nil {
		return match{}, err
	}
	defer resp.Body.Close()

	var matches []match
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		return match{}, err
	}

	var foundMatch bool
	var m match
	for _, m = range matches {
		if m.Top.Info.PID == teamID {
			m.Pos = "bottom"
			foundMatch = true
			break
		} else if m.Bottom.Info.PID == teamID {
			m.Pos = "top"
			foundMatch = true
			break
		}
	}
	if !foundMatch {
		return match{}, errors.New("match not found")
	}

	matchLink := fmt.Sprintf(cloudfront+"matches/%s?extend[%s.team][players][users]", m.ID, m.Pos)
	resp, err = http.Get(matchLink)
	if err != nil {
		return match{}, nil
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		return match{}, err
	}

	pos := m.Pos
	m = matches[0]
	m.Pos = pos
	return m, nil
}
