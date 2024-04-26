package graphlib

import(
	"testing"
	"fmt"
)

func TestCostQueue(t *testing.T){
	p:=newCostQueue[int]()
	c1:=&item[int]{
		key:1,
		value:3.0,
	}
	c2:=&item[int]{
		key:2,
		value:1.0,
	}
	c3:=&item[int]{
		key:3,
		value:-3.0,
	}
	c4:=&item[int]{
		key:4,
		value:13.0,
	}
	c5:=&item[int]{
		key:5,
		value:MaxFloatDistance,
	}
	c6:=&item[int]{
		key:6,
		value:0.0,
	}
	p.Push(c1)
	p.Push(c2)
	p.Push(c3)
	p.Push(c4)
	p.Push(c5)
	p.Push(c6)
	if p.Len()!=6{
		fmt.Println("queue length error")
		return 
	}
	for p.Len() > 0 {
		v:=p.Pop()
		fmt.Printf("key:%d,value:%f\n",v.key,v.value)
	}

	fmt.Println("=========")
	p.Push(c1)
	p.Push(c2)
	p.Push(c3)
	p.Push(c4)
	p.Update(1,100.0)
	p.Update(4,-100.0)
	for p.Len() > 0 {
		v:=p.Pop()
		fmt.Printf("key:%d,value:%f\n",v.key,v.value)
	}
}

