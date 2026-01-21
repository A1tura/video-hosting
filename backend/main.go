package main

import (
	"context"
	"net/http"
	"server/controllers"
	"server/db"
	"server/storage"
)

var dbCon = db.Connect("user=admin password=admin dbname=stream host=db port=5432 sslmode=disable")
var s = storage.CreateStorage("./videos", dbCon)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // change this later
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "s", &s)
        ctx = context.WithValue(ctx, "db", dbCon)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {

	http.Handle("/uploadMetadata", CORSMiddleware(CtxMiddleware(http.HandlerFunc(controllers.UploadMetadata))))
	http.Handle("/uploadVideoChunk", CORSMiddleware(CtxMiddleware(http.HandlerFunc(controllers.UploadVideoChunk))))
	http.Handle("/video/{id}", CORSMiddleware(CtxMiddleware(http.HandlerFunc(controllers.GetPlaylist))))
	http.Handle("/segment/{id}/{segment}", CORSMiddleware(CtxMiddleware(http.HandlerFunc(controllers.GetSegment))))

	http.Handle("/videos", CORSMiddleware(CtxMiddleware(http.HandlerFunc(controllers.GetVideos))))

	http.ListenAndServe("0.0.0.0:8081", nil)
}
