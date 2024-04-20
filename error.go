package graphlib

import (
	"errors"
	"strings"
)

var (
	errVertexNotExists = errors.New("vertex not exists")
	errEdgeNotExists   = errors.New("edge not exists")
	errUnknownProperty = errors.New("unknown graph property")
	errNotDigraph      = errors.New("the graph is not digraph")
)

var (
	errNotImplement   = errors.New("not implement the method now")
	errExistsCycle    = errors.New("there are cycles in the current execgraph")
	errAlreadyRunning = errors.New("the current ExecGraph is already running")
	errJobNotExists   = errors.New("the job not exists in current graph")
	errExecCanceled   = errors.New("the execgraph has been canceled")
	errJobCanceled    = errors.New("the job has been canceled")
	errForbidModify   = errors.New("current status is not waiting,cannot modify execgraph structure")
	errJobIsNull      = errors.New("the job is null")
	errNoEntrypoint   = errors.New("not found entrypoint node in current execgraph object")
)

func IsNotExists(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not exists")
}

func IsAlreadyExists(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
