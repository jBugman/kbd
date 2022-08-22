package main

// spellchecker:ignore Keydb, Keybd
// spellchecker:words LLKHF

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/stephen-fox/user32util"
)

const f7 = 0x76
const f8 = 0x77

const LLKHF_INJECTED = 0x00000010

func main() {
	dll, err := user32util.LoadUser32DLL()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer dll.Release()

	rand.Seed(time.Now().UnixNano())

	on := true

	fn := func(event user32util.LowLevelKeyboardEvent) {
		if event.Struct.Flags&LLKHF_INJECTED > 0 {
			return
		}
		switch event.KeyboardButtonAction() {
		case user32util.WMKeyUp:
			switch event.Struct.VkCode {
			case f8:
				on = !on

			case f7:
				fmt.Println("exiting")
				os.Exit(0)

			case '2':
				if on {
					go send('3', rng(0, 25), rng(25, 75), dll)
					go send('4', rng(0, 25), rng(25, 75), dll)
				}
			case '3':
				if on {
					go send('2', rng(0, 25), rng(25, 75), dll)
					go send('4', rng(0, 25), rng(25, 75), dll)
				}
			case '4':
				if on {
					go send('2', rng(0, 25), rng(25, 75), dll)
					go send('3', rng(0, 25), rng(25, 75), dll)
				}
			}
		}
	}

	listener, err := user32util.NewLowLevelKeyboardListener(fn, dll)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	select {
	case err := <-listener.OnDone():
		fmt.Fprintln(os.Stderr, err)
	case <-interrupts:
	}
}

func send(key uint16, down time.Duration, up time.Duration, dll *user32util.User32DLL) {
	time.Sleep(down)
	if err := user32util.SendKeydbInput(user32util.KeybdInput{
		WVK: key,
	}, dll); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	time.Sleep(up)
	if err := user32util.SendKeydbInput(user32util.KeybdInput{
		WVK:     key,
		DwFlags: user32util.KeyEventFKeyUp,
	}, dll); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func rng(from int, to int) time.Duration {
	return time.Duration(
		(from + rand.Intn(to-from+1)) * int(time.Millisecond),
	)
}
