# Accio

<p align="center">
  <img src="https://raw.githubusercontent.com/yourusername/accio/main/assets/logo.png" alt="Accio Logo" width="200"/>
</p>

Accio is a powerful command-line tool written in Go for searching usernames across multiple websites. It's inspired by the [Sherlock](https://github.com/sherlock-project/sherlock) OSINT tool but built with Go's performance and concurrency features.

## Features

- Search for usernames across 30+ social networks and websites
- Fast, concurrent checking using Go's goroutines
- Multiple output formats (text, JSON, CSV, Markdown)
- Colored terminal output
- Configurable concurrency and timeout settings
- Automatic retries for failed requests
- Detailed statistics

## Installation

### Using Go

```bash
go install github.com/yourusername/accio/cmd/accio@latest
```

### From Source

```bash
git clone https://github.com/yourusername/accio.git
cd accio
go build ./cmd/accio
```

### Binary Releases

Download the latest binary for your platform from the [Releases](https://github.com/yourusername/accio/releases) page.

## Usage

### Basic Usage

```bash
accio -username johndoe
```

### Command-Line Options

```
Options:
  -username string
        Username to search for
  -verbose
        Enable verbose output
  -timeout int
        Timeout in seconds for HTTP requests (default 10)
  -output string
        Output file to save results
  -format string
        Output format (text, json, csv, markdown) (default "text")
  -no-color
        Disable colored output
  -concurrency int
        Number of concurrent requests (default: number of CPU cores)
  -retries int
        Number of retries for failed requests (default 2)
  -version
        Show version information
  -list-sites
        List all available sites
```

### Examples

Search for a username and display results in the terminal:
```bash
accio -username johndoe
```

Search with verbose output and save results to a file:
```bash
accio -username johndoe -verbose -output results.txt
```

Output results in JSON format:
```bash
accio -username johndoe -format json > results.json
```

Output results in CSV format:
```bash
accio -username johndoe -format csv > results.csv
```

Output results in Markdown format:
```bash
accio -username johndoe -format markdown > results.md
```

List all available sites:
```bash
accio -list-sites
```

## How It Works

Accio checks for the existence of a given username across various websites by:

1. Formatting the site's URL pattern with the provided username
2. Making HTTP requests to each site concurrently
3. Analyzing the response to determine if the username exists
4. Presenting the results in the specified format

## Supported Sites

Accio currently supports 30+ websites including:

- GitHub
- Twitter
- Instagram
- Facebook
- YouTube
- Pinterest
- Reddit
- Twitch
- Medium
- Quora
- Flickr
- Steam
- Vimeo
- SoundCloud
- And many more...

Run `accio -list-sites` to see the full list of supported sites.

## Contributing

Contributions are welcome! Here are some ways you can contribute:

- Add support for new sites
- Improve detection accuracy
- Fix bugs
- Add new features
- Improve documentation

Please feel free to submit a Pull Request.

## License

MIT

## Acknowledgements

- Inspired by [Sherlock](https://github.com/sherlock-project/sherlock)
- Built with Go's powerful concurrency features