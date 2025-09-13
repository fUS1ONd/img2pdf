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
	Path         string
	CreationTime time.Time
}

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(input, output string) error {
	images, err := c.collectImages(input)
	if err != nil {
		return fmt.Errorf("failed to collect images: %w", err)
	}

	if len(images) == 0 {
		return fmt.Errorf("no images found")
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].CreationTime.Before(images[j].CreationTime)
	})

	return c.createPDF(images, output)
}

func (c *Converter) collectImages(input string) ([]ImageInfo, error) {
	var images []ImageInfo

	if isDirectory(input) {
		return c.collectFromDirectory(input)
	}

	files := strings.Split(input, ",")
	for _, file := range files {
		file = strings.TrimSpace(file)
		if !hasImageExtension(file) {
			continue
		}

		info, err := c.getImageInfo(file)
		if err != nil {
			fmt.Printf("Warning: skipping %s: %v\n", file, err)
			continue
		}
		images = append(images, info)
	}

	return images, nil
}

func (c *Converter) collectFromDirectory(dir string) ([]ImageInfo, error) {
	var images []ImageInfo

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
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
		Path:         path,
		CreationTime: stat.ModTime(),
	}, nil
}

func (c *Converter) createPDF(images []ImageInfo, output string) error {
	if len(images) == 0 {
		return fmt.Errorf("no images to convert")
	}

	imagePaths := make([]string, len(images))
	for i, img := range images {
		imagePaths[i] = img.Path
	}

	return api.ImportImagesFile(imagePaths, output, nil, nil)
}

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
