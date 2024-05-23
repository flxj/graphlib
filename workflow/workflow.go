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
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/flxj/graphlib"
)

var (
	errWFModify = errors.New("workflow is running,not support dynamic modification")
	errRef      = errors.New("unknown ref format of parameter")
)

type Parameter struct {
	// paraneter name
	Name string
	//
	// format is 'workflow.task.input|output.parameter',
	// for example "myWorkflow.task1.output.x","myWorkflow.task1.input.y"
	Ref string
	// value
	Value interface{}
	// path to file
	//saveTo string
}

// workflow node.
type Task interface {
	// task instance name.
	Name() string
	// run task.
	Run() error
	// set input parameters.
	Input([]Parameter) error
	// get output patameters.
	Output() ([]Parameter, error)
}

type TaskInfo struct {
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	Err       string          `json:"err"`
	StartAt   string          `json:"start_at"`
	EndAt     string          `json:"end_at"`
	Successor map[string]bool `json:"successor"`
	Precursor map[string]bool `json:"precursor"`
}

type WorkflowInfo struct {
	Name    string     `json:"name"`
	Corn    string     `json:"corn"`
	Status  string     `json:"status"`
	StartAt string     `json:"start_at"`
	EndAt   string     `json:"end_at"`
	Err     string     `json:"err"`
	Tasks   []TaskInfo `json:"tasks"`
}

type WorkflowOption func(*Workflow) error

func WfNameOption(name string) WorkflowOption {
	return func(wf *Workflow) error {
		wf.name = name
		return nil
	}
}

func WfCornOption(corn string) WorkflowOption {
	return func(wf *Workflow) error {
		wf.corn = corn
		return nil
	}
}

// create a new workflow.
func NewWorkflow(ops ...WorkflowOption) (*Workflow, error) {
	wf := &Workflow{
		status: graphlib.Waiting,
		steps:  make(map[string]*step),
		infos:  make(map[string]*TaskInfo),
	}
	for _, op := range ops {
		if err := op(wf); err != nil {
			return nil, err
		}
	}

	return wf, nil
}

type Workflow struct {
	name string
	corn string

	mu      sync.RWMutex
	count   int
	status  graphlib.State
	err     string
	startAt time.Time
	endAt   time.Time
	steps   map[string]*step
	infos   map[string]*TaskInfo

	//waitCh chan struct{}
	eg graphlib.ExecGraph[string, graphlib.Job]
}

func (w *Workflow) Name() string {
	return w.name
}

func (w *Workflow) Corn() string {
	return w.corn
}

func (w *Workflow) Status() string {
	return string(w.status)
}

// run the workflow async.
func (w *Workflow) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return nil
	}

	switch w.status {
	case graphlib.Running:
		return nil
	default:
	}

	w.status = graphlib.Running
	w.startAt = time.Now()
	w.err = ""

	var err error
	defer func() {
		if err != nil {
			w.status = graphlib.Failed
			w.err = err.Error()
			w.endAt = time.Now()
			return
		}
	}()

	for _, tk := range w.infos {
		tk.Status = string(graphlib.Waiting)
		tk.StartAt = ""
		tk.EndAt = ""
		tk.Err = ""
	}
	if w.eg != nil {
		_ = w.eg.Stop()
		w.eg = nil
	}
	w.eg, err = graphlib.NewExecGraph[string, graphlib.Job](w.name)
	if err != nil {
		return err
	}
	for name, tk := range w.steps {
		if err = w.eg.AddJob(name, tk.job); err != nil {
			return err
		}
	}
	for _, tk := range w.infos {
		for s := range tk.Successor {
			if err = w.eg.AddDependency(tk.Name, s); err != nil {
				return err
			}
		}
	}

	// TODO: clean all outputs value

	w.count++
	if err = w.eg.Start(); err != nil {
		return err
	}
	return nil
}

// stop the workflow.
func (w *Workflow) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != graphlib.Running {
		return nil
	}
	if w.eg != nil {
		_ = w.eg.Stop()
	}
	w.status = graphlib.Stopped
	w.err = "canceled"

	return nil
}

// query workflow info.
func (w *Workflow) Info() (*WorkflowInfo, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	wf := &WorkflowInfo{
		Name:    w.name,
		StartAt: w.startAt.Format(time.RFC3339),
		EndAt:   w.endAt.Format(time.RFC3339),
		Err:     w.err,
	}
	if w.status == graphlib.Running {
		w.status = w.eg.Status()
		for name, tk := range w.infos {
			if tk.Status == string(graphlib.Waiting) || tk.Status == string(graphlib.Running) {
				job, _ := w.eg.Job(name)
				tk.Status = string(job.Status)
				tk.StartAt = job.StartAt.Format(time.RFC3339)
				tk.EndAt = job.EndAt.Format(time.RFC3339)
				if job.Error != nil {
					tk.Err = job.Error.Error()
				}
				if job.EndAt.After(w.endAt) {
					w.endAt = job.EndAt
				}
			}
		}
	}
	for _, tk := range w.infos {
		task := *tk
		wf.Tasks = append(wf.Tasks, task)
	}
	wf.Status = string(w.status)

	return wf, nil
}

// add tasks to the workflow.
func (w *Workflow) AddTask(tasks ...Task) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}
	for _, task := range tasks {
		w.steps[task.Name()] = newStep(w, task)
		w.infos[task.Name()] = &TaskInfo{
			Name:   task.Name(),
			Status: string(graphlib.Waiting),
		}
	}
	return nil
}

