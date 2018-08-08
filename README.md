# initials-avatar

[![Build Status](https://travis-ci.org/holys/initials-avatar.svg?branch=master)](https://travis-ci.org/holys/initials-avatar)
[![Coverage Status](https://coveralls.io/repos/holys/initials-avatar/badge.svg?branch=master&service=github)](https://coveralls.io/github/holys/initials-avatar?branch=master)
[![GoDoc](https://godoc.org/github.com/holys/initials-avatar/avatar?status.svg)](https://godoc.org/github.com/holys/initials-avatar)


Generate an avatar image from a user's initials. Image background color depends on  name hashes(consistent hashing).


## Online Demo

https://initials.herokuapp.com 

You may switch to [heroku-branch](https://github.com/holys/initials-avatar/tree/feature-heroku) to see how to deploy to heroku.

## Installation

*VERSION REQUIRED* **GO 1.3 or greater**

```
$ go get github.com/holys/initials-avatar/...
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
// run the http server. The port is :3000 by default. Assumes $GOBIN is in your $PATH.
$ avatar server

// try it on your browser
// http://127.0.0.1:3000/hello 

// to view avaliable options
$ avatar server --help

```

![](./resource/images/hello.png)

```
// Chinese example
http://127.0.0.1:3000/%E5%AD%94
```

![](./resource/images/kong.png)

## HTTP Benchmark

Environment:

```
MacBook Pro (Retina, 13-inch, Mid 2014)
Processor: 2.6 GHz Intel Core i5
Memory: 8 GB 1600 MHz DDR3
```


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


## LICENSE 
MIT LICENSE, see [LICENSE](./LICENSE) for details.

