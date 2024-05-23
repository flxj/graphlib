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

// Shell Task
type Command struct {
	Cmd     string
	Shell   string
	Args    []string
	Remote  bool
	SSHost  string
	SSHPort int
	SSHPwd  string
	SSHUser string
}

func NewShellTask(name string, cmd *Command) Task {
	return &shellTask{
		name:    name,
		cmd:     cmd,
		inputs:  make(map[string]Parameter),
		outputs: make(map[string]Parameter),
	}
}

type shellTask struct {
	name string
	cmd  *Command

	inputs  map[string]Parameter
	outputs map[string]Parameter
}

func (s *shellTask) Name() string {
	return s.name
}

func (s *shellTask) Run() error {
	return nil
}

func (s *shellTask) Input(ps []Parameter) error {
	for _, p := range ps {
		s.inputs[p.Name] = p
	}
	return nil
}

func (s *shellTask) Output() ([]Parameter, error) {
	var ps []Parameter
	for _, p := range s.outputs {
		ps = append(ps, p)
	}
	return ps, nil
}
