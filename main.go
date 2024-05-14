package main

import (
	"log"

	"github.com/cloyop/locker-desktop/desktop"
	"github.com/cloyop/locker-desktop/pkg"
)

func main() {
	if !pkg.Config() {
		log.Fatal("Missing path please set: \n'export LOCKER_PATH=path/to/lockerDir' \n'export PATH=$PATH:$LOCKER_PATH'\n")
	}
	locker := desktop.NewLocker()
	if pkg.ShouldInit() {
		locker.W.SetContent(desktop.MakeLocker(locker))
	} else {
		locker.W.SetContent(desktop.UnLockView(locker))
	}
	locker.W.ShowAndRun()
}
