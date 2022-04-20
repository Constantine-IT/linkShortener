package storage

import (
	"errors"
	"sync"
)

//	ErrEmptyNotAllowed - ошибка возникающая при попытке вставить пустое значение в любое поле структуры хранения URL в оперативной памяти
var ErrEmptyNotAllowed = errors.New("ram-storage: empty value is not allowed")

//	rowStorage - структура записи в хранилище URL в оперативной памяти
type rowStorage struct {
	longURL string
	userID  string
}

//	Storage - структура хранилища URL
type Storage struct {
	data  map[string]rowStorage
	mutex sync.Mutex
}

//	Констуктор хранилища URL в оперативной памяти
func NewStorage() *Storage {
	return &Storage{data: make(map[string]rowStorage)}
}

// Методы работы с хранилищем URL

// Insert - Метод для сохранения связки HASH и (<original_URL> + UserID)
func (s *Storage) Insert(shortURL, longURL, userID, filePath string) error {
	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if shortURL == "" || longURL == "" || userID == "" {
		return ErrEmptyNotAllowed
	}

	//	Блокируем структуру храниения в опративной памяти на время записи информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	сохраняем URL в оперативной памяти в структуре Storage
	//	каждая запись - это сопоставленная с HASH структура из (URL + UserID) - rowStorage
	s.data[shortURL] = rowStorage{longURL, userID}
	//	если файл для хранения URL не задан, то храним список только в оперативной памяти
	if filePath != "" {
		//	создаем экземпляр структуры хранения связки HASH<==>URL+UserID
		shortenURL := shortenURL{
			HashURL: shortURL,
			LongURL: longURL,
			UserID:  userID,
		}
		//	создаем экземпляр writer для файла-хранилища
		writtenURL, err := newWriter(filePath)
		if err != nil {
			return err
		}
		defer writtenURL.close()
		//	производим сохранение в файл связки HASH<==>URL+UserID
		if err := writtenURL.write(&shortenURL); err != nil {
			return err
		}
	}
	return nil
}

// Get - Метод для нахождения <original_URL> по HASH
func (s *Storage) Get(hash string) (longURL string, userID string, flag bool) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	если записи с запрашиваемым HASH нет в базе, то выставялем FLAG в положение FALSE
	if _, ok := s.data[hash]; !ok {
		return "", "", false
	}
	return s.data[hash].longURL, s.data[hash].userID, true
}

//	Структура записи связки HASH и <original_URL> для выдачи по запросу всех строк с одинаковым UserID
type HashURLrow struct {
	Hash    string `json:"short_url"`
	LongURL string `json:"original_url"`
}

// GetByUserID - Метод для нахождения списка сохраненных пар HASH и <original_URL> по UserID
func (s *Storage) GetByUserID(userID string) ([]HashURLrow, bool) {

	hashRows := make([]HashURLrow, 0)

	// блокируем хранилище URL на время считывания информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	отбираем строки с указанным UserID и добавляем пару HASH и <original_URL> из них в исходящий слайс
	for hash, row := range s.data {
		if row.userID == userID {
			hashRows = append(hashRows, HashURLrow{hash, row.longURL})
		}
	}

	//	если записей с таким UserID не найдено - выставляем FLAG в положение FALSE
	if len(hashRows) == 0 {
		return nil, false
	} else {
		//	если строки найдены, то возвращаем список пар HASH и <original_URL> для запрашенмого UserID
		return hashRows, true
	}
}
