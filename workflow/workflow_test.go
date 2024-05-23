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

package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// go test -run TestWF -timeout 10m
func TestWF(t *testing.T) {
	wf, err := NewWorkflow(WfNameOption("test"))
	if err != nil {
		fmt.Println("craete wf err ", err)
		return
	}

	g := `
t1---> t2 ---+
             |
             v
t3--->t4 --->t5--->t6
`
	//  n-> t1 --> n+1 --> t2 --> (n+1)*2 --> t5
	//
	//  m -> t3 --> m*3 -->t4 --> 3*m+4 --> t5 --> (n+1)*2 + 3*m+4 --> t6 --output
	//
	// 2n+3m+6

	t1 := NewGeneralTask("t1", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t1 start")
		p, ok := ps["n"]
		if !ok {
			return nil, fmt.Errorf("not found input n")
		}
		n := p.Value.(int)
		out := Parameter{
			Name:  "n",
			Value: n + 1,
		}
		time.Sleep(2 * time.Second)
		fmt.Println("task t1 end")
		return []Parameter{out}, nil
	})

	t2 := NewGeneralTask("t2", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t2 start")
		p, ok := ps["n"]
		if !ok {
			return nil, fmt.Errorf("not found input n")
		}
		n1 := p.Value.(int)
		out := Parameter{
			Name:  "n",
			Value: n1 * 2,
		}
		time.Sleep(2 * time.Second)
		fmt.Println("task t2 end")
		return []Parameter{out}, nil
	})

	t3 := NewGeneralTask("t3", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t3 start")
		p, ok := ps["m"]
		if !ok {
			return nil, fmt.Errorf("not found input m")
		}
		m := p.Value.(int)
		out := Parameter{
			Name:  "m",
			Value: m * 3,
		}
		time.Sleep(2 * time.Second)
		fmt.Println("task t3 end")
		return []Parameter{out}, nil
	})

	t4 := NewGeneralTask("t4", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t4 start")
		p, ok := ps["m"]
		if !ok {
			return nil, fmt.Errorf("not found input m")
		}
		m := p.Value.(int)
		out := Parameter{
			Name:  "m",
			Value: m + 4,
		}
		time.Sleep(2 * time.Second)
		fmt.Println("task t4 end")
		return []Parameter{out}, nil
	})

	t5 := NewGeneralTask("t5", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t5 start")
		pn, ok := ps["n"]
		if !ok {
			return nil, fmt.Errorf("not found input n")
		}
		n := pn.Value.(int)

		pm, ok := ps["m"]
		if !ok {
			return nil, fmt.Errorf("not found input m")
		}
		m := pm.Value.(int)

		out := Parameter{
			Name:  "res",
			Value: m + n,
		}
		time.Sleep(2 * time.Second)
		fmt.Println("task t5 end")
		return []Parameter{out}, nil
	})

	t6 := NewGeneralTask("t6", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t6 start")
		s, ok := ps["sum"]
		if !ok {
			return nil, fmt.Errorf("not found input sum")
		}
		sum := s.Value.(int)

		c, ok := ps["check"]
		if !ok {
			return nil, fmt.Errorf("not found input check")
		}
		chk := c.Value.(int)

		if chk != sum {
			return nil, fmt.Errorf("get sum %d, but expect %d", sum, chk)
		}
		fmt.Println("task t6 end")
		return []Parameter{}, nil
	})

	_ = wf.AddTask(t1, t2, t3, t4, t5, t6)
	//
	_ = wf.AddDependency(t1.Name(), t2.Name())
	_ = wf.AddDependency(t2.Name(), t5.Name())
	_ = wf.AddDependency(t3.Name(), t4.Name())
	_ = wf.AddDependency(t4.Name(), t5.Name())
	_ = wf.AddDependency(t5.Name(), t6.Name())
	//
	n := 100
	m := 200
	chk := 2*n + 3*m + 6

	_ = wf.SetInput(t1.Name(), &Parameter{Name: "n", Value: n})
	_ = wf.SetInput(t3.Name(), &Parameter{Name: "m", Value: m})
	_ = wf.SetInput(t2.Name(), &Parameter{Name: "n", Ref: fmt.Sprintf("%s.%s.output.n", wf.Name(), t1.Name())})
	_ = wf.SetInput(t4.Name(), &Parameter{Name: "m", Ref: fmt.Sprintf("%s.%s.output.m", wf.Name(), t3.Name())})
	_ = wf.SetInput(t5.Name(), []*Parameter{
		{Name: "m", Ref: fmt.Sprintf("%s.%s.output.m", wf.Name(), t4.Name())},
		{Name: "n", Ref: fmt.Sprintf("%s.%s.output.n", wf.Name(), t2.Name())},
	}...)
	_ = wf.SetInput(t6.Name(), []*Parameter{
		{Name: "sum", Ref: fmt.Sprintf("%s.%s.output.res", wf.Name(), t5.Name())},
		{Name: "check", Value: chk},
	}...)

	_ = wf.SetOutput(t1.Name(), &Parameter{Name: "n"})
	_ = wf.SetOutput(t3.Name(), &Parameter{Name: "m"})
	_ = wf.SetOutput(t2.Name(), &Parameter{Name: "n"})
	_ = wf.SetOutput(t4.Name(), &Parameter{Name: "m"})
	_ = wf.SetOutput(t5.Name(), &Parameter{Name: "res"})

	fmt.Println(g)
	fmt.Printf("start to run g. n=%d,m=%d,chk=%d\n", n, m, chk)

	if err := wf.Start(); err != nil {
		fmt.Println("[Err] start wf err ", err)
		return
	}

	fmt.Println("Waiting...")
	time.Sleep(10 * time.Second)
	fmt.Println("Completed")

	info, _ := wf.Info()

	bs, _ := json.Marshal(info)
	var str bytes.Buffer
	_ = json.Indent(&str, bs, "", "    ")
	fmt.Println(str.String())
}
