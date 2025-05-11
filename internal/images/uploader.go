// internal/images/uploader.go

package images

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var imgRe = regexp.MustCompile(`!\[[^\]]*]\(([^)]+)\)`)

type Service struct {
	BaseURL  string
	Client   *http.Client
	AdminJWT string
	cache    map[string]string // sha1 → remoteURL
}

func New(base string, jwt string, c *http.Client) *Service {
	return &Service{
		BaseURL:  base,
		Client:   c,
		AdminJWT: jwt,
		cache:    make(map[string]string),
	}
}

func (s *Service) Rewrite(md []byte, root string) ([]byte, error) {
	return imgRe.ReplaceAllFunc(md, func(m []byte) []byte {
		match := imgRe.FindSubmatch(m)
		locPath := string(match[1])
		full := filepath.Join(root, locPath)

		remote, err := s.upload(full)
		if err != nil {
			return m // keep local ref if upload fails
		}
		return bytes.Replace(m, match[1], []byte(remote), 1)
	}), nil
}

func (s *Service) upload(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", sha1.Sum(raw))
	if url, ok := s.cache[sum]; ok {
		return url, nil
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("file", filepath.Base(path))
	io.Copy(part, bytes.NewReader(raw))
	w.Close()

	req, _ := http.NewRequest("POST", s.BaseURL+"images/upload/", body)
	req.Header.Set("Authorization", "Ghost "+s.AdminJWT)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := s.Client.Do(req)
	if err != nil || resp.StatusCode >= 300 {
		return "", fmt.Errorf("upload failed %v", err)
	}
	// parse {"images":[{"url":"https://…"}]}
	var r struct {
		Images []struct {
			URL string `json:"url"`
		}
	}
	json.NewDecoder(resp.Body).Decode(&r)
	remote := r.Images[0].URL
	s.cache[sum] = remote
	return remote, nil
}
