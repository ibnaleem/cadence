package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const SimilarityThresholdWarn  = float32(0.85)
const SimilarityThresholdMatch = float32(0.65)

type ollamaEmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type ollamaEmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

var ollamaClient = &http.Client{Timeout: 10 * time.Second}

func Embed(text string) ([]float32, error) {
	body, _ := json.Marshal(ollamaEmbedRequest{Model: "embeddinggemma", Input: []string{text}})
	resp, err := ollamaClient.Post("http://localhost:11434/api/embed", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama not reachable: %w", err)
	} // if
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama: %s", resp.Status)
	} // if

	var result ollamaEmbedResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10<<20)).Decode(&result); err != nil {
		return nil, err
	} // if
	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("ollama: no embedding returned")
	} // if

	return result.Embeddings[0], nil
} // Embed

func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	} // for
	if normA == 0 || normB == 0 {
		return 0
	} // if
	return dot / float32(math.Sqrt(float64(normA))*math.Sqrt(float64(normB)))
} // CosineSimilarity
