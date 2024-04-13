package graphlib

import (
	"io"
)

type JobState uint16

const (
	Waiting JobState = iota
	Running
	Success
	Stopped
	Failed
)

type Job struct {
	Name string // unique name

	Err   error
	State JobState
}

type ExecGraph interface {
	Digraph[string, Job, int]
}

func NewExecGraph() (ExecGraph, error) {
	return nil, nil
}
func NewExecGraphFromFile(r io.Reader) (ExecGraph, error) {
	return nil, nil
}
func NewExecGraphFromDAG[K comparable, V any, W number](g Graph[K, V, W]) (ExecGraph, error) {
	return nil, nil
}

type dagExecGraph struct {
}
