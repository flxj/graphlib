## graphlib workflow 

Graphlib workflow is a lightweight workflow engine.

### Workflow 

A Workflow object is a collection of Task instances, each of which represents a pending operation, The execution sequence of tasks conforms to the predefined DAG structure.

Here is an example to illustrate the use of workflow. Assuming we need to calculate the value of the integer expression 2*(n+1)+3*m+4, (n,m is unknown variable), and we can decompose the calculation process of the modifier into several basic sets of binary operations.

```java
input n --> (n+1) ---> [2*(n+1)]
                                \
                                 }---> [2*(n+1)] + [(3*m)+4] ---> output
                                /
input m --> (3*m)----> [(3*m)+4]
```
It can be seen that the above decomposition is a directed acyclic graph, so workflow can be used for modeling and processing.

We represent each calculation step as a Task, and then define the execution order of the task and the transmission relationship of the output using workflow.

```python
t1      t3
|       |
v       v
t2      t4
|       |
v       |
t5 <----+
|
v
t6
```

t1 receives parameter n, calculates n+1, and passes the result to t2, t2 multiplies the input by 2 and passes the result to t5; t3 receives parameter m, calculates 3*m, and passes the result to t4, t4 adds 4 to the input and passes the result to t5; t5 adds up all input results and passes them to t6, t6 output result.

First, create an empty Workflow
```golang
wf,_:=NewWorkflow(WfNameOption("test"))
```

Then define several tasks, each of which performs some operations of calculating the input numbers before outputting them:
```golang
t1:=NewGeneralTask("t1",func(ps map[string]Parameter)([]Parameter,error){
		p,ok:= ps["n"]
		if !ok {
			return nil,fmt.Errorf("not found input n")
		}
		n := p.Value.(int)
		out:=Parameter{
			Name:"n",
			Value: n+1,
		}
		return []Parameter{out},nil 
	})

	t2:=NewGeneralTask("t2",func(ps map[string]Parameter)([]Parameter,error){
		p,ok:= ps["n"]
		if !ok {
			return nil,fmt.Errorf("not found input n")
		}
		n1 := p.Value.(int)
		out:=Parameter{
			Name:"n",
			Value: n1*2,
		}
		return []Parameter{out},nil 
	})

	t3:=NewGeneralTask("t3",func(ps map[string]Parameter)([]Parameter,error){
		p,ok:= ps["m"]
		if !ok {
			return nil,fmt.Errorf("not found input m")
		}
		m := p.Value.(int)
		out:=Parameter{
			Name:"m",
			Value: m*3,
		}
		return []Parameter{out},nil 
	})

	t4:=NewGeneralTask("t4",func(ps map[string]Parameter)([]Parameter,error){
		p,ok:= ps["m"]
		if !ok {
			return nil,fmt.Errorf("not found input m")
		}
		m := p.Value.(int)
		out:=Parameter{
			Name:"m",
			Value: m+4,
		}
		return []Parameter{out},nil 
	})

	t5:=NewGeneralTask("t5",func(ps map[string]Parameter)([]Parameter,error){
		pn,ok:= ps["n"]
		if !ok {
			return nil,fmt.Errorf("not found input n")
		}
		n := pn.Value.(int)

		pm,ok:= ps["m"]
		if !ok {
			return nil,fmt.Errorf("not found input m")
		}
		m := pm.Value.(int)


		out:=Parameter{
			Name:"res",
			Value: m+n,
		}
		return []Parameter{out},nil 
	})

	t6 := NewGeneralTask("t6",func(ps map[string]Parameter)([]Parameter,error){
		s,ok:= ps["sum"]
		if !ok {
			return nil,fmt.Errorf("not found input sum")
		}
		sum := s.Value.(int)

        fmt.Println("result: ",sum)

		return []Parameter{},nil 
	})
```

Then add the above task to the workflow.
```golang
_ = wf.AddTask(t1,t2,t3,t4,t5,t6)
```
Next, configure the dependency relationships of each task
```shell
t1->t2
t2->t5
t3->t4
t4->t5
t5->t6
```

```golang
_ = wf.AddDependency(t1.Name(),t2.Name())
_ = wf.AddDependency(t2.Name(),t5.Name())
_ = wf.AddDependency(t3.Name(),t4.Name())
_ = wf.AddDependency(t4.Name(),t5.Name())
_ = wf.AddDependency(t5.Name(),t6.Name())
```

Then configure the dependencies of the parameters, such as setting the inputs of t1 and t3 to n and m, and declaring the corresponding outputs.
```golang
n:=100
m:=200

_ = wf.SetInput(t1.Name(),&Parameter{Name:"n",Value:n})
_ = wf.SetInput(t3.Name(),&Parameter{Name:"m",Value:m,})

_ = wf.SetOutput(t1.Name(),&Parameter{Name:"n"})
_ = wf.SetOutput(t3.Name(),&Parameter{Name:"m"})
```
Note that the declaration of the above output is still "n" and "m", which is determined by the implementation of the specific task.

Similarly, set the input and output of t2 and t4, and note that the input of t2 depends on the output of t1. Therefore, set the Ref field (dependency format:`workflowName.taskName.output.parameterName`）
```golang
_ = wf.SetInput(t2.Name(),&Parameter{Name:"n",Ref:fmt.Sprintf("%s.%s.output.n",wf.Name(),t1.Name())})
_ = wf.SetInput(t4.Name(),&Parameter{Name:"m",Ref:fmt.Sprintf("%s.%s.output.m",wf.Name(),t3.Name())})

_ = wf.SetOutput(t2.Name(),&Parameter{Name:"n"})
_ = wf.SetOutput(t4.Name(),&Parameter{Name:"m"})
```

Finally, set the input and output of t5 and t6
```golang
_ = wf.SetInput(t5.Name(),[]*Parameter{
		{Name:"m",Ref:fmt.Sprintf("%s.%s.output.m",wf.Name(),t4.Name())},
		{Name:"n",Ref:fmt.Sprintf("%s.%s.output.n",wf.Name(),t2.Name())},
	}...)
_ = wf.SetInput(t6.Name(),[]*Parameter{
	    {Name:"sum",Ref:fmt.Sprintf("%s.%s.output.res",wf.Name(),t5.Name())},
    }...)

_ = wf.SetOutput(t5.Name(),&Parameter{Name:"res"})
```


After all the above configurations are completed, The workflow has been created and can now be run
```golang
if err:= wf.Start();err!=nil{
	return err 
}
```
The operation is asynchronous, and the real-time status of the workflow can be viewed through the Info() method.



### Management Workflow by http 

Currently supporting the use of HTTP interface to manage workflow.

Firstly, the user needs to prepare several workflow instances as described in the previous section.

Then create an HTTP service and add the workflow object to the service
```golang
port:=8080
svc:=NewService("localhost",port)

svc.Register(wf)
	
```

Finally, run the service
```golang
if err:=svc.Run();err!=nil{
	return err 
}
```

After the service runs successfully, tools such as curl can be used to manage the workflow instances included in the service through the HTTP interface.

View all workflows
```shell
curl -X GET http://localhost:8080/workflows/
```

View detailed information about a specific workflow
```shell
curl -X GET http://localhost:8080/workflows/name | python -m json.tool
```

Run a specific workflow
```shell
curl -X PATCH  http://localhost:8080/workflows/name -G -d "action=run"
```

Stop a specific workflow
```shell
curl -X PATCH  http://localhost:8080/workflows/name  -G -d "action=stop"
```

### dashboard

TODO