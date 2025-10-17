package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type ImageInfo struct {
	Path    string
	ModTime time.Time
}

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(input, output, order string) error {
	if strings.TrimSpace(input) == "" {
		return ErrInvalidInput
	}

	images := c.collectImages(input)

	if len(images) == 0 {
		return ErrNoImagesFound
	}

	return c.createPDF(images, output, order)
}

func (c *Converter) collectImages(input string) []ImageInfo {
	var images []ImageInfo

	files := strings.SplitSeq(input, ",")

	for file := range files {
		file = strings.TrimSpace(file)

		// Если -i afadf.jpg,,afafdadsf.jpg
		if file == "" {
			continue
		}

		if isDirectory(file) {
			imagesFromDir, err := c.collectFromDirectory(file)
			if err != nil {
				fmt.Printf("Warning: skipping %s: %v\n", file, err)
				continue
			} else {
				images = append(images, imagesFromDir...)
			}
			continue
		}

		if !hasImageExtension(file) {
			fmt.Printf("Warning: %v\n", &InvalidExtensionError{
				Path:      file,
				Extension: filepath.Ext(file),
			})

			continue
		}

		info, err := c.getImageInfo(file)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Warning: skipping %s: %v\n", file, &FileNotFoundError{Path: file})
			} else {
				fmt.Printf("Warning: skipping %s: %v\n", file, err)
			}
			continue
		}
		images = append(images, info)
	}

	return images
}

func (c *Converter) collectFromDirectory(dir string) ([]ImageInfo, error) {
	var images []ImageInfo

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return &DirectoryError{
				Path:   dir,
				Reason: err.Error(),
			}
		}

		if d.IsDir() || !hasImageExtension(path) {
			return nil
		}

		info, err := c.getImageInfo(path)
		if err != nil {
			fmt.Printf("Warning: skipping %s: %v\n", path, err)
			return nil
		}

		images = append(images, info)
		return nil
	})

	return images, err
}

func (c *Converter) getImageInfo(path string) (ImageInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{
		Path:    path,
		ModTime: stat.ModTime(),
	}, nil
}

func (c *Converter) createPDF(images []ImageInfo, output string, order string) error {
	switch order {
	case "mod":
		sort.Slice(images, func(i, j int) bool {
			return images[i].ModTime.Before(images[j].ModTime)
		})
	case "nam":
		sort.Slice(images, func(i, j int) bool {
			return filepath.Base(images[i].Path) < filepath.Base(images[j].Path)
		})
	}

	imagePaths := make([]string, len(images))
	for i, img := range images {
		imagePaths[i] = img.Path
	}

	if err := api.ImportImagesFile(imagePaths, output, nil, nil); err != nil {
		return &ConversionError{
			Output: output,
			Reason: err.Error(),
		}
	}
	return nil
}

// TODO: А что если дадут dir и обычные файлы? как тогда?
func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func hasImageExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".tif" || ext == ".tiff"
}
