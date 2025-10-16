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

func (c *Converter) Convert(input, output, order string) error {
	images := c.collectImages(input)

	if len(images) == 0 {
		return fmt.Errorf("no images found")
	}

	// По последовательности переданной порядок добавления в pdf

	// -order=mod
	// -order=nam
	// -order=seq (default)

	// TODO: Реализовать порядок по имени, по дате создания с отдельными флагами
	return c.createPDF(images, output, order)
}

func (c *Converter) collectImages(input string) []ImageInfo {
	var images []ImageInfo

	files := strings.Split(input, ",")

	for _, file := range files {
		file = strings.TrimSpace(file)

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
			fmt.Printf("Warning: skipping %s: No image's extension.\n", file)
			continue
			// TODO: Ошибку сделать кастомную
			// TODO: -i afadf.jpg,,afafdadsf.jpg TEST
		}

		info, err := c.getImageInfo(file)
		if err != nil {
			fmt.Printf("Warning: skipping %s: %v\n", file, err)
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

func (c *Converter) createPDF(images []ImageInfo, output string, order string) error {

	// TODO: Двойной перебор пофиксить
	imagePaths := make([]string, len(images))
	for i, img := range images {
		imagePaths[i] = img.Path
	}

	var err error
	switch order {
	case "mod":
		sort.Slice(images, func(i, j int) bool {
			return images[i].CreationTime.Before(images[j].CreationTime)
		})

		for i, img := range images {
			imagePaths[i] = img.Path
		}

		err = api.ImportImagesFile(imagePaths, output, nil, nil)
		return err
	case "nam":
		sort.Strings(imagePaths)
		err = api.ImportImagesFile(imagePaths, output, nil, nil)
		return err
	// Default = sequently
	default:
		err = api.ImportImagesFile(imagePaths, output, nil, nil)
		return err
	}
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
