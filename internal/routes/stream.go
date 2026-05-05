package routes

import (
	"net/http"
	"quiz/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func handleAudioRequest(c *gin.Context) {
	fileID := strings.TrimSpace(c.Param("id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	url := utils.BuildDriveAudioURL(fileID)

	client := &http.Client{
		Timeout: 2 * time.Minute,
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear la petición"})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "no se pudo obtener el audio remoto"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":  "google drive devolvió error",
			"status": resp.StatusCode,
		})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "audio/mpeg"
	}

	if strings.Contains(strings.ToLower(contentType), "text/html") {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "google drive no devolvió un archivo de audio válido",
		})
		return
	}

	extraHeaders := map[string]string{
		"Cache-Control":               "public, max-age=3600",
		"Access-Control-Allow-Origin": "*",
	}

	if resp.Header.Get("Accept-Ranges") != "" {
		extraHeaders["Accept-Ranges"] = resp.Header.Get("Accept-Ranges")
	}

	if resp.Header.Get("Content-Disposition") != "" {
		extraHeaders["Content-Disposition"] = resp.Header.Get("Content-Disposition")
	}

	c.DataFromReader(
		http.StatusOK,
		resp.ContentLength,
		contentType,
		resp.Body,
		extraHeaders,
	)
}

func StreamRoutes(rg *gin.RouterGroup) {
	streamGroup := rg.Group("/stream")
	{
		streamGroup.GET("/audio/:id", handleAudioRequest)
	}
}
