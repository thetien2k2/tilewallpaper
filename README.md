# tilewallpaper
Simple program to create wallpaper tiled by input image.
Input image will be tile from right.
In case heigh of image smaller than wallpaper, gaussian blur background will be created.

# build requirements:
- go
- libvips

# usage:
set your screen width and heigh in main.go
go mod tidy
go build
./tilewallpaper your_input_image