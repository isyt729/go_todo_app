package main

import (
	"context"
	"log"
	"net/http"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct{
	srv *http.Server 
	l net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server{
	return &Server{
		srv: &http.Server{Handler: mux}, 
		l: l,
	}
}

func (s *Server) Run(ctx context.Context) error{
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error{
		if err := s.srv.Serve(s.l); err != nil &&
			err != http.ErrServerClosed {
				log.Printf("faild to close server: %+v",err)
				return err
		}
		return nil
	})

	// チャネルからの終了通知を待機する
	<-ctx.Done()

	if err := s.srv.Shutdown(context.Background()); err != nil{
		log.Printf("faild to shutdown : %v",err)
	}

	return eg.Wait()
}