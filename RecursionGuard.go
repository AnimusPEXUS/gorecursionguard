package gorecursionguard

import (
	"fmt"
	"log"
	"sync"
)

const RG_MSG = "recursionGuard: recursion detected"

type RGMode uint

const (
	Panic        RGMode = iota //panic()
	LogPrint                   // log.Println + allow recursion
	FmtPrint                   // fmt.Println + allow recursion
	SilentReturn               // cancel recursion
	SilentPass                 // allow recursion
)

type RecursionGuard struct {
	being_called       bool
	mutex              *sync.Mutex
	mode               RGMode
	cb_if_being_called func(RGMode) RGMode
}

func NewRecursionGuard(
	mode RGMode,

	// if not nil called in case of recursive call detected, return forces mode change
	cb_if_being_called func(RGMode) RGMode,
) *RecursionGuard {
	self := new(RecursionGuard)
	self.mutex = new(sync.Mutex)
	self.mode = mode
	self.cb_if_being_called = cb_if_being_called
	return self
}

func (self *RecursionGuard) Do(fn func()) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	mode := self.mode

	if self.being_called {
		if self.cb_if_being_called != nil {
			mode = self.cb_if_being_called(self.mode)
		}
		switch mode {
		case Panic:
			panic(RG_MSG)
		case FmtPrint:
			fmt.Println(RG_MSG)
		case LogPrint:
			log.Println(RG_MSG)
		case SilentReturn:
			return
		case SilentPass:
			break
		}
	}
	self.being_called = true
	fn()
}
