package lib

import (
	"github.com/codegangsta/negroni"
	"log"
	"os"
)

func NewLogger() *negroni.Logger {
	return &negroni.Logger{
		log.New(os.Stdout, "[http] ", log.LstdFlags),
	}
}
