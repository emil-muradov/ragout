package aiassistant

import (
	"bytes"
	"context"
	"html/template"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func ParsePrompt(path string, data any) (string, error) {
	promptContent, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("temp").Parse(string(promptContent))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ConvertFloat64VectorToFloat32Vector(embedding []float64) []float32 {
	float32Vector := make([]float32, len(embedding))
	for i, val := range embedding {
		float32Vector[i] = float32(val)
	}
	return float32Vector
}

func TextToChunks(ctx context.Context, pathToTextFile string) ([]schema.Document, error) {
	f, err := os.Open(pathToTextFile)
	if err != nil {
		return nil, err
	}
	p := documentloaders.NewText(f)
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 200
	split.ChunkOverlap = 20
	split.Separators = []string{"\n\n"}
	docs, err := p.LoadAndSplit(ctx, split)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
