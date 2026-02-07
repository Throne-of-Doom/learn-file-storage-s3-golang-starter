package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(name string, mediaType string) string {
	ext := getExtensionFromMediaType(mediaType)
	filename := name + ext
	return "/assets/" + filename
}

func getExtensionFromMediaType(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}

	ext := parts[1]
	if ext == "jpeg" {
		return ".jpg"
	}
	return "." + ext
}

// Map URL path -> disk path
func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	// assetPath is "/assets/<file>", we want "./assets/<file>"
	// so strip the leading "/assets" and prepend "./assets"
	const prefix = "/assets"
	if strings.HasPrefix(assetPath, prefix) {
		return cfg.assetsRoot + assetPath[len(prefix):]
	}
	return cfg.assetsRoot + assetPath
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return "http://localhost:" + cfg.port + assetPath
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var output ffprobeOutput
	if err := json.Unmarshal(out.Bytes(), &output); err != nil {
		return "", err
	}
	if len(output.Streams) == 0 {
		return "", fmt.Errorf("no video streams found")
	}
	width := output.Streams[0].Width
	height := output.Streams[0].Height

	if width > height {
		return "16:9", nil
	}
	if height > width {
		return "9:16", nil
	}
	return "other", nil
}
