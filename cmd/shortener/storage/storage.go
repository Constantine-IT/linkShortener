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
	//	Блокируем структуру храниения в оперативной памяти на время записи информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	сохраняем URL в оперативной памяти в структуре Storage
	//	каждая запись - это сопоставленная с HASH структура из (URL + UserID) - rowStorage
	s.Data[hash] = RowStorage{longURL, userID, false}
	//	если файл для хранения URL не задан, то храним список только в оперативной памяти
	if fileWriter != nil {
		//	создаем экземпляр структуры хранения связки HASH<==>URL+UserID
		shortURL := shortenURL{
			HashURL:   hash,
			LongURL:   longURL,
			UserID:    userID,
			IsDeleted: true,
		}
		//	производим сохранение в файл связки HASH<==>URL+UserID
		if err := fileWriter.Write(&shortURL); err != nil {
			return err
		}
	}
	return nil
}

// Delete - метод помечает записи в базе данных, как удаленные по их HASH и UserID
func (s *Storage) Delete(hashes []string, userID string) error {
	//	Блокируем структуру храниения в оперативной памяти на время записи информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// помечаем строки с HASH из входящего среза как "удалённые", если они принадлежат пользователю с указанным UserID
	for _, hash := range hashes {
		if s.Data[hash].userID == userID {
			s.Data[hash] = RowStorage{s.Data[hash].longURL, userID, true}
			//	если файл для хранения URL не задан, то храним изменения только в оперативной памяти
			if fileWriter != nil {
				//	создаем экземпляр структуры хранения связки HASH<==>URL+UserID
				shortURL := shortenURL{
					HashURL:   hash,
					LongURL:   s.Data[hash].longURL,
					UserID:    userID,
					IsDeleted: true,
				}
				//	производим сохранение в файл, помеченной как "удалённая", связки HASH<==>URL+UserID
				if err := fileWriter.Write(&shortURL); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get - метод для нахождения <original_URL> и UserID по HASH
func (s *Storage) Get(hash string) (longURL string, flg int) {
	// блокируем хранилище URL на время считывания информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	если записи с запрашиваемым HASH нет в базе, то выставялем FLAG в положение 0
	if _, ok := s.Data[hash]; !ok {
		return "", 0
	}
	if s.Data[hash].isDeleted == true {
		return "", 2
	}
	return s.Data[hash].longURL, 1
}

// GetByLongURL - метод для нахождения HASH по <original_URL>
func (s *Storage) GetByLongURL(longURL string) (hash string, flg bool) {
	// блокируем хранилище URL на время считывания информации
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//	если записи с запрашиваемым URL нет в базе, то выставялем FLAG в положение FALSE
	for hash, rowStorage := range s.Data {
		if rowStorage.longURL == longURL && rowStorage.isDeleted == false {
			return hash, true
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
	for hash, rowStorage := range s.Data {
		if rowStorage.userID == userID && rowStorage.isDeleted == false {
			hashRows = append(hashRows, HashURLrow{hash, rowStorage.longURL})
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
	//	при остановке сервера закрываем reader и writer для файла-хранилища URL
	if err := fileReader.Close(); err != nil {
		return err
	}
	if err := fileWriter.Close(); err != nil {
		return err
	}
	return nil
}
