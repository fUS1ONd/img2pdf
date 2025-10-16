package main

import (
	"errors"
	"fmt"
)

// Custom error types
var (
	ErrNoImagesFound = errors.New("no images found")
	ErrInvalidInput  = errors.New("invalid input")
)

// ImageError представляет ошибку при обработке изображения
type ImageError struct {
	Path   string
	Reason string
}

func (e *ImageError) Error() string {
	return fmt.Sprintf("image error for %q: %s", e.Path, e.Reason)
}

// InvalidExtensionError для файлов с неправильным расширением
type InvalidExtensionError struct {
	Path      string
	Extension string
}

func (e *InvalidExtensionError) Error() string {
	return fmt.Sprintf("invalid extension %q for file %q", e.Extension, e.Path)
}

// FileNotFoundError для несуществующих файлов
type FileNotFoundError struct {
	Path string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found: %q", e.Path)
}

// DirectoryError для ошибок при работе с директориями
type DirectoryError struct {
	Path   string
	Reason string
}

func (e *DirectoryError) Error() string {
	return fmt.Sprintf("directory error for %q: %s", e.Path, e.Reason)
}

// ConversionError для ошибок при конвертации в PDF
type ConversionError struct {
	Output string
	Reason string
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("conversion error for output %q: %s", e.Output, e.Reason)
}

// Helper функции для проверки типов ошибок
func IsNoImagesFound(err error) bool {
	return errors.Is(err, ErrNoImagesFound)
}

func IsInvalidExtension(err error) bool {
	var ie *InvalidExtensionError
	return errors.As(err, &ie)
}

func IsFileNotFound(err error) bool {
	var fe *FileNotFoundError
	return errors.As(err, &fe)
}

func IsDirectoryError(err error) bool {
	var de *DirectoryError
	return errors.As(err, &de)
}

func IsConversionError(err error) bool {
	var ce *ConversionError
	return errors.As(err, &ce)
}
