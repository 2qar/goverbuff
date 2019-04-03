package main

import (
    "fmt"
    "golang.org/x/net/html"
    "io"
    //"os"
    "errors"
    "strings"
    "net/http"
    "strconv"
    "regexp"
)

type Player struct {
    SR      int
    Roles   map[string]int
}

func (p *Player) GetMain() string {
    var topRole string
    var topWins int
    for role, wins := range p.Roles {
        if wins > topWins {
            topWins = wins
            topRole = role
        }
    }
    return topRole
}

func parseRole(tokenizer *html.Tokenizer, roles map[string]int) {
    var currentRole string
    for {
        tt := tokenizer.Next()

        if tt == html.EndTagToken {
            t := tokenizer.Token()
            if t.Data == "tr" {
                return
            }
        }

        if tt == html.StartTagToken {
            t := tokenizer.Token()
            if t.Data == "td" && len(t.Attr) > 0{
                firstVal := t.Attr[0].Val

                if firstVal != "" {
                    games, err := strconv.Atoi(firstVal)
                    if err == nil {
                        roles[currentRole] = games
                    }
                }
            } else if t.Data == "a" {
                if t.Attr[1].Val == "color-white" {
                    tt = tokenizer.Next()
                    currentRole = tokenizer.Token().Data
                }
            }
        }
    }
}

func parsePlayer(r io.Reader) (p Player) {
    tokenizer := html.NewTokenizer(r)

    var sr string
    p.Roles = map[string]int{
        "Offense": 0,
        "Defense": 0,
        "Support": 0,
        "Tank": 0,
    }

    for {
        tt := tokenizer.Next()

        if tt == html.ErrorToken {
            break
        }

        if tt == html.StartTagToken {
            t := tokenizer.Token()
            if t.Data == "span" && len(t.Attr) == 1 { // check for sr
                if t.Attr[0].Val == "player-skill-rating" {
                    tt = tokenizer.Next()
                    t = tokenizer.Token()
                    if tt == html.TextToken {
                        sr = strings.Replace(t.Data, " ", "", -1)
                    }
                }
            } else if t.Data == "section" { // check for roles
                for i := 0; i < 2; i++ {
                    tt = tokenizer.Next()
                }
                if tt == html.TextToken {
                    t = tokenizer.Token()
                    if t.Data == "Roles" {
                        for {
                            tt = tokenizer.Next()

                            if tt == html.StartTagToken {
                                t = tokenizer.Token()
                                if t.Data == "tbody" && len(t.Attr) == 1 {
                                    for i := 0; i < 4; i++ {
                                        parseRole(tokenizer, p.Roles)
                                    }
                                    break
                                }
                            }
                        }
                    }
                }
                
            } 
        }
    }

    if sr == "" {
        p.SR = -1
    } else {
        p.SR, _ = strconv.Atoi(sr)
    }

    return
}

func GetPlayer(btag string) (Player, error) {
    if match, _ := regexp.MatchString("\\w{1,}#\\d{3,5}", btag); !match {
        return Player{}, errors.New("invalid btag")
    } 

    validTag := strings.Replace(btag, "#", "-", 1)
    resp, err := http.Get(fmt.Sprintf("https://www.overbuff.com/players/pc/%s", validTag))
    if err != nil {
        return Player{}, err
    }
    defer resp.Body.Close()

    return parsePlayer(resp.Body), nil
}

/*
func main() {
    docFile, _ := os.Open("Tydra-11863")
    defer docFile.Close()

    player := parsePlayer(docFile)
    
    fmt.Println(player)
    fmt.Println(player.GetMain())
}
*/
