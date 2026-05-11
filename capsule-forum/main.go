package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	gemini "git.sr.ht/~adnano/go-gemini"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Hostname string `toml:"hostname"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
	DBPath   string `toml:"db_path"`
}

func main() {
	cfgPath := "config.toml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	var cfg Config
	if _, err := toml.DecodeFile(cfgPath, &cfg); err != nil {
		log.Fatalf("config: %v", err)
	}
	if cfg.Port == 0 {
		cfg.Port = 1965
	}
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}

	db, err := openDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatalf("tls: %v", err)
	}

	mux := buildMux(db)

	var server gemini.Server
	server.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server.Handler = mux
	server.GetCertificate = func(_ string) (*tls.Certificate, error) {
		return &cert, nil
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("capsule-forum listening on %s:%d for %s", cfg.Host, cfg.Port, cfg.Hostname)
	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatal(err)
	}
}
