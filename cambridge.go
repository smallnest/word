package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// <span class="uk"><span class="region">uk</span>
//                   <span title="semaphore: listen to British English pronunciation" data-src-mp3="/media/english/uk_pron/u/uks/uksel/ukselfs013.mp3" data-src-ogg="/media/english/uk_pron_ogg/u/uks/uksel/ukselfs013.ogg" class="circle circle-btn sound audio_play_button">
//             <i class="fcdo fcdo-volume-up">​</i>
//         </span>
// 			   <span class="pron">/<span class="ipa">ˈsem.ə.fɔː<sup class="sp">r</sup></span>/</span>  </span>

// <span class="us"><span class="region">us</span>
//                   <span title="semaphore: listen to American pronunciation" data-src-mp3="/media/english/us_pron/e/eus/eus74/eus74137.mp3" data-src-ogg="/media/english/us_pron_ogg/e/eus/eus74/eus74137.ogg" class="circle circle-btn sound audio_play_button">
//             <i class="fcdo fcdo-volume-up">​</i>
//         </span>
// 			   <span class="pron">/<span class="ipa">ˈsem.ə.fɔːr</span>/</span>  </span>

// <div class="pos-body">

//             <div class="sense-block" id="cald4-1-1-1">  <div class="sense-body">
//             <div class="def-block pad-indent" data-wl-senseid="ID_00028775_01"><a class="wordlist-button circle circle-btn circle-btn--sml circle-btn--alt fav-entry wordlist-form" title="Add this meaning to a word list"><i class="fcdo fcdo-star" aria-hidden="true">​</i></a><p class="def-head semi-flush"><span class="def-info"><span class="freq">›</span> </span><b class="def">a <a class="query" href="https://dictionary.cambridge.org/dictionary/english/system" title="system">system</a> of <a class="query" href="https://dictionary.cambridge.org/dictionary/english/communication" title="communication">communication</a> using two <a class="query" href="https://dictionary.cambridge.org/dictionary/english/flag" title="flags">flags</a> <a class="query" href="https://dictionary.cambridge.org/dictionary/english/held" title="held">held</a> in <a class="query" href="https://dictionary.cambridge.org/dictionary/english/your" title="your">your</a> <a class="query" href="https://dictionary.cambridge.org/dictionary/english/hand" title="hands">hands</a> that are <a class="query" href="https://dictionary.cambridge.org/dictionary/english/moved" title="moved">moved</a> into different <a class="query" href="https://dictionary.cambridge.org/dictionary/english/position" title="positions">positions</a> to <a class="query" href="https://dictionary.cambridge.org/dictionary/english/represent" title="represent">represent</a> different <a class="query" href="https://dictionary.cambridge.org/dictionary/english/capital" title="letters">letters</a>, <a class="query" href="https://dictionary.cambridge.org/dictionary/english/number" title="numbers">numbers</a>, or <a class="query" href="https://dictionary.cambridge.org/dictionary/english/symbol" title="symbols">symbols</a></b></p></div>
// 			</div>                    </div></div>

var cambridgeURL = "https://dictionary.cambridge.org/dictionary/english/"

type cambridge struct{}

func (e cambridge) audio(word string, us bool) (mp3, ipa, def string, err error) {
	u := cambridgeURL + word
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("failed to get audio from cambridge: %v\n", err)
		return mp3, ipa, def, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response from cambridge: %v\n", err)
		return mp3, ipa, def, err
	}

	if us {
		val, ok := doc.Find(".us span.audio_play_button").Attr("data-src-mp3")
		if !ok {
			fmt.Println("failed to get audio")
			return mp3, ipa, def, errors.New("not found")
		}

		mp3 = "https://dictionary.cambridge.org" + val

		ipa = doc.Find(".us span.ipa").Text()
	} else { //uk
		val, ok := doc.Find(".uk span.audio_play_button").Attr("data-src-mp3")
		if !ok {
			fmt.Println("failed to get audio")
			return mp3, ipa, def, errors.New("not found")
		}

		mp3 = "https://dictionary.cambridge.org" + val

		ipa = doc.Find(".uk span.ipa").Text()
	}

	def = doc.Find("p.def-head b.def").Text()

	return mp3, ipa, def, nil
}
