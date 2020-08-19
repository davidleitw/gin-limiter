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
- Download

    Type the following command in your terminal.
    ```bash
    go get github.com/go-redis/redis/v8
    go get github.com/davidleitw/gin-limiter 
    ``` 

- Import
    ```go 
    import "github.com/go-redis/redis/v8"
    import limiter "github.com/davidleitw/gin-limiter"
    ```

---

### Quickstart

- Create limit controller object
    ```go
    // Set redis client
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})


    // Debug mode, each 24 minutes can send 100 times request from single Ip.
    limitControl, err := limiter.DefaultController(rdb, "24-M", 100, "debug")

    if err != nil {
        log.Println(err)
    }

    ```

    Debug mode will display some information on the terminal.
    ![](https://imgur.com/KeZsQpQ.png)

- For each route, add a sub-limiter. 
    ```go
    server := gin.Default()

    // "/ExamplePost1" route, allows 20 requests every 4 minutes from a single IP address
    err = limitControl.Add("/ExamplePost1", "POST", "4-M", 20)

    if err != nil {
        log.Println(err)
    }


    // ".ExampleGet1" route, allows 40 ruquests every 20 hours from a single IP address.
    err = limitControl.Add("/ExampleGet1", "GET", "20-H", 40)

    if err != nil {
        log.Println(err)
    }


    // Create middleware
    server.Use(limitControl.GenerateLimitMiddleWare()) 

    server.POST("/ExamplePost1", func(ctx *gin.Context) {
        ctx.String(200, "Hello Example! In ExamplePost1")
    })

    server.GET("/ExampleGet1", func(ctx *gin.Context) {
        ctx.String(200, "Hello Example! In ExampleGet1")
    })
    ```


    See more examples [HERE](https://github.com/davidleitw/gin-limiter/blob/master/Example/example.go). 

---

### Response 
- When the total of request times is within limit, we will write data to header.
    ```
    Return header:

    X-RateLimit-Limit-global     -> Request limit of a single ip can send request for the server. 
    X-RateLimit-Remaining-global -> Remaining times which single ip can send request for the server.
    X-RateLimit-Reset-global     -> Time to global limit reset. 

    X-RateLimit-Limit-single     -> Request limit of a single ip can send request for the single route.
    X-RateLimit-Remaining-single -> Remaining times which single ip can send request for the single route.
    X-RateLimit-Reset-single     -> Time to single route limit reset. 

    ```

- When global limit or single route limit is reached, a `429` HTTP status code is sent.
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

If you want to know the updated progress, please check noSingleVersion branch. 

### License

All source code is licensed under the [MIT License](./LICENSE).

