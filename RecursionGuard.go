package gorecursionguard

import (
	"fmt"
	"log"

	"github.com/AnimusPEXUS/goreentrantlock"
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
	mutex              *goreentrantlock.ReentrantMutexCheckable
	mode               RGMode
	cb_if_being_called func(RGMode) RGMode
}

func NewRecursionGuard(
	mode RGMode,

	// if not nil, then cb_if_being_called is called in case of recursive call detection,
	// return forces mode change
	cb_if_being_called func(RGMode) RGMode,
) *RecursionGuard {
	self := new(RecursionGuard)
	self.mutex = goreentrantlock.NewReentrantMutexCheckable(false)
	self.mode = mode
	self.cb_if_being_called = cb_if_being_called
	return self
}

func (self *RecursionGuard) Do(fn func()) {

	// NOTE: usual golang sync.Mutex can't be used here, as it's not reentrant
	// 	@valyala
	// https://github.com/golang/go/issues/24192#issuecomment-369606420
	// rfyiamcool, reentrant locks usually mean bad code. See https://stackoverflow.com/a/14671462
	// реентрант локи, наверное, дебилы придумали.. куда им до Санька go-программиста
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
