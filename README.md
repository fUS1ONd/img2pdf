# img2pdf

Simple image to PDF converter written in Go.

## Supported formats

- **JPG/JPEG**
- **PNG**
- **WEBP** 
- **TIFF/TIF**

## Installation

### From source

```bash
git clone https://github.com/fUS1ONd/img2pdf.git
cd img2pdf
go build
```

### Binary release

Download the latest version from [Releases](https://github.com/fUS1ONd/img2pdf/releases).

## Usage

```bash
# Convert all images from directory
./img2pdf -i ./photos -o result.pdf

# Convert specific files
./img2pdf -i "photo1.jpg,image.png,scan.tiff" -o document.pdf

# Convert all in one line
./img2pdf -i "photo1.jpg,photos/,scan.tiff,image.png" -order nam
# Show help
./img2pdf -help
```

## Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `-i` | Directory or comma-separated list of files | required |
| `-o` | Output PDF file path | `output.pdf` |
| `-order` | Set order that pages are saving in pdf | `seq` |
| `-help` | Show help | - |

*NEW* order types:
 - sequently = `seq` (default)
 - naming = `nam`
 - modtime = `mod`

## Features

- Sorting by sequently\modtime\naming
- Support for different formats in one PDF
- Simple friendly command line interface
- Wonderful tests (converter cover ~90%)

## Development

### Tests

```bash
go test -v
```

### Build

```bash
go build
```
