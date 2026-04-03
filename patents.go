package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const patentsViewBaseURL = "https://api.patentsview.org/patents/query"
const assigneesViewBaseURL = "https://api.patentsview.org/assignees/query"

func getPatentsViewAPIKey() string {
	return os.Getenv("PATENTSVIEW_API_KEY")
}

func handleSearchPatents(ctx context.Context, req *mcp.CallToolRequest, input SearchPatentsInput) (*mcp.CallToolResult, any, error) {
	apiKey := getPatentsViewAPIKey()
	if apiKey == "" {
		return errorResult("PATENTSVIEW_API_KEY environment variable is not set. Get a free key at https://patentsview.org/apis"), nil, nil
	}

	limit := clampLimit(input.Limit, 25, 100)

	// Build query
	var conditions []interface{}

	if input.Query != "" {
		conditions = append(conditions, map[string]interface{}{
			"_text_any": map[string]string{
				"patent_abstract": input.Query,
			},
		})
	}
	if input.Assignee != "" {
		conditions = append(conditions, map[string]interface{}{
			"_text_any": map[string]string{
				"assignees.assignee_organization": input.Assignee,
			},
		})
	}
	if input.Inventor != "" {
		conditions = append(conditions, map[string]interface{}{
			"_text_any": map[string]string{
				"inventors.inventor_last_name": input.Inventor,
			},
		})
	}
	if input.DateFrom != "" {
		conditions = append(conditions, map[string]interface{}{
			"_gte": map[string]string{
				"patent_date": input.DateFrom,
			},
		})
	}
	if input.DateTo != "" {
		conditions = append(conditions, map[string]interface{}{
			"_lte": map[string]string{
				"patent_date": input.DateTo,
			},
		})
	}

	var query interface{}
	if len(conditions) == 0 {
		return errorResult("At least one search parameter is required"), nil, nil
	} else if len(conditions) == 1 {
		query = conditions[0]
	} else {
		query = map[string]interface{}{"_and": conditions}
	}

	requestBody := map[string]interface{}{
		"q": query,
		"f": []string{
			"patent_id", "patent_title", "patent_date", "patent_abstract",
			"patent_type", "patent_num_claims",
			"assignees.assignee_organization",
			"inventors.inventor_first_name", "inventors.inventor_last_name",
		},
		"o": map[string]interface{}{
			"per_page": limit,
		},
		"s": []map[string]string{
			{"patent_date": "desc"},
		},
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to build request: %v", err)), nil, nil
	}

	headers := map[string]string{
		"X-Api-Key": apiKey,
	}

	respBody, err := doPost(ctx, patentsViewBaseURL, "application/json", bodyJSON, headers)
	if err != nil {
		return errorResult(fmt.Sprintf("PatentsView API error: %v", err)), nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}

	if errMsg, ok := result["error"].(bool); ok && errMsg {
		msg, _ := result["message"].(string)
		return errorResult(fmt.Sprintf("PatentsView: %s", msg)), nil, nil
	}

	return textResult(formatPatentResults(result)), nil, nil
}

func handleSearchPatentAssignees(ctx context.Context, req *mcp.CallToolRequest, input SearchPatentAssigneesInput) (*mcp.CallToolResult, any, error) {
	apiKey := getPatentsViewAPIKey()
	if apiKey == "" {
		return errorResult("PATENTSVIEW_API_KEY environment variable is not set. Get a free key at https://patentsview.org/apis"), nil, nil
	}

	limit := clampLimit(input.Limit, 25, 100)

	requestBody := map[string]interface{}{
		"q": map[string]interface{}{
			"_text_any": map[string]string{
				"assignee_organization": input.Query,
			},
		},
		"f": []string{
			"assignee_id", "assignee_organization", "assignee_type",
			"assignee_total_num_patents",
		},
		"o": map[string]interface{}{
			"per_page": limit,
		},
		"s": []map[string]string{
			{"assignee_total_num_patents": "desc"},
		},
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to build request: %v", err)), nil, nil
	}

	headers := map[string]string{
		"X-Api-Key": apiKey,
	}

	respBody, err := doPost(ctx, assigneesViewBaseURL, "application/json", bodyJSON, headers)
	if err != nil {
		return errorResult(fmt.Sprintf("PatentsView API error: %v", err)), nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}

	return textResult(formatAssigneeResults(result)), nil, nil
}

