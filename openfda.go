package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Query Builders ---

func buildEnforcementSearch(input SearchFDAEnforcementInput) []string {
	var parts []string
	if input.Query != "" {
		q := strings.ReplaceAll(input.Query, " ", "+")
		parts = append(parts, fmt.Sprintf(`reason_for_recall:"%s"+product_description:"%s"`, q, q))
	}
	if input.Status != "" {
		parts = append(parts, fmt.Sprintf(`status:"%s"`, input.Status))
	}
	if input.Classification != "" {
		parts = append(parts, fmt.Sprintf(`classification:"%s"`, input.Classification))
	}
	if input.DateFrom != "" || input.DateTo != "" {
		from := input.DateFrom
		if from == "" {
			from = "19700101"
		}
		to := input.DateTo
		if to == "" {
			to = "21001231"
		}
		parts = append(parts, fmt.Sprintf("report_date:[%s+TO+%s]", from, to))
	}
	return parts
}

func enforcementEndpoint(productType string) string {
	switch strings.ToLower(productType) {
	case "device":
		return "device/enforcement"
	case "food":
		return "food/enforcement"
	default:
		return "drug/enforcement"
	}
}

func buildDrugLabelSearch(input SearchFDADrugLabelsInput) []string {
	var parts []string
	if input.Query != "" {
		q := strings.ReplaceAll(input.Query, " ", "+")
		parts = append(parts, q)
	}
	if input.BrandName != "" {
		parts = append(parts, fmt.Sprintf(`openfda.brand_name:"%s"`, input.BrandName))
	}
	if input.GenericName != "" {
		parts = append(parts, fmt.Sprintf(`openfda.generic_name:"%s"`, input.GenericName))
	}
	if input.Manufacturer != "" {
		parts = append(parts, fmt.Sprintf(`openfda.manufacturer_name:"%s"`, input.Manufacturer))
	}
	return parts
}

func buildAdverseEventSearch(input SearchFDAAdverseEventsInput) []string {
	var parts []string
	if input.DrugName != "" {
		q := strings.ReplaceAll(input.DrugName, " ", "+")
		parts = append(parts, fmt.Sprintf(`patient.drug.openfda.brand_name:"%s"+patient.drug.openfda.generic_name:"%s"`, q, q))
	}
	if input.Reaction != "" {
		parts = append(parts, fmt.Sprintf(`patient.reaction.reactionmeddrapt:"%s"`, input.Reaction))
	}
	if input.Serious != nil && *input.Serious {
		parts = append(parts, "serious:1")
	}
	if input.DateFrom != "" || input.DateTo != "" {
		from := input.DateFrom
		if from == "" {
			from = "19700101"
		}
		to := input.DateTo
		if to == "" {
			to = "21001231"
		}
		parts = append(parts, fmt.Sprintf("receivedate:[%s+TO+%s]", from, to))
	}
	return parts
}

func buildDrugApprovalSearch(input SearchFDADrugApprovalsInput) []string {
	var parts []string
	if input.Query != "" {
		q := strings.ReplaceAll(input.Query, " ", "+")
		parts = append(parts, q)
	}
	if input.BrandName != "" {
		parts = append(parts, fmt.Sprintf(`openfda.brand_name:"%s"`, input.BrandName))
	}
	if input.SponsorName != "" {
		parts = append(parts, fmt.Sprintf(`sponsor_name:"%s"`, input.SponsorName))
	}
	return parts
}

func buildFoodEventSearch(input SearchFDAFoodEventsInput) []string {
	var parts []string
	if input.Query != "" {
		q := strings.ReplaceAll(input.Query, " ", "+")
		parts = append(parts, q)
	}
	if input.Serious != nil && *input.Serious {
		parts = append(parts, `outcomes:"serious"`)
	}
	if input.DateFrom != "" || input.DateTo != "" {
		from := input.DateFrom
		if from == "" {
			from = "19700101"
		}
		to := input.DateTo
		if to == "" {
			to = "21001231"
		}
		parts = append(parts, fmt.Sprintf("date_started:[%s+TO+%s]", from, to))
	}
	return parts
}

