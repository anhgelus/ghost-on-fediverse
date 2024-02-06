package src

import (
	"errors"
	"github.com/McKael/madon"
	"os"
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
		return ErrEnvNotSet
	}
	mail, v := os.LookupEnv(EnvLoginMail)
	if !v {
		return ErrEnvNotSet
	}
	password, v := os.LookupEnv(EnvLoginPassword)
	if !v {
		return ErrEnvNotSet
	}
	token, v := os.LookupEnv(EnvLoginToken)
	if !v {
		return ErrEnvNotSet
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
