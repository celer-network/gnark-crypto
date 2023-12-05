// Copyright 2020 Consensys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package test_vector_utils

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bw6-756/fr"
	"github.com/consensys/gnark-crypto/ecc/bw6-756/fr/polynomial"
	"hash"
	"reflect"
	"strings"
)

func ToElement(i int64) *fr.Element {
	var res fr.Element
	res.SetInt64(i)
	return &res
}

type HashDescription map[string]interface{}

func HashFromDescription(d HashDescription) (hash.Hash, error) {
	if _type, ok := d["type"]; ok {
		switch _type {
		case "const":
			startState := int64(d["val"].(float64))
			return &MessageCounter{startState: startState, step: 0, state: startState}, nil
		default:
			return nil, fmt.Errorf("unknown fake hash type \"%s\"", _type)
		}
	}
	return nil, fmt.Errorf("hash description missing type")
}

type MessageCounter struct {
	startState int64
	state      int64
	step       int64
}

func (m *MessageCounter) Write(p []byte) (n int, err error) {
	inputBlockSize := (len(p)-1)/fr.Bytes + 1
	m.state += int64(inputBlockSize) * m.step
	return len(p), nil
}

func (m *MessageCounter) Sum(b []byte) []byte {
	inputBlockSize := (len(b)-1)/fr.Bytes + 1
	resI := m.state + int64(inputBlockSize)*m.step
	var res fr.Element
	res.SetInt64(int64(resI))
	resBytes := res.Bytes()
	return resBytes[:]
}

func (m *MessageCounter) Reset() {
	m.state = m.startState
}

func (m *MessageCounter) Size() int {
	return fr.Bytes
}

func (m *MessageCounter) BlockSize() int {
	return fr.Bytes
}

func NewMessageCounter(startState, step int) hash.Hash {
	transcript := &MessageCounter{startState: int64(startState), state: int64(startState), step: int64(step)}
	return transcript
}

func NewMessageCounterGenerator(startState, step int) func() hash.Hash {
	return func() hash.Hash {
		return NewMessageCounter(startState, step)
	}
}

type ListHash []fr.Element

func (h *ListHash) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (h *ListHash) Sum(b []byte) []byte {
	res := (*h)[0].Bytes()
	*h = (*h)[1:]
	return res[:]
}

func (h *ListHash) Reset() {
}

func (h *ListHash) Size() int {
	return fr.Bytes
}

func (h *ListHash) BlockSize() int {
	return fr.Bytes
}
func SetElement(z *fr.Element, value interface{}) (*fr.Element, error) {

	// TODO: Put this in element.SetString?
	switch v := value.(type) {
	case string:

		if sep := strings.Split(v, "/"); len(sep) == 2 {
			var denom fr.Element
			if _, err := z.SetString(sep[0]); err != nil {
				return nil, err
			}
			if _, err := denom.SetString(sep[1]); err != nil {
				return nil, err
			}
			denom.Inverse(&denom)
			z.Mul(z, &denom)
			return z, nil
		}

	case float64:
		asInt := int64(v)
		if float64(asInt) != v {
			return nil, fmt.Errorf("cannot currently parse float")
		}
		z.SetInt64(asInt)
		return z, nil
	}

	return z.SetInterface(value)
}

func SliceToElementSlice[T any](slice []T) ([]fr.Element, error) {
	elementSlice := make([]fr.Element, len(slice))
	for i, v := range slice {
		if _, err := SetElement(&elementSlice[i], v); err != nil {
			return nil, err
		}
	}
	return elementSlice, nil
}

func SliceEquals(a []fr.Element, b []fr.Element) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %d≠%d", len(a), len(b))
	}
	for i := range a {
		if !a[i].Equal(&b[i]) {
			return fmt.Errorf("at index %d: %s ≠ %s", i, a[i].String(), b[i].String())
		}
	}
	return nil
}

func SliceSliceEquals(a [][]fr.Element, b [][]fr.Element) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %d≠%d", len(a), len(b))
	}
	for i := range a {
		if err := SliceEquals(a[i], b[i]); err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}
	}
	return nil
}

func PolynomialSliceEquals(a []polynomial.Polynomial, b []polynomial.Polynomial) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %d≠%d", len(a), len(b))
	}
	for i := range a {
		if err := SliceEquals(a[i], b[i]); err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}
	}
	return nil
}

func ElementToInterface(x *fr.Element) interface{} {
	if i := x.BigInt(nil); i != nil {
		return i
	}
	return x.Text(10)
}

func ElementSliceToInterfaceSlice(x interface{}) []interface{} {
	if x == nil {
		return nil
	}

	X := reflect.ValueOf(x)

	res := make([]interface{}, X.Len())
	for i := range res {
		xI := X.Index(i).Interface().(fr.Element)
		res[i] = ElementToInterface(&xI)
	}
	return res
}

func ElementSliceSliceToInterfaceSliceSlice(x interface{}) [][]interface{} {
	if x == nil {
		return nil
	}

	X := reflect.ValueOf(x)

	res := make([][]interface{}, X.Len())
	for i := range res {
		res[i] = ElementSliceToInterfaceSlice(X.Index(i).Interface())
	}

	return res
}
