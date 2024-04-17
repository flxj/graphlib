package graphlib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

type (
	State string
	// Job represents a runnable function object,
	// which in graphlib is an abstraction of an independent task.
	Job func() error
	// ContextJob represents a runnable function object with a
	// context parameter, which differs from Job in that when
	// running ContextJob,the ExecGraph execution engine automatically
	// generates a context parameter for it and cancels the context
	// when the user performs a Stop operation.
	ContextJob func(context.Context) error
)

type job interface {
	Job | ContextJob
}

const (
	// Waiting represents the initial state of a job or
	// ExecGraph, in which the user can adjust the ExecGraph
	// structure to update the job.
	Waiting State = "waiting"
	// Running status indicates that Job or ExecGraph is running
	Running State = "running"
	// If the job runs without returning an error, it will
	// be in the Success state.After all jobs in the
	// ExecGraph run successfully, the ExecGraph will also
	// be set to the Success state
	Success State = "success"
	// Stopped status indicates that Job or ExecGraph
	// has been actively terminated by the user.
	Stopped State = "stopped"
	// Failed indicates that the job encountered an error
	// while running.If any job in ExecGraph fails to run
	// at the same time, the final status of ExecGraph will
	// also be set to Failed.
	Failed State = "failed"
	// Paused represents the state of ExecGraph
	// after being actively paused by the user.
	Paused State = "paused"

	//
	DefaultMaxConcurrencyJob = 1000000
)

var (
	ErrNotImplement   = errors.New("not implement the method now")
	ErrExistsCycle    = errors.New("there are cycles in the current execgraph")
	ErrAlreadyRunning = errors.New("the current ExecGraph is already running")
	ErrJobNotExists   = errors.New("the job not exists in current graph")
	ErrExecCanceled   = errors.New("the execgraph has been canceled")
	ErrJobCanceled    = errors.New("the job has been canceled")
	ErrForbidModify   = errors.New("current status is not waiting,cannot modify execgraph structure")
	ErrJobIsNull      = errors.New("the job is null")
	ErrNoEntrypoint   = errors.New("not found entrypoint node in current execgraph object")
)

func edgeFormat[K comparable](v1, v2 K) K {
	return any(fmt.Sprintf("%v-%v", v1, v2)).(K)
}

// JobInfo record the job's status of the ExecGraph node.
type JobInfo[K comparable] struct {
	// The unique identifier of the job.
	Key K
	// The error message returned by running the job.
	Error error
	// The current status of the job.
	Status State
	// Start time.
	StartAt time.Time
	// End time.
	EndAt time.Time
}

