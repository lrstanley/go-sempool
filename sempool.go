// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package sempool provides a very lightweight goroutine work pool
//
// Here is an example, which creates no more than 5 go-routines at a time. Once
// the as each go-routine finishes, a slot is opened, which is then occupied by
// the next available piece of work:
//
// 	concurrency := 5
// 	pool := sempool.New(concurrency)
// 	urls := []string{"url1", "url2"}
// 	for _, url := range urls {
// 	 	pool.Slot() // wait for an open slot
// 	 	go func(url) {
// 	 	 	defer pool.Free() // free the slot we're occupying
//
// 	 	 	// get url or other stuff
// 	 	}(url)
// 	}
//
// 	pool.Wait()
package sempool // import "github.com/lrstanley/go-sempool"

// Pool represents a go-routine worker pool. This does NOT manage the
// workers, only how many workers are running.
type Pool struct {
	total   int
	threads chan bool
	done    bool
}

// Slot is used to wait for an open slot to start processing
func (p *Pool) Slot() {
	if p.done {
		panic("Slot() called in go-routine on completed pool")
	}

	p.threads <- true
}

// Free is used to free the slot taken by Pool.Slot()
func (p *Pool) Free() {
	if p.done {
		panic("Free() called in go-routine on completed pool")
	}

	<-p.threads
}

// Wait is used to wait for all open Slot()'s to be Free()'d
func (p *Pool) Wait() {
	if p.done {
		panic("Wait() called on completed pool")
	}

	for i := 0; i < cap(p.threads); i++ {
		p.threads <- true
	}

	p.done = true
}

// WaitChan returns a channel that can be used to wait for a response to a
// channel.
func (p *Pool) WaitChan() chan struct{} {
	notify := make(chan struct{}, 1)

	go func() {
		p.Wait()

		notify <- struct{}{}
	}()

	return notify
}

// New returns a new Pool{} method
func New(count int) Pool {
	if count < 1 {
		count = 1
	}

	return Pool{total: count, threads: make(chan bool, count)}
}
