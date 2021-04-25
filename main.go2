package main

import (
	"fmt"
	"errors"
	"strings"
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

type DefaultIterator[T any] struct {
	next func() (T, error)
}

func (iter *DefaultIterator[T]) Next() (T, error) {
	return iter.next()
}

type Transducer[Tx, Ty any] struct {
	next func(iter Iterator[Tx]) (func() (Ty, error))
}

func (trans *Transducer[Tx, Ty]) Transduce(iter Iterator[Tx]) (Iterator[Ty]) {
	return &DefaultIterator[Ty]{
		next: trans.next(iter),
	}
}

// func (trans *Transducer[Tx, Ty, Tz]) Append(transAfter *Transducer[Ty, Tz, Tz]) *Transducer[Tx, Ty, Tz] {
// 	return &Transducer[Tx, Ty, Tz]{
// 		next: func(iter Iterator[Ty]) (func() (Tz, error)) {
// 			return transAfter.Transduce(trans.Transduce(iter)).Next
// 		},
// 	}
// }


// type Mapper[T1, Ty any] struct {
// 	apply func(element T1) Ty
// }

func Map[T1, T2 any](apply func(element T1) T2) *Transducer[T1, T2] {
	return &Transducer[T1, T2]{
		next: func(iter Iterator[T1]) (func() (T2, error)) {
			return func() (T2, error) {
				next, err := iter.Next()
				return apply(next), err
			}
		},
	}
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
	(*ints).i = (*ints).i+1
	return ints.data[ints.i-1], nil
}

// func Map[T1, T2 any](Iterator[T1], func(element T1) T2) Iterator[T2]


func Print[T any](s []T) {
	for _, v := range s {
		fmt.Print(v)
	}
}

func main() {
	Print([]string{"Hello, ", "go2\n"})

	iter := IterIntSlice([]int{1, 2, 3})
	trans := Map[int, int](func(e int) int { 
		return e * 2
	})
	trans2 := Map[int, string](func(e int) string { 
		return strings.Repeat("x", e)
	})
	[]interface{}{trans, trans2}
	iter2 := trans2.Transduce(trans.Transduce(&iter))
	for {
		next, err := iter2.Next()
		if err != nil {
			break
		}
		fmt.Println(next)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
