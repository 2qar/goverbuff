package main

import (
    "fmt"
)

func main() {
    c := make(chan Player)

    btags := [3]string{"Tydra#11863", "heckoffnerd#1772", "Thomas#11515"}

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
