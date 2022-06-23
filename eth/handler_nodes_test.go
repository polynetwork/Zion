/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package eth

import (
	"math/rand"
	"testing"
	"time"
)

func TestGoroutineManage(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	type Halt struct {
		Waiting chan struct{}
		Done    chan struct{}
	}
	type Signal struct {
		randNum int
	}

	sigCh := make(chan *Signal, 100)
	haltCh := make(chan *Halt, 10)
	stop := make(chan struct{})
	round := int64(0)

	process := func(halt *Halt, randNum int) {
		cur := round
		round += 1

		t.Log(cur, "rand number", randNum)
		done := make(chan struct{})
		timer := time.NewTimer(1 * time.Second)

		defer func() {
			timer.Stop()
			t.Log("-------------------------------------")
			halt.Done <- struct{}{}
		}()

		time.AfterFunc(5*time.Second, func() {
			close(done)
		})
		n := 0

		for {
			select {
			case <-timer.C:
				t.Log(cur, "number", n)
				timer.Reset(1 * time.Second)
				n += 1
			case <-halt.Waiting:
				t.Log(cur, "halt")
				return
			case <-done:
				t.Log(cur, "done")
				return
			case <-stop:
				t.Log(cur, "system stopped!")
				return
			}
		}
	}

	go func() {
		for {
			select {
			case sig := <-sigCh:
				if round > 0 {
					halt := <-haltCh
					close(halt.Waiting)
					<-halt.Done
				}
				halt := &Halt{
					Waiting: make(chan struct{}),
					Done:    make(chan struct{}),
				}
				haltCh <- halt
				go process(halt, sig.randNum)

			case <-stop:
				t.Log("system stopped!")
				return
			}
		}
	}()

	for i := 0; i < 50; i++ {
		num := rand.Intn(8)
		sigCh <- &Signal{
			randNum: num,
		}
		time.Sleep(time.Duration(num) * time.Second)
	}
}