func formatPatentResults(data map[string]interface{}) string {
	var sb strings.Builder

	totalHits := 0
	if th, ok := data["total_hits"].(float64); ok {
		totalHits = int(th)
	}
	count := 0
	if c, ok := data["count"].(float64); ok {
		count = int(c)
	}

	sb.WriteString(fmt.Sprintf("Found %d total patents (showing %d):\n\n", totalHits, count))

	patents, ok := data["patents"].([]interface{})
	if !ok {
		sb.WriteString("No patent data in response.\n")
		return sb.String()
	}

	for i, p := range patents {
		patent, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		sb.WriteString(fmt.Sprintf("--- Patent %d ---\n", i+1))
		writeField(&sb, "Patent Number", patent, "patent_id")
		writeField(&sb, "Title", patent, "patent_title")
		writeField(&sb, "Date", patent, "patent_date")
		writeField(&sb, "Type", patent, "patent_type")
		writeField(&sb, "Claims", patent, "patent_num_claims")

		// Assignees
		if assignees, ok := patent["assignees"].([]interface{}); ok {
			for _, a := range assignees {
				if assignee, ok := a.(map[string]interface{}); ok {
					if org, ok := assignee["assignee_organization"].(string); ok && org != "" {
						sb.WriteString(fmt.Sprintf("  Assignee: %s\n", org))
					}
				}
			}
		}

		// Inventors
		if inventors, ok := patent["inventors"].([]interface{}); ok {
			names := make([]string, 0, len(inventors))
			for _, inv := range inventors {
				if inventor, ok := inv.(map[string]interface{}); ok {
					first, _ := inventor["inventor_first_name"].(string)
					last, _ := inventor["inventor_last_name"].(string)
					if first != "" || last != "" {
						names = append(names, strings.TrimSpace(first+" "+last))
					}
				}
			}
			if len(names) > 0 {
				if len(names) > 5 {
					sb.WriteString(fmt.Sprintf("  Inventors: %s, and %d more\n", strings.Join(names[:5], ", "), len(names)-5))
				} else {
					sb.WriteString(fmt.Sprintf("  Inventors: %s\n", strings.Join(names, ", ")))
				}
			}
		}

		// Abstract
		if abstract, ok := patent["patent_abstract"].(string); ok && abstract != "" {
			sb.WriteString(fmt.Sprintf("  Abstract: %s\n", truncate(abstract, 500)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatAssigneeResults(data map[string]interface{}) string {
	var sb strings.Builder

	totalHits := 0
	if th, ok := data["total_hits"].(float64); ok {
		totalHits = int(th)
	}
	count := 0
	if c, ok := data["count"].(float64); ok {
		count = int(c)
	}

	sb.WriteString(fmt.Sprintf("Found %d total assignees (showing %d):\n\n", totalHits, count))

	assignees, ok := data["assignees"].([]interface{})
	if !ok {
		sb.WriteString("No assignee data in response.\n")
		return sb.String()
	}

	for i, a := range assignees {
		assignee, ok := a.(map[string]interface{})
		if !ok {
			continue
		}

		sb.WriteString(fmt.Sprintf("--- Assignee %d ---\n", i+1))
		writeField(&sb, "Organization", assignee, "assignee_organization")
		writeField(&sb, "Type", assignee, "assignee_type")
		writeField(&sb, "Total Patents", assignee, "assignee_total_num_patents")
		writeField(&sb, "ID", assignee, "assignee_id")
		sb.WriteString("\n")
	}
	return sb.String()
}
