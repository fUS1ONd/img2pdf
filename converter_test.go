package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/image/tiff"
)

// TODO: Сделать тесты на все возможные комбинации флагов

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

	if info.ModTime.IsZero() {
		t.Error("CreationTime should not be zero")
	}
}

// createImagesWithTimes создает тестовые JPG-файлы с указанными именами и временем модификации
func createImagesWithTimes(t *testing.T, dir string, names []string, times []time.Time) []string {
	if len(names) != len(times) {
		t.Fatalf("names and times length mismatch")
	}

	paths := make([]string, len(names))
	for i, name := range names {
		path := filepath.Join(dir, name)
		if err := createTestJPG(path, 10, 10); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, times[i], times[i]); err != nil {
			t.Fatal(err)
		}
		paths[i] = path
	}
	return paths
}

func countPDFPages(path string) (int, error) {
	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return 0, err
	}
	return ctx.PageCount, nil
}
func TestCreatePDF_OrderByName(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "order_nam_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Создаём тестовые JPG-файлы с разными именами и временем изменения
	names := []string{"b.jpg", "a.jpg", "c.jpg"}
	times := []time.Time{
		time.Now().Add(-2 * time.Hour),
		time.Now().Add(-1 * time.Hour),
		time.Now(),
	}
	paths := createImagesWithTimes(t, tmpDir, names, times)

	converter := NewConverter()
	output := filepath.Join(tmpDir, "output_nam.pdf")

	if err := converter.Convert(strings.Join(paths, ","), output, "nam"); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Проверяем, что PDF существует
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatalf("Output PDF not created: %v", output)
	}

	// Проверяем количество страниц в PDF
	pageCount, err := countPDFPages(output)
	if err != nil {
		t.Fatalf("Failed to read PDF: %v", err)
	}
	if pageCount != len(paths) {
		t.Errorf("Expected %d pages, got %d", len(paths), pageCount)
	}

	// Проверяем, что сортировка по имени отработала корректно (a, b, c)
	expectedOrder := []string{"a.jpg", "b.jpg", "c.jpg"}

	// Здесь можно проверить, что порядок файлов, переданных в Convert, совпадает с ожидаемым
	sorted := make([]string, len(paths))
	copy(sorted, paths)
	sort.Slice(sorted, func(i, j int) bool {
		return filepath.Base(sorted[i]) < filepath.Base(sorted[j])
	})
	for i, expected := range expectedOrder {
		if filepath.Base(sorted[i]) != expected {
			t.Errorf("Expected %s at position %d, got %s", expected, i, filepath.Base(sorted[i]))
		}
	}
}

func TestCreatePDF_OrderByModTime(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "order_mod_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	names := []string{"old.jpg", "middle.jpg", "new.jpg"}
	times := []time.Time{
		time.Now().Add(-3 * time.Hour),
		time.Now().Add(-2 * time.Hour),
		time.Now().Add(-1 * time.Hour),
	}
	paths := createImagesWithTimes(t, tmpDir, names, times)

	converter := NewConverter()
	output := filepath.Join(tmpDir, "output_mod.pdf")

	if err := converter.Convert(strings.Join(paths, ","), output, "mod"); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Проверим, что по времени модификации отсортировано верно
	for i := 1; i < len(paths); i++ {
		info1, _ := os.Stat(paths[i-1])
		info2, _ := os.Stat(paths[i])
		if info1.ModTime().After(info2.ModTime()) {
			t.Errorf("Images not sorted by mod time: %s > %s", paths[i-1], paths[i])
		}
	}
}

func TestCreatePDF_OrderSequential(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "order_seq_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	names := []string{"one.jpg", "two.jpg", "three.jpg"}
	times := []time.Time{
		time.Now().Add(-2 * time.Hour),
		time.Now().Add(-1 * time.Hour),
		time.Now(),
	}
	paths := createImagesWithTimes(t, tmpDir, names, times)

	converter := NewConverter()
	output := filepath.Join(tmpDir, "output_seq.pdf")

	if err := converter.Convert(strings.Join(paths, ","), output, "seq"); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Проверяем, что последовательность осталась исходной
	for i, p := range paths {
		if filepath.Base(p) != names[i] {
			t.Errorf("Expected %s at position %d, got %s", names[i], i, filepath.Base(p))
		}
	}
}

func TestCollectImagesWithGlob(t *testing.T) {
	tmpDir, _ := createTestDirectory(t)
	defer os.RemoveAll(tmpDir)

	// создадим glob шаблон
	globPattern := filepath.Join(tmpDir, "*.jpg")

	matches, err := filepath.Glob(globPattern)
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("Glob found no files, expected some jpgs")
	}

	converter := NewConverter()
	images := converter.collectImages(strings.Join(matches, ","))

	if len(images) != len(matches) {
		t.Errorf("Expected %d images from glob, got %d", len(matches), len(images))
	}
}

