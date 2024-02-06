package src

import (
	"errors"
	"fmt"
	"github.com/McKael/madon"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	EnvInstance      = "FED_INSTANCE"
	EnvLoginMail     = "FED_LOGIN_EMAIL"
	EnvLoginPassword = "FED_LOGIN_PASSWORD"
	EnvLoginToken    = "FED_LOGIN_TOKEN"
)

var (
	MastodonClient *madon.Client
	Scopes         = []string{"read", "write", "follow"}

	ErrCannotConnect = errors.New("cannot connect to Mastodon")
	ErrEnvNotSet     = errors.New("env var not set")
)

func ConnectMastodon() error {
	instance, v := os.LookupEnv(EnvInstance)
	if !v {
		return errors.Join(ErrEnvNotSet, errors.New(EnvInstance+" is not set"))
	}
	mail, v := os.LookupEnv(EnvLoginMail)
	if !v {
		return errors.Join(ErrEnvNotSet, errors.New(EnvLoginMail+" is not set"))
	}
	password, v := os.LookupEnv(EnvLoginPassword)
	if !v {
		return errors.Join(ErrEnvNotSet, errors.New(EnvLoginPassword+" is not set"))
	}
	token, v := os.LookupEnv(EnvLoginToken)
	if !v {
		return errors.Join(ErrEnvNotSet, errors.New(EnvLoginToken+" is not set"))
	}
	var err error
	MastodonClient, err = madon.NewApp(
		"Ghost Integration",
		"https://github.com/anhgelus/ghost-on-fediverse",
		Scopes,
		madon.NoRedirect,
		instance,
	)
	if err != nil {
		return errors.Join(err, ErrCannotConnect)
	}
	err = MastodonClient.SetUserToken(token, mail, password, Scopes)
	if err != nil {
		return errors.Join(err, ErrCannotConnect)
	}
	return nil
}

func TootPost(post CurrentPost) error {
	var id int64
	var err error
	if post.FeatureImage != "" {
		id, err = uploadMedia(post)
		if err != nil {
			return err
		}
	}
	_, err = MastodonClient.PostStatus(genMessage(post), 0, []int64{id}, false, "", "")
	return err
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
