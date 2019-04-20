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

func TestGetPlayerNoSR(t *testing.T) {
	p, err := GetPlayer("OGDog#1515")
	if err != nil {
		t.Error(err)
	} else if p.SR != -1 {
		t.Error("sr not none")
	}
}

func TestGetPlayerMain(t *testing.T) {
	p, err := GetPlayer("OGDog#1515")
	if err != nil {
		t.Error(err)
	} else if p.GetMain() != "Defense" {
		t.Error("wrong main")
	}
}

func TestGetFakePlayer(t *testing.T) {
	_, err := GetPlayer("TheresNoWayAnybodyHasThisName#1234")
	if !IsNotFound(err) {
		t.Error(err)
	}
}
