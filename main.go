package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func main() {

	var (
		heightPtr, widthPtr         *int
		cropHeightPtr, cropWidthPtr *int
	)

	heightPtr = flag.Int("height", 500, "maximum height of the resized picture")
	widthPtr = flag.Int("width", 320, "maximum width of the resized picture")

	cropHeightPtr = flag.Int("cropheight", 0, "cropped height of the picture")
	cropWidthPtr = flag.Int("cropwidth", 0, "cropped width of the picture")

	flag.Parse()

	files, err := ioutil.ReadDir("./in")

	messages := make(chan string, len(files))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(files), " files found")

	// Resize concurrently
	for _, file := range files {
		go convertAndResize(file, uint(*widthPtr), uint(*heightPtr), *cropWidthPtr, *cropHeightPtr, messages)
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

func convertAndResize(file os.FileInfo, width, height uint, cropWidth, cropHeight int, messages chan string) error {
	var (
		frs *os.File
	)

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
	// resized := resize.Resize(width, height, t, resize.Lanczos3)
	resized := resize.Thumbnail(width, height, t, resize.Lanczos3)

	newFileName := pngFileName(file.Name())

	var croppedImg image.Image

	// crop
	shouldCropImage := cropWidth > 0 && cropHeight > 0

	if shouldCropImage {
		croppedImg, err = cutter.Crop(resized, cutter.Config{
			Width:  cropWidth,
			Height: cropHeight,
			Mode:   cutter.Centered,
		})
	}

	// create
	f, err = os.Create("./out/" + newFileName)

	if err != nil {
		log.Fatal("Error creating the file: ", err)
		return err
	}

	frs, err = os.Create("./resized/" + newFileName)

	if err != nil {
		log.Fatal("Error creating the file: ", err)
		return err
	}

	// encode
	err = jpeg.Encode(f, t, nil)

	if err != nil {
		log.Fatal("Error encoding the file: ", err)
		return err
	}

	if shouldCropImage {
		err = jpeg.Encode(frs, croppedImg, nil)
	} else {
		err = jpeg.Encode(frs, resized, nil)
	}

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
	filename = strings.Replace(filename, ".jpg", ".jpg", -1)
	filename = strings.Replace(filename, ".JPG", ".jpg", -1)
	filename = strings.Replace(filename, ".jpeg", ".jpg", -1)
	filename = strings.Replace(filename, ".JPEG", ".jpg", -1)
	filename = strings.Replace(filename, ".PNG", ".jpg", -1)
	return filename
}
