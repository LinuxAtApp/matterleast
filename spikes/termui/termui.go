package main

import (
	"github.com/LinuxAtApp/matterleast/spikes/termui/datalog"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
	"strconv"
	"time"
)

type TextBox struct {
	text string
	cursor_pos int
	size int
}

var text_box = TextBox{ size: 1 }
var chat = &datalog.TownSquare

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

//TODO: Line Wrapping and scrolling
func redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	x, y := termbox.Size()

	if len(text_box.text) + len(chat.Name)+2 < x {
		text_box.size = 1
		tbprint(0,y-1,termbox.ColorGreen, termbox.ColorDefault, chat.Name+">")
		tbprint(len(chat.Name)+2, y-1,termbox.ColorWhite, termbox.ColorDefault, text_box.text)
		termbox.SetCursor(len(chat.Name)+2+text_box.cursor_pos,y-1)
	} else {
		text_box.size = 2
		tbprint(0,y-2,termbox.ColorGreen, termbox.ColorDefault, chat.Name+">")
		tbprint(len(chat.Name)+2, y-2,termbox.ColorWhite, termbox.ColorDefault,
						text_box.text[:x-len(chat.Name)-2])
		tbprint(4, y-1,termbox.ColorWhite, termbox.ColorDefault, text_box.text[x-len(chat.Name)-2:])
		termbox.SetCursor(len(chat.Name)+6+text_box.cursor_pos-x,y-1)
	}

	l := len(chat.History) - 1
	padding := 0
	for i := l; i >= 0; i-- {
		if len(chat.History[i]) + len(chat.Authors[i]) + 2 >= x {
			padding++
			tbprint(0, y-1-(l-i)-text_box.size-padding, termbox.ColorWhite, termbox.ColorDefault,
							chat.Authors[i]+": "+chat.History[i][:x-2-len(chat.Authors[i])])
			tbprint(0, y-(l-i)-text_box.size-padding, termbox.ColorWhite, termbox.ColorDefault,
							"    "+chat.History[i][x-2-len(chat.Authors[i]):])
		} else {
			tbprint(0, y-1-(l-i)-text_box.size-padding, termbox.ColorWhite, termbox.ColorDefault,
							chat.Authors[i]+": "+chat.History[i])
		}
	}
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
	// Added this to force utf8, however it may not be necessary
	buf := make([]byte, 1)
	utf8.EncodeRune(buf, ch)
	switch {
		case len(text_box.text) == text_box.cursor_pos:
			text_box.text += string(buf)
		case text_box.cursor_pos == 0:
			text_box.text = string(buf) + text_box.text
		default:
			text_box.text = text_box.text[:text_box.cursor_pos] + string(buf) + text_box.text[text_box.cursor_pos:]
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

func handle_enter() bool {
	if text_box.text == "" {
		return false
	}
	switch text_box.text[0] {
		case '/':
			if text_box.text == "/quit" {
				return true
			}
			if text_box.text[:6] == "/enter" {
				switch text_box.text[7:] {
					case "TownSquare":
						chat = &datalog.TownSquare
					case "OffTopic":
						chat = &datalog.OffTopic
					case "TermDev":
						chat = &datalog.TermDev
				}
				chat.History = append(chat.History, "Switched channels")
				chat.Authors = append(chat.Authors, "System")
				text_box.text = ""
				text_box.cursor_pos = 0
			}
		default:
			// Can also send message to back end here
			chat.Mux.Lock()
			chat.History = append(chat.History, text_box.text)
			chat.Authors = append(chat.Authors, "You")
			chat.Mux.Unlock()
			text_box.text = ""
			text_box.cursor_pos = 0
	}
	return false
}

// TODO: This function will receive messages to be added to the history
func Receive_Message(msg, author string) {
	chat.Mux.Lock()
	chat.History = append(chat.History, msg)
	chat.Authors = append(chat.Authors, author)
	chat.Mux.Unlock()
	termbox.Interrupt()
}

func someuser1() {
	count := 0
	for {
		time.Sleep(5 * time.Second)
		Receive_Message("Hello"+strconv.Itoa(count), "Guy1")
		count++
	}
}

func someuser2() {
	count := 0
	for {
		time.Sleep(4 * time.Second)
		Receive_Message("Hello"+strconv.Itoa(count), "Guy2")
		count++
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	redraw_all()

	go someuser1()
	go someuser2()

	exit := false
mainloop:
	for !exit {
		switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
					case termbox.KeyEsc:
						break mainloop
					case termbox.KeyEnter:
						exit = handle_enter()
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
			case termbox.EventInterrupt:
		}
		redraw_all()
	}
}