// DAG
type ExecGraph[K comparable, J job] interface { // TODO: change the name as GoGraph
	//
	// The Start method will attempt to run ExecGraph. If it is currently
	// in a Waiting, Stopped, or Failed state, the entire Graph will be run
	// from scratch.If it is in the Paused state, it will continue running
	// from the pause point.
	//
	// Until the status of the Graph changes to Success or Failed, it indicates
	// that the Graph has finished running.
	//
	// Running a Graph that is already in the Running state will return an
	// error indicating that it is already running.
	//
	// Currently, all Start calls to the Success state graph will be ignored.
	// If you want to run it again, you need to first call the Reset method
	// to reset the Graph to the Waiting state.
	Start() error
	//
	// The Stop method will attempt to terminate the current Graph run
	// (Cancel the running job and no longer run a new job).
	//
	// If calling the Graph method on the Success state, it will return an error.
	Stop() error
	//
	// Wait is used to wait for the Graph to finish running and will return
	// the latest error message.
	//
	// The gorouting calling this method will block until there are no
	// runnable jobs in the Graph.
	Wait() error
	//
	// The Reset method will reset the Graph to its initial state.
	// Calling this method on a Graph in Running state
	// will first stop the Graph and then reset it.
	Reset() error
	//
	// The Pause method will pause the execution of the Graph.
	//
	// Note that if there are jobs that happen to be running when calling
	// Pause, meaning they start before the current pause point, then these
	// jobs will not be canceled.
	Pause() error
	//
	// The Status method is used to view the global status of the current Graph.
	Status() State
	//
	// The Job method returns the basic information of the job corresponding
	// to the key value,including the current status, error information, etc.
	//
	// If the job does not exist, return NotExists error.
	Job(key K) (JobInfo[K], error)
	//
	// Add a new job to the Graph, and if the same key value already exists,
	// update the Job.
	//
	// Note that dynamic modification of the running Graph structure is
	// currently not allowed, so this method can only be called on graphs
	// in the Waiting state.
	AddJob(key K, job J) error
	//
	// AddTimeoutJob adds a new job to the Graph and sets the timeout for
	// the job to run.
	//
	// If the same key value already exists, the Job will be updated.
	//
	// Note that dynamic modification of the running Graph structure is
	// currently not allowed, so this method can only be called on graphs
	// in the Waiting state.
	AddTimeoutJob(key K, job J, timeout time.Duration) error
	//
	// AddRetryJob adds a new job to the Graph and sets the number of retries
	// for the job to run (n<=0 indicates no retry).
	//
	// If the same key value already exists, the Job will be updated.
	//
	// Note that dynamic modification of the running Graph structure is
	// currently not allowed,so this method can only be called on graphs
	// in the Waiting state.
	AddRetryJob(key K, job J, retry int) error
	//
	// RemoveJob removes a Job from the current Graph And remove the
	// relevant dependencies.
	//
	// If the job does not exist, return a NotExists error.
	//
	// Note that currently only calling this method on Graph in Waiting
	// state is supported.
	RemoveJob(key K) error
	//
	// AddDependency adds a dependency relationship between jobs in the current
	// Graph, which is unidirectional:Source job ->target job
	// indicates that source is the predecessor task of the target.
	//
	// Note that currently only calling this method on Graph in Waiting
	// state is supported.
	AddDependency(source, target K) error
	//
	// RemoveDependency removes dependencies between tasks from the current Graph.
	//
	// Note that currently only calling this method on Graph in Waiting
	// state is supported.
	RemoveDependency(source, target K) error
	//
	// Set the maximum number of concurrent job runs.
	// Note that this setting is currently not supported,
	// meaning there is no limit on the concurrency of jobs.
	SetMaxConcurrencyJob(n int)
	//
	// StoppJob stops a running job.
	//
	// Note that stopping the job may have an impact on its subsequent operations.
	//
	// If ignoreErr is set to false, the stopped job is considered a failure
	// (the reason for the failure is due to the user's active cancellation),
	// which will result in the final graph state being Failed;
	//
	// On the other hand, if ignoreErr is set to true, the job is considered
	// to have ended successfully.
	StopJob(key K, ignoreErr bool) error
	//
	// RunJob runs a specific job separately.
	// Note that this feature is currently not supported.
	RunJob(key K) error
	//
	// Detecting possible job ring dependencies in the Graph.
	DetectCycle() ([][]K, error)
}

// Create an empty ExecGraph.
func NewExecGraph[K comparable, J job]() (ExecGraph[K, J], error) {
	dag, err := NewDigraph[K, struct{}, int]()
	if err != nil {
		return nil, err
	}
	eg := &execGraph[K, J]{
		dag:        dag,
		complete:   make(chan struct{}),
		wait:       make(chan struct{}),
		status:     Waiting,
		nodes:      make(map[K]*execNode[K, J]),
		candicates: make(map[K]int),
		resCh:      make(chan *execResult[K], 256),
		finishes:   make(map[K]struct{}),
	}
	return eg, nil
}

// Load a DAG from text data and create an ExecGraph based on it.
func NewExecGraphFromFile[K comparable, J job](r io.Reader) (ExecGraph[K, J], error) {
	return nil, ErrNotImplement
}

