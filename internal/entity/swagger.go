package entity

// SwaggerInfo represents the top-level structure of the Swagger file
type SwaggerInfo struct {
	Swagger             string                    `json:"swagger"`
	Info                Info                      `json:"info"`
	Host                string                    `json:"host"`
	BasePath            string                    `json:"basePath"`
	Paths               map[string]PathItem       `json:"paths"`
	SecurityDefinitions map[string]SecurityScheme `json:"securityDefinitions"`
}

// Info represents the "info" section of the Swagger file
type Info struct {
	Description string  `json:"description"`
	Title       string  `json:"title"`
	Contact     Contact `json:"contact"`
	License     License `json:"license"`
	Version     string  `json:"version"`
}

// Contact represents the "contact" section of the Swagger file
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// License represents the "license" section of the Swagger file
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PathItem represents a path in the "paths" section of the Swagger file
type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

// Operation represents an operation (e.g., GET, POST) for a path
type Operation struct {
	Security []interface{} `json:"security"`
}

type SecurityScheme struct {
	In   string `json:"in"`
	Name string `json:"name"`
	Type string `json:"type"`
}
