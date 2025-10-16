package main

import (
	"flag"
	"fmt"
	"os"
)

// Прикол: можно наслаивать друг на друга их
// Идея: сделать так, чтобы утилита либо создавала с нуля (нужно проверять название output на предмет существования)
// либо чтобы утилита наслаивала на существующий output

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

// TODO: Добавить пояснение про order
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
