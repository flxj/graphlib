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

type generalTask struct {
	name    string
	inputs  map[string]Parameter
	outputs map[string]Parameter
	runner  func(map[string]Parameter) ([]Parameter, error)
}

func NewGeneralTask(name string, runner func(map[string]Parameter) ([]Parameter, error)) Task {
	return &generalTask{
		name:    name,
		inputs:  make(map[string]Parameter),
		outputs: make(map[string]Parameter),
		runner:  runner,
	}
}

func (t *generalTask) Name() string {
	return t.name
}

func (t *generalTask) Run() error {
	out, err := t.runner(t.inputs)
	if err != nil {
		return err
	}
	for _, p := range out {
		t.outputs[p.Name] = p
	}
	return nil
}

func (t *generalTask) Input(ps []Parameter) error {
	for _, p := range ps {
		t.inputs[p.Name] = p
	}
	return nil
}

func (t *generalTask) Output() ([]Parameter, error) {
	var ps []Parameter
	for _, p := range t.outputs {
		ps = append(ps, p)
	}
	return ps, nil
}
