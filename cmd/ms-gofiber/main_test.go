package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"ms-gofiber/internal/config"
)

type fakeServer struct {
	listenErr   error
	shutdownErr error
	listenBlock chan struct{}
}

func (s *fakeServer) Listen(string) error {
	if s.listenBlock != nil {
		<-s.listenBlock
	}
	return s.listenErr
}

func (s *fakeServer) ShutdownWithContext(context.Context) error {
	if s.listenBlock != nil {
		close(s.listenBlock)
		s.listenBlock = nil
	}
	return s.shutdownErr
}

func TestRunBranches(t *testing.T) {
	origLoad := loadConfig
	origBuild := buildApp
	origNotify := notifySignal
	origWithTimeout := withTimeout
	t.Cleanup(func() {
		loadConfig = origLoad
		buildApp = origBuild
		notifySignal = origNotify
		withTimeout = origWithTimeout
	})

	cfg := &config.Config{AppHost: "127.0.0.1", AppPort: 18080}

	loadConfig = func() (*config.Config, error) { return nil, errors.New("cfg") }
	if err := run(); err == nil || !strings.Contains(err.Error(), "config load error") {
		t.Fatalf("expected config load error, got %v", err)
	}

	loadConfig = func() (*config.Config, error) { return cfg, nil }
	buildApp = func(*config.Config) (server, func(), error) { return nil, nil, errors.New("build") }
	if err := run(); err == nil || !strings.Contains(err.Error(), "app build error") {
		t.Fatalf("expected build error, got %v", err)
	}

	buildApp = func(*config.Config) (server, func(), error) {
		return &fakeServer{listenErr: errors.New("listen")}, func() {}, nil
	}
	if err := run(); err == nil || !strings.Contains(err.Error(), "fiber listen error") {
		t.Fatalf("expected listen error, got %v", err)
	}

	buildApp = func(*config.Config) (server, func(), error) {
		return &fakeServer{shutdownErr: errors.New("shutdown"), listenBlock: make(chan struct{})}, func() {}, nil
	}
	notifySignal = func(c chan<- os.Signal, sig ...os.Signal) {
		go func() { c <- os.Interrupt }()
	}
	if err := run(); err == nil || !strings.Contains(err.Error(), "fiber shutdown error") {
		t.Fatalf("expected shutdown error, got %v", err)
	}

	buildApp = func(*config.Config) (server, func(), error) {
		return &fakeServer{listenBlock: make(chan struct{})}, func() {}, nil
	}
	withTimeout = func(parent context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
		return context.WithCancel(parent)
	}
	if err := run(); err != nil {
		t.Fatalf("expected success run, got %v", err)
	}
}

func TestMainFunction(t *testing.T) {
	origRunMain := runMain
	origFatalf := fatalf
	t.Cleanup(func() {
		runMain = origRunMain
		fatalf = origFatalf
	})

	fatalCalled := false
	fatalMsg := ""
	fatalf = func(format string, v ...any) {
		fatalCalled = true
		fatalMsg = format
	}

	runMain = func() error { return nil }
	main()
	if fatalCalled {
		t.Fatalf("fatal should not be called on nil error")
	}

	runMain = func() error { return errors.New("boom") }
	main()
	if !fatalCalled || !strings.Contains(fatalMsg, "%v") {
		t.Fatalf("expected fatal called with format, called=%v msg=%s", fatalCalled, fatalMsg)
	}
}
