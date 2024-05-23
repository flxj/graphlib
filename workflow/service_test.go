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
	"fmt"
	"testing"
	"time"
)

// go test -run TestSvc -timeout 10m
func TestSvc(t *testing.T) {
	wf, err := NewWorkflow(WfNameOption("test"))
	if err != nil {
		fmt.Println("craete wf err ", err)
		return
	}

	/*
		t1 ----> t2----+
		  \            |
		   \           |
			v          v
		t3-->t4------->t5--->t6
	*/

	t1 := NewGeneralTask("t1", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t1 start")
		time.Sleep(1 * time.Second)
		fmt.Println("task t1 end")
		return []Parameter{}, nil
	})

	t2 := NewGeneralTask("t2", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t2 start")
		time.Sleep(2 * time.Second)
		fmt.Println("task t2 end")
		return []Parameter{}, nil
	})

	t3 := NewGeneralTask("t3", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t3 start")
		time.Sleep(3 * time.Second)
		fmt.Println("task t3 end")
		return []Parameter{}, nil
	})

	t4 := NewGeneralTask("t4", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t4 start")
		time.Sleep(4 * time.Second)
		fmt.Println("task t4 end")
		return []Parameter{}, nil
	})

	t5 := NewGeneralTask("t5", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t5 start")
		time.Sleep(5 * time.Second)
		fmt.Println("task t5 end")
		return []Parameter{}, nil
	})

	t6 := NewGeneralTask("t6", func(ps map[string]Parameter) ([]Parameter, error) {
		fmt.Println("task t6 start")
		time.Sleep(time.Second)
		fmt.Println("task t6 end")
		return []Parameter{}, nil
	})

	_ = wf.AddTask(t1, t2, t3, t4, t5, t6)
	_ = wf.AddDependency(t1.Name(), t2.Name())
	_ = wf.AddDependency(t1.Name(), t4.Name())
	_ = wf.AddDependency(t2.Name(), t5.Name())
	_ = wf.AddDependency(t3.Name(), t4.Name())
	_ = wf.AddDependency(t4.Name(), t5.Name())
	_ = wf.AddDependency(t5.Name(), t6.Name())
	//

	port := 8080
	svc := NewService("localhost", port)

	svc.Register(wf)

	if err := svc.Run(); err != nil {
		fmt.Println("[Err] start svc err ", err)
		return
	}

	time.Sleep(20 * time.Second)

}
