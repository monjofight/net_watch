package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/monjofight/net_watch/pkg/netflix"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Initializing database...")

	// titlesテーブルの作成
	statement1, _ := database.Prepare(`
		CREATE TABLE IF NOT EXISTS titles (
			id INT PRIMARY KEY, 
			name TEXT, 
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`)
	statement1.Exec()

	// log.Println("1. titlesテーブルの作成完了")

	// seasonsテーブルの作成
	statement2, _ := database.Prepare(`
		CREATE TABLE IF NOT EXISTS seasons (
			id INT PRIMARY KEY, 
			title_id INT,
			name TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY(title_id) REFERENCES titles(id)
		);
	`)
	statement2.Exec()

	// log.Println("2. seasonsテーブルの作成完了")

	// episodesテーブルの作成
	statement3, _ := database.Prepare(`
		CREATE TABLE IF NOT EXISTS episodes (
			id SERIAL PRIMARY KEY, 
			title_id INT,
			season_id INT,
			name TEXT,
			image TEXT,
			watched BOOLEAN DEFAULT FALSE, 
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY(title_id) REFERENCES titles(id),
			FOREIGN KEY(season_id) REFERENCES seasons(id)
		);
	`)
	statement3.Exec()

	// log.Println("3. episodesテーブルの作成完了")

	return database
}

func SaveToDB(title *netflix.Title, database *sql.DB) {
	currentTime := time.Now()
	tx, err := database.Begin()
	if err != nil {
		log.Fatal("Failed to start transaction:", err)
		return
	}
	defer tx.Rollback()

	// 既存のタイトル、シーズン、エピソードのIDを取得
	existingTitleIDs := make(map[int]bool)
	existingSeasonIDs := make(map[int]bool)
	existingEpisodeIDs := make(map[int]bool)

	rows, _ := tx.Query("SELECT id FROM titles")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		existingTitleIDs[id] = true
	}

	rows, _ = tx.Query("SELECT id FROM seasons")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		existingSeasonIDs[id] = true
	}

	rows, _ = tx.Query("SELECT id FROM episodes")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		existingEpisodeIDs[id] = true
	}

	// タイトルの確認と挿入
	if !existingTitleIDs[title.ID] {
		_, err := tx.Exec("INSERT INTO titles (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)", title.ID, title.Name, currentTime, currentTime)
		if err != nil {
			log.Println("Failed to insert title:", err)
			return
		}
	}

	var seasonValues []string
	var seasonArgs []interface{}
	var episodeValues []string
	var episodeArgs []interface{}
	seasonIndex := 1
	episodeIndex := 1

	for _, season := range title.Seasons {
		if !existingSeasonIDs[season.ID] {
			seasonValues = append(seasonValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", seasonIndex, seasonIndex+1, seasonIndex+2, seasonIndex+3, seasonIndex+4))
			seasonArgs = append(seasonArgs, season.ID, title.ID, season.Name, currentTime, currentTime)
			seasonIndex += 5
		}

		for _, episode := range season.Episodes {
			if !existingEpisodeIDs[episode.ID] {
				episodeValues = append(episodeValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", episodeIndex, episodeIndex+1, episodeIndex+2, episodeIndex+3, episodeIndex+4, episodeIndex+5, episodeIndex+6, episodeIndex+7))
				episodeArgs = append(episodeArgs, episode.ID, title.ID, season.ID, episode.Name, episode.Image, false, currentTime, currentTime)
				episodeIndex += 8
			}
		}
	}

	// シーズンとエピソードをデータベースに追加
	if len(seasonValues) > 0 {
		seasonStmt := fmt.Sprintf("INSERT INTO seasons (id, title_id, name, created_at, updated_at) VALUES %s", strings.Join(seasonValues, ","))
		_, err = tx.Exec(seasonStmt, seasonArgs...)
		if err != nil {
			log.Println("Failed to insert seasons:", err)
			return
		}
	}

	if len(episodeValues) > 0 {
		episodeStmt := fmt.Sprintf("INSERT INTO episodes (id, title_id, season_id, name, image, watched, created_at, updated_at) VALUES %s", strings.Join(episodeValues, ","))
		_, err = tx.Exec(episodeStmt, episodeArgs...)
		if err != nil {
			log.Println("Failed to insert episodes:", err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}
}

func GetTitles(db *sql.DB) []netflix.Title {
	rows, err := db.Query("SELECT id, name FROM titles")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var titles []netflix.Title
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		titles = append(titles, netflix.Title{ID: id, Name: name})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return titles
}
