package main

import (
	"flag"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/vharitonsky/iniflags"
	"golang.org/x/crypto/bcrypt"
)

var defaultCost int

func main() {
	bind := flag.String("bind", ":8881", "address to bind baas (can be a unix domain socket: /var/run/baas.sock)")
	flag.IntVar(&defaultCost, "cost", 10, "bcrypt salt cost")
	threads := flag.Int("threads", runtime.NumCPU(), "max process (default: cpu count)")

	iniflags.Parse()

	runtime.GOMAXPROCS(*threads)

	r := router.New()
	r.GET("/hash", handleHash)
	r.GET("/verify", handleVerify)

	log.Printf("Start http server on %s", *bind)
	err := startServer(*bind, r.Handler)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func startServer(bind string, handler func(ctx *fasthttp.RequestCtx)) error {
	server := &fasthttp.Server{
		Handler: handler,
		Name:    "go-baas",
	}
	if strings.HasPrefix(bind, "/") {
		return server.ListenAndServeUNIX(bind, 0777)
	} else {
		return server.ListenAndServe(bind)
	}
}

func handleHash(ctx *fasthttp.RequestCtx) {
	raw := ctx.QueryArgs().Peek("raw")
	if raw == nil {
		ctx.Error("Invalid request", 400)
		return
	}
	cost := defaultCost
	if costArg := ctx.QueryArgs().Peek("cost"); costArg != nil {
		parsed, err := strconv.Atoi(string(costArg))
		if err != nil {
			ctx.Error("Invalid request", 400)
			return
		}
		cost = parsed
	}
	if cost > 15 || cost < 2 {
		ctx.Error("Invalid cost. It must be in range (2, 15)", 400)
		return
	}
	hash, err := bcrypt.GenerateFromPassword(raw, cost)
	if err != nil {
		ctx.Error("Could not hash password:"+err.Error(), 500)
		return
	}
	_, _ = ctx.Write(hash)
}

func handleVerify(ctx *fasthttp.RequestCtx) {
	raw := ctx.QueryArgs().Peek("raw")
	hash := ctx.QueryArgs().Peek("hash")
	if raw == nil || hash == nil {
		ctx.Error("Invalid request", 400)
		return
	}
	err := bcrypt.CompareHashAndPassword(hash, raw)
	if err == nil {
		_, _ = ctx.WriteString("OK")
	} else {
		_, _ = ctx.WriteString("FAIL")
	}
}
