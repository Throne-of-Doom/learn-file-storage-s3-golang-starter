package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	data, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid file upload", err)
		return
	}

	defer data.Close()

	contentType := header.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid content type", err)
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "invalid file type", nil)
		return
	}

	metaData, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Video not found", err)
		return
	}
	if userID != metaData.UserID {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this video", nil)
		return
	}

	pathR := make([]byte, 32)

	_, err = rand.Read(pathR)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating filename", err)
		return
	}

	videoPath := base64.RawURLEncoding.EncodeToString(pathR)
	assetPath := getAssetPath(videoPath, mediaType)
	diskPath := cfg.getAssetDiskPath(assetPath)

	dst, err := os.Create(diskPath)
	if err != nil {
		respondWithError(w, 500, "error creating path", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		respondWithError(w, 500, "error copying path", err)
		return
	}

	url := cfg.getAssetURL(assetPath)

	metaData.ThumbnailURL = &url

	err = cfg.db.UpdateVideo(metaData)
	if err != nil {
		respondWithError(w, 500, "error updating video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, metaData)
}
