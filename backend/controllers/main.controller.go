package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/db"
	"server/errors"
	"server/storage"
	"server/types"
	"strconv"
)

func UploadMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*") // change this later
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT")

		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		storage := r.Context().Value("s").(*storage.Storage)
		var req types.UploadMetadataRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			fmt.Fprint(w, "Error while reading body")
			return
		}

		id, created := storage.CreateVideo(req)
		if !created {
			fmt.Fprint(w, "Error while creating video medatada file")
			return
		}

		res := types.UploadMetadataResponse{
			Id: id,
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(503)
			fmt.Println("Internal server error")
			return
		}
	}
}

func UploadVideoChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*") // change this later
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT")

		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		storage := r.Context().Value("s").(*storage.Storage)
		var req types.UploadVideoChunkRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			fmt.Fprint(w, "Error while reading body")
			return
		}

		added := storage.AddChunk(req.Id, req.Chunk, req.ChunkId)
		if !added {
			fmt.Fprint(w, "Error while creating video chunk file")
			return
		}

		fmt.Fprint(w, types.UploadVideoChunkResponse{
			Status: true,
		})
	}
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		storage := r.Context().Value("s").(*storage.Storage)
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "Video do not exist")
			return
		}

		if storage.VideoAvaliable(id) == errors.NOT_EXIST {
			w.WriteHeader(404)
			fmt.Fprint(w, "Video do not exist")
			return
		} else if storage.VideoAvaliable(id) == errors.NOT_AVALIABE {
			w.WriteHeader(403)
			fmt.Fprint(w, "Video it's not avaliable now")
			return
		} else {
			var playlistContent []byte
			playlist, err := os.Open(storage.GetFileStorage() + "/" + idStr + "/" + "playlist.m3u8")
			if err != nil {
				w.WriteHeader(503)
				fmt.Fprint(w, "Internal server error")
				return
			}
			playlistStats, err := os.Stat(storage.GetFileStorage() + "/" + idStr + "/" + "playlist.m3u8")
			if err != nil {
				playlistContent = make([]byte, 5000)
			} else {
				playlistContent = make([]byte, playlistStats.Size())
			}

			if _, err := playlist.Read(playlistContent); err != nil {
				w.WriteHeader(503)
				fmt.Fprint(w, "Internal server error")
				return
			}

			w.Header().Add("Content-Type", "application/vnd.apple.mpegurl")
			w.Header().Add("Content-Length", strconv.Itoa(len(playlistContent)))
			fmt.Fprint(w, string(playlistContent))
		}
	}
}

func GetSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		storage := r.Context().Value("s").(*storage.Storage)
		videoIdString := r.PathValue("id")
		videoId, err := strconv.Atoi(videoIdString)
		segmentName := r.PathValue("segment")
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "Video do not exist")
		}

		if storage.VideoAvaliable(videoId) == errors.NOT_EXIST {
			w.WriteHeader(404)
			fmt.Fprint(w, "Video do not exist")
			return
		} else if storage.VideoAvaliable(videoId) == errors.NOT_AVALIABE {
			w.WriteHeader(403)
			fmt.Fprint(w, "Video it's not avaliable now")
			return
		} else {
			var segmentContent []byte
			segment, err := os.Open(storage.GetFileStorage() + "/" + videoIdString + "/" + segmentName + ".ts")
			if err != nil {
				w.WriteHeader(404)
				fmt.Println("ffff")
				fmt.Fprint(w, "Video do not exist")
				return
			}
			segmentStats, err := segment.Stat()
			if err != nil {
				segmentContent = make([]byte, 1*1024*1024)
			} else {
				segmentContent = make([]byte, segmentStats.Size())
			}

			if _, err := segment.Read(segmentContent); err != nil {
				w.WriteHeader(503)
				fmt.Fprint(w, "Internal server error")
				return
			}
			w.Header().Add("Content-Type", "video/MP2T")
			w.Header().Add("Content-Length", strconv.Itoa(len(segmentContent)))
			fmt.Fprint(w, string(segmentContent))
		}
	}
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var res []types.VideosResponse
		db := r.Context().Value("db").(*db.DB)

		videos := db.GetVideos()

		for _, video := range videos {
			res = append(res, types.VideosResponse{
				Id:    video.Id,
				Title: video.Title,
			})
		}

		jsonRes, err := json.Marshal(res)
		if err != nil {
			fmt.Fprint(w, "Something going wrong")
			return
		}

        w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(jsonRes))
	}
}
