package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/ncruces/go-sqlite3/driver" // Import SQLite3 driver
	_ "github.com/ncruces/go-sqlite3/embed"  // Import SQLite3 driver
)

type Item struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Rarity          string  `json:"rarity"`
	Details         string  `json:"details"`
	Percentage      float64 `json:"percentage"`
	ImageIdentifier string  `json:"image_identifier"`
}

// GachaRequest represents the request body for performing gacha.
type GachaRequest struct {
	Times int `json:"times"`
}

// GachaResponse represents the response body after performing gacha.
type GachaResponse struct {
	Results []Item `json:"results"`
}

var db *sql.DB
var rng *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng = rand.New(source)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./gacha.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブルが存在しない場合は作成する (image_path カラムがなければ追加)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			rarity TEXT NOT NULL,
			details TEXT,
			percentage REAL NOT NULL,
			image_identifier TEXT UNIQUE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// images フォルダが存在しない場合は作成
	os.MkdirAll("./images", 0755)

	r := gin.Default()

	r.StaticFS("/admin", http.Dir("./gacha-admin/dist"))

	// ガチャAPIのエンドポイント
	r.POST("/api/gacha", gachaHandler)

	// 画像取得API
	r.GET("/api/image", getImageHandler)

	// 管理画面用APIエンドポイント
	adminGroup := r.Group("/api/admin/items")
	{
		adminGroup.GET("", getItemListHandler)
		adminGroup.POST("", createItemHandler)
		adminGroup.PUT("/:id", updateItemHandler)
		adminGroup.DELETE("/:id", deleteItemHandler)
	}

	fmt.Println("Server started on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
