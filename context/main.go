package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		// 親コンテキストの生成
		ctx := context.Background()
		// キャンセルコンテキストの生成
		t := time.Now()
		t1 := t.Add(5 * time.Second)
		// 5秒経過したらキャンセルする
		deadlineCtx, cancel := context.WithDeadline(ctx, t1)
		defer cancel()

		wg := sync.WaitGroup{}

		for i := 1; i < 10; i++ {
			wg.Add(1)
			go func(deadlineCtx context.Context, num int) {
				defer wg.Done()
				select {
				case <-deadlineCtx.Done():
					return
				default:
					io.WriteString(w, fmt.Sprintf("Count: %d Second \n", num))
				}
			}(deadlineCtx, i)
			time.Sleep(1 * time.Second)
		}
		wg.Wait()
		io.WriteString(w, "time over")
	}

	http.HandleFunc("/hello", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