// --- Result Formatters ---

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: "Error: " + msg}},
		IsError: true,
	}
}

func formatEnforcementResults(data *OpenFDAResponse) string {
	var sb strings.Builder
	total := 0
	if data.Meta != nil {
		total = data.Meta.Results.Total
	}
	sb.WriteString(fmt.Sprintf("Found %d total results (showing %d):\n\n", total, len(data.Results)))

	for i, raw := range data.Results {
		var r map[string]interface{}
		json.Unmarshal(raw, &r)

		sb.WriteString(fmt.Sprintf("--- Result %d ---\n", i+1))
		writeField(&sb, "Recall Number", r, "recall_number")
		writeField(&sb, "Product", r, "product_description")
		writeField(&sb, "Reason", r, "reason_for_recall")
		writeField(&sb, "Classification", r, "classification")
		writeField(&sb, "Status", r, "status")
		writeField(&sb, "Company", r, "recalling_firm")
		writeField(&sb, "City", r, "city")
		writeField(&sb, "State", r, "state")
		writeField(&sb, "Distribution", r, "distribution_pattern")
		writeField(&sb, "Recall Date", r, "recall_initiation_date")
		writeField(&sb, "Report Date", r, "report_date")
		writeField(&sb, "Product Quantity", r, "product_quantity")
		writeField(&sb, "Code Info", r, "code_info")
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatDrugLabelResults(data *OpenFDAResponse) string {
	var sb strings.Builder
	total := 0
	if data.Meta != nil {
		total = data.Meta.Results.Total
	}
	sb.WriteString(fmt.Sprintf("Found %d total results (showing %d):\n\n", total, len(data.Results)))

	for i, raw := range data.Results {
		var r map[string]interface{}
		json.Unmarshal(raw, &r)

		sb.WriteString(fmt.Sprintf("--- Result %d ---\n", i+1))
		writeNestedField(&sb, "Brand Name", r, "openfda.brand_name")
		writeNestedField(&sb, "Generic Name", r, "openfda.generic_name")
		writeNestedField(&sb, "Manufacturer", r, "openfda.manufacturer_name")
		writeNestedField(&sb, "Product Type", r, "openfda.product_type")
		writeNestedField(&sb, "Route", r, "openfda.route")
		writeArrayField(&sb, "Indications", r, "indications_and_usage", 2000)
		writeArrayField(&sb, "Warnings", r, "warnings", 2000)
		writeArrayField(&sb, "Contraindications", r, "contraindications", 1500)
		writeArrayField(&sb, "Dosage", r, "dosage_and_administration", 1500)
		writeArrayField(&sb, "Adverse Reactions", r, "adverse_reactions", 1500)
		writeArrayField(&sb, "Drug Interactions", r, "drug_interactions", 1000)
		writeArrayField(&sb, "Description", r, "description", 1500)
		writeField(&sb, "Effective Date", r, "effective_time")
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatAdverseEventResults(data *OpenFDAResponse) string {
	var sb strings.Builder
	total := 0
	if data.Meta != nil {
		total = data.Meta.Results.Total
	}
	sb.WriteString(fmt.Sprintf("Found %d total results (showing %d):\n\n", total, len(data.Results)))

	for i, raw := range data.Results {
		var r map[string]interface{}
		json.Unmarshal(raw, &r)

		sb.WriteString(fmt.Sprintf("--- Report %d ---\n", i+1))
		writeField(&sb, "Safety Report ID", r, "safetyreportid")
		writeField(&sb, "Receive Date", r, "receivedate")
		writeField(&sb, "Serious", r, "serious")
		writeField(&sb, "Sender Organization", r, "primarysource.reportercountry")

		// Extract patient info
		if patient, ok := r["patient"].(map[string]interface{}); ok {
			writeField(&sb, "Patient Age", patient, "patientonsetage")
			writeField(&sb, "Patient Sex", patient, "patientsex")

			// Drugs
			if drugs, ok := patient["drug"].([]interface{}); ok {
				for j, d := range drugs {
					if drug, ok := d.(map[string]interface{}); ok {
						sb.WriteString(fmt.Sprintf("  Drug %d: ", j+1))
						if name, ok := drug["medicinalproduct"].(string); ok {
							sb.WriteString(name)
						}
						if char, ok := drug["drugcharacterization"].(string); ok {
							switch char {
							case "1":
								sb.WriteString(" (Suspect)")
							case "2":
								sb.WriteString(" (Concomitant)")
							case "3":
								sb.WriteString(" (Interacting)")
							}
						}
						sb.WriteString("\n")
					}
				}
			}

			// Reactions
			if reactions, ok := patient["reaction"].([]interface{}); ok {
				rxns := make([]string, 0, len(reactions))
				for _, r := range reactions {
					if rxn, ok := r.(map[string]interface{}); ok {
						if term, ok := rxn["reactionmeddrapt"].(string); ok {
							rxns = append(rxns, term)
						}
					}
				}
				if len(rxns) > 0 {
					sb.WriteString(fmt.Sprintf("  Reactions: %s\n", strings.Join(rxns, ", ")))
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatCountResults(data []json.RawMessage) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Count results (%d entries):\n\n", len(data)))

	for _, raw := range data {
		var entry map[string]interface{}
		json.Unmarshal(raw, &entry)
		term := ""
		count := 0.0
		if t, ok := entry["term"].(string); ok {
			term = t
		}
		if c, ok := entry["count"].(float64); ok {
			count = c
		}
		sb.WriteString(fmt.Sprintf("  %s: %.0f\n", term, count))
	}
	return sb.String()
}

func formatDrugApprovalResults(data *OpenFDAResponse) string {
	var sb strings.Builder
	total := 0
	if data.Meta != nil {
		total = data.Meta.Results.Total
	}
	sb.WriteString(fmt.Sprintf("Found %d total results (showing %d):\n\n", total, len(data.Results)))

	for i, raw := range data.Results {
		var r map[string]interface{}
		json.Unmarshal(raw, &r)

		sb.WriteString(fmt.Sprintf("--- Result %d ---\n", i+1))
		writeField(&sb, "Application Number", r, "application_number")
		writeField(&sb, "Sponsor", r, "sponsor_name")
		writeNestedField(&sb, "Brand Name", r, "openfda.brand_name")
		writeNestedField(&sb, "Generic Name", r, "openfda.generic_name")
		writeNestedField(&sb, "Manufacturer", r, "openfda.manufacturer_name")
		writeNestedField(&sb, "Substance", r, "openfda.substance_name")
		writeNestedField(&sb, "Route", r, "openfda.route")
		writeNestedField(&sb, "Product Type", r, "openfda.product_type")

		// Products/submissions
		if products, ok := r["products"].([]interface{}); ok {
			for j, p := range products {
				if prod, ok := p.(map[string]interface{}); ok {
					sb.WriteString(fmt.Sprintf("  Product %d:\n", j+1))
					if name, ok := prod["brand_name"].(string); ok {
						sb.WriteString(fmt.Sprintf("    Brand: %s\n", name))
					}
					if dosage, ok := prod["dosage_form"].(string); ok {
						sb.WriteString(fmt.Sprintf("    Form: %s\n", dosage))
					}
					if route, ok := prod["route"].(string); ok {
						sb.WriteString(fmt.Sprintf("    Route: %s\n", route))
					}
					if active, ok := prod["active_ingredients"].([]interface{}); ok {
						for _, ai := range active {
							if ingredient, ok := ai.(map[string]interface{}); ok {
								name, _ := ingredient["name"].(string)
								strength, _ := ingredient["strength"].(string)
								if name != "" {
									sb.WriteString(fmt.Sprintf("    Ingredient: %s %s\n", name, strength))
								}
							}
						}
					}
				}
				if j >= 2 {
					sb.WriteString(fmt.Sprintf("    ... and %d more products\n", len(products)-3))
					break
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatFoodEventResults(data *OpenFDAResponse) string {
	var sb strings.Builder
	total := 0
	if data.Meta != nil {
		total = data.Meta.Results.Total
	}
	sb.WriteString(fmt.Sprintf("Found %d total results (showing %d):\n\n", total, len(data.Results)))

	for i, raw := range data.Results {
		var r map[string]interface{}
		json.Unmarshal(raw, &r)

		sb.WriteString(fmt.Sprintf("--- Report %d ---\n", i+1))
		writeField(&sb, "Report Number", r, "report_number")
		writeField(&sb, "Date Started", r, "date_started")
		writeField(&sb, "Date Created", r, "date_created")

		// Products
		if products, ok := r["products"].([]interface{}); ok {
			for _, p := range products {
				if prod, ok := p.(map[string]interface{}); ok {
					writeField(&sb, "Product", prod, "name_brand")
					writeField(&sb, "Industry", prod, "industry_name")
					writeField(&sb, "Role", prod, "role")
				}
			}
		}

		// Reactions
		if reactions, ok := r["reactions"].([]interface{}); ok {
			rxns := make([]string, 0, len(reactions))
			for _, rx := range reactions {
				if s, ok := rx.(string); ok {
					rxns = append(rxns, s)
				}
			}
			if len(rxns) > 0 {
				sb.WriteString(fmt.Sprintf("  Reactions: %s\n", strings.Join(rxns, ", ")))
			}
		}

		// Outcomes
		if outcomes, ok := r["outcomes"].([]interface{}); ok {
			outs := make([]string, 0, len(outcomes))
			for _, o := range outcomes {
				if s, ok := o.(string); ok {
					outs = append(outs, s)
				}
			}
			if len(outs) > 0 {
				sb.WriteString(fmt.Sprintf("  Outcomes: %s\n", strings.Join(outs, ", ")))
			}
		}

		writeField(&sb, "Consumer Age", r, "consumer.age")
		writeField(&sb, "Consumer Gender", r, "consumer.gender")
		sb.WriteString("\n")
	}
	return sb.String()
}

// --- Tool Handlers ---

func handleSearchEnforcement(ctx context.Context, req *mcp.CallToolRequest, input SearchFDAEnforcementInput) (*mcp.CallToolResult, any, error) {
	limit := clampLimit(input.Limit, 10, 100)
	searchParts := buildEnforcementSearch(input)
	endpoint := enforcementEndpoint(input.ProductType)
	apiURL := buildOpenFDAURL(endpoint, searchParts, limit, 0, "")

	body, err := doGet(ctx, apiURL)
	if err != nil {
		return errorResult(fmt.Sprintf("openFDA API error: %v", err)), nil, nil
	}

	var resp OpenFDAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}
	if resp.Error != nil {
		return errorResult(fmt.Sprintf("openFDA: %s - %s", resp.Error.Code, resp.Error.Message)), nil, nil
	}

	return textResult(formatEnforcementResults(&resp)), nil, nil
}

func handleSearchDrugLabels(ctx context.Context, req *mcp.CallToolRequest, input SearchFDADrugLabelsInput) (*mcp.CallToolResult, any, error) {
	limit := clampLimit(input.Limit, 5, 100)
	searchParts := buildDrugLabelSearch(input)
	apiURL := buildOpenFDAURL("drug/label", searchParts, limit, 0, "")

	body, err := doGet(ctx, apiURL)
	if err != nil {
		return errorResult(fmt.Sprintf("openFDA API error: %v", err)), nil, nil
	}

	var resp OpenFDAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}
	if resp.Error != nil {
		return errorResult(fmt.Sprintf("openFDA: %s - %s", resp.Error.Code, resp.Error.Message)), nil, nil
	}

	return textResult(formatDrugLabelResults(&resp)), nil, nil
}

func handleSearchAdverseEvents(ctx context.Context, req *mcp.CallToolRequest, input SearchFDAAdverseEventsInput) (*mcp.CallToolResult, any, error) {
	limit := clampLimit(input.Limit, 10, 100)
	searchParts := buildAdverseEventSearch(input)

	apiURL := buildOpenFDAURL("drug/event", searchParts, limit, 0, input.CountField)

	body, err := doGet(ctx, apiURL)
	if err != nil {
		return errorResult(fmt.Sprintf("openFDA API error: %v", err)), nil, nil
	}

	// Count mode returns a different structure
	if input.CountField != "" {
		var countResp struct {
			Results []json.RawMessage `json:"results"`
		}
		if err := json.Unmarshal(body, &countResp); err != nil {
			return errorResult(fmt.Sprintf("Failed to parse count response: %v", err)), nil, nil
		}
		return textResult(formatCountResults(countResp.Results)), nil, nil
	}

	var resp OpenFDAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}
	if resp.Error != nil {
		return errorResult(fmt.Sprintf("openFDA: %s - %s", resp.Error.Code, resp.Error.Message)), nil, nil
	}

	return textResult(formatAdverseEventResults(&resp)), nil, nil
}

func handleSearchDrugApprovals(ctx context.Context, req *mcp.CallToolRequest, input SearchFDADrugApprovalsInput) (*mcp.CallToolResult, any, error) {
	limit := clampLimit(input.Limit, 10, 99)
	searchParts := buildDrugApprovalSearch(input)
	apiURL := buildOpenFDAURL("drug/drugsfda", searchParts, limit, 0, "")

	body, err := doGet(ctx, apiURL)
	if err != nil {
		return errorResult(fmt.Sprintf("openFDA API error: %v", err)), nil, nil
	}

	var resp OpenFDAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}
	if resp.Error != nil {
		return errorResult(fmt.Sprintf("openFDA: %s - %s", resp.Error.Code, resp.Error.Message)), nil, nil
	}

	return textResult(formatDrugApprovalResults(&resp)), nil, nil
}

func handleSearchFoodAdverseEvents(ctx context.Context, req *mcp.CallToolRequest, input SearchFDAFoodEventsInput) (*mcp.CallToolResult, any, error) {
	limit := clampLimit(input.Limit, 10, 100)
	searchParts := buildFoodEventSearch(input)
	apiURL := buildOpenFDAURL("food/event", searchParts, limit, 0, "")

	body, err := doGet(ctx, apiURL)
	if err != nil {
		return errorResult(fmt.Sprintf("openFDA API error: %v", err)), nil, nil
	}

	var resp OpenFDAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}
	if resp.Error != nil {
		return errorResult(fmt.Sprintf("openFDA: %s - %s", resp.Error.Code, resp.Error.Message)), nil, nil
	}

	return textResult(formatFoodEventResults(&resp)), nil, nil
}

func handleLookupRecall(ctx context.Context, req *mcp.CallToolRequest, input LookupFDARecallInput) (*mcp.CallToolResult, any, error) {
	endpoint := enforcementEndpoint(input.ProductType)
	searchParts := []string{fmt.Sprintf(`recall_number:"%s"`, input.RecallNumber)}
	apiURL := buildOpenFDAURL(endpoint, searchParts, 1, 0, "")

	body, err := doGet(ctx, apiURL)
	if err != nil {
		return errorResult(fmt.Sprintf("openFDA API error: %v", err)), nil, nil
	}

	var resp OpenFDAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return errorResult(fmt.Sprintf("Failed to parse response: %v", err)), nil, nil
	}
	if resp.Error != nil {
		return errorResult(fmt.Sprintf("openFDA: %s - %s", resp.Error.Code, resp.Error.Message)), nil, nil
	}
	if len(resp.Results) == 0 {
		return textResult(fmt.Sprintf("No recall found with number: %s. Try a different product_type (drug, device, food).", input.RecallNumber)), nil, nil
	}

	return textResult(formatEnforcementResults(&resp)), nil, nil
}
