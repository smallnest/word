package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

var us bool

type Engine interface {
	audio(word string, us bool) (mp3, ipa, def string, err error)
}

var engine Engine = &cambridge{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: word <your_word_to_pronounce>")
		return
	}

	if os.Args[1] == "-list" {
		list()
		return
	}

	if len(os.Args) == 3 && os.Args[2] == "-us" {
		us = true
	}
	word := os.Args[1]

	playWord(word)
}

func list() {
	if len(os.Args) == 4 && os.Args[3] == "-us" {
		us = true
	}
	file := os.Args[2]

	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("failed to get word list: %v\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	d := color.New(color.FgHiGreen, color.Bold)
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println()
		d.Println(word)
		playWord(word)
		time.Sleep(2 * time.Second)
	}
}

func playWord(word string) {
	// go 1.12 added os.UserHomeDir but I want it can support other go versions below 1.12
	d, err := homedir.Dir()
	if err != nil {
		fmt.Println("failed to get the user home dir")
		return
	}

	var path string
	if us {
		path = d + "/.words/us"
	} else {
		path = d + "/.words/uk"
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}

	// first check whether mp3 exists？
	done := make(chan bool)
	var mp3Exist bool
	if _, err := os.Stat(path + "/" + word + ".mp3"); !os.IsNotExist(err) {
		go play(path+"/"+word+".mp3", done)
		mp3Exist = true
	}

	mp3, ipa, def, err := engine.audio(word, us)
	if err != nil {
		return
	}

	fmt.Println(ipa)
	fmt.Println()
	fmt.Println(def)

	if !mp3Exist {
		go downloadAndPlay(path, word, mp3, done)
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
	}
}
