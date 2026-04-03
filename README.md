# govmcp

A Go-based MCP (Model Context Protocol) server that provides search access to open US government data — FDA recalls, drug labels, adverse events, Federal Register documents, and USPTO patents.

Compiles to a single binary with no runtime dependencies.

## Tools

### FDA (no API key required)

| Tool | Description |
|------|-------------|
| `search_fda_enforcement` | Search drug, device, and food recall enforcement actions |
| `search_fda_drug_labels` | Search drug labeling (SPL) — indications, warnings, dosage |
| `search_fda_adverse_events` | Search FAERS adverse event reports, with aggregation support |
| `search_fda_drug_approvals` | Search Drugs@FDA approval history |
| `search_fda_food_adverse_events` | Search food/dietary supplement adverse event reports |
| `lookup_fda_recall` | Look up a specific recall by number |

### Federal Register (no API key required)

| Tool | Description |
|------|-------------|
| `search_federal_register` | Search rules, notices, proposed rules — filterable by agency |

### USPTO (requires API key)

| Tool | Description |
|------|-------------|
| `search_patents` | Search US patents by text, assignee, inventor, date |
| `search_patent_assignees` | Search patent assignees/companies |

USPTO tools require a free PatentsView API key. Set the `PATENTSVIEW_API_KEY` environment variable to enable them. Get a key at [patentsview.org](https://patentsview.org/apis).

## Install

### Option 1: `go install`

```bash
go install github.com/thomastx05/govmcp@latest
```

The binary will be placed in your `$GOPATH/bin` (or `$HOME/go/bin`).

### Option 2: Build from source

```bash
git clone https://github.com/thomastx05/govmcp.git
cd govmcp
go build -o govmcp .
```

On Windows this produces `govmcp.exe`.

## Configure in Claude Code

Add to your Claude Code settings (`~/.claude/settings.json`):

```json
{
  "mcpServers": {
    "govmcp": {
      "command": "govmcp",
      "env": {
        "PATENTSVIEW_API_KEY": "your-key-here"
      }
    }
  }
}
```

If you built from source instead of using `go install`, use the full path to the binary:

```json
{
  "mcpServers": {
    "govmcp": {
      "command": "/path/to/govmcp"
    }
  }
}
```

The `PATENTSVIEW_API_KEY` env var is optional — omit it if you don't need USPTO patent search.

You can also set `OPENFDA_API_KEY` for higher rate limits on FDA queries (optional, not required).

## Configure in Claude Desktop

Add to your Claude Desktop config (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "govmcp": {
      "command": "govmcp",
      "args": [],
      "env": {
        "PATENTSVIEW_API_KEY": "your-key-here"
      }
    }
  }
}
```

Config file locations:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

## Example queries

Once installed, you can ask Claude things like:

- "Search for FDA recalls related to metformin"
- "What adverse events have been reported for Ozempic?"
- "Find FDA drug label info for lisinopril"
- "Search the Federal Register for FDA guidance on AI in medical devices"
- "Look up patents from Google related to large language models"

## Data sources

- [openFDA](https://open.fda.gov/) — Drug labels, adverse events, recalls, approvals, food events
- [Federal Register API](https://www.federalregister.gov/developers/documentation/api/v1) — Rules, notices, proposed rules
- [PatentsView API](https://patentsview.org/apis) — US patent data

## License

MIT
