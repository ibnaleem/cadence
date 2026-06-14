package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
)

const SimilarityThresholdWarn  = float32(0.85)
const SimilarityThresholdMatch = float32(0.65)

type voyageRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type voyageResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func Embed(text string) ([]float32, error) {
	key := os.Getenv("VOYAGE_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("VOYAGE_API_KEY not set")
	} // if

	body, _ := json.Marshal(voyageRequest{Input: []string{text}, Model: "voyage-3-lite"})
	req, _ := http.NewRequest("POST", "https://api.voyageai.com/v1/embeddings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} // if
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("voyage API: %s", resp.Status)
	} // if

	var result voyageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	} // if
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("voyage API: no embedding returned")
	} // if

	return result.Data[0].Embedding, nil
} // Embed

func CosineSimilarity(a, b []float32) float32 {
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
