package tools

// Tests for the ai.vision plugin tool handlers.
//
// All tools call vision.NewClient().Analyze which requires ANTHROPIC_API_KEY.
// In CI this will return analysis_error ("ANTHROPIC_API_KEY not set").
// Tests focus on:
//   1. Validation errors (missing required args) — no API key needed.
//   2. API-unavailable path — tools with valid args return analysis_error.

import (
	"context"
	"os"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- helpers ----------

func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var s *structpb.Struct
	if args != nil {
		var err error
		s, err = structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("NewStruct: %v", err)
		}
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

func errorCode(resp *pluginv1.ToolResponse) string {
	if resp == nil {
		return ""
	}
	return resp.GetErrorCode()
}

// apiKeySet returns true when ANTHROPIC_API_KEY is configured.
func apiKeySet() bool {
	return os.Getenv("ANTHROPIC_API_KEY") != ""
}

// makeTempImage creates a minimal PNG-like temp file.
func makeTempImage(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "image-*.png")
	if err != nil {
		t.Fatal(err)
	}
	// Minimal PNG header bytes.
	_, _ = f.Write([]byte("\x89PNG\r\n\x1a\n"))
	_ = f.Close()
	return f.Name()
}

// ---------- analyze_image ----------

func TestAnalyzeImage_MissingImagePath(t *testing.T) {
	resp := callTool(t, AnalyzeImage(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing image_path")
	}
	if errorCode(resp) != "validation_error" {
		t.Errorf("expected validation_error, got %q", errorCode(resp))
	}
}

func TestAnalyzeImage_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, AnalyzeImage(), map[string]any{"image_path": img})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
	if errorCode(resp) != "analysis_error" {
		t.Errorf("expected analysis_error, got %q", errorCode(resp))
	}
}

func TestAnalyzeImage_WithCustomPrompt_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, AnalyzeImage(), map[string]any{
		"image_path": img,
		"prompt":     "What color is the background?",
	})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}

// ---------- extract_text ----------

func TestExtractText_MissingImagePath(t *testing.T) {
	resp := callTool(t, ExtractText(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing image_path")
	}
}

func TestExtractText_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, ExtractText(), map[string]any{"image_path": img})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}

// ---------- find_elements ----------

func TestFindElements_MissingImagePath(t *testing.T) {
	resp := callTool(t, FindElements(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing image_path")
	}
}

func TestFindElements_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, FindElements(), map[string]any{"image_path": img})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}

// ---------- describe_screen ----------

func TestDescribeScreen_MissingImagePath(t *testing.T) {
	resp := callTool(t, DescribeScreen(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing image_path")
	}
}

func TestDescribeScreen_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, DescribeScreen(), map[string]any{"image_path": img})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}

// ---------- compare_images ----------

func TestCompareImages_MissingBothPaths(t *testing.T) {
	resp := callTool(t, CompareImages(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing image paths")
	}
}

func TestCompareImages_MissingSecondPath(t *testing.T) {
	resp := callTool(t, CompareImages(), map[string]any{"image_path_1": "/tmp/a.png"})
	if !isError(resp) {
		t.Error("expected validation_error for missing image_path_2")
	}
}

func TestCompareImages_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img1 := makeTempImage(t)
	img2 := makeTempImage(t)
	resp := callTool(t, CompareImages(), map[string]any{
		"image_path_1": img1,
		"image_path_2": img2,
	})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}

// ---------- extract_data ----------

func TestExtractData_MissingImagePath(t *testing.T) {
	resp := callTool(t, ExtractData(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing image_path")
	}
}

func TestExtractData_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, ExtractData(), map[string]any{"image_path": img})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}

func TestExtractData_WithDataType_NoAPIKey(t *testing.T) {
	if apiKeySet() {
		t.Skip("ANTHROPIC_API_KEY is set; skipping no-key path")
	}
	img := makeTempImage(t)
	resp := callTool(t, ExtractData(), map[string]any{
		"image_path": img,
		"data_type":  "table",
	})
	if !isError(resp) {
		t.Error("expected analysis_error without API key")
	}
}
