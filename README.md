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

### Response 
- 請求符合規定, 無論是global limit，或者是該次request所造訪的single route limit皆沒有違反。 
    ```shell
    Return header:

    X-RateLimit-Limit-global     -> 單一Ip在期間內總共能造訪幾次 (global)
    X-RateLimit-Remaining-global -> 單一Ip剩餘造訪次數
    X-RateLimit-Reset-global     -> 下次重製剩餘次數的時間

    X-RateLimit-Limit-single     -> 單一Ip對於造訪route內總共能造訪幾次 (single)
    X-RateLimit-Remaining-single -> 單一Ip對於造訪route剩餘造訪次數
    X-RateLimit-Reset-single     -> 此route下次重製時間
    ```


<br>

- 請求違反的globol limit，或者是single route limit其中一項。
    回傳Http429(Too many Requests) 
    ```shell
    若是global的造訪次數已經用完, 則會回傳
        X-RateLimit-Reset-global     -> 下次重製剩餘次數的時間

    若是single route的造訪次數已經用完, 則會回傳
        X-RateLimit-Reset-single     -> 此route下次重製時間
    ```

<hr>

### Licenses

All source code is licensed under the [MIT License](https://github.com/davidleitw/gin-limiter/blob/master/LICENSE).

