package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
)

func Serve() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthz", healthz)

	pusher := router.PathPrefix("/push").Methods("POST").Subrouter()
	pusher.Use(authMiddleware)
	pusher.HandleFunc("/txt", pushTxt)
	pusher.HandleFunc("/file", pushFile)
	pusher.HandleFunc("/image", pushImage)

	srv := &http.Server{
		Handler: router,
		Addr:    conf.Addr,
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		var shutdownOnce sync.Once
		for range sigs {
			logger.Warn("Receiving the termination signal, graceful shutdown...")
			shutdownOnce.Do(func() {
				err := srv.Shutdown(context.Background())
				if err != nil {
					logger.Error(err)
				}
			})
		}
	}()
	go bot.ID()

	logger.Infof("Server Listen at %s", conf.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("server shutdown success.")
		} else {
			logger.Fatal(err)
		}
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil || r.Header.Get("X-Token") != conf.Token {
			logger.Warnf("Forbidden: IP[%s], URI[%]", getIP(r), r.RequestURI)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func pushTxt(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Errorf("[txt] %s", err)
		return
	}
	if strings.TrimSpace(r.FormValue("type")) != "txt" {
		http.Error(w, "bad request type", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(r.FormValue("message"))
	if msg == "" {
		http.Error(w, "message is empty", http.StatusBadRequest)
		return
	}

	to := strings.TrimSpace(r.FormValue("recipient"))
	if to == "" {
		http.Error(w, "recipient is empty", http.StatusBadRequest)
		return
	}

	var md bool
	if strings.TrimSpace(r.FormValue("markdown")) == "true" {
		md = true
	}
	recipient := strings.Split(to, ",")
	for _, t := range recipient {
		id, err := strconv.Atoi(t)
		if err != nil {
			logger.Error("[txt] invalid recipient ID: %s", t)
			continue
		}
		go func(bot *Telegram, id int64, msg string, md bool) {
			err := bot.SendMessage(msg, id, md)
			if err != nil {
				logger.Errorf("[txt] message send failed, id: %d", id)
			}
		}(bot, int64(id), msg, md)
	}
}

func pushFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Errorf("[file] %s", err)
		return
	}
	if strings.TrimSpace(r.FormValue("type")) != "file" {
		http.Error(w, "bad request type", http.StatusBadRequest)
		return
	}

	filename := strings.TrimSpace(r.FormValue("filename"))
	if filename == "" {
		http.Error(w, "filename is empty", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Errorf("[file] receive file failed: %s", err)
		http.Error(w, "receive file failed", http.StatusInternalServerError)
		return
	}
	defer func() { _ = file.Close() }()

	mime := strings.TrimSpace(r.FormValue("mime"))
	caption := strings.TrimSpace(r.FormValue("caption"))

	to := strings.TrimSpace(r.FormValue("recipient"))
	if to == "" {
		http.Error(w, "recipient is empty", http.StatusBadRequest)
		return
	}

	recipient := strings.Split(to, ",")
	for _, t := range recipient {
		id, err := strconv.Atoi(t)
		if err != nil {
			logger.Error("[file] invalid recipient ID: %s", t)
			continue
		}
		go func(bot *Telegram, id int64, file io.Reader, filename, mime, caption string) {
			err := bot.SendFile(file, filename, mime, caption, id)
			if err != nil {
				logger.Errorf("[file] file send failed, id: %d", id)
			}
		}(bot, int64(id), file, filename, mime, caption)
	}
}

func pushImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Errorf("[image] %s", err)
		return
	}
	if strings.TrimSpace(r.FormValue("type")) != "file" {
		http.Error(w, "bad request type", http.StatusBadRequest)
		return
	}

	image, _, err := r.FormFile("image")
	if err != nil {
		logger.Errorf("[image] receive image failed: %s", err)
		http.Error(w, "receive image failed", http.StatusInternalServerError)
		return
	}
	defer func() { _ = image.Close() }()

	caption := strings.TrimSpace(r.FormValue("caption"))

	to := strings.TrimSpace(r.FormValue("recipient"))
	if to == "" {
		http.Error(w, "recipient is empty", http.StatusBadRequest)
		return
	}

	recipient := strings.Split(to, ",")
	for _, t := range recipient {
		id, err := strconv.Atoi(t)
		if err != nil {
			logger.Error("[image] invalid recipient ID: %s", t)
			continue
		}
		go func(bot *Telegram, id int64, image io.Reader, caption string) {
			err := bot.SendImage(image, caption, id)
			if err != nil {
				logger.Errorf("[image] image send failed, id: %d", id)
			}
		}(bot, int64(id), image, caption)
	}
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
