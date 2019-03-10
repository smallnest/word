package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func downloadAndPlay(dir string, name string, u string) {
	fileName := filepath.Join(dir, name+".mp3")
	if _, err := os.Stat(filepath.Join(dir, name+".mp3")); !os.IsNotExist(err) {
		play(fileName)
		return
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("failed to create file %s: %v\n", fileName, err)
		return
	}
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("failed to downlaod mp3 %s: %v\n", u, err)
		return
	}

	io.Copy(f, resp.Body)
	resp.Body.Close()
	f.Close()

	play(fileName)
}

func play(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("failed to read mp3: %v\n", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Printf("failed to decode mp3: %v\n", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	time.Sleep(200 * time.Millisecond)
}
