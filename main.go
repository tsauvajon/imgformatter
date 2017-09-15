package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	files, err := ioutil.ReadDir("./in")

	messages := make(chan string, len(files))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(files), " files found")

	// Resize concurrently
	for _, file := range files {
		go convertAndResize(file, messages)
	}

	for range files {
		fmt.Println(<-messages)
	}

	fmt.Println()
	fmt.Println("Resizing OK")
	fmt.Println()
	fmt.Println()

	messages = make(chan string, len(files))

	// Compress concurrently
	for _, file := range files {
		go compress(file, messages)
	}

	for range files {
		fmt.Println(<-messages)
	}

	fmt.Println()
	fmt.Println("Compressing OK")

	fmt.Println()
	fmt.Println()
	fmt.Println("Done")

	fmt.Scanln()
}

func convertAndResize(file os.FileInfo, messages chan string) error {
	// open
	f, err := os.Open("./in/" + file.Name())

	if err != nil {
		fmt.Println("Error opening the file: ", err)
		return err
	}

	// decode
	t, format, err := image.Decode(f)

	if err != nil {
		fmt.Println("Error decoding the file: ", err)
		return err
	}

	if format != "jpeg" && format != "png" {
		err := errors.New("Unkown format: " + format)
		fmt.Println(err)
		return err
	}

	// resize
	resized := resize.Resize(0, 160, t, resize.Lanczos3)

	newFileName := pngFileName(file.Name())

	// create
	f, err = os.Create("./out/" + newFileName)

	if err != nil {
		log.Fatal("Error creating the file: ", err)
		return err
	}

	frs, err := os.Create("./resized/" + newFileName)

	if err != nil {
		log.Fatal("Error creating the file: ", err)
		return err
	}

	// encode
	err = png.Encode(f, t)

	if err != nil {
		log.Fatal("Error encoding the file: ", err)
		return err
	}

	err = png.Encode(frs, resized)

	if err != nil {
		log.Fatal("Error encoding the file: ", err)
		return err
	}

	messages <- file.Name()

	return nil
}

func compress(file os.FileInfo, messages chan string) error {
	newFileName := pngFileName(file.Name())
	_, err := exec.Command("cmd", "/C", "pngcrush_1_8_11_w64.exe", "./resized/"+newFileName, "./compressed/"+newFileName).Output()

	messages <- newFileName

	return err
}

func pngFileName(filename string) string {
	filename = strings.Replace(filename, ".jpg", ".png", -1)
	filename = strings.Replace(filename, ".JPG", ".png", -1)
	filename = strings.Replace(filename, ".jpeg", ".png", -1)
	filename = strings.Replace(filename, ".JPEG", ".png", -1)
	filename = strings.Replace(filename, ".PNG", ".png", -1)
	return filename
}
