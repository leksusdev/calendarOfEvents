package storage

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/leksusdev/calendarOfEvents/config"
)

type ZipStorage struct {
	*Storage
}

func NewZipStorage(filename string) *ZipStorage {
	return &ZipStorage{
		&Storage{filename: filename},
	}
}

func (z *ZipStorage) Save(data []byte) error {
	f, err := os.Create(z.GetFilename())
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	zw := zip.NewWriter(f)
	defer func(zw *zip.Writer) {
		err := zw.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(zw)

	w, err := zw.Create(config.ZipLogEntryName)
	if err != nil {
		return fmt.Errorf("ошибка создания файла в архиве: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("ошибка записи данных в архив: %w", err)
	}

	return nil
}

func (z *ZipStorage) Load() ([]byte, error) {
	r, err := zip.OpenReader(z.GetFilename())
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия архива: %w", err)
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	if len(r.File) == 0 {
		return nil, errors.New("архив пуст")
	}

	file := r.File[0]
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла в архиве: %w", err)
	}
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rc)

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения содержимого архива: %w", err)
	}
	return data, nil
}
