package graphlib

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type intORStr interface {
	int | string
}

type A[T intORStr] struct {
	Data T
	Mp   map[int]int
}

func (a *A[T]) update(d T) {
	a.Data = d
}

func (a *A[T]) print() {
	n, ok := any(a.Data).(int)
	if ok {
		fmt.Println("data is int: ", n)
	} else {
		s, ok := any(a.Data).(string)
		if ok {
			fmt.Println("data is string: ", s)
		} else {
			fmt.Println("unknown")
		}
	}
}

func Test1(t *testing.T) {
	a := &A[int]{Data: 5}

	a.print()

	a.update(7)

}

func Test2(t *testing.T) {
	ctx := context.Background()

	ctx1, cal1 := context.WithCancel(ctx)

	go func() {
		<-ctx1.Done()
		fmt.Println("ctx1 done 1")
	}()

	fmt.Println("cal1")
	time.Sleep(2 * time.Second)
	cal1()
	time.Sleep(2 * time.Second)
	go func() {
		<-ctx1.Done()
		fmt.Println("ctx1 done 2")
	}()
	time.Sleep(2 * time.Second)
	fmt.Println("cal2")
	cal1()
	go func() {
		<-ctx1.Done()
		fmt.Println("ctx1 done 3")
	}()
	time.Sleep(2 * time.Second)
}

func Test3(t *testing.T) {
	ap := make(map[int]*A[string])
	ap[1] = &A[string]{Data: "a"}
	ap[2] = &A[string]{Data: "b", Mp: make(map[int]int)}
	ap[3] = &A[string]{Data: "c", Mp: make(map[int]int)}
	ap[2].Mp[2] = 200
	ap[3].Mp[3] = 300

	bp := make(map[int]*A[string])
	for k, v := range ap {
		//var p A[string]
		p := *v
		bp[k] = &p
	}
	fmt.Println("copy bp =============")
	for k, v := range bp {
		fmt.Printf("(%d,%s,%v)\n", k, v.Data, v.Mp)
	}
	fmt.Println("change ap =============")
	ap[2].Data = "bbbbb"
	ap[2].Mp[2] = 222
	for k, v := range bp {
		fmt.Printf("(%d,%s,%v)\n", k, v.Data, v.Mp)
	}

	fmt.Println("change bp =============")
	bp[2].Data = "dddddd"
	bp[3].Mp[3] = 2333
	for k, v := range bp {
		fmt.Printf("(%d,%s,%v)\n", k, v.Data, v.Mp)
	}
}
