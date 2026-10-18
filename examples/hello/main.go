package main

import (
	"context"
	"os"

	"github.com/NLipatov/tuigo/color"
	"github.com/NLipatov/tuigo/compose"
	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/terminal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := terminal.NewSession(ctx, os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
	events, err := session.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()
	if err := render(&session); err != nil {
		panic(err)
	}
	for event := range events {
		switch event.Type {
		case terminal.EventKey:
			if quitRequested(event.Key) {
				cancel()
				return
			}
		case terminal.EventResize:
			if err := render(&session); err != nil {
				panic(err)
			}
		case terminal.EventError:
			panic(event.Err)
		}
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

func render(session *terminal.Session) error {
	w, h, err := session.Size()
	if err != nil {
		return err
	}
	message := "Hello World!"
	exitTip := "Exit: q | Ctrl + C | Esc"
	root := compose.Div{
		Width:  w,
		Height: h,
		Fg:     color.FgWhite,
		Bg:     color.BgBlack,
		Children: []compose.Node{
			&compose.Div{
				Width:  w,
				Height: max(h-1, 1),
				Fg:     color.FgWhite,
				Bg:     color.BgBlack,
				Children: []compose.Node{
					&compose.Div{
						Width:  w,
						Height: max((h/2)-1, 1) + 1,
						Fg:     color.FgWhite,
						Bg:     color.BgBlack,
						Padding: compose.Padding{
							Up:   max((h/2)-1, 1),
							Left: max((w/2)-(len(message)/2), 1),
						},
						Children: []compose.Node{
							&compose.Text{
								Content: message,
								Fg:      color.FgGreen,
								Bg:      color.BgBlack,
							},
						},
						Direction: compose.Row,
					},
					&compose.Div{
						Width:  w,
						Height: max(h/16, 1) + 1,
						Fg:     color.FgWhite,
						Bg:     color.BgBlack,
						Padding: compose.Padding{
							Up:   max(h/16, 1),
							Left: max((w/2)-(len(exitTip)/2), 1),
						},
						Children: []compose.Node{
							&compose.Text{
								Content: exitTip,
								Fg:      color.FgWhite,
								Bg:      color.BgBlack,
							},
						},
						Direction: compose.Row,
					},
				},
				Direction: compose.Column,
			},
		},
		Direction: compose.Column,
	}
	frame, err := root.Frame()
	if err != nil {
		return err
	}
	return session.Render(frame)
}
