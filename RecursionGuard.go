package gorecursionguard

import (
	"fmt"
	"log"

	"github.com/AnimusPEXUS/goreentrantlock"
)

const RG_MSG = "recursionGuard: recursion detected"

type RGMode uint

const (
	RGM_Panic        RGMode = iota //panic()
	RGM_LogPrint                   // log.Println + allow recursion
	RGM_FmtPrint                   // fmt.Println + allow recursion
	RGM_SilentReturn               // cancel recursion
	RGM_SilentPass                 // allow recursion
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
		case RGM_Panic:
			panic(RG_MSG)
		case RGM_FmtPrint:
			fmt.Println(RG_MSG)
		case RGM_LogPrint:
			log.Println(RG_MSG)
		case RGM_SilentReturn:
			return
		case RGM_SilentPass:
			break
		}
	}
	self.being_called = true
	defer func() { self.being_called = false }()
	fn()
}
