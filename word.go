package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

var us bool

var engine = &cambridge{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: word <your_word_to_pronounce>")
		return
	}

	if len(os.Args) == 3 && os.Args[2] == "-us" {
		us = true
	}

	word := os.Args[1]

	mp3, ipa, def, err := engine.audio(word, us)
	if err != nil {
		return
	}

	fmt.Println(ipa)
	fmt.Println()
	fmt.Println(def)

	// go 1.12 added os.UserHomeDir but I want it can support other go versions below 1.12
	d, err := homedir.Dir()
	if err != nil {
		fmt.Println("failed to get the user home dir")
		return
	}

	if _, err := os.Stat(d + "/.words"); os.IsNotExist(err) {
		os.Mkdir(d+"/.words", 0644)
	}

	downloadAndPlay(d+"/.words", word, mp3)
}
