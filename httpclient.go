package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

const userAgent = "govmcp/0.1.0 (MCP Server)"

func doGet(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 500)]))
	}
	return body, nil
}

func doPost(ctx context.Context, rawURL string, contentType string, bodyData []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(string(bodyData)))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 500)]))
	}
	return body, nil
}

// buildOpenFDAURL constructs an openFDA API URL.
// searchParts are joined with "+AND+". The openFDA API uses these as literal strings in the query.
func buildOpenFDAURL(endpoint string, searchParts []string, limit, skip int, countField string) string {
	base := "https://api.fda.gov/" + endpoint + ".json"

	// Build query string manually since openFDA uses literal + in search syntax
	var params []string

	if len(searchParts) > 0 {
		search := strings.Join(searchParts, "+AND+")
		params = append(params, "search="+search)
	}
	if countField != "" {
		params = append(params, "count="+url.QueryEscape(countField))
	} else {
		if limit > 0 {
			params = append(params, "limit="+strconv.Itoa(limit))
		}
		if skip > 0 {
			params = append(params, "skip="+strconv.Itoa(skip))
		}
	}

	if key := os.Getenv("OPENFDA_API_KEY"); key != "" {
		params = append(params, "api_key="+url.QueryEscape(key))
	}

	if len(params) > 0 {
		return base + "?" + strings.Join(params, "&")
	}
	return base
}

func clampLimit(limit, defaultVal, maxVal int) int {
	if limit <= 0 {
		return defaultVal
	}
	if limit > maxVal {
		return maxVal
	}
	return limit
}

// writeField writes a formatted field to the string builder if the value exists in the map.
func writeField(sb *strings.Builder, label string, m map[string]interface{}, key string) {
	v, ok := m[key]
	if !ok || v == nil {
		return
	}
	switch val := v.(type) {
	case string:
		if val != "" {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", label, val))
		}
	case float64:
		sb.WriteString(fmt.Sprintf("  %s: %g\n", label, val))
	case bool:
		sb.WriteString(fmt.Sprintf("  %s: %t\n", label, val))
	default:
		sb.WriteString(fmt.Sprintf("  %s: %v\n", label, val))
	}
}

// writeNestedField writes a field from a nested map path like "openfda.brand_name".
func writeNestedField(sb *strings.Builder, label string, m map[string]interface{}, path string) {
	parts := strings.Split(path, ".")
	current := m
	for i, part := range parts {
		v, ok := current[part]
		if !ok || v == nil {
			return
		}
		if i == len(parts)-1 {
			// Last part - write the value
			switch val := v.(type) {
			case string:
				if val != "" {
					sb.WriteString(fmt.Sprintf("  %s: %s\n", label, val))
				}
			case []interface{}:
				if len(val) > 0 {
					strs := make([]string, 0, len(val))
					for _, item := range val {
						if s, ok := item.(string); ok {
							strs = append(strs, s)
						}
					}
					if len(strs) > 0 {
						sb.WriteString(fmt.Sprintf("  %s: %s\n", label, strings.Join(strs, ", ")))
					}
				}
			default:
				sb.WriteString(fmt.Sprintf("  %s: %v\n", label, val))
			}
			return
		}
		// Not the last part - descend
		if nested, ok := v.(map[string]interface{}); ok {
			current = nested
		} else {
			return
		}
	}
}

// truncate shortens a string to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// writeArrayField writes the first element of a string array field.
func writeArrayField(sb *strings.Builder, label string, m map[string]interface{}, key string, maxLen int) {
	v, ok := m[key]
	if !ok || v == nil {
		return
	}
	switch val := v.(type) {
	case []interface{}:
		if len(val) > 0 {
			if s, ok := val[0].(string); ok && s != "" {
				sb.WriteString(fmt.Sprintf("  %s: %s\n", label, truncate(s, maxLen)))
			}
		}
	case string:
		if val != "" {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", label, truncate(val, maxLen)))
		}
	}
}
