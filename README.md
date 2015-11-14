#initials-avatar

Generate an avatar image from a user's initials. Image background color depends on  name hashes(consistent hashing).

[![GoDoc](https://godoc.org/github.com/holys/initials-avatar/avatar?status.svg)](https://godoc.org/github.com/holys/initials-avatar/avatar)



## Installation

```
$ go get  github.com/holys/initials-avatar/avatar
```


## Usage
### Lib Example 

```
import  "github.com/holys/initials-avatar"

a := avatar.New("/path/to/fontfile")
b, _ := a.DrawToBytes("David", 128)
// now `b` is image data which you can write to file or http stream.
```


### HTTP Example
```
// run it at :3000 by default. Assumes $GOBIN is in your $PATH.
$ avatar

// try it on your browser
// http://127.0.0.1:3000/hello 

```

![](./resource/images/hello.png)

```
// Chinese example
http://127.0.0.1:3000/%E5%AD%94
```

![](./resource/images/kong.png)

## HTTP Benchmark

### With logger

```
$ wrk -t12 -c400 -d30s http://10.20.142.147:3000/a
Running 30s test @ http://10.20.142.147:3000/a
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.18ms   20.56ms 127.10ms   94.18%
    Req/Sec     2.36k     1.34k    5.92k    62.66%
  766657 requests in 30.04s, 1.32GB read
  Socket errors: connect 0, read 150, write 6, timeout 1915
Requests/sec:  25519.86
Transfer/sec:     45.00MB
```


### Without logger

```
$ wrk -t12 -c400 -d30s http://10.20.142.147:3000/a
Running 30s test @ http://10.20.142.147:3000/a
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     7.53ms   11.28ms 195.48ms   96.19%
    Req/Sec     3.66k     2.29k   47.47k    75.96%
  1285314 requests in 30.00s, 589.60MB read
  Socket errors: connect 0, read 276, write 0, timeout 1872
  Non-2xx or 3xx responses: 1
Requests/sec:  42849.74
Transfer/sec:     19.66MB

```

Thanks [@lixiaojun](https://github.com/lixiaojun) for his work.



