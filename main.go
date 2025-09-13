package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		input  = flag.String("input", "", "Input directory or JPG files (comma-separated (,))")
		output = flag.String("output", "output.pdf", "Output PDF file path")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if (*help) || (*input == "") {
		printUsage()
		return
	}

	converter := NewConverter()
	if err := converter.Convert(*input, *output); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted to %s\n", *output)
}

func printUsage() {
	fmt.Println("Image to PDF Converter")
	fmt.Println("\nSupported formats: JPG, JPEG, PNG, WEBP, TIFF")
	fmt.Println("\nUsage:")
	fmt.Println("  jpg2pdf -input <directory|files> -output <pdf_file>")
	fmt.Println("\nExamples:")
	fmt.Println("  jpg2pdf -input ./images -output result.pdf")
	fmt.Println("  jpg2pdf -input \"image1.jpg,photo.png,scan.tiff\" -output result.pdf")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}
