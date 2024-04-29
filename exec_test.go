/*
Copyright (C) 2023 flxj(https://github.com/flxj)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graphlib

import (
	//"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestExecJob1(t *testing.T) {
	g, err := NewExecGraph[int, Job]("exec")
	if err != nil {
		fmt.Printf("[ERR] create exec graph error: %v\n", err)
		return
	}
	var (
		v1 int
		v2 int
		v3 int
	)
	// input:  v1 <- x, v2 <- y
	// output: v3 <- 2*(x+100) + 3*x-10

	job1 := func() error {
		fmt.Println("job1 start...")
		v1 += 100
		time.Sleep(5 * time.Second)
		fmt.Println("job1 completed")
		return nil
	}
	job2 := func() error {
		fmt.Println("job2 start...")
		v1 = 2 * v1
		time.Sleep(2 * time.Second)
		fmt.Println("job2 completed")
		return nil
	}
	job3 := func() error {
		fmt.Println("job3 start...")
		v2 = 3 * v2
		time.Sleep(time.Second)
		fmt.Println("job3 completed")
		return nil
	}
	job4 := func() error {
		fmt.Println("job4 start...")
		v2 = v2 - 10
		time.Sleep(3 * time.Second)
		fmt.Println("job4 completed")
		return nil
	}
	job5 := func() error {
		fmt.Println("job5 start...")
		v3 = v1 + v2
		time.Sleep(2 * time.Second)
		fmt.Println("job6 completed")
		return nil
	}

	jobs := map[int]Job{
		1: job1,
		2: job2,
		3: job3,
		4: job4,
		5: job5,
	}

	for k, j := range jobs {
		if err := g.AddJob(k, j); err != nil {
			fmt.Printf("[ERR] add job error: %v\n", err)
			return
		}
	}

	deps := [][]int{
		{1, 2},
		{3, 4},
		{2, 5},
		{4, 5},
	}
	for _, d := range deps {
		if err := g.AddDependency(d[0], d[1]); err != nil {
			fmt.Printf("[ERR] add dep error: %v\n", err)
			return
		}
	}

	v1 = 100
	v2 = 200

	var val = 2*(v1+100) + 3*v2 - 10

	gs := `
job1---->job2--.
                \
                 \
                  V
job3---->job4--->job5
`

	fmt.Println(gs)
	fmt.Println("expectation result is: ", val)

	fmt.Println("exec graph status=>", g.Status())

	if err := g.Start(); err != nil {
		fmt.Printf("[ERR] start graph error: %v\n", err)
		return
	}
	fmt.Println("exec graph status=>", g.Status())

	if err := g.Wait(); err != nil {
		fmt.Printf("[ERR] wait graph error: %v\n", err)
		fmt.Println("exec graph status=>", g.Status())
		_ = g.Stop()
		return
	}

	fmt.Println("exec graph status=>", g.Status())

	if v3 != val {
		fmt.Printf("exec err: expect %d, actual get %d\n", val, v3)
	} else {
		fmt.Println("success")
	}

	fmt.Println("============== jobs status")
	for _, jb := range g.AllJobs() {
		fmt.Printf("job %d status is %s\n", jb.Key, jb.Status)
	}
}

func TestExecJob2(t *testing.T) {
	g, err := NewExecGraph[int, Job]("exec")
	if err != nil {
		fmt.Printf("[ERR] create exec graph error: %v\n", err)
		return
	}
	var (
		v1 int
		v2 int
		v3 int
	)
	// input:  v1 <- x, v2 <- y
	// output: v3 <- 2*(x+100) + 3*x-10

	job1 := func() error {
		fmt.Println("job1 start...")
		v1 += 100
		time.Sleep(5 * time.Second)
		fmt.Println("job1 completed")
		return nil
	}
	job2 := func() error {
		fmt.Println("job2 start...")
		time.Sleep(2 * time.Second)
		fmt.Println("job2 err")
		return errors.New("job failed")
	}
	job3 := func() error {
		fmt.Println("job3 start...")
		v2 = 3 * v2
		time.Sleep(time.Second)
		fmt.Println("job3 completed")
		return nil
	}
	job4 := func() error {
		fmt.Println("job4 start...")
		v2 = v2 - 10
		time.Sleep(3 * time.Second)
		fmt.Println("job4 completed")
		return nil
	}
	job5 := func() error {
		fmt.Println("job5 start...")
		v3 = v1 + v2
		time.Sleep(2 * time.Second)
		fmt.Println("job6 completed")
		return nil
	}

	jobs := map[int]Job{
		1: job1,
		2: job2,
		3: job3,
		4: job4,
		5: job5,
	}

	for k, j := range jobs {
		if err := g.AddJob(k, j); err != nil {
			fmt.Printf("[ERR] add job error: %v\n", err)
			return
		}
	}

	deps := [][]int{
		{1, 2},
		{3, 4},
		{2, 5},
		{4, 5},
	}
	for _, d := range deps {
		if err := g.AddDependency(d[0], d[1]); err != nil {
			fmt.Printf("[ERR] add dep error: %v\n", err)
			return
		}
	}

	v1 = 100
	v2 = 200

	var val = 2*(v1+100) + 3*v2 - 10

	gs := `
job1---->job2--.
                \
                 \
                  V
job3---->job4--->job5
`
	fmt.Println(gs)
	fmt.Println("expectation result is: ", val)

	fmt.Println("exec graph status=>", g.Status())

	if err := g.Start(); err != nil {
		fmt.Printf("[ERR] start graph error: %v\n", err)
		return
	}
	fmt.Println("exec graph status=>", g.Status())

	if err := g.Wait(); err != nil {
		fmt.Printf("[ERR] wait graph error: %v\n", err)
		fmt.Println("exec graph status=>", g.Status())
	}

	fmt.Println("exec graph status=>", g.Status())

	if v3 != val {
		fmt.Printf("exec err: expect %d, actual get %d\n", val, v3)
	} else {
		fmt.Println("success")
	}

	fmt.Println("============== jobs status")
	for _, jb := range g.AllJobs() {
		fmt.Printf("job %d status is %s\n", jb.Key, jb.Status)
	}
}

// stop
func TestExecJob3(t *testing.T) {
	g, err := NewExecGraph[int, Job]("exec")
	if err != nil {
		fmt.Printf("[ERR] create exec graph error: %v\n", err)
		return
	}
	job1 := func() error {
		fmt.Println("job1 start...")
		time.Sleep(5 * time.Second)
		fmt.Println("job1 completed")
		return nil
	}
	job2 := func() error {
		fmt.Println("job2 start...")
		time.Sleep(2 * time.Second)
		fmt.Println("job2 completed")
		return nil
	}
	job3 := func() error {
		fmt.Println("job3 start...")
		time.Sleep(time.Second)
		fmt.Println("job3 completed")
		return nil
	}
	job4 := func() error {
		fmt.Println("job4 start...")
		time.Sleep(3 * time.Second)
		fmt.Println("job4 completed")
		return nil
	}
	job5 := func() error {
		fmt.Println("job5 start...")
		time.Sleep(2 * time.Second)
		fmt.Println("job6 completed")
		return nil
	}

	jobs := map[int]Job{
		1: job1,
		2: job2,
		3: job3,
		4: job4,
		5: job5,
	}

	for k, j := range jobs {
		if err := g.AddJob(k, j); err != nil {
			fmt.Printf("[ERR] add job error: %v\n", err)
			return
		}
	}

	deps := [][]int{
		{1, 2},
		{3, 4},
		{2, 5},
		{4, 5},
	}
	for _, d := range deps {
		if err := g.AddDependency(d[0], d[1]); err != nil {
			fmt.Printf("[ERR] add dep error: %v\n", err)
			return
		}
	}

	gs := `
job1---->job2--.
                \
                 \
                  V
job3---->job4--->job5
`
	fmt.Println(gs)
	fmt.Println("==============> start")
	if err := g.Start(); err != nil {
		fmt.Printf("[ERR] start graph error: %v\n", err)
		return
	}
	fmt.Println("exec graph status==>", g.Status())
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("==============> stop")
		if err := g.Stop(); err != nil {
			fmt.Printf("[ERR] stop graph error: %v\n", err)
		}
	}()
	fmt.Println("==============> waiting...")
	if err := g.Wait(); err != nil {
		fmt.Printf("[ERR] wait graph error: %v\n", err)
		fmt.Println("exec graph status=>", g.Status())
	}

	fmt.Println("exec graph status==>", g.Status())

	fmt.Println("============== jobs status")
	for _, jb := range g.AllJobs() {
		fmt.Printf("job %d status is %s\n", jb.Key, jb.Status)
	}
}

// paused
func TestExecJob4(t *testing.T) {
	g, err := NewExecGraph[int, Job]("exec")
	if err != nil {
		fmt.Printf("[ERR] create exec graph error: %v\n", err)
		return
	}
	job1 := func() error {
		fmt.Println("job1 start...")
		time.Sleep(5 * time.Second)
		fmt.Println("job1 completed")
		return nil
	}
	job2 := func() error {
		fmt.Println("job2 start...")
		time.Sleep(2 * time.Second)
		fmt.Println("job2 completed")
		return nil
	}
	job3 := func() error {
		fmt.Println("job3 start...")
		time.Sleep(time.Second)
		fmt.Println("job3 completed")
		return nil
	}
	job4 := func() error {
		fmt.Println("job4 start...")
		time.Sleep(3 * time.Second)
		fmt.Println("job4 completed")
		return nil
	}
	job5 := func() error {
		fmt.Println("job5 start...")
		time.Sleep(2 * time.Second)
		fmt.Println("job6 completed")
		return nil
	}
	job6 := func() error {
		fmt.Println("job6 start...")
		time.Sleep(2 * time.Second)
		fmt.Println("job6 completed")
		return nil
	}

	jobs := map[int]Job{
		1: job1,
		2: job2,
		3: job3,
		4: job4,
		5: job5,
		6: job6,
	}

	for k, j := range jobs {
		if err := g.AddJob(k, j); err != nil {
			fmt.Printf("[ERR] add job error: %v\n", err)
			return
		}
	}

	deps := [][]int{
		{1, 2},
		{3, 4},
		{2, 5},
		{4, 5},
		{5, 6},
	}
	for _, d := range deps {
		if err := g.AddDependency(d[0], d[1]); err != nil {
			fmt.Printf("[ERR] add dep error: %v\n", err)
			return
		}
	}

	gs := `
job1---->job2--.
                \
                 \
                  V
job3---->job4--->job5---->job6
`
	fmt.Println(gs)
	fmt.Println("==============> start")
	if err := g.Start(); err != nil {
		fmt.Printf("[ERR] start graph error: %v\n", err)
		return
	}
	fmt.Println("exec graph status==>", g.Status())
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("==============> pause 20s")
		if err := g.Pause(); err != nil {
			fmt.Printf("[ERR] stop graph error: %v\n", err)
			return
		}
		fmt.Println("exec graph status==>", g.Status())
		time.Sleep(20 * time.Second)
		fmt.Println("==============> continue running")
		if err := g.Start(); err != nil {
			fmt.Printf("[ERR] continue start graph error: %v\n", err)
			return
		}
		fmt.Println("exec graph status==>", g.Status())
	}()
	fmt.Println("==============> waiting...")
	if err := g.Wait(); err != nil {
		fmt.Printf("[ERR] wait graph error: %v\n", err)
		fmt.Println("exec graph status=>", g.Status())
	}

	fmt.Println("exec graph status==>", g.Status())

	fmt.Println("============== jobs status")
	for _, jb := range g.AllJobs() {
		fmt.Printf("job %d status is %s\n", jb.Key, jb.Status)
	}
}
