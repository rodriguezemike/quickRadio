package main

import (
	"context"
	"fmt"
	"quickRadio/radioErrors"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func GetGameHtml(linksMap map[string]interface{}) string {
	var html string
	sleepTimer := linksMap["load_sleep_timer"].(float64)
	baseUrl := fmt.Sprintf("%v", linksMap["base"])

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(baseUrl),
		chromedp.Sleep(time.Duration(sleepTimer)*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	radioErrors.ErrorCheck(err)
	return html
}

func main() {
	/*
		linksMap := internals.GetLinksJson()
		gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
		gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
		team := fmt.Sprintf("%v", linksMap["team_abbrev"])
		gameRegexs := []string{fmt.Sprintf("%v", linksMap["home_game_regex"]), fmt.Sprintf("%v", linksMap["away_game_regex"])}
		html := GetGameHtml(linksMap)
		landingLink, err := internals.GetGameLandingLink(html, gamecenterBase, gamecenterLanding, gameRegexs)
		internals.ErrorCheck(err)
		gameDataObject := internals.GetGameDataObjectFromResponse(landingLink)
		radioLink, err := internals.GetRadioLink(gameDataObject, team)
		internals.ErrorCheck(err)
		log.Println(radioLink)
	*/
}
