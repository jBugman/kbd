package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	reghotkey := user32.MustFindProc("RegisterHotKey")
	getmsg := user32.MustFindProc("GetMessageW")

	reg := reg(reghotkey)

	keys := map[int]int{
		1: '2',
		2: '3',
		3: '4',
	}

	for id, code := range keys {
		if err := reg(id, code); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	for {
		var msg message
		getmsg.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 1)

		if id := msg.WPARAM; id != 0 {
			fmt.Printf("pressed: %c\n", keys[int(id)])
			switch id {
			case 2, 3, 4:
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func reg(reghotkey *syscall.Proc) func(int, int) error {
	return func(id int, code int) error {
		r1, _, err := reghotkey.Call(
			0, uintptr(id), 0, uintptr(code))
		if r1 == 1 {
			fmt.Printf("registered: %c\n", code)
			return nil
		} else {
			return fmt.Errorf("failed to register: %w", err)
		}
	}
}

type hotkey struct {
	id   int
	code int
}

type message struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}