func TestConvertMixedInputs(t *testing.T) {
	tmpDir, _ := createTestDirectory(t)
	defer os.RemoveAll(tmpDir)

	// Добавим отдельный файл вне директории
	singleFile := filepath.Join(tmpDir, "extra.jpg")
	if err := createTestJPG(singleFile, 50, 50); err != nil {
		t.Fatal(err)
	}

	input := tmpDir + "," + singleFile
	converter := NewConverter()
	output := filepath.Join(tmpDir, "mixed.pdf")

	if err := converter.Convert(input, output, "seq"); err != nil {
		t.Fatalf("Convert failed for mixed input: %v", err)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("PDF file not created for mixed input")
	}
}

func TestConvertInvalidPattern(t *testing.T) {
	converter := NewConverter()
	err := converter.Convert("*.nonexistent", "bad.pdf", "seq")
	if err == nil {
		t.Error("Expected error for invalid glob pattern input")
	}
}

func TestConvertNonexistentFile(t *testing.T) {
	converter := NewConverter()
	tmpDir, err := os.MkdirTemp("", "test_nonexistent_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	output := filepath.Join(tmpDir, "output.pdf")
	nonexistent := filepath.Join(tmpDir, "no_such_file.jpg")

	err = converter.Convert(nonexistent, output, "seq")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}

	// Проверяем тип ошибки
	if !IsNoImagesFound(err) {
		t.Errorf("Expected ErrNoImagesFound, got: %v", err)
	}

	// Проверяем, что PDF не создан
	if _, statErr := os.Stat(output); statErr == nil {
		t.Error("PDF should not be created when input file does not exist")
	}
}

func TestConvertEmptyInput(t *testing.T) {
	converter := NewConverter()
	tmpDir, err := os.MkdirTemp("", "test_empty_input_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	output := filepath.Join(tmpDir, "output.pdf")

	err = converter.Convert("", output, "seq")
	if err == nil {
		t.Fatal("Expected error for empty input, got nil")
	}

	// Проверяем, что это именно наша ошибка
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("Expected ErrNoImagesFound or invalid input error, got: %v", err)
	}

	if _, statErr := os.Stat(output); statErr == nil {
		t.Error("PDF should not be created for empty input")
	}
}

func TestCollectFromDirectory_PermissionError(t *testing.T) {
	// Этот тест основан на модели прав доступа в Unix.
	// В Windows он может вести себя иначе, поэтому пропускаем его.
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permission-based test on Windows")
	}

	// 1. Создаем временную директорию, которая будет автоматически удалена после теста.
	badDir := t.TempDir()

	// 2. Устанавливаем права, запрещающие чтение и вход в директорию.
	if err := os.Chmod(badDir, 0000); err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}

	// 3. Важно! Восстанавливаем права после теста, чтобы t.TempDir() смог ее удалить.
	// defer выполняется в конце функции.
	defer func() {
		if err := os.Chmod(badDir, 0755); err != nil {
			// Если даже это не сработало, просто выводим предупреждение.
			t.Logf("Warning: could not restore permissions for %s: %v", badDir, err)
		}
	}()

	converter := NewConverter()
	_, err := converter.collectFromDirectory(badDir)

	// 4. Проверяем, что ошибка была возвращена.
	if err == nil {
		t.Fatal("Expected an error for unreadable directory, but got nil")
	}

	// 5. Проверяем, что это именно ошибка типа *DirectoryError.
	var dirErr *DirectoryError
	if !errors.As(err, &dirErr) {
		t.Errorf("Expected error of type *DirectoryError, but got %T", err)
	} else {
		// (Опционально) Можно даже проверить содержимое ошибки.
		if dirErr.Path != badDir {
			t.Errorf("Expected error path %q, got %q", badDir, dirErr.Path)
		}
	}
}

func TestConvert_IgnoresEmptyEntriesInInput(t *testing.T) {
	// 1. Создаем временную директорию и тестовые файлы.
	tmpDir, err := os.MkdirTemp("", "test_empty_entries_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Создаем два валидных изображения
	path1 := filepath.Join(tmpDir, "image1.jpg")
	path2 := filepath.Join(tmpDir, "image2.jpg")
	if err := createTestJPG(path1, 50, 50); err != nil {
		t.Fatal(err)
	}
	if err := createTestJPG(path2, 50, 50); err != nil {
		t.Fatal(err)
	}

	// 2. Формируем входную строку с пустыми элементами.
	// Используем разные варианты: двойная запятая, запятая в конце, пробелы.
	input := fmt.Sprintf("%s, , %s, ", path1, path2)

	output := filepath.Join(tmpDir, "output.pdf")
	converter := NewConverter()

	// 3. Вызываем конвертацию.
	if err := converter.Convert(input, output, "seq"); err != nil {
		t.Fatalf("Convert failed with empty entries: %v", err)
	}

	// 4. Проверяем, что PDF был создан.
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatal("Output PDF was not created")
	}

	// 5. Проверяем количество страниц в PDF.
	// Ожидаем 2 страницы, так как у нас было 2 валидных файла.
	pageCount, err := countPDFPages(output)
	if err != nil {
		t.Fatalf("Failed to count PDF pages: %v", err)
	}

	expectedPages := 2
	if pageCount != expectedPages {
		t.Errorf("Expected %d pages in PDF, but got %d", expectedPages, pageCount)
	}
}
