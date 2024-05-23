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
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	errWfNotExists = errors.New("workflow not exists")
)

// http service.
type Service struct {
	host string
	port int

	mu      sync.RWMutex
	running bool
	wfs     map[string]*Workflow
	svc     *gin.Engine
}

func NewService(host string, port int) *Service {
	return &Service{
		host: host,
		port: port,
		wfs:  make(map[string]*Workflow),
	}
}

func (s *Service) Run() error {
	s.mu.Lock()

	if s.running {
		s.mu.Unlock()
		return nil
	}

	s.svc = gin.Default()
	s.router()

	var err error
	go func() {
		time.Sleep(2 * time.Second)
		if err == nil {
			s.running = true
			s.mu.Unlock()
		}
	}()

	err = s.svc.Run(fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.mu.Unlock()
		return err
	}
	return nil
}

func (s *Service) router() {
	wf := s.svc.Group("/workflows")

	// list all wf.
	wf.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"workflows": s.listWorkflow(),
		})
	})

	// get wf by name.
	wf.GET("/:name", func(c *gin.Context) {
		info, err := s.getWorkflow(c.Param("name"))
		if err != nil {
			c.JSON(404, err)
			return
		}

		c.JSON(200, gin.H{
			"workflow": info,
		})
	})

	// create a wf.
	wf.POST("/", func(c *gin.Context) {
		c.JSON(500, errors.New("not support now"))
	})

	// start/stop wf.
	wf.PATCH("/:name", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Query("action")

		var err error
		switch action {
		case "run", "start":
			err = s.runWorkflow(name)
		case "stop":
			err = s.stopWorkflow(name)
		default:
			c.JSON(400, fmt.Errorf("not support action %s", action))
			return
		}
		if err != nil {
			c.JSON(500, err)
			return
		}
		c.JSON(200, gin.H{})
	})

	// delete wf.
	wf.DELETE("/:name", func(c *gin.Context) {
		if err := s.deleteWorkflow(c.Param("name")); err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, gin.H{})
	})
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, wf := range s.wfs {
		_ = wf.Stop()
	}
	s.running = false
	return nil
}

func (s *Service) Register(wf *Workflow) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wfs[wf.Name()] = wf
}

func (s *Service) runWorkflow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wf, ok := s.wfs[name]
	if !ok {
		return errWfNotExists
	}
	return wf.Start()
}

func (s *Service) stopWorkflow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wf, ok := s.wfs[name]
	if !ok {
		return errWfNotExists
	}
	return wf.Stop()
}

func (s *Service) deleteWorkflow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wf, ok := s.wfs[name]
	if !ok {
		return errWfNotExists
	}
	_ = wf.Stop()
	delete(s.wfs, name)
	return nil
}

type workflow struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (s *Service) listWorkflow() []*workflow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ws []*workflow
	for name, wf := range s.wfs {
		ws = append(ws, &workflow{
			Name:   name,
			Status: wf.Status(),
		})
	}
	return ws
}

func (s *Service) getWorkflow(name string) (*WorkflowInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wf, ok := s.wfs[name]
	if !ok {
		return nil, errWfNotExists
	}
	return wf.Info()
}
