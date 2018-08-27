// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2012 The Gorilla Authors. All rights reserved.
// Copyright 2018 Aleksei Kovrizhkin <lekovr+apisite@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"net/http"
	"reflect"
)

// Call RPC method without JSON encoding and method's structs knowledge.
// Request args are given via map and decoded into struct by decoder func.
func (s *Server) Call(r *http.Request,
	decoder func(input interface{}, output interface{}) error,
	method string,
	args *map[string]interface{},
) (interface{}, error) {
	// Lookup method.
	serviceSpec, methodSpec, errGet := s.services.get(method)
	if errGet != nil {
		return nil, errGet
	}
	// Decode the args.
	argsStruct := reflect.New(methodSpec.argsType)
	if errRead := decoder(args, argsStruct.Interface()); errRead != nil {
		return nil, errRead
	}
	// Call the service method.
	reply := reflect.New(methodSpec.replyType)
	// omit the HTTP request if the service method doesn't accept it
	var errValue []reflect.Value
	if serviceSpec.passReq {
		errValue = methodSpec.method.Func.Call([]reflect.Value{
			serviceSpec.rcvr,
			reflect.ValueOf(r),
			argsStruct,
			reply,
		})
	} else {
		errValue = methodSpec.method.Func.Call([]reflect.Value{
			serviceSpec.rcvr,
			argsStruct,
			reply,
		})
	}
	// Cast the result to error if needed.
	var errResult error
	errInter := errValue[0].Interface()
	if errInter != nil {
		errResult = errInter.(error)
	}
	return reply.Interface(), errResult
}
