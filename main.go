package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		input  = flag.String("i", "", "Input directory or JPG files (space-separated)")
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

func printInputHelp() {
	fmt.Println("Image to PDF Converter")

	fmt.Println("\nUsage:")
	fmt.Println("  ./jpg2pdf -i \"image1.jpg,photo.png,scan.tiff\" -o result.pdf")
	fmt.Println("  ./jpg2pdf -i ./images -o result.pdf")
}
