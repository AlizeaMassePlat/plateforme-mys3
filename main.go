package main

import (
	"net/http"
	"myS3/storage"

	"github.com/gin-gonic/gin"
)

type CreateBucketRequest struct {
	BucketName string `json:"bucket_name" binding:"required"`
}

func main() {
	// Initialiser le routeur Gin
	r := gin.Default()

	// Route pour créer un bucket
	r.POST("/create-bucket", func(c *gin.Context) {
		var req CreateBucketRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bucketName := req.BucketName

		// Créer le bucket en envoyant une requête HTTP
		err := storage.AjoutBucket(bucketName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bucket created successfully", "bucket": bucketName})
	})

	// Démarrer le serveur
	r.Run(":8080")
}
