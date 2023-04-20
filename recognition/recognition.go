package recognition

import (
	"bytes"
	"fmt"
	"gb-detection/base"
	"gb-detection/model"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/Kagami/go-face"
)

const dataDir = "testdata"

var (
	modelsDir = filepath.Join(dataDir, "models")
	imagesDir = filepath.Join(dataDir, "images")
)

type Recognizer struct {
	recognizer
}

func (r *Recognizer) InitRecognizer() {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("Can`t init face recognizer: %v", err)
	}
	return &Recognizer{recognizer: rec}, nil
}

func (r *Recognizer) Close() {
	r.recognizer.Close()
}

func (r *Recognizer) RecognizeSingleFile(imagePath string) (*[128]float32, error) {
	face, err := r.recognazer.RecognizeSingleFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("Can't recognize face: %v", err)
	}
	if face == nil {
		return nil, fmt.Errorf("Can't recognize face")
	}
	result := [128]float32(face.Descriptor)
	return &result, nil
}

func (r *Recognizer) RecognizeBuf(buf io.Reader) (*[128]float32, error) {
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("Can't read bytes from file: %v", err)
	}
	face, err := r.recognazer.RecognizeSingle(bytes)
	if err != nil {
		return nil, fmt.Errorf("Can't recognize face: %v", err)
	}
	if face == nil {
		return nil, fmt.Errorf("Can't recognize face")
	}
	result := [128]float32(face.Descriptor)
	return &result, nil
}

func (r *Recognizer) Classify(db base.ImageDb, descriptor [128]float32) (model.Image, error) {
	all, err := db.AllWithDescriptor()
	if err != nil {
		return model.Image{}, fmt.Errorf("Can`t get all images: %v", err)
	}
	var descriptors []face.Descriptor
	var ids []int32
	for _, image := range all {
		imgDescriptor := image.Descriptor
		if imgDescriptor != nil {
			descriptors = append(descriptors, face.Descriptor(*imgDescriptor))
			ids = append(ids, int32(image.ID))
		}
	}
	r.recognazer.SetSamples(descriptors, ids)
	classifyResult := r.recognazer.Classify(face.Descriptor(descriptor))
	for _, image := range all {
		if image.ID == int(classifyResult) {
			return image, nil
		}
	}
	return model.Image{}, fmt.Errorf("Can`t find image with id %v", classifyResult)
}

func (r *Recognizer) JpegFormatter(buf io.Reader) (io.Reader, error) {
	jpegImage, format, err := image.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("Can't decode image: %v", err)
	}
	log.Printf("Image format: %v", format)
	bufferBytes := []byte{}
	jpegBuf := bytes.NewBuffer(bufferBytes)
	err = jpeg.Encode(jpegBuf, jpegImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, fmt.Errorf("Can't encode jpeg image: %v", err)
	}
	return jpegBuf, nil
}
