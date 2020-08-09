<h1 align="center">gin- limiter</h1>
<h5 align="center">A simple gin middleware for IP limiter based on redis.</h5>

<p align="center">
    <a href="https://www.gnu.org/licenses/"> 
        <img src="https://img.shields.io/github/license/davidleitw/goGamer.svg" alt="License">
    </a>
    <a href="http://hits.dwyl.com/davidleitw/gin-limiter">
        <img src=http://hits.dwyl.com/davidleitw/gin-limiter.svg alt="HitCount">
    </a>
    <a href="https://github.com/davidleitw/gin-limiter/stargazers"> 
        <img src="https://img.shields.io/github/stars/davidleitw/gin-limiter" alt="GitHub stars">
    </a>
</p>
<hr>
### Installation
```go 
go get github.com/davidleitw/gin-limiter 
``` 

**Import**
```go 
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

- ##### 搭配gin，針對每個不同的route做出獨立的限制
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
完整程式碼請看 [Example](https://github.com/davidleitw/gin-limiter/tree/master/Example)
<hr>
### Licenses

All source code is licensed under the [MIT License](https://github.com/davidleitw/gin-limiter/blob/master/LICENSE).

