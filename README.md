<h1 align="center">gin- limiter</h1>
<h5 align="center">A simple gin middleware for IP limiter based on redis.</h5>

<p align="center">
    <a href="https://www.gnu.org/licenses/"> 
        <img src="https://img.shields.io/github/license/davidleitw/goGamer.svg" alt="License">
    </a>
    <a href="http://hits.dwyl.io/davidleitw/gin-limiter">
        <img src=http://hits.dwyl.io/davidleitw/gin-limiter.svg alt="HitCount">
    </a>
    <a href="https://github.com/davidleitw/gin-limiter/stargazers"> 
        <img src="https://img.shields.io/github/stars/davidleitw/gin-limiter" alt="GitHub stars">
    </a>
</p>

### Installation

```go 
go get github.com/go-redis/redis/v8
go get github.com/davidleitw/gin-limiter 
``` 

**Import**
```go 
import "github.com/go-redis/redis/v8"
import limiter "github.com/davidleitw/gin-limiter"
```
<hr>

### Quickstart

- ##### Create limit controller object
```go
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0}) // set redis Client

    limitControl, err := limiter.DefaultController(rdb, "24-M", 100, "debug") // Debug mode, each 24 minutes can send 100 times request from single Ip.
    if err != nil {
        log.Println(err)
    }

```

- ##### Debug mode can show some information on the command.
![](https://imgur.com/KeZsQpQ.png)

- ##### For each route add a sub-limiter. 
```go
    server := gin.Default()

    err = limitControl.Add("/ExamplePost1", "POST", "4-M", 20) // "/ExamplePost1" route, each 4 minutes can send 20 times request from single Ip.
    if err != nil {
        log.Println(err)
    }

    err = limitControl.Add("/ExampleGet1", "GET", "20-H", 40) // ".ExampleGet1" route, each 20 hours can send 40 times request from single Ip.
    if err != nil {
        log.Println(err)
    }

    server.Use(limitControl.GenerateLimitMiddleWare()) // Create middleWare

    server.POST("/ExamplePost1", func(ctx *gin.Context) {
        ctx.String(200, "Hello Example! In ExamplePost1")
    })

    server.GET("/ExampleGet1", func(ctx *gin.Context) {
        ctx.String(200, "Hello Example! In ExampleGet1")
    })
```


See more [Example](https://github.com/davidleitw/gin-limiter/blob/master/Example/example.go) and full code. 

<hr>

### Response 
- Request is pass, global limit or single route limit is legal. Then we return header with some limiter information. 
    ```shell
    Return header:

    X-RateLimit-Limit-global     -> limit request time which single ip can send request for the server. 
    X-RateLimit-Remaining-global -> remaining time which single ip can send request for the server.
    X-RateLimit-Reset-global     -> time for global limit reset. 

    X-RateLimit-Limit-single     -> limit request time which single ip can send request for the single route.
    X-RateLimit-Remaining-single -> remaining time which single ip can send request for the single route.
    X-RateLimit-Reset-single     -> time for single route limit reset. 

    ```


<br>

- When the global limit or single route limit is reached, a `429` HTTP status code is sent.
    and add the header with:
    ```shell
    Return header:
    
    If global remaining request time < 0
        X-RateLimit-Reset-global     -> return global limit reset time. 

    If single remaining request time < 0
        X-RateLimit-Reset-single     -> return this single route limit reset time.
    ```

<hr>

### Reference
- https://github.com/ulule/limiter
- https://github.com/jpillora/ipfilter
- https://github.com/KennyChenFight/dcard-simple-demo

<hr>

### Licenses

All source code is licensed under the [MIT License](https://github.com/davidleitw/gin-limiter/blob/master/LICENSE).

