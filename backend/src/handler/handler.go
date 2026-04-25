package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/algorithm"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/parser"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/scraper"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/traverse", TraverseHandler)
	mux.HandleFunc("/api/lca", LCAHandler)
	mux.HandleFunc("/api/upload", UploadHTMLHandler)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func TraverseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method tidak didukung"))
		return
	}

	var req model.TraversalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, fmt.Errorf("request JSON tidak valid: %w", err))
		return
	}

	html, err := resolveHTML(req.InputMode, req.URL, req.HTML)
	if err != nil {
		writeError(w, err)
		return
	}

	root, err := parser.ParseHTML(html)
	if err != nil {
		writeError(w, fmt.Errorf("gagal parse HTML: %w", err))
		return
	}

	sel, err := selector.StringToSelector(req.Selector)
	if err != nil {
		writeError(w, fmt.Errorf("selector tidak valid: %w", err))
		return
	}

	steps, matchedIDs, stats, err := generateSteps(root, sel, req.Algorithm)
	if err != nil {
		writeError(w, err)
		return
	}

	elapsed, err := measureTraversalElapsed(root, sel, req.Algorithm, req.Parallel)
	if err != nil {
		writeError(w, err)
		return
	}
	stats.Elapsed = elapsed

	if !req.AllResult && req.Limit > 0 {
		matchedSeen := 0
		cutIndex := len(steps)

		for i, step := range steps {
			if step.IsMatch {
				matchedSeen++
				if matchedSeen >= req.Limit {
					cutIndex = i + 1
					break
				}
			}
		}

		steps = steps[:cutIndex]
		if len(matchedIDs) > req.Limit {
			matchedIDs = matchedIDs[:req.Limit]
		}
		if len(steps) > 0 {
			lastStep := steps[len(steps)-1]
			stats.Visited = lastStep.VisitedCount
			stats.Matched = lastStep.MatchedCount
		}
	}

	resp := model.TraverseResponse{
		Success:        true,
		Tree:           root,
		Steps:          steps,
		MatchedNodeIDs: matchedIDs,
		Stats:          stats,
		Algorithm:      req.Algorithm,
		Parallel:       req.Parallel,
		Selector:       req.Selector,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func LCAHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method tidak didukung"))
		return
	}

	var req model.LCARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, fmt.Errorf("request JSON tidak valid: %w", err))
		return
	}

	html, err := resolveHTML(req.InputMode, req.URL, req.HTML)
	if err != nil {
		writeError(w, err)
		return
	}

	root, err := parser.ParseHTML(html)
	if err != nil {
		writeError(w, fmt.Errorf("gagal parse HTML: %w", err))
		return
	}

	lca, lcaSteps := algorithm.FindLCAWithSteps(root, req.NodeA, req.NodeB)
	if lca == nil {
		writeError(w, fmt.Errorf("LCA tidak ditemukan untuk node %d dan %d", req.NodeA, req.NodeB))
		return
	}

	resp := model.LCAResponse{
		Success: true,
		NodeID:  lca.Id,
		Tag:     lca.TagName,
		Label:   fmt.Sprintf("%s#%d", lca.TagName, lca.Id),
		Steps:   lcaSteps,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func generateSteps(root *model.DOMNode, sel selector.Selector, algo string) ([]model.TraversalStep, []int, model.TraversalStats, error) {
	switch algo {
	case "bfs":
		steps, matchedIDs, stats := algorithm.BFSWithSteps(root, sel)
		return steps, matchedIDs, stats, nil

	case "dfs":
		steps, matchedIDs, stats := algorithm.DFSWithSteps(root, sel)
		return steps, matchedIDs, stats, nil

	default:
		return nil, nil, model.TraversalStats{}, fmt.Errorf("algorithm tidak valid: %q", algo)
	}
}

func UploadHTMLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method tidak didukung, gunakan POST"))
		return
	}

	const maxSize = 10 << 20

	if err := r.ParseMultipartForm(maxSize); err != nil {
		writeError(w, fmt.Errorf("gagal parse multipart form: %w", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, fmt.Errorf("field \"file\" tidak ditemukan: %w", err))
		return
	}
	defer file.Close()

	name := header.Filename
	if len(name) < 5 || name[len(name)-5:] != ".html" {
		writeError(w, fmt.Errorf("file harus berekstensi .html, diterima: %q", name))
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		writeError(w, fmt.Errorf("gagal membaca file: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"filename": header.Filename,
		"html":     string(content),
	})
}

func resolveHTML(inputMode, url, html string) (string, error) {
	switch inputMode {
	case "url":
		res, err := scraper.FetchHTML(url)
		if err != nil {
			return "", fmt.Errorf("gagal fetch URL: %w", err)
		}
		return res.HTML, nil

	case "html":
		if html == "" {
			return "", fmt.Errorf("field html kosong")
		}
		return html, nil

	default:
		return "", fmt.Errorf("inputMode tidak valid: %q (gunakan \"url\" atau \"html\")", inputMode)
	}
}

func measureTraversalElapsed(root *model.DOMNode, sel selector.Selector, algo string, parallel bool) (float64, error) {
	start := time.Now()

	switch algo {
	case "bfs":
		if parallel {
			algorithm.FastBFS(root, sel)
		} else {
			algorithm.BFS(root, sel)
		}

	case "dfs":
		if parallel {
			algorithm.FastDFS(root, sel)
		} else {
			algorithm.DFS(root, sel)
		}

	default:
		return 0, fmt.Errorf("algorithm tidak valid: %q", algo)
	}

	return float64(time.Since(start).Microseconds()) / 1000, nil
}
