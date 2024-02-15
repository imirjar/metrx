package mock

import (
	"encoding/json"
	"os"
)

func (m *MemStorage) Export(path string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	file.Write(data)

	return nil
}

func (m *MemStorage) Import(path string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, m); err != nil {
		return err
	}
	return nil
}
