package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/smtp"
	"os"
	"strings"
)

//go:embed facts.json
var factsJSON []byte

func main() {
	loadDotEnv()

	gmail := strings.TrimSpace(mustEnv("GMAIL_ADDRESS"))
	appPW := strings.TrimSpace(strings.ReplaceAll(mustEnv("GMAIL_APP_PASSWORD"), " ", ""))
	appPW = strings.TrimSpace(strings.ReplaceAll(appPW, "\n", ""))
	rawRcpts := mustEnv("RECIPIENTS")

	var facts []string
	if err := json.Unmarshal(factsJSON, &facts); err != nil {
		log.Fatalf("facts.json: %v", err)
	}
	if len(facts) == 0 {
		log.Fatal("facts.json is empty")
	}

	fact := facts[rand.IntN(len(facts))]

	var to []string
	for _, p := range strings.Split(rawRcpts, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, g, ok := strings.Cut(p, ":")
		if !ok {
			log.Fatalf("bad recipient %q (want NUMBER:GATEWAY, e.g. 5551234567:vtext.com)", p)
		}
		addr := strings.TrimSpace(n) + "@" + strings.TrimSpace(g)
		to = append(to, addr)
	}
	if len(to) == 0 {
		log.Fatal("RECIPIENTS is empty")
	}

	headers := strings.Join([]string{
		"From: " + gmail,
		"To: " + strings.Join(to, ", "),
		"Subject: ",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		fact,
	}, "\r\n")

	auth := smtp.PlainAuth("", gmail, appPW, "smtp.gmail.com")
	if err := smtp.SendMail("smtp.gmail.com:587", auth, gmail, to, []byte(headers)); err != nil {
		log.Fatalf("send: %v", err)
	}
	preview := fact
	if len(preview) > 60 {
		preview = preview[:60] + "…"
	}
	fmt.Printf("Sent to %d recipient(s): %s\n", len(to), preview)
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

// loadDotEnv reads a local .env for personal testing only.
// Existing environment variables win (GitHub Actions secrets unchanged).
func loadDotEnv() {
	data, err := os.ReadFile(".env")
	if err != nil || len(data) == 0 {
		return
	}
	s := strings.ReplaceAll(strings.TrimPrefix(string(data), "\ufeff"), "\r\n", "\n")
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
}
