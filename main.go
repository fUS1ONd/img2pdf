package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		input  = flag.String("i", "", "Input directory or JPG files (space-separated)")
		order  = flag.String("order", "seq", "default sequently order")
		output = flag.String("o", "output.pdf", "Output PDF file path")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help && *input == "" {
		printUsage()
		return
	}

	if !*help && *input == "" {
		printInputHelp()
		return
	}

	converter := NewConverter()
	if err := converter.Convert(*input, *output, *order); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted to %s\n", *output)
}

func printUsage() {
	fmt.Println("Image to PDF Converter")
	fmt.Println("\nSupported formats: JPG, JPEG, PNG, WEBP, TIFF")
	fmt.Println("\nUsage:")
	fmt.Println("  ./img2pdf -i <directory|files> -o <pdf_file> -order <order type>")
	fmt.Println("\nExamples:")
	fmt.Println("  ./img2pdf -i images/")
	fmt.Println("  ./img2pdf -i \"image1.jpg,photo.png,scan.tiff\" -o result.pdf")
	fmt.Println("  ./img2pdf -i \"images/,photo.jpg,scan.png\" -o result.pdf -order mod")
	fmt.Println("\nNote: The -i flag accepts both directories and individual files (comma-separated)")
	fmt.Println("\nOptions:")
	fmt.Println("  -i string")
	fmt.Println("    \tInput directory or JPG files (space-separated)")
	fmt.Println("  -o string")
	fmt.Println("    \tOutput PDF file path (default \"output.pdf\")")
	fmt.Println("  -order string")
	fmt.Println("    \tSorting order for images: seq (sequential), nam (by name), mod (by modification time) (default \"seq\")")
	fmt.Println("  -help")
	fmt.Println("    \tShow this help message")
}

func printInputHelp() {
	fmt.Println("Usage:")
	fmt.Println("  ./img2pdf -i <directory|files> -o <pdf_file> -order <order type>")
}
