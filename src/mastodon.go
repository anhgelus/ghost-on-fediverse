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
