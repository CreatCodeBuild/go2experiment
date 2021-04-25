package main

import (
	"fmt"
	"errors"
	"time"
)

// The playground now uses square brackets for type parameters. Otherwise,
// the syntax of type parameter lists matches the one of regular parameter
// lists except that all type parameters must have a name, and the type
// parameter list cannot be empty. The predeclared identifier "any" may be
// used in the position of a type parameter constraint (and only there);
// it indicates that there are no constraints.

type Iterator[T1 any] interface {
	Next() (T1, error)
}

type Mapper[T1, Ty any] struct {
	iter Iterator[T1]
	apply func(element T1) Ty
}

func (iter *Mapper[T1, Ty]) Next() (Ty, error) {
	ele, err := iter.iter.Next()
	if err != nil {
		return iter.apply(ele), err
	}
	return iter.apply(ele), nil
}

func Map[T1, T2 any](iter Iterator[T1], apply func(element T1) T2) Mapper[T1, T2] {
	iter2 := Mapper[T1, T2]{
		iter: iter,
		apply: apply,
	}
	return iter2
}

type iterIntSlice struct {
	i int
	data []int
}

func IterIntSlice(ints []int) iterIntSlice {
	return iterIntSlice{
		i: 0,
		data: ints,
	}
}

func (ints *iterIntSlice) Next() (int, error) {
	if ints.i >= len(ints.data) {
		return 0, errors.New("");
	}
	// fmt.Println(ints.i)
	(*ints).i = (*ints).i+1
	fmt.Println(ints.i, "x")
	return ints.data[ints.i-1], nil
}

// func Map[T1, T2 any](Iterator[T1], func(element T1) T2) Iterator[T2]


func Print[T any](s []T) {
	for _, v := range s {
		fmt.Print(v)
	}
}

func main() {
	Print([]string{"Hello, ", "playground\n"})

	iter := IterIntSlice([]int{1, 2, 3})
	
	trans := Map[int, int](func(e int) int { 
		return e * 2
	}).Filter(func(e int) bool {
		return e % 2 == 0
	}).Reduce[int, int](ToIntSlice)

	for {
		next, err := iter2.Next()
		if err != nil {
			break
		}
		fmt.Println(next, iter2)
	}
	// Map[[]int]([]int{1,2,3}, func(ele int) int {
	// 	return ele * 2
	// }).Map(func(ele int) {
	// 	fmt.Println(ele)
	// })

	// task1 := Task(1, func(t *AsyncTask) {
	// 	fmt.Println("task1 starts")
	// 	t2 := t.Task(2, func(t2 *AsyncTask) {
	// 		fmt.Println("task2 starts")
	// 		t2.Sleep(time.Second * 10)
	// 		fmt.Println("task2 is ended")
	// 	})
	// 	t2.Go()
	// 	t3 := t.Task(3, func(t3 *AsyncTask) {
	// 		fmt.Println("task3 starts")
	// 	}, func(at *AsyncTask) {
	// 		t3.Sleep(time.Second * 5)
	// 	}, func(at *AsyncTask) {
	// 		fmt.Println("task3 is ended")
	// 	})
	// 	t3.Go()
	// 	fmt.Println("task1 done")
	// })
	// task1.Go()
	// time.Sleep(time.Second)
	// task1.Cancel()
	// task1.Wait()
	// time.Sleep(time.Second * 10)
	// task1.Cancel()
}

func Task(id int, f func(*AsyncTask)) *AsyncTask {
	return &AsyncTask{
		id: id,
		f: f,
		state: created,
		doneCh: make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
}

type State string

const (
	running = "running"
	created = "created"
	done = "done"
)

type AsyncTask struct {
	id int
	f func(*AsyncTask)
	children []*AsyncTask
	state State
	doneCh chan struct{}
	cancelCh chan struct{}
}

func (t *AsyncTask) Go() {
	go func() {
		t.state = running
		t.f(t)
		t.state = done
		close(t.doneCh)
	}()
}

func (t *AsyncTask) Task(id int, f func(*AsyncTask)) *AsyncTask {
	child := Task(id, f)
	t.children = append(t.children, child)
	return child
}

func (t *AsyncTask) Wait() {
	for i, child := range t.children {
		fmt.Println("wait", i)
		child.Wait()
	}
	<- t.doneCh
	fmt.Println("x")
}

func (t *AsyncTask) Cancel() {
	fmt.Println(t.id, "is being cancelled")
	close(t.cancelCh)
	for _, child := range t.children {
		// fmt.Println("wait", i)
		child.Cancel()
	}
}

func (t *AsyncTask) Sleep(d time.Duration) {
	after := time.After(d)
	select {
	case <- t.cancelCh:
		fmt.Println(t.id, "is cancelled")
		break
	case <- after:
		break
	}
}