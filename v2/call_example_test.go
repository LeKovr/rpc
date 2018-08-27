// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2012 The Gorilla Authors. All rights reserved.
// Copyright 2018 Aleksei Kovrizhkin <lekovr+apisite@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/mitchellh/mapstructure"

	"github.com/LeKovr/rpc/v2"
)

type Service1Request struct {
	A int
	B int
}

type Service1Response struct {
	Result int
}

type Service1 struct {
}

func (t *Service1) Multiply(r *http.Request, req *Service1Request, res *Service1Response) error {
	res.Result = req.A * req.B
	return nil
}

func ExampleServer_Call() {

	const tpl = `{{ $x := api "Service1.Multiply" "a" 2 "b" 3 }}{{ $x.Result }}`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// api placeholder for startup template parsing
	funcs := template.FuncMap{
		"api": func(method string, dict ...interface{}) (interface{}, error) {
			return "", nil
		},
	}

	t, err := template.New("page").Funcs(funcs).Parse(tpl)
	check(err)

	s := rpc.NewServer()
	s.RegisterService(new(Service1), "")

	// Attach api to real service.
	funcs["api"] = func(method string, dict ...interface{}) (interface{}, error) {
		args, err := makeMap(dict...)
		if err != nil {
			return nil, err
		}
		decoder := func(input interface{}, output interface{}) error {
			config := &mapstructure.DecoderConfig{
				Result: output,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return err
			}
			return decoder.Decode(input)
		}
		// In http handler we should pass Request object here
		return s.Call(nil /*ctx.Request*/, decoder, method, args)
	}

	err = t.Funcs(funcs).Execute(os.Stdout, "")
	check(err)
	// Output:
	// 6
}

// make map from {name, value} pairs
func makeMap(dict ...interface{}) (*map[string]interface{}, error) {

	if len(dict)%2 != 0 {
		// log.Printf("Args: %+v", args)
		return nil, errors.New("arg count must be even")
	}

	args := make(map[string]interface{})
	for i := 0; i < len(dict); i += 2 {
		key, isset := dict[i].(string)
		if !isset {
			return nil, fmt.Errorf("not string key in position %d", i)
		}
		args[key] = dict[i+1]
	}
	return &args, nil
}