// delete tasks from the workflow.
func (w *Workflow) RemoveTask(names ...string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	for _, name := range names {
		if _, ok := w.steps[name]; !ok {
			return fmt.Errorf("task %s not exists", name)
		}
		delete(w.infos, name)
		for _, tk := range w.infos {
			delete(tk.Successor, name)
			delete(tk.Precursor, name)
		}
	}
	return nil
}

func (w *Workflow) AddDependency(precursor, successor string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	if _, ok := w.steps[precursor]; !ok {
		return fmt.Errorf("task %s not exists", precursor)
	}
	if _, ok := w.steps[successor]; !ok {
		return fmt.Errorf("task %s not exists", successor)
	}

	p := w.infos[precursor]
	if p.Successor == nil {
		p.Successor = make(map[string]bool)
	}
	p.Successor[successor] = true

	s := w.infos[successor]
	if s.Precursor == nil {
		s.Precursor = make(map[string]bool)
	}
	s.Precursor[precursor] = true

	return nil
}

func (w *Workflow) RemoveDependency(precursor, successor string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	if _, ok := w.steps[precursor]; !ok {
		return fmt.Errorf("task %s not exists", precursor)
	}
	if _, ok := w.steps[successor]; !ok {
		return fmt.Errorf("task %s not exists", successor)
	}

	delete(w.infos[precursor].Successor, successor)
	delete(w.infos[successor].Precursor, precursor)

	return nil
}

func (w *Workflow) SetInput(task string, p ...*Parameter) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	tk, ok := w.steps[task]
	if !ok {
		return fmt.Errorf("task %s not exists", task)
	}

	tk.setInput(p)

	return nil
}

func (w *Workflow) SetOutput(task string, p ...*Parameter) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	tk, ok := w.steps[task]
	if !ok {
		return fmt.Errorf("task %s not exists", task)
	}

	tk.setOutput(p)

	return nil
}

func (w *Workflow) SetInputs(inputs map[string][]*Parameter) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	for task, p := range inputs {
		tk, ok := w.steps[task]
		if !ok {
			return fmt.Errorf("task %s not exists", task)
		}
		tk.setInput(p)
	}

	return nil
}

func (w *Workflow) SetOutputs(outputs map[string][]*Parameter) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == graphlib.Running {
		return errWFModify
	}

	for task, p := range outputs {
		tk, ok := w.steps[task]
		if !ok {
			return fmt.Errorf("task %s not exists", task)
		}

		tk.setOutput(p)
	}

	return nil
}

func (w *Workflow) GetOutput(task string) ([]Parameter, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	tk, ok := w.steps[task]
	if !ok {
		return nil, fmt.Errorf("task %s not exists", task)
	}

	return tk.getOutput()
}

type step struct {
	wf      *Workflow
	name    string
	inputs  map[string]*Parameter
	outputs map[string]*Parameter

	task Task
	job  func() error
}

func newStep(wf *Workflow, task Task) *step {
	s := &step{
		wf:   wf,
		name: task.Name(),
		task: task,
	}
	job := func() error {
		if s.task != nil {
			if s.inputs != nil {
				var ps []Parameter
				for _, p := range s.inputs {
					pp := *p
					if p.Ref != "" {
						pr, err := s.getOutputByRef(p.Ref)
						if err != nil {
							return err
						}
						pp.Value = pr.Value
					}
					ps = append(ps, pp)
				}
				if err := s.task.Input(ps); err != nil {
					return err
				}
			}
			if err := s.task.Run(); err != nil {
				return err
			}
			if s.outputs != nil {
				outs, err := s.task.Output()
				if err != nil {
					return err
				}
				for _, p := range outs {
					if _, ok := s.outputs[p.Name]; ok {
						// TODO save value to file
						s.outputs[p.Name].Value = p.Value
					}
				}
			}
		}
		return nil
	}
	s.job = job
	return s
}

func (s *step) setInput(ps []*Parameter) {
	if s.inputs == nil {
		s.inputs = make(map[string]*Parameter)
	}
	for _, p := range ps {
		s.inputs[p.Name] = p
	}
}

func (s *step) setOutput(ps []*Parameter) {
	if s.outputs == nil {
		s.outputs = make(map[string]*Parameter)
	}
	for _, p := range ps {
		s.outputs[p.Name] = p
	}
}

func (s *step) getOutput() ([]Parameter, error) {
	if s.task != nil {
		return s.task.Output()
	}
	return []Parameter{}, nil
}

func (s *step) getOutputByRef(ref string) (Parameter, error) {
	var p Parameter

	r := strings.Split(ref, ".")
	if len(r) != 4 {
		return p, errRef
	}
	if r[0] != s.wf.Name() {
		return p, fmt.Errorf("not found workflow %s", r[0])
	}
	tk, ok := s.wf.steps[r[1]]
	if !ok {
		return p, fmt.Errorf("not found task %s", r[1])
	}
	switch strings.ToLower(r[2]) {
	case "input", "inputs":
		if tk.inputs != nil {
			pp, ok := tk.inputs[r[3]]
			if !ok {
				return p, fmt.Errorf("not found parameter %s", r[3])
			}
			p = *pp
		} else {
			return p, fmt.Errorf("not found parameter %s", r[3])
		}
	case "output", "outputs":
		if tk.outputs != nil {
			pp, ok := tk.outputs[r[3]]
			if !ok {
				return p, fmt.Errorf("not found parameter %s", r[3])
			}
			p = *pp
		} else {
			return p, fmt.Errorf("not found parameter %s", r[3])
		}
	default:
		return p, fmt.Errorf("ref must be input or output %s", r[2])
	}

	return p, nil
}
