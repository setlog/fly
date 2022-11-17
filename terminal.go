package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/setlog/panik"
)

func selectOption(title string, options []string) int {
	screen, err := tcell.NewScreen()
	panik.OnError(err)
	panik.OnError(screen.Init())
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)
	screen.Clear()
	drawText(screen, 0, 0, defStyle, title)
	for i, option := range options {
		drawText(screen, 0, i+1, defStyle, fmt.Sprintf("%d. %s", i+1, option))
	}
	for {
		screen.Show()
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				screen.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyRune {
				if ev.Rune() >= '1' && ev.Rune() <= rune('1'+len(options)-1) {
					screen.Fini()
					return int(ev.Rune() - '1')
				}
			}
		}
	}
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		s.SetContent(x, y, r, nil, style)
		x++
	}
}
