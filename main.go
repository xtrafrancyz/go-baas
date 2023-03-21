package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"

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

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", handleHash)
	mux.HandleFunc("/verify", handleVerify)

	log.Printf("Start http server on %s", *bind)
	var err error
	if strings.HasPrefix(*bind, "/") {
		unixListener, err := net.Listen("unix", *bind)
		if err == nil {
			err = http.Serve(unixListener, mux)
		}
	} else {
		err = http.ListenAndServe(*bind, mux)
	}
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func handleHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	cost := defaultCost
	raw := r.Form.Get("raw")
	if r.Form.Has("cost") {
		cost0, err := strconv.Atoi(r.Form.Get("cost"))
		if err != nil {
			http.Error(w, "Invalid cost", 400)
			return
		}
		cost = cost0
	}
	if raw == "" {
		http.Error(w, "Missing raw password", 400)
		return
	}
	if cost > 16 || cost < 5 {
		http.Error(w, "Invalid cost. It must be in range (5, 16)", 400)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), cost)
	if err != nil {
		http.Error(w, "Could not hash password:"+err.Error(), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(hash)
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	raw := r.Form.Get("raw")
	hash := r.Form.Get("hash")
	if raw == "" || hash == "" {
		http.Error(w, "Invalid request", 400)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	w.WriteHeader(200)
	if err == nil {
		_, _ = w.Write([]byte("OK"))
	} else {
		_, _ = w.Write([]byte("FAIL"))
	}
}
