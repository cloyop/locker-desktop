package storage

import (
	"bytes"
	"encoding/gob"
	"os"
	"time"

	"github.com/cloyop/locker-desktop/pkg"
)

type Metadata struct {
	LastAction    *time.Timer
	ChangesMade   bool
	Password, Pin string
	Data          StoreData
}

func NewMetaData() *Metadata {
	return &Metadata{
		Data: StoreData{},
	}
}
func (m *Metadata) Save() error {
	layerOneBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(layerOneBuffer).Encode(&m.Data); err != nil {
		return err
	}
	bytes := layerOneBuffer.Bytes()
	SafeStoreData, err := pkg.Cipher([]byte(m.Password+m.Pin), &bytes)
	if err != nil {
		return err
	}
	SafeToWrite, err := pkg.Cipher([]byte(m.Password), &SafeStoreData)
	if err != nil {
		return err
	}
	if err := os.WriteFile(os.Getenv("LOCKER_PATH")+"/locker.txt", SafeToWrite, os.ModePerm); err != nil {
		return err
	}
	return nil
}
