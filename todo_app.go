package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net"
	"os"

	"golang.org/x/sync/errgroup"
)

func main(){
	// コマンドライン引数にてポート番号を指定する必要があるため、未指定の場合終了
	if len(os.Args) != 2{
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil{
		log.Fatalf("faild to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil{
		log.Printf("faild to terminate server: %v",err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error{
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Hello, %s!",r.URL.Path[1:])
		}),
	}

	eg, ctx :=errgroup.WithContext(ctx)
	eg.Go(func() error{
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
				log.Printf("faild to close server: %+v",err)
				return err
		}
		return nil
	})

	// チャネルからの終了通知を待機する
	<-ctx.Done()

	if err := 	s.Shutdown(context.Background()); err != nil{
		log.Printf("faild to shutdown : %v",err)
	}

	return eg.Wait()
}