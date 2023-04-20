package main

import (
	"gb-detection/api"
	"image"
	"image/jpeg"
	"image/png"
	"log"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

}

func main() {
	log.Println("Start...")
	api.ServerRun()
}
