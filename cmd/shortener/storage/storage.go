package storage

import (
	"sync"
)

//	Storage - структура хранилища URL для работы в оперативной памяти
type Storage struct {
	Data  map[string]RowStorage
	mutex sync.Mutex
}

// Insert - Метод для сохранения связки HASH и (<original_URL> + UserID)
func (s *Storage) Insert(hash, longURL, userID string) error {

	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if hash == "" || longURL == "" || userID == "" {
		return ErrEmptyNotAllowed
	}

	//	Блокируем структуру храниения в опративной памяти на время записи информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	сохраняем URL в оперативной памяти в структуре Storage
	//	каждая запись - это сопоставленная с HASH структура из (URL + UserID) - rowStorage
	s.Data[hash] = RowStorage{longURL, userID}
	//	если файл для хранения URL не задан, то храним список только в оперативной памяти
	if fileWriter != nil {
		//	создаем экземпляр структуры хранения связки HASH<==>URL+UserID
		shortURL := shortenURL{
			HashURL: hash,
			LongURL: longURL,
			UserID:  userID,
		}
		//	производим сохранение в файл связки HASH<==>URL+UserID
		if err := fileWriter.Write(&shortURL); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByHashes - метод для пометки записей в базе данных как удаленные по их HASH и UserID
func (s *Storage) DeleteByHashes(hash []string, userID string) error {
	return nil
}

// Get - Метод для нахождения <original_URL> и UserID по HASH
func (s *Storage) Get(hash string) (longURL string, userID string, flg bool) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	если записи с запрашиваемым HASH нет в базе, то выставялем FLAG в положение FALSE
	if _, ok := s.Data[hash]; !ok {
		return "", "", false
	}
	return s.Data[hash].longURL, s.Data[hash].userID, true
}

// GetByLongURL - Метод для нахождения HASH по <original_URL>
func (s *Storage) GetByLongURL(longURL string) (hash string, flg bool) {

	// блокируем хранилище URL на время считывания информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	если записи с запрашиваемым URL нет в базе, то выставялем FLAG в положение FALSE
	for h, row := range s.Data {
		if row.longURL == longURL {
			return h, true
		}
	}
	return "", false
}

// GetByUserID - Метод для нахождения списка сохраненных пар HASH и <original_URL> по UserID
func (s *Storage) GetByUserID(userID string) ([]HashURLrow, bool) {

	hashRows := make([]HashURLrow, 0)

	// блокируем хранилище URL на время считывания информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	отбираем строки с указанным UserID и добавляем пару HASH и <original_URL> из них в исходящий слайс
	for hash, row := range s.Data {
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

func (s *Storage) Close() error {
	var err error

	//	при остановке сервера отложенно закрываем reader и writer для файла-хранилища URL
	err = fileReader.Close()
	if err != nil {
		return err
	}

	err = fileWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
