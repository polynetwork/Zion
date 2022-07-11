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

	type task struct {
		id      int
		randnum int
		halt    chan struct{}
		done    chan struct{}
	}
	type signal struct {
		id      int
		randNum int
	}

	sigCh := make(chan *signal, 100)
	taskCh := make(chan *task, 10)
	stop := make(chan struct{})

	process := func(task *task) {
		t.Logf("start the %d task, rand number %d", task.id, task.randnum)
		done := make(chan struct{})
		timer := time.NewTimer(1 * time.Second)

		defer func() {
			timer.Stop()
			t.Log("-------------------------------------")
			task.done <- struct{}{}
		}()

		time.AfterFunc(100*time.Second, func() {
			close(done)
		})
		n := 1
		for {
			select {
			case <-timer.C:
				t.Log("number", n)
				timer.Reset(1 * time.Second)
				n += 1
			case <-task.halt:
				t.Log("halt")
				return
			case <-done:
				t.Log("done")
				return
			case <-stop:
				t.Log("system stopped!")
				return
			}
		}
	}

	go func() {
		for {
			select {
			case sig := <-sigCh:
				if len(taskCh) > 0 {
					halt := <-taskCh
					close(halt.halt)
					<-halt.done
				}
				task := &task{
					id:      sig.id,
					randnum: sig.randNum,
					halt:    make(chan struct{}),
					done:    make(chan struct{}),
				}
				taskCh <- task
				go process(task)

			case <-stop:
				t.Log("system stopped!")
				return
			}
		}
	}()

	for i := 0; i < 50; i++ {
		num := rand.Intn(6)
		sig := &signal{
			id:      i,
			randNum: 3 + num,
		}
		sigCh <- sig
		time.Sleep(time.Duration(sig.randNum) * time.Second)
	}
}
