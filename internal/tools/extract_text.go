package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-vision/internal/vision"
	"google.golang.org/protobuf/types/known/structpb"
)

func ExtractTextSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path": map[string]any{
				"type":        "string",
				"description": "Path to the image file to extract text from",
			},
		},
		"required": []any{"image_path"},
	})
	return s
}

func ExtractText() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "image_path"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		imagePath := helpers.GetString(req.Arguments, "image_path")
		prompt := "Extract all text visible in this image. Return only the text content."

		client := vision.NewClient()
		result, err := client.Analyze(ctx, imagePath, prompt)
		if err != nil {
			return helpers.ErrorResult("analysis_error", fmt.Sprintf("failed to extract text: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
