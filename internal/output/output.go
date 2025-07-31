package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Result represents the result of checking a username on a site
type Result struct {
	Site     string `json:"site"`
	URL      string `json:"url"`
	Exists   bool   `json:"exists"`
	Error    error  `json:"-"`
	Response string `json:"-"`
}

// MarshalJSON custom JSON marshaling to handle the error field
func (r Result) MarshalJSON() ([]byte, error) {
	type Alias Result
	var errStr string
	if r.Error != nil {
		errStr = r.Error.Error()
	}
	return json.Marshal(&struct {
		Alias
		Error string `json:"error,omitempty"`
	}{
		Alias: Alias(r),
		Error: errStr,
	})
}

// FormatType represents the output format type
type FormatType string

const (
	// FormatText is plain text output
	FormatText FormatType = "text"
	// FormatJSON is JSON output
	FormatJSON FormatType = "json"
	// FormatCSV is CSV output
	FormatCSV FormatType = "csv"
	// FormatMarkdown is Markdown output
	FormatMarkdown FormatType = "markdown"
)

// Formatter handles the formatting and output of results
type Formatter struct {
	Verbose bool
	Format  FormatType
	Color   bool
}

// NewFormatter creates a new Formatter instance
func NewFormatter(verbose bool) *Formatter {
	return &Formatter{
		Verbose: verbose,
		Format:  FormatText,
		Color:   true,
	}
}

// WithFormat sets the output format
func (f *Formatter) WithFormat(format FormatType) *Formatter {
	f.Format = format
	return f
}

// WithColor enables or disables colored output
func (f *Formatter) WithColor(color bool) *Formatter {
	f.Color = color
	return f
}

// PrintResult prints a single result
func (f *Formatter) PrintResult(result Result) {
	switch f.Format {
	case FormatJSON:
		// For JSON format, we don't print individual results
		// They will be collected and printed at the end
		return
	case FormatCSV:
		// For CSV format, we don't print individual results
		// They will be collected and printed at the end
		return
	case FormatMarkdown:
		if result.Exists {
			fmt.Printf("- [x] %s: [%s](%s)\n", result.Site, result.Site, result.URL)
		} else if f.Verbose {
			fmt.Printf("- [ ] %s\n", result.Site)
		}
	default: // FormatText
		if result.Exists {
			if f.Color {
				fmt.Printf("\033[32m[+]\033[0m %s: %s\n", result.Site, result.URL)
			} else {
				fmt.Printf("[+] %s: %s\n", result.Site, result.URL)
			}
		} else if f.Verbose {
			if f.Color {
				fmt.Printf("\033[31m[-]\033[0m %s: Not Found\n", result.Site)
			} else {
				fmt.Printf("[-] %s: Not Found\n", result.Site)
			}
		}

		if result.Error != nil && f.Verbose {
			if f.Color {
				fmt.Printf("    \033[33mError: %v\033[0m\n", result.Error)
			} else {
				fmt.Printf("    Error: %v\n", result.Error)
			}
		}
	}
}

// PrintSummary prints a summary of all results
func (f *Formatter) PrintSummary(results []Result) {
	switch f.Format {
	case FormatJSON:
		// Print all results as JSON
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Error generating JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
		return
	case FormatCSV:
		// Print all results as CSV to stdout
		writer := csv.NewWriter(os.Stdout)
		writer.Write([]string{"Site", "URL", "Exists", "Error"})
		for _, result := range results {
			var errStr string
			if result.Error != nil {
				errStr = result.Error.Error()
			}
			writer.Write([]string{
				result.Site,
				result.URL,
				fmt.Sprintf("%t", result.Exists),
				errStr,
			})
		}
		writer.Flush()
		return
	case FormatMarkdown:
		var found int
		for _, result := range results {
			if result.Exists {
				found++
			}
		}
		fmt.Printf("\n## Summary\n\n")
		fmt.Printf("- **Found**: %d\n", found)
		fmt.Printf("- **Total**: %d\n", len(results))
		fmt.Printf("- **Time**: %s\n", time.Now().Format(time.RFC3339))
		return
	default: // FormatText
		var found int
		for _, result := range results {
			if result.Exists {
				found++
			}
		}

		if f.Color {
			fmt.Printf("\n\033[1mFound %d results out of %d sites\033[0m\n", found, len(results))
		} else {
			fmt.Printf("\nFound %d results out of %d sites\n", found, len(results))
		}
	}
}

// SaveToFile saves results to a file
func (f *Formatter) SaveToFile(results []Result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Determine format based on file extension
	format := f.Format
	if strings.HasSuffix(filename, ".json") {
		format = FormatJSON
	} else if strings.HasSuffix(filename, ".csv") {
		format = FormatCSV
	} else if strings.HasSuffix(filename, ".md") {
		format = FormatMarkdown
	}

	switch format {
	case FormatJSON:
		// Save as JSON
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(results)
	case FormatCSV:
		// Save as CSV
		writer := csv.NewWriter(file)
		writer.Write([]string{"Site", "URL", "Exists", "Error"})
		for _, result := range results {
			var errStr string
			if result.Error != nil {
				errStr = result.Error.Error()
			}
			writer.Write([]string{
				result.Site,
				result.URL,
				fmt.Sprintf("%t", result.Exists),
				errStr,
			})
		}
		writer.Flush()
		return writer.Error()
	case FormatMarkdown:
		// Save as Markdown
		fmt.Fprintf(file, "# Accio Results\n\n")
		fmt.Fprintf(file, "Username search results generated on %s\n\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(file, "## Found Accounts\n\n")

		for _, result := range results {
			if result.Exists {
				fmt.Fprintf(file, "- [%s](%s)\n", result.Site, result.URL)
			}
		}

		var found int
		for _, result := range results {
			if result.Exists {
				found++
			}
		}

		fmt.Fprintf(file, "\n## Summary\n\n")
		fmt.Fprintf(file, "- **Found**: %d\n", found)
		fmt.Fprintf(file, "- **Total**: %d\n", len(results))

		return nil
	default: // FormatText
		// Save as plain text
		for _, result := range results {
			if result.Exists {
				fmt.Fprintf(file, "[+] %s: %s\n", result.Site, result.URL)
			} else if f.Verbose {
				fmt.Fprintf(file, "[-] %s: Not Found\n", result.Site)
			}

			if result.Error != nil && f.Verbose {
				fmt.Fprintf(file, "    Error: %v\n", result.Error)
			}
		}

		var found int
		for _, result := range results {
			if result.Exists {
				found++
			}
		}

		fmt.Fprintf(file, "\nFound %d results out of %d sites\n", found, len(results))
		return nil
	}
}
