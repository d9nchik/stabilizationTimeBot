package puller

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"stabilizationTimeBot/pkg/core"
)

type Puller struct {
	sender           core.Sender
	previousFileHash string
}

func NewPuller(sender core.Sender) *Puller {
	return &Puller{sender: sender}
}

func (p *Puller) Run(ctx context.Context) {
	ticker := time.NewTicker(1)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Reset(time.Minute * 30)
			p.pull()

		case <-ctx.Done():
			return
		}
	}
}

func (p *Puller) pull() {
	file, err := getFile()
	if err != nil {
		log.Printf("Problem with retrieving file: %v", err.Error())
		return
	}
	defer os.Remove(file)

	hash, err := getHashOfFile(file)
	if err != nil {
		log.Printf("Problem with retrieving file hash: %v", err.Error())
		return
	}

	if p.previousFileHash == hash {
		return
	}

	if ok := p.sender.SendFile(file); ok {
		p.previousFileHash = hash
	}
}

func getFile() (string, error) {
	resp, err := http.Get("https://energy.volyn.ua/spozhyvacham/perervy-u-elektropostachanni/hrafik-vidkliuchen/!files/hsv-flutskm.pdf")
	if err != nil {
		return "", err
	}

	file, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func getHashOfFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
