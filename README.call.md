## Notice

This repo contains original gorilla/rpc repository with addition - function
```
// Call RPC method without JSON encoding and method's structs knowledge.
// Request args are given via map and decoded into struct by decoder func.
func (s *Server) Call(r *http.Request,
    decoder func(input interface{}, output interface{}) error,
    method string,
    args *map[string]interface{},
) (interface{}, error) {

```
for rpc and rpc/v2 packages.

Puspose of this addition is to allow calling RPC method from templates.

### Addition's copyright

Copyright (c) 2018 Aleksei Kovrizhkin <lekovr+apisite@gmail.com>
