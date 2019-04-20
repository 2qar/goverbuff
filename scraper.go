package odscraper

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

func FindTeam(tournamentID, name string) (TeamInfo, error) {
	url := cloudfront + "tournaments/" + tournamentID + "/" + "teams?name=" + name
	resp, err := http.Get(url)
	if err != nil {
		return TeamInfo{}, err
	}
	defer resp.Body.Close()

	var teams []teamData
	err = json.NewDecoder(resp.Body).Decode(&teams)
	if err != nil {
		return TeamInfo{}, err
	}
	if len(teams) == 0 {
		return TeamInfo{}, errors.New(fmt.Sprintf("unable to find team \"%s\"", name))
	}

	info, err := getTeamInfo(teams[0])
	if err != nil {
		return TeamInfo{}, nil
	}

	return info, nil
}

// Get information on the enemy team in a round of Open Division
func GetOtherTeam(tournamentLink, teamID string, round int) (TeamInfo, error) {
	cutIndex := strings.LastIndex(tournamentLink, "/") + 1
	stageID := tournamentLink[cutIndex:]

	m, err := GetMatch(stageID, teamID, round)
	if err != nil {
		return TeamInfo{}, err
	}

	info, err := getTeamInfo(m.Team().Info)
	if err != nil {
		return info, err
	}

	info.MatchLink = tournamentLink + "/match/" + m.ID
	return info, nil
}

func getPlayerInfo(activeIDs []string, p player, captain bool) PlayerInfo {
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
	return PlayerInfo{active, p.User.Name, stats}
}

func getTeamInfo(t teamData) (TeamInfo, error) {
	resp, err := http.Get(cloudfront + "persistent-teams/" + t.PID)
	if err != nil {
		return TeamInfo{}, err
	}
	defer resp.Body.Close()

	var pts [1]persistentTeam
	err = json.NewDecoder(resp.Body).Decode(&pts)
	if err != nil {
		return TeamInfo{}, err
	}
	pt := pts[0]

	var activeIDs []string
	ids := t.ActiveIDS
	for _, p := range t.Players {
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
	pCh := make(chan PlayerInfo, len(pt.Players)+1)

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

	var players []PlayerInfo
	for i := 0; i < len(pt.Players)+1; i++ {
		players = append(players, <-pCh)
	}

	return TeamInfo{"", pt.Name, pt.Logo, players}, nil
}

type match struct {
	ID     string `json:"_id"`
	Pos    string
	Top    team `json:"top"`
	Bottom team `json:"bottom"`
}

func (m *match) Team() team {
	var t team
	if m.Pos == "top" {
		t = m.Top
	} else {
		t = m.Bottom
	}
	return t
}

type team struct {
	ID             string         `json:"teamID"`
	Info           teamData       `json:"team"`
	PersistentTeam persistentTeam `json:"persistentTeam"`
}

type teamData struct {
	ActiveIDS []string `json:"playerIDs"`
	PID       string   `json:"persistentTeamID"`
	Players   []struct {
		ID  string `json:"_id"`
		PID string `json:"persistentPlayerID"`
	} `json:"players"`
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

type TeamInfo struct {
	MatchLink string
	Name      string
	Logo      string
	Players   []PlayerInfo
}

type PlayerInfo struct {
	Active bool
	Name   string
	Stats  PlayerStats
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
