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
go build -o img2pdf
```

### Binary release

Download the latest version from [Releases](https://github.com/fUS1ONd/img2pdf/releases).

## Usage

```bash
# Convert all images from directory
./img2pdf -input ./photos -output result.pdf

# Convert specific files
./img2pdf -input "photo1.jpg,image.png,scan.tiff" -output document.pdf

# Show help
./img2pdf -help
```

## Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `-input` | Directory or comma-separated list of files | required |
| `-output` | Output PDF file path | `output.pdf` |
| `-help` | Show help | - |

## Features

- Automatic sorting by creation time
- Support for different formats in one PDF
- Simple command line interface

## Development

### Tests

```bash
go test -v
```

### Build

```bash
go build
```
