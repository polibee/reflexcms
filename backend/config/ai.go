package config

import (
	"github.com/goravel/framework/contracts/ai"
	openaifacades "github.com/goravel/openai/facades"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("ai", map[string]any{
		// Default AI Provider. Left empty until the AI module (M2+) is enabled;
		// see docs/工程化开发计划.md §O8 — the pipeline is event-driven and
		// stays inert without a provider configured.
		"default": config.Env("AI_PROVIDER"),
		"providers": map[string]any{
			"openai": map[string]any{
				"key": config.Env("OPENAI_API_KEY", ""),
				"models": map[string]any{
					"text":          map[string]any{"default": ""},
					"audio":         map[string]any{"default": ""},
					"transcription": map[string]any{"default": ""},
					"image":         map[string]any{"default": ""},
				},
				"failover": map[string][]string{},
				"url":      config.Env("OPENAI_BASE_URL", ""),
				"via": func() (ai.Provider, error) {
					return openaifacades.OpenAI("openai")
				},
			},
		},
	})
}
