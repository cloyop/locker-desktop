package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type StoreData map[string]*KeyValueStore

type KeyValueStore struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *Metadata) SetToJson(fn string) (string, error) {
	t := strings.Split(strings.Split(time.Now().String(), ".")[0], " ")
	s := fmt.Sprintf("%v/locker%v_%v.json", fn, t[0], t[1])
	f, err := os.OpenFile(s, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fmt.Println("File Created")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(m.Data); err != nil {
		return "", err

	}
	return fmt.Sprintf("Data writted on\n%v", f.Name()), err
}
func (m *Metadata) SetFromJson(r io.Reader) (int, error) {
	var thing map[string]KeyValueStore
	if err := json.NewDecoder(r).Decode(&thing); err != nil {
		return 0, err
	}
	c := 0
	for name, v := range thing {
		name = strings.ToLower(name)
		if name == "" || v.Key == "" {
			continue
		}
		m.Data[name] = &KeyValueStore{Key: v.Key, Value: v.Value}
		c++
	}
	m.ChangesMade = true
	return c, nil
}
