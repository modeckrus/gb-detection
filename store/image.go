package store

import "gb-detection/model"

type ImageDb interface {
	Add(name string)
	Update(name string)
	Get()
	File()
}

type ImageDBmem struct {
	Image model.Image `json:"image"`
}

func NewImageDBmem() *ImageDBmem {

}
