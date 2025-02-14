package main

import (
	"fmt"
	"sync"
)

func worker(t *Tile, wg *sync.WaitGroup) {
	defer wg.Done()
	t.resolveCombat()
	fmt.Println("Worker started with ", t.toString())
}
