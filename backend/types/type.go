package types

import (
	"fmt"
	"os"
	"strconv"
)

type UploadMetadataRequest struct {
	Title    string `json:"title"`
	FileSize int    `json:"fileSize"`
	Chunks   int    `json:"chunks"`
}

type UploadMetadataResponse struct {
	Id int `json:"id"`
}

type UploadVideoChunkRequest struct {
	Id      int    `json:"id"`
	ChunkId int    `json:"chunkId"`
	Chunk   []byte `json:"chunk"`
}

type UploadVideoChunkResponse struct {
	Status bool `json:"status"`
}

type VideosResponse struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Video struct {
	Title    string
	Filesize int

	Chunks      int
	TotalChunks int

	IsCompiled  bool
	IsSegmented bool
}

type DbVideo struct {
	Id       int
	Title    string
	Filesize int

	Chunks      int
	TotalChunks int

	IsCompiled  bool
	IsSegmented bool
}

func (v *Video) Compile(path string) bool {
	file, err := os.Create(path + "/" + "video.mp4")
	defer file.Close()
	if err != nil {
		return false
	}

	for i := 0; i < v.Chunks; i++ {
		chunkPath := path + "/" + strconv.Itoa(i) + ".chunk"
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return false
		}
		chunkStats, err := os.Stat(chunkPath)
		if err != nil {
			return false
		}

		chunk := make([]byte, chunkStats.Size())

		if _, err := chunkFile.Read(chunk); err != nil {
			return false
		}

		file.Write(chunk)
		chunkFile.Close()
		os.Remove(chunkPath)
	}

	fmt.Println("File has been created")

	return true
}
