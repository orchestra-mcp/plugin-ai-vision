package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-vision/internal/vision"
	"google.golang.org/protobuf/types/known/structpb"
)

func FindElementsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path": map[string]any{
				"type":        "string",
				"description": "Path to the image file to search for elements",
			},
			"element_type": map[string]any{
				"type":        "string",
				"description": "Type of element to find (e.g. button, input, link). Optional.",
			},
		},
		"required": []any{"image_path"},
	})
	return s
}

func FindElements() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "image_path"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		imagePath := helpers.GetString(req.Arguments, "image_path")
		elementType := helpers.GetString(req.Arguments, "element_type")

		var prompt string
		if elementType != "" {
			prompt = fmt.Sprintf("Find and list all %s elements in this image.", elementType)
		} else {
			prompt = "List all UI elements in this image."
		}

		client := vision.NewClient()
		result, err := client.Analyze(ctx, imagePath, prompt)
		if err != nil {
			return helpers.ErrorResult("analysis_error", fmt.Sprintf("failed to find elements: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
