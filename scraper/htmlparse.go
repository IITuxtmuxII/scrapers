package scraper

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//getHrefMain scrapes the main index listing of the page for links to torrent descriptions
func getHrefMain(tok html.Token) (ok bool, href string) {
	for _, a := range tok.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return //implicit signature return
}

func parsePageMain(tokenizer html.Tokenizer, chURL chan<- string, chFin chan<- bool, chCount chan<- int) {
	var leave bool

	for {
		if leave {
			break
		}
		tokType := tokenizer.Next()
		switch tokType {
		case html.ErrorToken:
			//EOF
			leave = true
			break
		case html.StartTagToken:
			tok := tokenizer.Token()
			isAnchor := tok.Data == "a"
			if !isAnchor {
				continue
			}
			ok, url := getHrefMain(tok)

			//If the URL ends in t, it's a nyaa.si torrent quicklink
			//If the URL ends in s, it's a nyaa.si comment quicklick
			//We don't want either of those right now so skip them
			if !ok || url[len(url)-1:] == "t" || url[len(url)-1:] == "s" {
				continue
			}

			//TODO: Factor this more cleanly
			//var names correspond to their respective sites
			nyaaSi := strings.Index(url, "/view") == 0
			anidex := strings.Index(url, "?page=torrent&id=") == 0
			if nyaaSi {
				for len(chUrl) == cap(chUrl) {
					fmt.Println("Nyaa channel full, sleeping 3 seconds")
					time.Sleep(time.Millisecond * 3000)
				}
				chUrl <- baseURL + url
				chURLCount <- 1
				childPageCount++
				continue
			}
			if anidex {
				for len(chURL) == cap(chURL) {
					fmt.Println("Anidex channel full, sleeping 3 seconds")
					time.Sleep(time.Millisecond * 3000)
				}
				chAnidexURL <- url[17:]
				chURLCount <- 1
				childPageCount++
				continue
			}
		}
	}
}
