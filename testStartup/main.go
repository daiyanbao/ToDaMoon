package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

func main() {
	content := time.Now().String()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	content += "\n" + wd + "\n"

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	content += fmt.Sprintln("Name     :", usr.Name)
	content += fmt.Sprintln("username :", usr.Username)
	content += fmt.Sprintln("home     :", usr.HomeDir)
	content += fmt.Sprintln("UID      :", usr.Uid)
	content += fmt.Sprintln("GID      :", usr.Gid)

	err = ioutil.WriteFile("test.txt", []byte(content), 0667)
	if err != nil {
		log.Fatal(err)
	}
}
