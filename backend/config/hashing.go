package config

import (
	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("hashing", map[string]any{
		// Supported Drivers: "bcrypt", "argon2id"
		"driver": "bcrypt",
		"bcrypt": map[string]any{
			// Cost >= 12 per security baseline (plan §10.3).
			"rounds": 12,
		},
		"argon2id": map[string]any{
			"memory":  65536,
			"time":    4,
			"threads": 1,
		},
	})
}
