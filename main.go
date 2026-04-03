package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

func main() {
	log.SetOutput(os.Stderr)

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "govmcp",
			Version: version,
		},
		nil,
	)

	// --- FDA Tools (openFDA, no API key required) ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_fda_enforcement",
		Description: "Search FDA drug, device, and food recall enforcement actions via openFDA. Returns recall details including product, reason, classification, and status. Use product_type to filter by drug, device, or food.",
	}, handleSearchEnforcement)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_fda_drug_labels",
		Description: "Search FDA drug labeling (SPL/structured product labeling) data via openFDA. Returns drug label information including indications, warnings, dosage, contraindications, and manufacturer details.",
	}, handleSearchDrugLabels)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_fda_adverse_events",
		Description: "Search FDA drug adverse event reports (FAERS) via openFDA. Returns individual case safety reports with drug names, reactions, and patient info. Use count_field parameter (e.g. patient.reaction.reactionmeddrapt.exact) to get aggregated reaction counts instead of individual reports.",
	}, handleSearchAdverseEvents)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_fda_drug_approvals",
		Description: "Search FDA drug application and approval data (Drugs@FDA) via openFDA. Returns drug approval history, application numbers, sponsor information, and product details.",
	}, handleSearchDrugApprovals)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_fda_food_adverse_events",
		Description: "Search FDA food and dietary supplement adverse event reports via openFDA. Returns reports of adverse events associated with food products including reactions and outcomes.",
	}, handleSearchFoodAdverseEvents)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "lookup_fda_recall",
		Description: "Look up a specific FDA recall by its recall number (e.g. D-0572-2024). Returns detailed information about the recall including product, reason, classification, and distribution.",
	}, handleLookupRecall)

	// --- Federal Register Tool (no API key required) ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_federal_register",
		Description: "Search the Federal Register for rules, proposed rules, notices, and presidential documents. Filter by agency (e.g. food-and-drug-administration for FDA, patent-and-trademark-office for USPTO), document type, and date range. Useful for finding FDA guidance documents, warning letter notices, and regulatory actions.",
	}, handleSearchFederalRegister)

	// --- USPTO/PatentsView Tools (require PATENTSVIEW_API_KEY env var) ---

	if getPatentsViewAPIKey() != "" {
		log.Println("PatentsView API key found, registering patent search tools")

		mcp.AddTool(server, &mcp.Tool{
			Name:        "search_patents",
			Description: "Search US patents via the PatentsView API. Search by text in title/abstract, assignee company, inventor name, and date range. Returns patent numbers, titles, dates, assignees, inventors, and abstracts.",
		}, handleSearchPatents)

		mcp.AddTool(server, &mcp.Tool{
			Name:        "search_patent_assignees",
			Description: "Search patent assignees (companies/organizations) via the PatentsView API. Returns assignee names, types, and total patent counts. Useful for competitive patent landscape analysis.",
		}, handleSearchPatentAssignees)
	} else {
		log.Println("No PATENTSVIEW_API_KEY found, patent search tools not registered. Set this env var to enable USPTO patent search.")
	}

	// Run on stdio
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
