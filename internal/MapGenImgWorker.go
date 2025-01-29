package internal

import (
	"image"
	"image/png"
	"log"
	"os"
)

func GenerateImg(val [][]MapPoint, name string) {
	m := image.NewRGBA(image.Rect(0, 0, len(val), len(val[0])))
	for i := 0; i < len(val); i++ {
		for j := 0; j < len(val[0]); j++ {
			m.Set(i, j, val[i][j].C)
		}
	}
	filename := name
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, m)
}
