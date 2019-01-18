package main

import (
	"flag"
	"log"
	"runtime"
	"strings"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/vharitonsky/iniflags"
	"golang.org/x/crypto/bcrypt"
)

var cost int

func main() {
	bind := flag.String("bind", ":8881", "address to bind baas (can be a unix domain socket: /var/run/baas.sock)")
	flag.IntVar(&cost, "cost", 10, "bcrypt salt cost")
	threads := flag.Int("threads", runtime.NumCPU(), "max process (default: cpu count)")

	iniflags.Parse()

	runtime.GOMAXPROCS(*threads)

	router := routing.New()
	router.Get("/hash", handleHash)
	router.Get("/verify", handleVerify)

	log.Printf("Start http server on %s", *bind)
	err := startServer(*bind, router.HandleRequest)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func startServer(bind string, handler func(ctx *fasthttp.RequestCtx)) error {
	server := &fasthttp.Server{
		Handler: handler,
		Name:    "valencia",
	}
	if strings.HasPrefix(bind, "/") {
		return server.ListenAndServeUNIX(bind, 0777)
	} else {
		return server.ListenAndServe(bind)
	}
}

func handleHash(c *routing.Context) error {
	raw := c.QueryArgs().Peek("raw")
	if raw == nil {
		c.Error("Invalid request", 400)
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword(raw, cost)
	if err != nil {
		return err
	}
	_, _ = c.Write(hash)
	return nil
}

func handleVerify(c *routing.Context) error {
	raw := c.QueryArgs().Peek("raw")
	hash := c.QueryArgs().Peek("hash")
	if raw == nil || hash == nil {
		c.Error("Invalid request", 400)
		return nil
	}
	err := bcrypt.CompareHashAndPassword(hash, raw)
	if err == nil {
		_, _ = c.WriteString("OK")
	} else {
		_, _ = c.WriteString("FAIL")
	}
	return nil
}
