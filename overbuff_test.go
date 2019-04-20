package odscraper

import (
	"testing"
)

func TestGetPlayer(t *testing.T) {
	p, err := GetPlayer("Tydra#11863")
	if err != nil {
		t.Error(err)
	} else if p.SR == 0 {
		t.Error("no sr")
	}
}

func TestGetInvalidPlayer(t *testing.T) {
	_, err := GetPlayer("ogdog")
	if err.Error() != "invalid btag" {
		t.Error("invalid btag somehow worked")
	}
}
