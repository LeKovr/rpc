// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2012 The Gorilla Authors. All rights reserved.
// Copyright 2018 Aleksei Kovrizhkin <lekovr+apisite@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mitchellh/mapstructure"
)

var decoder = func(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result: output,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func TestCall(t *testing.T) {
	args := map[string]interface{}{
		"a": 2,
		"b": 3,
	}
	expected := args["a"].(int) * args["b"].(int)

	s := NewServer()
	s.RegisterService(new(Service1), "")

	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	resIface, err := s.Call(r, decoder, "Service1.Multiply", &args)
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resIface.(*Service1Response)
	if !ok {
		t.Errorf("Response type of %+v should be *Service1Response.", res)
	}
	if res.Result != expected {
		t.Errorf("Response was %d, should be %d.", res.Result, expected)
	}

}

func TestCallTCP(t *testing.T) {
	args := map[string]interface{}{
		"a": 2,
		"b": 3,
	}
	expected := args["a"].(int) + args["b"].(int)

	s := NewServer()
	s.RegisterTCPService(new(Service1), "")

	resIface, err := s.Call(nil, decoder, "Service1.Add", &args)
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resIface.(*Service1Response)
	if !ok {
		t.Errorf("Response type of %+v should be *Service1Response.", res)
	}
	if res.Result != expected {
		t.Errorf("Response was %d, should be %d.", res.Result, expected)
	}
}

func (t *Service1) Raise(r *http.Request, req *Service1Request, res *Service1Response) error {
	return errors.New("error raised")
}

func TestCallErrorMethodNotFound(t *testing.T) {
	args := map[string]interface{}{"a": 2}
	s := NewServer()
	s.RegisterService(new(Service1), "")

	_, err := s.Call(nil, decoder, "Service1.Unknown", &args)
	if err == nil {
		t.Error("Call of unknown method should throw error.")
	}
}

func TestCallErrorIncorrectMethodArgs(t *testing.T) {
	args := map[string]interface{}{"a": "2", "b": 3}
	s := NewServer()
	s.RegisterService(new(Service1), "")

	_, err := s.Call(nil, decoder, "Service1.Multiply", &args)
	if err == nil {
		t.Error("Call for incorrect request type should throw error.")
	}
}

func TestCallErrorRaised(t *testing.T) {
	args := map[string]interface{}{}
	s := NewServer()
	s.RegisterService(new(Service1), "")

	_, err := s.Call(nil, decoder, "Service1.Raise", &args)
	if err == nil {
		t.Error("Call with error should throw error.")
	}
}
