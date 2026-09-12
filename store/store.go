package store

import (
	"encoding/json"
	"fmt"
	"os"
)

const fileName = "store.json"

func LoadData() ([]byte, error) {
	if data, err := os.ReadFile(fileName); err != nil {
		fmt.Println("DDDDDDD")
		return nil, err
	} else {
		return data, nil
	}
}

func SaveData[T any](data T) {
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	converted, err := json.Marshal(any(data).(T))

	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := file.Write(converted); err != nil {
		fmt.Println(err)
		return
	}
}
