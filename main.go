/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	srcDir := flag.String("source", "", "a string")
	outDir := flag.String("dest", "", "a string")
	width := flag.Uint("width", 0, "a int")

	flag.Parse()

	fmt.Println(fmt.Sprintf("Scan %s directory...", *srcDir))

	files, err := ioutil.ReadDir(*srcDir)
	if err != nil {
		log.Fatal(err)
	}
	nbFiles := len(files)
	var wg sync.WaitGroup
	wg.Add(nbFiles)

	for _, file := range files {
		go func(file fs.FileInfo) {
			defer wg.Done()
			f, err := os.Open(fmt.Sprintf("%v/%s", *srcDir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(file.Name())

			// decode jpeg into image.Image
			img, err := jpeg.Decode(f)
			if err != nil {
				return
			}
			f.Close()

			resized := resize.Resize(*width, 0, img, resize.Lanczos3)

			out, err := os.Create(fmt.Sprintf("%s/%s", *outDir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// write new image to file
			jpeg.Encode(out, resized, nil)
		}(file)
	}

	wg.Wait()
}
