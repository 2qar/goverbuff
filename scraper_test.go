package odscraper

import (
	"strings"
	"testing"
)

var (
	stageID        = "5bf8b741b06aae03a9f18385"
	tournamentID   = "5c7ccfe88d004d0345bbd0cd"
	teamID         = "5bfe1b9418ddd9114f14efb0"
	tournamentLink = "https://battlefy.com/overwatch-open-division-north-america/2019-overwatch-open-division-season-2-north-america/5c7ccfe88d004d0345bbd0cd/stage/5c929d720bc67d0345180aa6"
)

func TestFindOneTeam(t *testing.T) {
	var team TeamInfo
	_, err := FindTeam(tournamentID, "Vixen", &team)
	if err != nil {
		t.Error(err)
	} else if team.Name == "" && len(team.Players) == 0 {
		t.Error("team empty")
	}
}

func TestFindManyTeams(t *testing.T) {
	_, err := FindTeam(tournamentID, "The", &TeamInfo{})
	if err != nil {
		t.Error(err)
	}
}

func TestFindInvalidTeam(t *testing.T) {
	var team TeamInfo
	_, err := FindTeam(tournamentID, "ThisTeamDoesNotExist", &team)
	if !strings.HasPrefix(err.Error(), "unable to find team") {
		t.Error(err)
	}
}

func TestGetOtherTeam(t *testing.T) {
	e, err := GetOtherTeam(tournamentLink, teamID, 1)
	if err != nil {
		t.Error(err)
	} else if e.Link == "" {
		t.Error("no match link")
	}
}
