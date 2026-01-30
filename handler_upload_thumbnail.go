package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

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

	mediaType := header.Header.Get("Content-Type")

	metaData, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Video not found", err)
		return
	}
	if userID != metaData.UserID {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this video", nil)
		return
	}

	assetPath := getAssetPath(videoID, mediaType)

	filename := path.Base(assetPath)

	diskPath := filepath.Join(cfg.assetsRoot, filename)

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

	url := "http://localhost:" + cfg.port + assetPath
	metaData.ThumbnailURL = &url

	err = cfg.db.UpdateVideo(metaData)
	if err != nil {
		respondWithError(w, 500, "error updating video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, metaData)
}
