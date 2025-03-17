package main

import (
	"errors"
	"sync"
)

type Opertation string

const (
	Sum Opertation = "sum"
	Sub Opertation = "sub"
	Mul Opertation = "mul"
	Div Opertation = "div"
)

type (
	Package struct {
		Value     int
		Op        Opertation
		PackageID string
		StreamID  string
	}

	StreamContext struct {
		Value    int
		StreamID string

		mu sync.Mutex
	}
)

var ErrDivByZero = errors.New("division by zero")
