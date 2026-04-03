package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func handleSearchFederalRegister(ctx context.Context, req *mcp.CallToolRequest, input SearchFederalRegisterInput) (*mcp.CallToolResult, any, error) {
	base := "https://www.federalregister.gov/api/v1/documents.json"
	u, err := url.Parse(base)
	if err != nil {
		return errorResult("Failed to build URL"), nil, nil
	}

	q := u.Query()
	q.Set("conditions[term]", input.Term)

	if input.Agency != "" {
		q.Set("conditions[agencies][]", input.Agency)
	}
	if input.DocType != "" {
		q.Set("conditions[type][]", input.DocType)
	}
	if input.DateFrom != "" {
		q.Set("conditions[publication_date][gte]", input.DateFrom)
	}
	if input.DateTo != "" {
		q.Set("conditions[publication_date][lte]", input.DateTo)
	}

	perPage := clampLimit(input.PerPage, 10, 100)
	q.Set("per_page", strconv.Itoa(perPage))

	if input.Page > 0 {
		q.Set("page", strconv.Itoa(input.Page))
	}

	q.Set("fields[]", "title,type,abstract,document_number,html_url,pdf_url,publication_date,agencies,excerpts")
	u.RawQuery = q.Encode()

	body, err := doGet(ctx, u.String())
	if err != nil {
		return errorResult(fmt.Sprintf("Federal Register API error: %v", err)), nil, nil
	}

	var resp FederalRegisterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}

	return textResult(formatFederalRegisterResults(&resp)), nil, nil
}

func formatFederalRegisterResults(data *FederalRegisterResponse) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d total results (page of %d total pages, showing %d):\n\n",
		data.Count, data.TotalPages, len(data.Results)))

	for i, doc := range data.Results {
		sb.WriteString(fmt.Sprintf("--- Result %d ---\n", i+1))
		sb.WriteString(fmt.Sprintf("  Title: %s\n", doc.Title))
		sb.WriteString(fmt.Sprintf("  Type: %s\n", doc.Type))
		sb.WriteString(fmt.Sprintf("  Published: %s\n", doc.PublicationDate))
		sb.WriteString(fmt.Sprintf("  Document Number: %s\n", doc.DocumentNumber))

		if len(doc.Agencies) > 0 {
			names := make([]string, 0, len(doc.Agencies))
			for _, a := range doc.Agencies {
				if a.RawName != "" {
					names = append(names, a.RawName)
				} else if a.Name != "" {
					names = append(names, a.Name)
				}
			}
			if len(names) > 0 {
				sb.WriteString(fmt.Sprintf("  Agencies: %s\n", strings.Join(names, ", ")))
			}
		}

		if doc.Abstract != "" {
			sb.WriteString(fmt.Sprintf("  Abstract: %s\n", truncate(doc.Abstract, 1000)))
		}
		if doc.Excerpts != "" {
			sb.WriteString(fmt.Sprintf("  Excerpt: %s\n", truncate(doc.Excerpts, 500)))
		}
		if doc.HTMLURL != "" {
			sb.WriteString(fmt.Sprintf("  URL: %s\n", doc.HTMLURL))
		}
		if doc.PDFURL != "" {
			sb.WriteString(fmt.Sprintf("  PDF: %s\n", doc.PDFURL))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
