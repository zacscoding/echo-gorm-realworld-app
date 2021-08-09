package e2e

const (
	UserJsonSchema = `
	{
	  "$schema": "http://json-schema.org/draft-04/schema#",
	  "type": "object",
	  "properties": {
		"user": {
		  "type": "object",
		  "properties": {
			"email": {
			  "type": "string"
			},
			"token": {
			  "type": "string"
			},
			"username": {
			  "type": "string"
			},
			"bio": {
			  "type": "string"
			},
			"image": {
			  "type": "string"
			}
		  },
		  "required": [
			"email",
			"token",
			"username",
			"bio",
			"image"
		  ]
		}
	  },
	  "required": [
		"user"
	  ]
	}`

	ProfileJsonSchema = `
	{
	  "$schema": "http://json-schema.org/draft-04/schema#",
	  "type": "object",
	  "properties": {
		"profile": {
		  "type": "object",
		  "properties": {
			"username": {
			  "type": "string"
			},
			"bio": {
			  "type": "string"
			},
			"image": {
			  "type": "string"
			},
			"following": {
			  "type": "boolean"
			}
		  },
		  "required": [
			"username",
			"bio",
			"image",
			"following"
		  ]
		}
	  },
	  "required": [
		"profile"
	  ]
	}`
)
