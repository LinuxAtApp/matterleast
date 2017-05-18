package main

import "github.com/nsf/termbox-go"

type TextBox struct {
	text string
}

type ChatHistory struct {
	history []string
}

var text_box = TextBox{}
var chat = ChatHistory{}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	_, y := termbox.Size()
	tbprint(0,y-1,termbox.ColorWhite, termbox.ColorDefault, "Enter Message>")
	tbprint(15, y-1,termbox.ColorWhite, termbox.ColorDefault, text_box.text)
	l := len(chat.history) - 1
	for i := l; i >= 0; i-- {
		tbprint(0, y-2-(l-i), termbox.ColorWhite, termbox.ColorDefault, "You: "+chat.history[i])
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	redraw_all()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
					case termbox.KeyEsc:
						break mainloop
					case termbox.KeyEnter:
						chat.history = append(chat.history, text_box.text)
						text_box.text = ""
					default:
						text_box.text += string(ev.Ch)
				}
			case termbox.EventError:
				panic(ev.Err)
		}
		redraw_all()
	}
}
