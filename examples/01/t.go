package main

import (
	"log"

	"github.com/AnimusPEXUS/gorecursionguard"
)

func main() {
	var counter int = 5
	rg := gorecursionguard.NewRecursionGuard(
		gorecursionguard.RGM_SilentPass,
		func(m gorecursionguard.RGMode) gorecursionguard.RGMode {
			log.Println("recursion detected:", counter)
			counter--
			if counter == 0 {
				counter = 0
				log.Println("counter == 0. terminating")
				return gorecursionguard.RGM_SilentReturn

			}
			return m
		},
	)

	var rf func()

	rf = func() {
		rg.Do(rf)
	}

	rf()
}
