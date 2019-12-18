package goverbuff

import (
	"net/http"
	"os"
	"testing"
)

var client *http.Client

func init() {
	client = DefaultClient()
}

func TestParsePlayer(t *testing.T) {
	f, err := os.Open("example_player")
	if err != nil {
		t.Error(err)
	}

	p := parsePlayer(f)
	if p.SR != 3992 {
		t.Errorf("%d != 3992", p.SR)
	} else if r := p.Main(); r != "Offense" {
		t.Errorf("%s != Offense", r)
	}
}

func TestGetPlayer(t *testing.T) {
	p, err := GetPlayer(client, "Tydra#11863")
	if err != nil {
		t.Error(err)
	} else if p.SR == 0 {
		t.Error("no sr")
	}
}

func TestGetInvalidPlayer(t *testing.T) {
	_, err := GetPlayer(client, "ogdog")
	if err.Error() != "invalid btag" {
		t.Error("invalid btag somehow worked")
	}
}

func TestGetPlayerNoSR(t *testing.T) {
	p, err := GetPlayer(client, "OGDog#1515")
	if err != nil {
		t.Error(err)
	} else if p.SR != -1 {
		t.Error("sr not none")
	}
}

func TestGetPlayerMain(t *testing.T) {
	p, err := GetPlayer(client, "OGDog#1515")
	if err != nil {
		t.Error(err)
	} else if m := p.Main(); m != "Defense" {
		t.Errorf("wrong main: %q != %q", m, "Defense")
	}
}

func TestGetFakePlayer(t *testing.T) {
	_, err := GetPlayer(client, "TheresNoWayAnybodyHasThisName#1234")
	if !IsNotFound(err) {
		t.Errorf("not a NotFound error: %s", err)
	}
}
