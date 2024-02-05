package src

import (
	"log"
)

const (
	AnsiReset       = "\033[0m"
	AnsiRed         = "\033[31m"
	AnsiGreen       = "\033[32m"
	AnsiYellow      = "\033[33m"
	AnsiBlue        = "\033[34m"
	AnsiMagenta     = "\033[35m"
	AnsiCyan        = "\033[36m"
	AnsiWhite       = "\033[37m"
	AnsiBlueBold    = "\033[34;1m"
	AnsiMagentaBold = "\033[35;1m"
	AnsiRedBold     = "\033[31;1m"
	AnsiYellowBold  = "\033[33;1m"
)

func LogInfo(v string) {
	log.Default().Println(AnsiGreen, v, AnsiReset)
}

func LogWarn(v string) {
	log.Default().Println(AnsiYellow, v, AnsiReset)
}

func LogError(v error) {
	log.Default().Println(AnsiYellowBold, v, AnsiReset)
}

func Crash(err error) {
	log.Default().Panic(err)
}
