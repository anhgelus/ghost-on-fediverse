package src

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Body struct {
	Post struct {
		Current struct {
			ID            int       `json:"id"`
			Title         string    `json:"title"`
			Slug          string    `json:"slug"`
			FeatureImage  string    `json:"feature_image"`
			Status        string    `json:"status"`
			Visibility    string    `json:"visibility"`
			CreatedAt     time.Time `json:"created_at"`
			UpdatedAt     time.Time `json:"updated_at"`
			PublishedAt   time.Time `json:"published_at"`
			CustomExcerpt string    `json:"customExcerpt"`
			Authors       []Author  `json:"authors"`
			Tags          []Tag     `json:"tags"`
			PrimaryTag    Tag       `json:"primary_tag"`
			EmailSegment  string    `json:"email_segment"`
			Url           string    `json:"url"`
			Excerpt       string    `json:"excerpt"`
		} `json:"current"`
	} `json:"post"`
}

type Author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Email string `json:"email"`
}

type Tag struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	FeatureImage string `json:"feature_image"`
	Visibility   string `json:"visibility"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		LogWarn("invalid Content-Type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !strings.Contains(r.Header.Get("User-Agent"), "Ghost/5") {
		LogWarn("invalid User-Agent")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		LogWarn("wrong method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		LogError(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	//post := body.Post.Current
}