// Create an ExecGraph based on an existing DAG object.
func NewExecGraphFromDAG[K comparable, V any, W number, J job](g Digraph[K, V, W]) (ExecGraph[K, J], error) {
	if !g.IsAcyclic() {
		return nil, ErrExistsCycle
	}
	var (
		err error
		dag Digraph[K, struct{}, int]
		vs  []Vertex[K, V]
		es  []Edge[K, W]
	)
	if dag, err = NewDigraph[K, struct{}, int](); err != nil {
		return nil, err
	}
	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}
	for _, v := range vs {
		nv := Vertex[K, struct{}]{Key: v.Key}
		if err = dag.AddVertex(nv); err != nil {
			return nil, err
		}
	}
	for _, e := range es {
		ne := Edge[K, int]{
			Key:  edgeFormat(e.Head, e.Tail),
			Head: e.Head,
			Tail: e.Tail,
		}
		if err = dag.AddEdge(ne); err != nil {
			return nil, err
		}
	}

	eg := &execGraph[K, J]{
		dag:        dag,
		complete:   make(chan struct{}),
		wait:       make(chan struct{}),
		status:     Waiting,
		nodes:      make(map[K]*execNode[K, J]),
		candicates: make(map[K]int),
		resCh:      make(chan *execResult[K], 256),
		finishes:   make(map[K]struct{}),
	}
	return eg, nil
}

type execResult[K comparable] struct {
	key   K
	err   error
	endAt time.Time
}

type execNode[K comparable, J job] struct {
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	runner     J
	info       *JobInfo[K]
	retryLimit int
	timeout    time.Duration
	version    int // The value of this field is monotonically increasing。
}

func newExecNode[K comparable, J job](key K, jb J, d time.Duration, n int) *execNode[K, J] {
	return &execNode[K, J]{
		runner: jb,
		info: &JobInfo[K]{
			Key:    key,
			Status: Waiting,
		},
		timeout:    d,
		retryLimit: n,
	}
}

func (e *execNode[K, J]) getInfo() JobInfo[K] {
	//e.mu.RLock()
	//defer e.mu.RUnlock()

	return JobInfo[K]{
		Key:     e.info.Key,
		Error:   e.info.Error,
		Status:  e.info.Status,
		StartAt: e.info.StartAt,
		EndAt:   e.info.EndAt,
	}
}

func (e *execNode[K, J]) updateJob(runner J) {
	e.runner = runner
}

func (e *execNode[K, J]) updateTimeout(d time.Duration) {
	e.timeout = d
}

func (e *execNode[K, J]) updateRetry(n int) {
	e.retryLimit = n
}

func (e *execNode[K, J]) run(ch chan *execResult[K]) error {
	// check status first.
	e.mu.RLock()
	if e.info.Status == Running {
		e.mu.RUnlock()
		return nil
	}
	e.mu.RUnlock()
	// run the job.
	go func() {
		var (
			err     error
			version int
		)
		e.mu.Lock()
		e.version++
		version = e.version // record the current version
		e.info.Status = Running
		e.info.StartAt = time.Now()
		if e.runner != nil {
			job, ok := any(e.runner).(func() error)
			if !ok {
				ctxJob, _ := any(e.runner).(func(context.Context) error)
				e.ctx, e.cancel = context.WithCancel(e.ctx)
				job = func() error {
					return ctxJob(e.ctx)
				}
			}
			e.mu.Unlock()

			err = runWithRetry(e.retryLimit, e.timeout, job)
		} else {
			e.mu.Unlock()
			err = ErrJobIsNull
		}
		end := time.Now()

		// TODO maybe we should use atomic operation for version.
		e.mu.Lock()
		defer e.mu.Unlock()
		if version == e.version && e.info.Status == Running {
			e.info.EndAt = end
			if err != nil {
				e.info.Status = Failed
				e.info.Error = err
			} else {
				e.info.Status = Success
			}
			ch <- &execResult[K]{
				key:   e.info.Key,
				err:   err,
				endAt: e.info.EndAt,
			}
		}
		// if old version value not equal current value, we just ignore the result.
		// if the node status is not running, which means maybe someone rest/stop the job,so we also need ignore the job result.
	}()

	return nil
}

