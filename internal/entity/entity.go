package entity

import (
	"encoding/json"
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

func (sc *StreamContext) ApplyOperation(p Package) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc, err := applyOpertaion(p, sc)
	if err != nil {
		return err
	}

	return nil
}

func applyOpertaion(p Package, streamCtx *StreamContext) (*StreamContext, error) {
	switch p.Op {
	case Sum:
		streamCtx.Value += p.Value
	case Sub:
		streamCtx.Value -= p.Value
	case Mul:
		streamCtx.Value *= p.Value
	case Div:
		if p.Value == 0 {
			return nil, ErrDivByZero
		}
	}

	return streamCtx, nil
}

func (self StreamContext) MarshalJSON() ([]byte, error) {
	ctx := struct {
		Value    int    `json:"value"`
		StreamID string `json:"streamId"`
	}{self.Value, self.StreamID}

	return json.Marshal(ctx)
}

func (self *StreamContext) UnmarshalJSON(b []byte) error {
	ctx := struct {
		Value    int    `json:"value"`
		StreamID string `json:"streamId"`
	}{}

	if err := json.Unmarshal(b, &ctx); err != nil {
		return err
	}

	self.StreamID = ctx.StreamID
	self.Value = ctx.Value

	return nil
}
