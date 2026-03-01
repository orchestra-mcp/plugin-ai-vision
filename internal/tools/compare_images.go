package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-vision/internal/vision"
	"google.golang.org/protobuf/types/known/structpb"
)

func CompareImagesSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path_1": map[string]any{
				"type":        "string",
				"description": "Path to the first image file",
			},
			"image_path_2": map[string]any{
				"type":        "string",
				"description": "Path to the second image file",
			},
		},
		"required": []any{"image_path_1", "image_path_2"},
	})
	return s
}

func CompareImages() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "image_path_1", "image_path_2"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		imagePath1 := helpers.GetString(req.Arguments, "image_path_1")
		imagePath2 := helpers.GetString(req.Arguments, "image_path_2")

		client := vision.NewClient()

		desc1, err := client.Analyze(ctx, imagePath1, "Describe this image in detail, noting all key visual elements, layout, content, and any text present.")
		if err != nil {
			return helpers.ErrorResult("analysis_error", fmt.Sprintf("failed to analyze first image: %v", err)), nil
		}

		desc2, err := client.Analyze(ctx, imagePath2, "Describe this image in detail, noting all key visual elements, layout, content, and any text present.")
		if err != nil {
			return helpers.ErrorResult("analysis_error", fmt.Sprintf("failed to analyze second image: %v", err)), nil
		}

		comparisonPrompt := fmt.Sprintf(
			"Compare these two image descriptions and identify similarities and differences:\n\nImage 1:\n%s\n\nImage 2:\n%s\n\nProvide a detailed comparison covering: visual differences, layout changes, content differences, and any notable changes.",
			desc1, desc2,
		)

		comparison, err := client.Analyze(ctx, imagePath1, comparisonPrompt)
		if err != nil {
			return helpers.ErrorResult("analysis_error", fmt.Sprintf("failed to compare images: %v", err)), nil
		}

		output := fmt.Sprintf("## Image Comparison\n\n### Image 1\n%s\n\n### Image 2\n%s\n\n### Comparison\n%s", desc1, desc2, comparison)
		return helpers.TextResult(output), nil
	}
}