func (e *execNode[K, J]) stop(ignoreErr bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.info.Status != Running {
		return fmt.Errorf("job status is not running,cannot stop it")
	}

	if e.cancel != nil {
		e.cancel()
	}
	e.info.EndAt = time.Now()
	if ignoreErr {
		e.info.Status = Success
	} else {
		e.info.Error = ErrJobCanceled
		e.info.Status = Stopped
	}

	return nil
}

func (e *execNode[K, J]) reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.cancel != nil {
		e.cancel()
	}
	e.info.Error = nil
	e.info.Status = Waiting
	e.info.StartAt = time.Time{}
	e.info.EndAt = time.Time{}
	e.version++
}

func (e *execNode[K, J]) isRunning() bool {
	return e.info.Status == Running
}

type execGraph[K comparable, J job] struct {
	mu    sync.RWMutex
	limit int
	// using DAG to orchestrate the execution workflow of Jobs.
	dag Digraph[K, struct{}, int]
	// to stop the start() goroutinue.
	complete  chan struct{}
	completed bool
	// to signal all waiters when the graph status switch to Stopped/Success/Failed.
	wait chan struct{}
	// record the latest error info generated by jobs.
	err error
	// global state.
	status State
	// used to prevent start gorutinue from running new jobs.
	suspend bool
	// the set all jobs.
	nodes map[K]*execNode[K, J]
	// trace the dependent jobs count for a waiting job,
	// when the count become zero,which means the candicate job ready to run.
	candicates map[K]int
	// recieve results from exec node.
	resCh chan *execResult[K]
	//
	finishes map[K]struct{}
	// This collection is used to cache all ready jobs,
	// and this field is used to control the number of concurrent executed jobs。
	//ready map[K]bool
}

