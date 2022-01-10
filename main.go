package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const serveIp = "0.0.0.0"
const serveDefaultPort = "8080"

const videoUrlPrefix = "/video/"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "A non-recommending and private alternative frontend for Bilibili (B站)")
}

//go:embed templates/*
var resources embed.FS
var t = template.Must(template.ParseFS(resources, "templates/*"))

func videoHandler(w http.ResponseWriter, r *http.Request) {
	bvid := r.URL.Path[len(videoUrlPrefix):]

	videoInfo, err := GetVideoInfo(bvid)
	if err != nil {
		fmt.Fprintf(w, "Failed to get video info for %s, %v", bvid, err)
		return
	}

	cid := videoInfo.Cid
	videoStream, err := GetVideoStream(bvid, cid)
	if err != nil {
		fmt.Fprintf(w, "Failed to get video stream for %s/%v, %v", bvid, cid, err)
		return
	}

	data := map[string]string{
		"Title":    videoInfo.Title,
		"VideoUrl": videoStream.urls[0],
	}
	t.ExecuteTemplate(w, "index.html", data)
}

func main() {
	http.HandleFunc(videoUrlPrefix, videoHandler)
	http.HandleFunc("/", indexHandler)
	servePort := os.Getenv("PORT")
	if servePort == "" {
		servePort = serveDefaultPort
	}
	listenOn := fmt.Sprintf("%s:%v", serveIp, servePort)
	log.Println(fmt.Sprintf("Listening on %s", listenOn))
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
