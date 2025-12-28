# go-lisp

basic lisp implementation. Syntax:

```
(lambda fibo (x)(
    (if (= x 0)
        0 
        (if (= x 1)
            1 
            (+ (fibo (- x 1)) (fibo (- x 2)))))	
))

(fibo 10)
```

# configuration server

app showcases using lisp as a configuration DSL. This can be used as a server to fetch params or feature flags. Example:

when init server with this lisp code:
```
(set port 8080)
(set host "localhost")
(lambda foo (param) (
    (if (= param 1) true false)
))
```
you can query these values with REST request:

```
curl localhost:8080/api/flag/host
{"str_value":"localhost"}
```
```
curl localhost:8080/api/flag/port
{"num_value":8080}
```
```
 curl -XPOST localhost:8080/api/flag/foo -d '{"args": [1]}'
{"bool_value":true}
```
```
curl localhost:8080/api/flag/hostaasd
{"error":"failed to get hostaasd: hostaasd not found"}
```