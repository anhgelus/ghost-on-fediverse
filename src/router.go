package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Body struct {
	Post struct {
		Current CurrentPost `json:"current"`
	} `json:"post"`
}

type CurrentPost struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	FeatureImage        string    `json:"feature_image"`
	FeatureImageCaption string    `json:"feature_image_caption"`
	Status              string    `json:"status"`
	Visibility          string    `json:"visibility"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	PublishedAt         time.Time `json:"published_at"`
	CustomExcerpt       string    `json:"customExcerpt"`
	Authors             []Author  `json:"authors"`
	Tags                []Tag     `json:"tags"`
	PrimaryTag          Tag       `json:"primary_tag"`
	EmailSegment        string    `json:"email_segment"`
	Url                 string    `json:"url"`
	Excerpt             string    `json:"excerpt"`
}

type Author struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Email string `json:"email"`
}

type Tag struct {
	ID           string `json:"id"`
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
		LogWarn(fmt.Sprintf("invalid User-Agent %s", r.Header.Get("User-Agent")))
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
	post := body.Post.Current
	var id int64
	if post.FeatureImage != "" {
		id, err = uploadMedia(post)
		if err != nil {
			LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	_, err = MastodonClient.PostStatus(genMessage(post), 0, []int64{id}, false, "", "")
	if err != nil {
		LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func uploadMedia(post CurrentPost) (int64, error) {
	imageUrl := post.FeatureImage
	alt := cleanText(post.FeatureImageCaption)

	s := strings.Split(imageUrl, ".")
	ext := s[len(s)-1]

	res, err := http.Get(imageUrl)
	if err != nil {
		return 0, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	file := fmt.Sprintf("tmp.%s", ext)
	err = os.WriteFile(file, b, 0666)
	if err != nil {
		return 0, err
	}

	att, err := MastodonClient.UploadMedia(file, alt, "")
	if err != nil {
		return 0, err
	}

	return att.ID, os.Remove(file)
}

func cleanText(t string) string {
	reg1 := regexp.MustCompile("<[^/][^<]+>")
	reg2 := regexp.MustCompile("</[^<]+>")
	b := []byte(t)
	v := []byte("")
	b = reg1.ReplaceAll(b, v)
	b = reg2.ReplaceAll(b, v)
	return string(b)
}

func genMessage(post CurrentPost) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", post.Title, post.CustomExcerpt, post.Url, genTags(post))
}

func genTags(post CurrentPost) string {
	tags := ""
	for _, t := range post.Tags {
		tags += "#" + t.Slug
	}
	return tags
}
