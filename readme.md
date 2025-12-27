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
