package storage

import (
	"fmt"
	"os"

	"github.com/leksusdev/calendarOfEvents/logger"
)

type JsonStorage struct {
	*Storage
}

func NewJsonStorage(filename string) *JsonStorage {
	logger.Info(fmt.Sprintf("Инициализация JSON хранилища: %s", filename))
	return &JsonStorage{
		&Storage{filename: filename},
	}
}

func (s *JsonStorage) Save(data []byte) error {
	filename := s.GetFilename()
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка сохранения в %s: %v", filename, err))
		return err
	}
	logger.Info(fmt.Sprintf("Данные успешно сохранены в %s", filename))
	return nil
}

func (s *JsonStorage) Load() ([]byte, error) {
	filename := s.GetFilename()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			logger.Error(fmt.Sprintf("Ошибка создания файла %s: %v", filename, err))
			return nil, err
		}
		err = file.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Ошибка закрытия файла %s после создания: %v", filename, err))
			return nil, err
		}
		logger.Info(fmt.Sprintf("Файл %s создан", filename))
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка загрузки из %s: %v", filename, err))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Данные успешно загружены из %s", filename))
	return data, nil
}
