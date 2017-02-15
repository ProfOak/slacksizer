package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide a file")
		return
	}

	filename := os.Args[1]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", filename)
	}

	slacksizer(filename)
}

func slacksizer(filename string) image.Image {
	filename = strings.ToLower(filename)

	// probably doesn't matter about output type, but why not
	if endswith(filename, ".jpg") || endswith(filename, ".jpeg") {
		slacksizer_jpg(filename)

	} else if endswith(filename, ".png") {
		slacksizer_png(filename)

	} else {
		fmt.Println("jpg/jpeg and png only")
	}
	return nil

}

func slacksizer_png(input_filename string) {
	// png helper function

	var (
		err         error
		input_img   image.Image
		output_img  image.Image
		output_file *os.File
	)

	input_img, err = png.Decode(open_image_file(input_filename))
	if err != nil {
		fmt.Printf("Not a valid png: %s\n", input_filename)
	}

	output_img = slack_size(input_img)

	output_file = create_file(input_filename)
	png.Encode(output_file, output_img)
}

func slacksizer_jpg(input_filename string) {
	// jpeg helper function

	var (
		err         error
		input_img   image.Image
		opt         jpeg.Options
		output_img  image.Image
		output_file *os.File
	)

	input_img, err = jpeg.Decode(open_image_file(input_filename))
	if err != nil {
		fmt.Printf("Not a valid jpg: %s\n", input_filename)
	}

	output_img = slack_size(input_img)

	output_file = create_file(input_filename)

	opt.Quality = 100
	jpeg.Encode(output_file, output_img, &opt)

}

func slack_size(img image.Image) image.Image {
	// slack images can only be a square of max height/width 128 pixels

	const MAX_SIZE = 128.0
	w := float64(img.Bounds().Max.X)
	h := float64(img.Bounds().Max.Y)

	percent_scale := float64(100 * min(MAX_SIZE/w, MAX_SIZE/h))

	width := uint(w * percent_scale / 100.0)
	height := uint(h * percent_scale / 100.0)

	return resize.Resize(width, height, img, resize.Bilinear)
}

func generate_filename(input_filename string) string {
	// eg: file.jpg -> file_output.jpg

	filename_parts := strings.Split(input_filename, ".")

	last_place := len(filename_parts) - 1

	extension := filename_parts[last_place]
	front_part := strings.Join(filename_parts[:last_place], "")

	return front_part + "_output." + extension
}

func create_file(filename string) *os.File {
	// totally overwrites files of the same name
	// which is unlikely to be the case

	output_filename := generate_filename(filename)
	output_file, err := os.Create(output_filename)
	if err != nil {
		fmt.Printf("Unable to create file: %s\n", output_file)
	}
	return output_file
}

func open_image_file(filename string) io.Reader {

	opened_file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to create file")
	}
	return opened_file
}

func min(x, y float64) float64 {
	// used in calculating scaling percentage

	if x < y {
		return x
	}
	return y
}

func endswith(s, substr string) bool {
	// i like python

	return strings.HasSuffix(s, substr)
}
