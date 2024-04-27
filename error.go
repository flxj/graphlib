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
	"errors"
	"strings"
)

var (
	errVertexNotExists  = errors.New("vertex not exists")
	errVertexExists     = errors.New("vertex already exists")
	errEdgeNotExists    = errors.New("edge not exists")
	errEdgeExists       = errors.New("edge already exists")
	errUnknownProperty  = errors.New("unknown graph property")
	errNotDigraph       = errors.New("the graph is not digraph")
	errHasNegativeCycle = errors.New("found negative cycle")
	errNotDAG           = errors.New("current digraph is not DAG")
	errNotConnected     = errors.New("current graph is not connected")
	errEmptyGraph       = errors.New("current graph is empty")
	errNotSimple        = errors.New("current graph is not simple")
	errViolateBipartite = errors.New("violate the definition of bipartite")
	errCloneFailed      = errors.New("clone current graph failed")
	errNone             = errors.New("")
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
