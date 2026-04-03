package main

import "encoding/json"

// --- openFDA input types ---

type SearchFDAEnforcementInput struct {
	Query          string `json:"query" jsonschema:"Search term (product name, company, reason, etc.)"`
	ProductType    string `json:"product_type,omitempty" jsonschema:"Filter by type: drug, device, or food. Defaults to all."`
	Status         string `json:"status,omitempty" jsonschema:"Filter by status: Ongoing, Completed, or Terminated"`
	Classification string `json:"classification,omitempty" jsonschema:"Recall class: Class I, Class II, or Class III"`
	DateFrom       string `json:"date_from,omitempty" jsonschema:"Start date in YYYYMMDD format"`
	DateTo         string `json:"date_to,omitempty" jsonschema:"End date in YYYYMMDD format"`
	Limit          int    `json:"limit,omitempty" jsonschema:"Max results 1-100, default 10"`
}

type SearchFDADrugLabelsInput struct {
	Query        string `json:"query" jsonschema:"Search term (drug name, active ingredient, indication, etc.)"`
	BrandName    string `json:"brand_name,omitempty" jsonschema:"Filter by brand name"`
	GenericName  string `json:"generic_name,omitempty" jsonschema:"Filter by generic name"`
	Manufacturer string `json:"manufacturer,omitempty" jsonschema:"Filter by manufacturer name"`
	Limit        int    `json:"limit,omitempty" jsonschema:"Max results 1-100, default 5"`
}

type SearchFDAAdverseEventsInput struct {
	DrugName   string `json:"drug_name" jsonschema:"Drug brand or generic name to search"`
	Reaction   string `json:"reaction,omitempty" jsonschema:"Filter by adverse reaction MedDRA term"`
	Serious    *bool  `json:"serious,omitempty" jsonschema:"Filter for serious events only"`
	DateFrom   string `json:"date_from,omitempty" jsonschema:"Start date in YYYYMMDD format"`
	DateTo     string `json:"date_to,omitempty" jsonschema:"End date in YYYYMMDD format"`
	Limit      int    `json:"limit,omitempty" jsonschema:"Max results 1-100, default 10"`
	CountField string `json:"count_field,omitempty" jsonschema:"Field to aggregate/count instead of returning records (e.g. patient.reaction.reactionmeddrapt.exact)"`
}

type SearchFDADrugApprovalsInput struct {
	Query       string `json:"query" jsonschema:"Search term (drug name, sponsor, application number)"`
	BrandName   string `json:"brand_name,omitempty" jsonschema:"Filter by brand name"`
	SponsorName string `json:"sponsor_name,omitempty" jsonschema:"Filter by sponsor/company name"`
	Limit       int    `json:"limit,omitempty" jsonschema:"Max results 1-99, default 10"`
}

type SearchFDAFoodEventsInput struct {
	Query    string `json:"query" jsonschema:"Search term (product name, symptom, industry)"`
	DateFrom string `json:"date_from,omitempty" jsonschema:"Start date in YYYYMMDD format"`
	DateTo   string `json:"date_to,omitempty" jsonschema:"End date in YYYYMMDD format"`
	Serious  *bool  `json:"serious,omitempty" jsonschema:"Filter for serious outcomes only"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Max results 1-100, default 10"`
}

type LookupFDARecallInput struct {
	RecallNumber string `json:"recall_number" jsonschema:"FDA recall number (e.g. D-0572-2024)"`
	ProductType  string `json:"product_type,omitempty" jsonschema:"Product type: drug, device, or food. Default: drug"`
}

// --- Federal Register input types ---

type SearchFederalRegisterInput struct {
	Term     string `json:"term" jsonschema:"Search term or phrase"`
	Agency   string `json:"agency,omitempty" jsonschema:"Agency slug (e.g. food-and-drug-administration, patent-and-trademark-office)"`
	DocType  string `json:"doc_type,omitempty" jsonschema:"Document type: RULE, PRORULE, NOTICE, or PRESDOCU"`
	DateFrom string `json:"date_from,omitempty" jsonschema:"Start date in YYYY-MM-DD format"`
	DateTo   string `json:"date_to,omitempty" jsonschema:"End date in YYYY-MM-DD format"`
	PerPage  int    `json:"per_page,omitempty" jsonschema:"Results per page 1-100, default 10"`
	Page     int    `json:"page,omitempty" jsonschema:"Page number, default 1"`
}

// --- PatentsView input types ---

type SearchPatentsInput struct {
	Query    string `json:"query" jsonschema:"Search text to match in patent title and abstract"`
	Assignee string `json:"assignee,omitempty" jsonschema:"Filter by assignee/company name"`
	Inventor string `json:"inventor,omitempty" jsonschema:"Filter by inventor last name"`
	DateFrom string `json:"date_from,omitempty" jsonschema:"Start date in YYYY-MM-DD format"`
	DateTo   string `json:"date_to,omitempty" jsonschema:"End date in YYYY-MM-DD format"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Max results 1-100, default 25"`
}

type SearchPatentAssigneesInput struct {
	Query string `json:"query" jsonschema:"Search assignee/organization name"`
	Limit int    `json:"limit,omitempty" jsonschema:"Max results 1-100, default 25"`
}

// --- openFDA response types ---

type OpenFDAResponse struct {
	Meta    *OpenFDAMeta      `json:"meta,omitempty"`
	Results []json.RawMessage `json:"results,omitempty"`
	Error   *OpenFDAError     `json:"error,omitempty"`
}

type OpenFDAMeta struct {
	Results OpenFDAResultsMeta `json:"results"`
}

type OpenFDAResultsMeta struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type OpenFDAError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// --- Federal Register response types ---

type FederalRegisterResponse struct {
	Count      int                       `json:"count"`
	TotalPages int                       `json:"total_pages"`
	Results    []FederalRegisterDocument  `json:"results"`
}

type FederalRegisterDocument struct {
	Title           string                  `json:"title"`
	Type            string                  `json:"type"`
	Abstract        string                  `json:"abstract"`
	DocumentNumber  string                  `json:"document_number"`
	HTMLURL         string                  `json:"html_url"`
	PDFURL          string                  `json:"pdf_url"`
	PublicationDate string                  `json:"publication_date"`
	Agencies        []FederalRegisterAgency `json:"agencies"`
	Excerpts        string                  `json:"excerpts"`
}

type FederalRegisterAgency struct {
	RawName string `json:"raw_name"`
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
}
