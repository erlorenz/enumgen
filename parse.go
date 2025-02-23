package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// MapParseResult holds both the parsed map and metadata about where it was found
type MapParseResult struct {
	Name   string            // Name of the variable containing the map
	Values map[string]string // The actual key-value pairs from the map
}

// ParseMapValues searches for and parses map declarations ending in "Values"
func ParseMapValues(node *ast.File) ([]MapParseResult, error) {
	// Store all discovered maps and their values
	var results []MapParseResult

	// Traverse the AST looking for variable declarations
	ast.Inspect(node, func(n ast.Node) bool {
		// First, check if we're looking at a variable declaration
		valueSpec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true // Continue searching if this isn't a value specification
		}

		// For each name in the value specification
		for i, name := range valueSpec.Names {
			// Check if the variable name ends with "Values"
			if !strings.HasSuffix(name.Name, "Values") {
				continue
			}

			// Verify we have a corresponding value and it's a map
			if i >= len(valueSpec.Values) {
				continue
			}

			// Try to extract the map value
			if mapValues, err := extractMapValues(valueSpec.Values[i]); err == nil {
				results = append(results, MapParseResult{
					Name:   name.Name,
					Values: mapValues,
				})
			}
		}

		return true
	})

	// Return an error if we didn't find any valid maps
	if len(results) == 0 {
		return nil, fmt.Errorf("no valid map variables ending in 'Values' were found")
	}

	return results, nil
}

// extractMapValues handles the actual parsing of map contents
func extractMapValues(expr ast.Expr) (map[string]string, error) {
	// Convert the expression to a composite literal (map declaration)
	mapLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("expected map literal, got %T", expr)
	}

	// Check that this is actually a map type
	mapType, ok := mapLit.Type.(*ast.MapType)
	if !ok {
		return nil, fmt.Errorf("expected map type, got %T", mapLit.Type)
	}

	// Verify that both key and value types are strings
	if !isStringType(mapType.Key) || !isStringType(mapType.Value) {
		return nil, fmt.Errorf("map must have string keys and values")
	}

	// Initialize our result map
	result := make(map[string]string)

	// Process each key-value pair in the map
	for _, elt := range mapLit.Elts {
		kvExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue // Skip non key-value pairs
		}

		// Extract the key string
		key, ok := kvExpr.Key.(*ast.BasicLit)
		if !ok || key.Kind != token.STRING {
			continue // Skip non-string keys
		}

		// Extract the value string
		value, ok := kvExpr.Value.(*ast.BasicLit)
		if !ok || value.Kind != token.STRING {
			continue // Skip non-string values
		}

		// Remove the surrounding quotes and add to our result map
		keyStr := strings.Trim(key.Value, `"`)
		valueStr := strings.Trim(value.Value, `"`)
		result[keyStr] = valueStr
	}

	return result, nil
}

// isStringType checks if an AST expression represents a string type
func isStringType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "string"
}
