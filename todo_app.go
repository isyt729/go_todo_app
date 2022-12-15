package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"golang.org/x/sync/errgroup"
	"github.com/isyt729/go_todo_app/config"
)

func main(){
	if err := run(context.Background()); err != nil{
		log.Printf("faild to terminate server: %v",err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error{
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := config.New()
	if err != nil{
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil{
		log.Fatalf("faild to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %v", url)

	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			// リクエストを遅延させて、
			// time.Sleep(10 * time.Second)
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