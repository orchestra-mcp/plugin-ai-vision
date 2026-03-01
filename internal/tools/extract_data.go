package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-vision/internal/vision"
	"google.golang.org/protobuf/types/known/structpb"
)

func ExtractDataSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path": map[string]any{
				"type":        "string",
				"description": "Path to the image file to extract data from",
			},
			"data_type": map[string]any{
				"type":        "string",
				"description": "Type of data to extract (e.g. table, form, chart). Optional.",
			},
		},
		"required": []any{"image_path"},
	})
	return s
}

func ExtractData() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "image_path"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		imagePath := helpers.GetString(req.Arguments, "image_path")
		dataType := helpers.GetString(req.Arguments, "data_type")

		var prompt string
		if dataType != "" {
			prompt = fmt.Sprintf("Extract %s data from this image as structured JSON.", dataType)
		} else {
			prompt = "Extract all data from this image as structured JSON."
		}

		client := vision.NewClient()
		result, err := client.Analyze(ctx, imagePath, prompt)
		if err != nil {
			return helpers.ErrorResult("analysis_error", fmt.Sprintf("failed to extract data: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
