package main

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

func TestFindTeam(t *testing.T) {
	_, err := FindTeam(tournamentID, "Vixen")
	if err != nil {
		t.Error(err)
	}
}

func TestFindInvalidTeam(t *testing.T) {
	_, err := FindTeam(tournamentID, "Vixen Gaming")
	if !strings.HasPrefix(err.Error(), "unable to find team") {
		t.Error(err)
	}
}

func TestGetOtherTeam(t *testing.T) {
	_, err := GetOtherTeam(tournamentLink, teamID, 1)
	if err != nil {
		t.Error(err)
	}
}
