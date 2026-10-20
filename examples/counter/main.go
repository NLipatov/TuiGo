package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NLipatov/tuigo/color"
	"github.com/NLipatov/tuigo/compose"
	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/terminal"
)

type counter struct {
	val uint
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := terminal.NewSession(ctx, os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()
	events, err := session.Start()
	if err != nil {
		panic(err)
	}
	counter := &counter{}
	render(session, counter)
	for event := range events {
		switch event.Type {
		case terminal.EventError:
			panic(event.Err)
		case terminal.EventKey:
			if quitRequested(event.Key) {
				cancel()
				return
			}
			if event.Key.Code == keyboard.KeyRune && event.Key.Text == " " {
				counter.val++
			}
			if event.Key.Code == keyboard.KeyRune && event.Key.Text == "r" {
				counter.val = 0
			}
			render(session, counter)
		default:
			render(session, counter)
		}
	}
}

func render(session terminal.Session, counter *counter) {
	w, h, err := session.Size()
	if err != nil {
		panic(err)
	}
	counterMessage := fmt.Sprintf("Counter: %d", counter.val)
	controlMessage := "Press 'space' to increment, 'r' to reset"
	exitMessage := "Exit: q | Ctrl + C | Esc"
	root := &compose.Div{
		Width:  w,
		Height: h,
		Padding: compose.Padding{
			Up:   max(1, (h/2)-3),
			Left: max(1, (w/2)-len(exitMessage)),
		},
		Fg: color.FgWhite,
		Bg: color.BgBlack,
		Children: []compose.Node{
			&compose.Text{
				Content: counterMessage,
				Fg:      color.FgGreen,
				Bg:      color.BgBlack,
			},
			&compose.Text{
				Content: controlMessage,
				Fg:      color.FgWhite,
				Bg:      color.BgBlack,
			},
			&compose.Text{
				Content: exitMessage,
				Fg:      color.FgWhite,
				Bg:      color.BgBlack,
			},
		},
		Direction: compose.Column,
	}
	frame, err := root.Frame()
	if err != nil {
		panic(err)
	}
	if err := session.Render(frame); err != nil {
		panic(err)
	}
}

func quitRequested(event keyboard.KeyEvent) bool {
	if event.Code == keyboard.KeyEsc {
		return true
	}
	if event.Code != keyboard.KeyRune {
		return false
	}
	if event.Text == "q" {
		return true
	}
	return event.Text == "c" && event.Mod&keyboard.ModCtrl != 0
}
