package oas

import (
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

// RegisterFromSpec reads an OpenAPI spec file and registers routes on the provided Fiber app.
// handlers should map OpenAPI operationId -> fiber.Handler.
func RegisterFromSpec(app *fiber.App, handlers map[string]fiber.Handler, specPath string) error {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return err
	}

	for rawPath, item := range doc.Paths {
		// convert OpenAPI path params {id} to Fiber :id
		path := convertPath(rawPath)

		if item.Get != nil {
			op := item.Get
			registerHandler(app, "GET", path, op.OperationID, handlers)
		}
		if item.Post != nil {
			op := item.Post
			registerHandler(app, "POST", path, op.OperationID, handlers)
		}
		if item.Put != nil {
			op := item.Put
			registerHandler(app, "PUT", path, op.OperationID, handlers)
		}
		if item.Delete != nil {
			op := item.Delete
			registerHandler(app, "DELETE", path, op.OperationID, handlers)
		}
		if item.Patch != nil {
			op := item.Patch
			registerHandler(app, "PATCH", path, op.OperationID, handlers)
		}
	}

	return nil
}

func registerHandler(app *fiber.App, method, path, operationID string, handlers map[string]fiber.Handler) {
	if operationID == "" {
		log.Printf("oas: skipping %s %s; no operationId", method, path)
		return
	}
	h, ok := handlers[operationID]
	if !ok {
		log.Printf("oas: no handler registered for operationId %s (%s %s)", operationID, method, path)
		return
	}

	switch strings.ToUpper(method) {
	case "GET":
		app.Get(path, h)
	case "POST":
		app.Post(path, h)
	case "PUT":
		app.Put(path, h)
	case "DELETE":
		app.Delete(path, h)
	case "PATCH":
		app.Patch(path, h)
	default:
		log.Printf("oas: unsupported method %s for %s", method, path)
	}
}

func convertPath(p string) string {
	// replace {param} with :param for Fiber
	p = strings.ReplaceAll(p, "{", ":")
	p = strings.ReplaceAll(p, "}", "")
	return p
}
