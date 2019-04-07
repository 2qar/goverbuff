package main

import (
	"fmt"
)

func testGetTeam() {
	c := make(chan Player)

	btags := []string{
		"Tydra#11863",
		"heckoffnerd#1772",
		"Thomas#11515",
		"Mason#12245",
		"Mason#12841",
		"Mason#12454",
		"HoopoeVX#1905",
		"Blake#13335",
		"AVA#11577",
	}

	for _, btag := range btags {
		go func(btag string) {
			player, _ := GetPlayer(btag)
			c <- player
		}(btag)
	}

	for _ = range btags {
		fmt.Println(<-c)
	}
}

func main() {
	//fmt.Println(GetMatch("5bf8b741b06aae03a9f18385", "5bfe1b9418ddd9114f14efb0", 1))
	teamInfo, err := GetOtherTeam("https://battlefy.com/overwatch-open-division-north-america/2019-overwatch-open-division-season-2-north-america/5c7ccfe88d004d0345bbd0cd/stage/5c929d720bc67d0345180aa6", "5bfe1b9418ddd9114f14efb0", 1)
	if err != nil {
		panic(err)
	}

	for _, p := range teamInfo.Players {
		fmt.Println(p)
	}
}