func (g *execGraph[K, J]) Start() error {
	if !g.dag.IsAcyclic() {
		return ErrExistsCycle
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.status {
	case Running:
		return ErrAlreadyRunning
	case Paused:
		// continue run from the last suspened point.
		g.suspend = false
		g.status = Running
	case Waiting, Stopped, Failed:
		// restart running the entire graph.
		g.reset()
		g.start()
	case Success:
		return nil // TODO maybe should set a once parameter.
	default:
		return fmt.Errorf("unknown graph status %s", g.status)
	}
	return nil
}

func (g *execGraph[K, J]) scheduledError(key K, err error) {
	g.status = Failed
	g.err = fmt.Errorf("scheduled job %v failed:%v", key, err)
	g.suspend = true
	close(g.wait)
}

func (g *execGraph[K, J]) start() {
	g.status = Running
	// load source vertecis.
	sources, err := g.dag.Sources()
	if err != nil {
		g.scheduledError(any("").(K), ErrNoEntrypoint)
		return
	}
	for _, v := range sources {
		if _, ok := g.nodes[v.Key]; ok {
			g.candicates[v.Key] = 0
		}
	}
	// start the main loop.
	go func() {
		finish := false
		hasReady := true
		for {
			select {
			case <-g.complete: // Stop() maybe close the channel.
				return
			case res, ok := <-g.resCh:
				if !ok {
					// TODO: teardown
				}
				if res != nil {
					g.mu.Lock()
					// if find error,should change status as Failed,and signal waiters.
					if res.err != nil {
						g.err = res.err
						if g.status == Running || g.status == Paused {
							g.status = Failed
							close(g.wait)
							// we dont set suspend=true,
							// because the default design principle of execGraph is to make all jobs run as much as possible.
							// If users want fast failure mode, they need to explicitly set parameters.
						}
					} else {
						// update candicates refCount
						outs, err := g.dag.OutNeighbours(res.key) // TODO:we need ignore some err, for example when res.key sink node,it not has successors.
						if err != nil {
							g.scheduledError(res.key, err)
							g.mu.Unlock()
							return
						}
						for _, v := range outs {
							// if key in candicates,then decrement the value of the corresponding record.
							dep, ok := g.candicates[v.Key]
							if ok {
								g.candicates[v.Key] = dep - 1
								if dep == 1 {
									hasReady = true
								} else if dep <= 0 {
									// TODO: panic
								}
							} else {
								// if the successor not in nodes set,then just ignore it,
								// This situation usually occurs when creating a Graph based on an existing DAG
								// and not adding a Job to certain DAG vertices.
								if _, ok := g.nodes[v.Key]; !ok {
									continue
								}
								// if successors not in candicates set, then add it
								if _, ok := g.finishes[v.Key]; !ok {
									n, err := g.dag.InDegree(v.Key)
									if err != nil {
										g.scheduledError(v.Key, err)
										g.mu.Unlock()
										return
									}
									if n != 0 {
										g.candicates[v.Key] = n - 1
										if n == 1 {
											hasReady = true
										}
									}
								}
								// if v.Key in finishes set:
								// This means that the successor of the job ends before the current job,
								// which occurs when the user triggers a retry operation on a particular job separately.
								// For this situation, the current approach is to ignore the impact of the job's operation on its successors (already completed)。
							}
						}
					}
					g.finishes[res.key] = struct{}{}
					g.mu.Unlock()
				}
			default:
				// scheduled new job to run
				if hasReady {
					g.mu.Lock()
					if !g.suspend {
						var ready []K
						for key, dep := range g.candicates {
							if dep == 0 {
								ready = append(ready, key)
							}
						}
						// shchedule ready jobs to run.
						for _, key := range ready {
							if _, ok := g.finishes[key]; !ok {
								node := g.nodes[key]
								if node != nil {
									_ = node.run(g.resCh)
								}
							}
							// move it from candicates
							delete(g.candicates, key)
						}
						hasReady = false
					}
					g.mu.Unlock()
				} else {
					// check whether the termination conditions are met:
					// i) there aren't any jobs running
					// ii) there aren't any jobs that are ready
					g.mu.RLock()
					finish = true
					for _, node := range g.nodes {
						if node.isRunning() {
							finish = false
							break
						}
					}
					for _, v := range g.candicates {
						if v == 0 {
							finish = false
							break
						}
					}
					g.mu.RUnlock()
				}
				// terminate the min loop
				if finish {
					g.mu.Lock()
					if g.status == Running {
						g.status = Success
						g.err = nil
						close(g.wait)
					}
					g.mu.Unlock()
					return
				}
			}
		}
	}()
}

// setting the all nodes to initialization status(waiting).
func (g *execGraph[K, J]) reset() {
	for _, node := range g.nodes {
		node.reset()
	}
	g.candicates = make(map[K]int)
	g.finishes = make(map[K]struct{})
	g.complete = make(chan struct{})
	g.wait = make(chan struct{})
	g.suspend = false
	g.completed = false
	g.err = nil
	g.status = Waiting
}

// stop the current running graph.
func (g *execGraph[K, J]) Stop() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	defer func() {
		// to inform all waiters that the graph 'completed'.
		if !g.completed {
			close(g.complete)
			g.completed = true
		}
	}()

	switch g.status {
	case Stopped:
	case Failed:
		// stop the any running jobs, and keep the Failed status.
		g.suspend = true
		if err := g.stop(); err != nil {
			g.err = err
			g.status = Failed
			return err
		}
	case Waiting, Running, Paused:
		// should stop any running jobs,and change the status as stopped.
		// first setting a stop falg,to signal start goroutinue to stop run new job.
		// then close all running jobs.
		g.suspend = true
		if err := g.stop(); err != nil {
			g.scheduledError(any("").(K), err)
			return err
		}
		g.status = Stopped
		g.err = ErrExecCanceled
		close(g.wait)
	case Success:
		return fmt.Errorf("current status is %s,which means no jobs running,needont stop anything", g.status)
	default:
		return fmt.Errorf("unknown graph status %s", g.status)
	}
	return nil
}

// stop all running jobs.
func (g *execGraph[K, J]) stop() error {
	for _, node := range g.nodes {
		if node.isRunning() {
			_ = node.stop(false)
		}
	}
	return nil
}

func (g *execGraph[K, J]) Pause() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.status {
	case Paused:
	case Running, Waiting:
		g.suspend = true
		g.status = Paused
	case Failed, Stopped:
		g.suspend = true
	case Success:
		return fmt.Errorf("current status is success,which means no jobs running,needont pause anything")
	default:
		return fmt.Errorf("unknown graph status %s", g.status)
	}
	return nil
}

