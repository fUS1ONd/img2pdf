package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/image/tiff"
)

// createTestJPG создает тестовое JPG изображение
func createTestJPG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Заполняем изображение градиентом
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := color.RGBA{
				R: uint8((x * 255) / width),
				G: uint8((y * 255) / height),
				B: 100,
				A: 255,
			}
			img.Set(x, y, c)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
}

// createTestImage создает тестовое изображение в указанном формате
func createTestImage(path string, width, height int, format string) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Заполняем изображение градиентом
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := color.RGBA{
				R: uint8((x * 255) / width),
				G: uint8((y * 255) / height),
				B: 100,
				A: 255,
			}
			img.Set(x, y, c)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(file, img)
	case "tiff", "tif":
		return tiff.Encode(file, img, nil)
	case "webp":
		// Для webp можно использовать базовое кодирование
		// В реальном проекте нужна библиотека для webp
		return fmt.Errorf("webp encoding not implemented in test")
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// createTestDirectory создает временную директорию с тестовыми изображениями
func createTestDirectory(t *testing.T) (string, []string) {
	tmpDir, err := os.MkdirTemp("", "jpg2pdf_test_")
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовые файлы с разным временем модификации
	files := []string{"image1.jpg", "image2.jpg", "image3.jpg"}
	
	for i, filename := range files {
		path := filepath.Join(tmpDir, filename)
		if err := createTestJPG(path, 100, 100); err != nil {
			t.Fatal(err)
		}
		
		// Устанавливаем разное время модификации для тестирования сортировки
		modTime := time.Now().Add(time.Duration(i) * time.Hour)
		if err := os.Chtimes(path, modTime, modTime); err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir, files
}

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	if converter == nil {
		t.Error("NewConverter returned nil")
	}
}

func TestHasImageExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"image.jpg", true},
		{"image.jpeg", true},
		{"image.JPG", true},
		{"image.JPEG", true},
		{"image.png", true},
		{"image.PNG", true},
		{"image.webp", true},
		{"image.WEBP", true},
		{"image.tiff", true},
		{"image.tif", true},
		{"image.TIFF", true},
		{"image.gif", false},
		{"image.bmp", false},
		{"image.txt", false},
		{"image", false},
		{"", false},
	}

	for _, test := range tests {
		result := hasImageExtension(test.filename)
		if result != test.expected {
			t.Errorf("hasImageExtension(%q) = %v; want %v", test.filename, result, test.expected)
		}
	}
}

func TestIsDirectory(t *testing.T) {
	// Тестируем с существующей директорией
	tmpDir, err := os.MkdirTemp("", "test_dir_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if !isDirectory(tmpDir) {
		t.Error("isDirectory should return true for existing directory")
	}

	// Тестируем с несуществующим путем
	if isDirectory("/nonexistent/path") {
		t.Error("isDirectory should return false for nonexistent path")
	}

	// Создаем файл и тестируем
	tmpFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if isDirectory(tmpFile) {
		t.Error("isDirectory should return false for file")
	}
}

func TestGetImageInfo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_image_info_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	converter := NewConverter()
	imagePath := filepath.Join(tmpDir, "test.jpg")
	
	if err := createTestJPG(imagePath, 50, 50); err != nil {
		t.Fatal(err)
	}

	info, err := converter.getImageInfo(imagePath)
	if err != nil {
		t.Errorf("getImageInfo failed: %v", err)
	}

	if info.Path != imagePath {
		t.Errorf("Expected path %q, got %q", imagePath, info.Path)
	}

	if info.CreationTime.IsZero() {
		t.Error("CreationTime should not be zero")
	}
}

func TestCollectFromDirectory(t *testing.T) {
	tmpDir, expectedFiles := createTestDirectory(t)
	defer os.RemoveAll(tmpDir)

	converter := NewConverter()
	images, err := converter.collectFromDirectory(tmpDir)
	
	if err != nil {
		t.Errorf("collectFromDirectory failed: %v", err)
	}

	if len(images) != len(expectedFiles) {
		t.Errorf("Expected %d images, got %d", len(expectedFiles), len(images))
	}

	// Проверяем, что изображения отсортированы по времени создания
	for i := 1; i < len(images); i++ {
		if images[i-1].CreationTime.After(images[i].CreationTime) {
			t.Error("Images are not sorted by creation time")
		}
	}
}

func TestCollectImages(t *testing.T) {
	tmpDir, _ := createTestDirectory(t)
	defer os.RemoveAll(tmpDir)

	converter := NewConverter()

	// Тест с директорией
	images, err := converter.collectImages(tmpDir)
	if err != nil {
		t.Errorf("collectImages with directory failed: %v", err)
	}
	if len(images) != 3 {
		t.Errorf("Expected 3 images from directory, got %d", len(images))
	}

	// Тест со списком файлов
	files := make([]string, 0)
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && hasImageExtension(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	filesList := strings.Join(files, ",")
	images, err = converter.collectImages(filesList)
	if err != nil {
		t.Errorf("collectImages with file list failed: %v", err)
	}
	if len(images) != 3 {
		t.Errorf("Expected 3 images from file list, got %d", len(images))
	}
}

func TestConvertEmptyInput(t *testing.T) {
	converter := NewConverter()
	
	// Тест с несуществующей директорией
	err := converter.Convert("/nonexistent", "output.pdf")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}

	// Тест с пустой директорией
	tmpDir, err := os.MkdirTemp("", "empty_dir_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = converter.Convert(tmpDir, "output.pdf")
	if err == nil {
		t.Error("Expected error for empty directory")
	}
}

func TestConvertSuccess(t *testing.T) {
	tmpDir, _ := createTestDirectory(t)
	defer os.RemoveAll(tmpDir)

	converter := NewConverter()
	outputPath := filepath.Join(tmpDir, "output.pdf")

	err := converter.Convert(tmpDir, outputPath)
	if err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	// Проверяем, что файл PDF создан
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("PDF file was not created")
	}

	// Проверяем, что файл PDF не пустой
	stat, err := os.Stat(outputPath)
	if err != nil {
		t.Error("Cannot stat PDF file")
	}
	if stat.Size() == 0 {
		t.Error("PDF file is empty")
	}
}

// Бенчмарк для тестирования производительности
func BenchmarkCreateTestImage(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	formats := []string{"jpg", "png", "tiff"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		format := formats[i%len(formats)]
		path := filepath.Join(tmpDir, fmt.Sprintf("bench_%d.%s", i, format))
		if err := createTestImage(path, 200, 200, format); err != nil {
			// Пропускаем неподдерживаемые форматы
			continue
		}
		os.Remove(path)
	}
}
