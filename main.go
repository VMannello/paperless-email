package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vmannello/paperless-email/pmail"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: pmail <path/to/config.yaml>\n")
		os.Exit(2)
	}

	cfg, err := pmail.LoadConfig(os.Args[1])
	if err != nil {
		log.Fatal("[ERROR]  could not load configuration:", err)
	}

	documentPath := os.Getenv("DOCUMENT_SOURCE_PATH")
	tags := strings.Split(os.Getenv("DOCUMENT_TAGS"), ",")

	pmail.SendEmail(cfg, tags, documentPath)
}
