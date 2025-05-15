package main

import (
	"database/sql"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ガチャ.
func gachaHandler(c *gin.Context) {
	var req GachaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Times <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number of times must be greater than 0"})
		return
	}

	items, err := getAllItemsFromDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	results := []Item{}
	for i := 0; i < req.Times; i++ {
		item := drawItem(items)
		if item != nil {
			results = append(results, *item)
		}
	}

	c.JSON(http.StatusOK, GachaResponse{Results: results})
}

// アイテムリスト.
func getItemListHandler(c *gin.Context) {
	items, err := getAllItemsFromDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// アイテム作成.
func createItemHandler(c *gin.Context) {
	var newItem Item
	newItem.Name = c.PostForm("name")
	newItem.Rarity = c.PostForm("rarity")
	newItem.Details = c.PostForm("details")
	percentageStr := c.PostForm("percentage")
	if percentageStr != "" {
		percentage, err := strconv.ParseFloat(percentageStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid percentage value"})
			return
		}
		newItem.Percentage = percentage
	}

	file, err := c.FormFile("image")
	var imageIdentifier string
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer src.Close()

		img, _, err := image.Decode(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			return
		}

		id := uuid.New().String()
		imageIdentifier = id
		pngFilename := filepath.Join("./images", id+".png")
		err = imaging.Save(img, pngFilename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save as PNG"})
			return
		}
	}

	result, err := db.Exec(
		"INSERT INTO items (name, rarity, details, percentage, image_identifier) VALUES (?, ?, ?, ?, ?)",
		newItem.Name, newItem.Rarity, newItem.Details, newItem.Percentage, imageIdentifier,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
	}
	newItem.ID = int(insertedID)
	newItem.ImageIdentifier = imageIdentifier

	c.JSON(http.StatusCreated, newItem)
}

// アイテム更新.
func updateItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var updatedItem Item
	updatedItem.ID = itemID
	updatedItem.Name = c.PostForm("name")
	updatedItem.Rarity = c.PostForm("rarity")
	updatedItem.Details = c.PostForm("details")
	percentageStr := c.PostForm("percentage")
	if percentageStr != "" {
		percentage, err := strconv.ParseFloat(percentageStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid percentage value"})
			return
		}
		updatedItem.Percentage = percentage
	}

	file, err := c.FormFile("image")
	var imageIdentifier string
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer src.Close()

		img, _, err := image.Decode(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			return
		}

		newID := uuid.New().String()
		imageIdentifier = newID
		pngFilename := filepath.Join("./images", newID+".png")
		err = imaging.Save(img, pngFilename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save as PNG"})
			return
		}

		// 古い画像を削除 (任意)
		var oldImageIdentifier string
		err = db.QueryRow("SELECT image_identifier FROM items WHERE id = ?", itemID).Scan(&oldImageIdentifier)
		if err == nil && oldImageIdentifier != "" {
			oldImagePath := filepath.Join("./images", oldImageIdentifier+".png")
			os.Remove(oldImagePath)
		}
	} else {
		err := db.QueryRow("SELECT image_identifier FROM items WHERE id = ?", itemID).Scan(&imageIdentifier)
		if err != nil && err != sql.ErrNoRows {
			log.Println("Error fetching existing image identifier:", err)
		}
	}

	_, err = db.Exec(
		"UPDATE items SET name = ?, rarity = ?, details = ?, percentage = ?, image_identifier = ? WHERE id = ?",
		updatedItem.Name, updatedItem.Rarity, updatedItem.Details, updatedItem.Percentage, imageIdentifier, itemID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	updatedItem.ImageIdentifier = imageIdentifier
	c.JSON(http.StatusOK, updatedItem)
}

// アイテム削除
func deleteItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	_, err = db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Item with ID %d deleted", id)})
}

// アイテム情報をすべて取得する.
func getAllItemsFromDB() ([]Item, error) {
	rows, err := db.Query("SELECT id, name, rarity, details, percentage, image_identifier FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Rarity, &item.Details, &item.Percentage, &item.ImageIdentifier); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// Image取得.
func getImageHandler(c *gin.Context) {
	imageIdentifier := c.Query("id")
	if imageIdentifier == "" {
		c.String(http.StatusBadRequest, "Image ID is required")
		return
	}

	imagePath := filepath.Join("./images", imageIdentifier+".png")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "Image not found")
		return
	}

	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read image: %s", err.Error()))
		return
	}

	c.Data(http.StatusOK, "image/png", imageData)
}

// アイテム抽選.
func drawItem(items []Item) *Item {
	totalWeight := 0.0
	for _, item := range items {
		totalWeight += item.Percentage
	}

	if totalWeight <= 0 {
		return nil
	}

	randomNumber := rng.Float64() * totalWeight
	currentWeight := 0.0
	for _, item := range items {
		currentWeight += item.Percentage
		if randomNumber < currentWeight {
			return &item
		}
	}
	return nil // Should not reach here if percentages are valid
}
