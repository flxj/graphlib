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
