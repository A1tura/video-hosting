package storage

import (
	"fmt"
	"os"
	video "server/Video"
	"server/db"
	"server/errors"
	"server/types"
	"strconv"
)

type Storage struct {
	db      *db.DB
	storage string
}

func CreateStorage(path string, db *db.DB) Storage {
	return Storage{
		storage: path,
		db:      db,
	}
}

func (s *Storage) CreateVideo(m types.UploadMetadataRequest) (int, bool) {
	var video types.Video

	video.Title = m.Title
	video.Filesize = m.FileSize
	video.IsCompiled = false
	video.IsSegmented = false
	video.Chunks = 0
	video.TotalChunks = m.Chunks

	return s.db.CreateVideo(video)
}

func DirExist(path string, create bool) bool {
	_, err := os.ReadDir(path)
	if err != nil {
		if create {
			if err := os.Mkdir(path, 0750); err != nil {
				return false
			}

			return true
		}
		return false
	}

	return true
}

func (s *Storage) AddChunk(id int, chunk []byte, chunkId int) bool {
	videoDB, _ := s.db.GetVideo(id)
	fmt.Println(videoDB)
	DirExist(s.storage+"/"+strconv.Itoa(id), true)
	chunkFile, err := os.Create(s.storage + "/" + strconv.Itoa(id) + "/" + strconv.Itoa(chunkId) + ".chunk")
	defer chunkFile.Close()
	if err != nil {
		return false
	}

	if _, err := chunkFile.Write(chunk); err != nil {
		fmt.Println(err)
		return false
	}

	videoDB.Chunks++

	if videoDB.Chunks == videoDB.TotalChunks {
		if videoDB.Compile(s.storage + "/" + strconv.Itoa(id)) {
			videoDB.IsCompiled = true

			videoSegments := video.Video{
				Input:  s.storage + "/" + strconv.Itoa(id) + "/" + "video.mp4",
				Output: s.storage + "/" + strconv.Itoa(id),
			}

			if videoSegments.Segment() {
				videoDB.IsSegmented = true
				s.db.UpdateVideo(videoDB, id)
				return true
			} else {
				s.db.UpdateVideo(videoDB, id)
				return false
			}
		}
		return false
	}

	s.db.UpdateVideo(videoDB, id)
	return true
}

func (s *Storage) GetFileStorage() string {
	return s.storage
}

func (s *Storage) VideoAvaliable(id int) int {
	video, exist := s.db.GetVideo(id)
	if !exist {
		return errors.NOT_EXIST
	}

	if !video.IsSegmented {
		return errors.NOT_AVALIABE
	}

	return 0
}
