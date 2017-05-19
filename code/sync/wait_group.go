var wg sync.WaitGroup
var urls = []string{
        "http://www.golang.org/",
        "http://www.google.com/",
        "http://www.somestupidname.com/",
}
for _, url := range urls {
        // 对 WaitGroup 计数器执行加一操作
        wg.Add(1)
        // 启动一个 goroutine ，用于获取给定的 URL 
        go func(url string) {
            // 在 goroutine 执行完毕时，对计数器执行减一操作
            defer wg.Done()
            // 获取 URL
            http.Get(url)
        }(url)
}

// 等待直到所有 HTTP 获取操作都执行完毕为止
wg.Wait()
