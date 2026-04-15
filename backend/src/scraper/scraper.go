package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ScrapeRes struct{
	HTML string
	URL string
	Elapsed time.Duration
}

func FetchHTML( url string) (*ScrapeRes, error){
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://"){
		url = "https://" + url
	}

	client := &http.Client{
		Timeout: 15*time.Second,
	}

	start := time.Now()

	resp, err := client.Get(url)

	if err != nil{
		return nil, fmt.Errorf("Gagal fetch URL: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("status tidak ok: %w", resp.StatusCode)
	
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		return nil, fmt.Errorf("gagal mengambil body dari response: %w", err )
	}

	return &ScrapeRes{
		HTML: string(body),
		URL: url,
		Elapsed: time.Since(start),
	}, nil
}