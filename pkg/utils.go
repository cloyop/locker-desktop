package pkg

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
)

func RandomPassword() string {
	var characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890<>,.?//''::;;{}[]_-+=~!@#$%^&*()"
	s := []byte{}
	for range 16 {
		s = append(s, characters[rand.Intn(len(characters)-1)])
	}
	return string(s)
}
func Config() (found bool) {
	if os.Getenv("LOCKER_PATH") != "" {
		found = true
	}
	return
}
func ShouldInit() bool {
	fileStats, err := os.Stat(os.Getenv("LOCKER_PATH") + "/locker.txt")
	if err != nil {
		return true
	}
	if fileStats.IsDir() {
		if err := os.RemoveAll(os.Getenv("LOCKER_PATH") + "/locker.txt"); err != nil {
			log.Fatal(err)
		}
		return true
	}
	if fileStats.Size() < 115 {
		return true
	}
	return false
}
func ClearTerminal() {
	rn := runtime.GOOS
	var clearCMD string
	if rn == "windows" {
		clearCMD = "cls"
	} else {
		clearCMD = "clear"
	}
	cmd := exec.Command(clearCMD)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
