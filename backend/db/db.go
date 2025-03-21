package db

import (
	"database/sql"
	"fmt"
	"log"
	"server/types"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Connect(host string) *DB {
	con, err := sql.Open("postgres", host)
	if err != nil {
		panic(err)
	}

	if err := con.Ping(); err != nil {
		panic(err)
	}

	return &DB{
		con,
	}
}

func (db *DB) CreateVideo(v types.Video) (int, bool) {
	var id int
	row := db.QueryRow(`INSERT INTO "video" (title, fileSize, chunks, totalChunks, isCompiled, isSegmented) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, v.Title, v.Filesize, v.Chunks, v.TotalChunks, v.IsCompiled, v.IsSegmented)

	err := row.Scan(&id)
	if err != nil {
		log.Println(err)
		return -1, false
	}

	return int(id), true
}

func (db *DB) VideoExist(id int) bool {
	row := db.QueryRow(`SELECT * FROM "video" WHERE id = $1`, id)

	if row.Err() == sql.ErrNoRows {
		return false
	}

	return true
}

func (db *DB) GetVideo(id int) (types.Video, bool) {
	var video types.Video

	var chunks int64
	var fileSize int64
	var totalChunks int64

	row := db.QueryRow(`SELECT title, filesize, chunks, totalchunks, iscompiled, issegmented FROM "video" WHERE id = $1`, id)

	if err := row.Scan(&video.Title, &fileSize, &chunks, &totalChunks, &video.IsCompiled, &video.IsSegmented); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Video do not exist")
		}
		return video, false
	}

	video.Chunks = int(chunks)
	video.Filesize = int(fileSize)
	video.TotalChunks = int(totalChunks)

	return video, true
}

func (db *DB) GetVideos() []types.DbVideo {
	var videos []types.DbVideo
	rows, err := db.Query(`SELECT id, title, filesize, chunks, totalchunks, iscompiled, issegmented FROM "video";`)
	if err != nil {
		fmt.Println("Some error while quering videos")
	}

	for {

		if rows.Next() == false {
			break
		}
		var video types.DbVideo
		rows.Scan(&video.Id, &video.Title, &video.Filesize, &video.Chunks, &video.TotalChunks, &video.IsCompiled, &video.IsSegmented)

		videos = append(videos, video)
	}

	return videos
}

func (db *DB) UpdateVideo(v types.Video, id int) bool {
	_, err := db.Exec(`UPDATE "video" SET title = $1, fileSize = $2, chunks = $3, totalChunks = $4, isCompiled = $5, isSegmented = $6 WHERE Id = $7`, v.Title, v.Filesize, v.Chunks, v.TotalChunks, v.IsCompiled, v.IsSegmented, id)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
