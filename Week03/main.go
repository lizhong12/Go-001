package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world！")
}

func main()  {
	sign := make(chan os.Signal)
	g, ctx := errgroup.WithContext(context.Background())
	signal.Notify(sign)
	g.Go(func() error {
		s := http.Server{
			Addr:              ":9090",
			Handler:           nil,
		}
		http.HandleFunc("/", sayHello)
		go func() {
			select {
			case <- ctx.Done():
				s.Shutdown(ctx)
				fmt.Println("ctx done...", ctx.Err())
			case <-sign:
				s.Shutdown(ctx)
				fmt.Println("signal kill", ctx.Err())

			}

		}()
		return s.ListenAndServe()

	})

	err := g.Wait()
	if err != nil {
		fmt.Println("err:",err)
	}
}
