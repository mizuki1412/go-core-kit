package snippet

import (
	"context"
	"errors"
	"github.com/mizuki1412/go-core-kit/v2/library/timekit"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func ContextCancelDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	// context 携带val
	ctx = context.WithValue(ctx, "k", "v")
	// 调用接口a
	err := childProcess(1, ctx)
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	// 调用接口b
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := childProcess(2, ctx)
		if err != nil {
			cancel()
		}
	}()
	// 调用接口c
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := childProcess(3, ctx)
		if err != nil {
			cancel()
		}
	}()
	wg.Wait()
}

func childProcess(id int, ctx context.Context) error {
	result := make(chan int)
	err := make(chan error)
	go func() {
		timekit.Sleep(3000)
		log.Println(ctx.Value("k"))
		if id == 3 {
			err <- errors.New("")
		} else {
			result <- 1
		}
	}()

	select {
	case <-ctx.Done():
		// 其他请求失败
		return ctx.Err()
	case e := <-err:
		// 本次请求失败，返回错误信息
		return e
	case <-result:
		// 本此请求成功，不返回错误信息
		return nil
	}
}

func ContextTimeoutDemo() {
	req, err := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	// 构造一个超时间为50毫秒的Context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
}