// wait the graph completed.
func (g *execGraph[K, J]) Wait() error {
	<-g.wait

	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.err
}

func (g *execGraph[K, J]) Status() State {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.status
}

func (g *execGraph[K, J]) Job(key K) (JobInfo[K], error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if job, ok := g.nodes[key]; ok {
		return job.getInfo(), nil
	}
	return JobInfo[K]{}, ErrJobNotExists
}

func (g *execGraph[K, J]) Reset() error {
	ch := make(chan error)
	go func() {
		ch <- g.Stop()
	}()

	defer close(ch)

	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	//
	for _, node := range g.nodes {
		node.reset()
	}
	g.reset()

	return nil
}

func (g *execGraph[K, J]) addJob(key K, job J, d time.Duration, n int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	//
	if g.status != Waiting {
		return ErrForbidModify
	}
	if jb, ok := g.nodes[key]; ok {
		jb.updateJob(job)
		if d != time.Duration(0) {
			jb.updateTimeout(d)
		}
		if n != 0 {
			jb.updateRetry(n)
		}
		return nil
	}

	v := Vertex[K, struct{}]{
		Key: key,
	}
	if err := g.dag.AddVertex(v); err != nil {
		return err
	}
	g.nodes[key] = newExecNode(key, job, d, n)

	return nil
}

func (g *execGraph[K, J]) AddJob(key K, job J) error {
	return g.addJob(key, job, time.Duration(0), 0)
}

func (g *execGraph[K, J]) AddTimeoutJob(key K, job J, timeout time.Duration) error {
	return g.addJob(key, job, timeout, 0)
}

func (g *execGraph[K, J]) AddRetryJob(key K, job J, retry int) error {
	return g.addJob(key, job, time.Duration(0), retry)
}

func (g *execGraph[K, J]) RemoveJob(key K) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	//
	if g.status != Waiting {
		return ErrForbidModify
	}
	//
	_, ok := g.nodes[key]
	if !ok {
		return ErrJobNotExists
	}
	if err := g.dag.RemoveVertex(key); err != nil {
		return err
	}
	delete(g.nodes, key)
	delete(g.candicates, key)

	return nil
}

func (g *execGraph[K, J]) AddDependency(source, target K) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	//
	if g.status != Waiting {
		return ErrForbidModify
	}

	if _, ok := g.nodes[source]; !ok {
		return ErrJobNotExists
	}
	if _, ok := g.nodes[target]; !ok {
		return ErrJobNotExists
	}

	edge := Edge[K, int]{
		Key:  edgeFormat(source, target),
		Head: source,
		Tail: target,
	}
	return g.dag.AddEdge(edge)
}

func (g *execGraph[K, J]) RemoveDependency(source, target K) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	//
	if g.status != Waiting {
		return ErrForbidModify
	}

	if _, ok := g.nodes[source]; !ok {
		return ErrJobNotExists
	}
	if _, ok := g.nodes[target]; !ok {
		return ErrJobNotExists
	}

	return g.dag.RemoveEdgeByKey(edgeFormat(source, target))
}

func (g *execGraph[K, J]) SetMaxConcurrencyJob(n int) {}

// stop a running execNode, and send a result to start goroutinue,if ignore err then should
func (g *execGraph[K, J]) StopJob(key K, ignoreErr bool) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	node, ok := g.nodes[key]
	if !ok {
		return ErrJobNotExists
	}
	if err := node.stop(ignoreErr); err != nil {
		return err
	}
	// signal execGraph that the job has completed.
	res := &execResult[K]{
		key:   key,
		endAt: time.Now(),
	}
	if !ignoreErr {
		res.err = ErrJobCanceled
	}
	g.resCh <- res

	return nil
}

func (g *execGraph[K, J]) RunJob(key K) error {
	return ErrNotImplement
}

func (g *execGraph[K, J]) DetectCycle() ([][]K, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.dag.DetectCycle()
}
