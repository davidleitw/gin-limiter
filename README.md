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

- Create a limit middleware dispatcher object
    ```go
    // Set redis client
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

    dispatcher, err := limiter.LimitDispatcher("24-M", 100, rdb)

    if err != nil {
        log.Println(err)
    }

    ```

- Add a middleware to controlling each route. 
    ```go
    server := gin.Default()

    server.POST("/ExamplePost1", dispatcher.MiddleWare("4-M", 20), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello ExamplePost1")
	})

	server.GET("/ExampleGet1", dispatcher.MiddleWare("5-M", 10), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello ExampleGet1")
	})

	err = server.Run(":8080")
	if err != nil {
		log.Println("gin server error = ", err)
	}
    ```


    See more examples [HERE](https://github.com/davidleitw/gin-limiter/blob/master/Example). 

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
        return global limit reset time. 

    If single remaining request time < 0
        return this single route limit reset time.
    ```

<hr>

### Reference
- https://github.com/ulule/limiter
- https://github.com/jpillora/ipfilter
- https://github.com/KennyChenFight/dcard-simple-demo

<hr>

If you want to know the updated progress, please check `Ipfilter` branch. 

### License

All source code is licensed under the [MIT License](./LICENSE).

