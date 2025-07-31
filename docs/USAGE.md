# Accio Usage Guide

This document provides detailed information on how to use Accio effectively.

## Basic Usage

The most basic usage is to search for a username:

```bash
accio -username johndoe
```

This will search for "johndoe" across all supported sites and display the results in the terminal.

## Command-Line Options

### Required Options

- `-username string`: The username to search for (required)

### Output Options

- `-verbose`: Enable verbose output, showing more details including "not found" results
- `-output string`: Save results to a file
- `-format string`: Output format (text, json, csv, markdown) (default "text")
- `-no-color`: Disable colored output in the terminal

### Performance Options

- `-timeout int`: Timeout in seconds for HTTP requests (default 10)
- `-concurrency int`: Number of concurrent requests (default: number of CPU cores)
- `-retries int`: Number of retries for failed requests (default 2)

### Informational Options

- `-version`: Show version information
- `-list-sites`: List all available sites

## Output Formats

Accio supports multiple output formats:

### Text Format (Default)

```bash
accio -username johndoe
```

Output example:
```
[+] GitHub: https://github.com/johndoe
[+] Twitter: https://twitter.com/johndoe
[-] Instagram: Not Found
```

### JSON Format

```bash
accio -username johndoe -format json
```

Output example:
```json
[
  {
    "site": "GitHub",
    "url": "https://github.com/johndoe",
    "exists": true
  },
  {
    "site": "Twitter",
    "url": "https://twitter.com/johndoe",
    "exists": true
  },
  {
    "site": "Instagram",
    "url": "https://www.instagram.com/johndoe",
    "exists": false
  }
]
```

### CSV Format

```bash
accio -username johndoe -format csv
```

Output example:
```
Site,URL,Exists,Error
GitHub,https://github.com/johndoe,true,
Twitter,https://twitter.com/johndoe,true,
Instagram,https://www.instagram.com/johndoe,false,
```

### Markdown Format

```bash
accio -username johndoe -format markdown
```

Output example:
```markdown
# Accio Results

Username search results generated on 2023-04-01T12:34:56Z

## Found Accounts

- [GitHub](https://github.com/johndoe)
- [Twitter](https://twitter.com/johndoe)

## Summary

- **Found**: 2
- **Total**: 3
```

## Saving Results to a File

You can save the results to a file using the `-output` flag:

```bash
accio -username johndoe -output results.txt
```

The file format will be determined by the file extension:
- `.json`: JSON format
- `.csv`: CSV format
- `.md`: Markdown format
- Other extensions: Text format

You can also explicitly specify the format:

```bash
accio -username johndoe -format json -output results.json
```

## Advanced Usage

### Combining Options

You can combine multiple options:

```bash
accio -username johndoe -verbose -timeout 20 -concurrency 50 -retries 3 -format json -output results.json
```

### Piping Output

You can pipe the output to other tools:

```bash
accio -username johndoe -format json | jq '.[] | select(.exists == true)'
```

### Redirecting Output

You can redirect the output to a file:

```bash
accio -username johndoe -format json > results.json
```

## Troubleshooting

### Rate Limiting

If you're getting a lot of errors, you might be getting rate limited. Try:
- Reducing concurrency: `-concurrency 5`
- Increasing timeout: `-timeout 30`
- Increasing retries: `-retries 5`

### False Positives/Negatives

Some sites may return false positives or negatives. Use the `-verbose` flag to see more details about the responses.