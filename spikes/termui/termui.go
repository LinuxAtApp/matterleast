package main

import "github.com/nsf/termbox-go"

type TextBox struct {
	text string
	cursor_pos int
}

type ChatHistory struct {
	history []string
	authors []string
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
	termbox.SetCursor(15+text_box.cursor_pos,y-1)
	termbox.Flush()
}

func handle_backspace() {
	if len(text_box.text) > 0 && text_box.cursor_pos > 0 {
		text := text_box.text[:text_box.cursor_pos-1]
		if text_box.cursor_pos != len(text_box.text) {
			text += text_box.text[text_box.cursor_pos:]
		}
		text_box.text = text
		text_box.cursor_pos--
	}
}

func handle_textinput(ch rune) {
	switch {
		case len(text_box.text) == text_box.cursor_pos:
			text_box.text += string(ch)
		case text_box.cursor_pos == 0:
			text_box.text = string(ch) + text_box.text
		default:
			text_box.text = text_box.text[:text_box.cursor_pos] + string(ch) + text_box.text[text_box.cursor_pos:]
	}
	text_box.cursor_pos++
}

func handle_cursorleft() {
	if text_box.cursor_pos > 0 {
		text_box.cursor_pos--
	}
}

func handle_cursorright() {
	if text_box.cursor_pos < len(text_box.text) {
		text_box.cursor_pos++
	}
}

func handle_enter() {
	// Can also send message to back end here
	chat.history = append(chat.history, text_box.text)
	text_box.text = ""
	text_box.cursor_pos = 0
}

// TODO: This function will receive messages to be added to the history
func Receive_Message(message string) {
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
						handle_enter()
					case termbox.KeyBackspace, termbox.KeyBackspace2:
						handle_backspace()
					case termbox.KeyArrowLeft:
						handle_cursorleft()
					case termbox.KeyArrowRight:
						handle_cursorright()
					case termbox.KeyArrowUp:
					case termbox.KeyArrowDown:
					//TODO: Find better way to handle this
					default:
						handle_textinput(ev.Ch)
				}
			case termbox.EventError:
				panic(ev.Err)
		}
		redraw_all()
	}
}
