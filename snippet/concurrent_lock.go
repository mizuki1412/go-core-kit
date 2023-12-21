package snippet

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/concurrentkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"sync"
)

/**
简单的秒杀
*/

var stock = 100 // 商品库存
var mu sync.Mutex

func MiaoSha() {
	g := concurrentkit.NewGroup()
	for i := 0; i < 1000; i++ {
		g.Add(func() {
			user := i
			for {
				if mu.TryLock() {
					handle(user)
					break
				} else {
					timekit.Sleep(10)
					continue
				}
			}
		}, false)
	}
	g.Process()
}

func handle(user int) {
	defer mu.Unlock()
	if stock > 0 {
		stock--
		fmt.Printf("用户%d秒杀成功，剩余库存：%d\n", user, stock)
	}
}
