package main

import (
	"fmt"
	"sync"
)

func tileWorker(t *Tile, wg *sync.WaitGroup) {
	defer wg.Done()
	t.resolveCombat()
	fmt.Println("Worker started with ", t.toString())
}
