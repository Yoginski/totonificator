package face

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"

	"totonificator/bindata"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	fontNames = []string{"Roboto-Regular.ttf"}
)

type FaceMaker struct {
	Original draw.Image
	Fonts    map[string]*truetype.Font
}

func NewFaceMaker(imageBytes []byte) (*FaceMaker, error) {
	fonts, err := loadFonts()
	if err != nil {
		return nil, fmt.Errorf("failed to get font: %s", err)
	}
	original, err := png.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to load original image: %s", err)
	}
	rwImage, ok := original.(draw.Image)
	if !ok {
		return nil, fmt.Errorf("loaded image is not drawable")
	}
	return &FaceMaker{Original: rwImage, Fonts: fonts}, nil
}

func (fm *FaceMaker) Make(caption, font, color string, size float64) ([]byte, error) {
	c, dst, err := fm.getContext(font, color, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get context: %s", err)
	}
	pt := freetype.Pt(200, 200+int(c.PointToFixed(size)>>6))
	_, err = c.DrawString(strings.ToUpper(caption), pt)
	if err != nil {
		return nil, fmt.Errorf("failed to draw a caption: %s", err)
	}
	pt.Y += c.PointToFixed(size * 1.5)
	var resultBuffer bytes.Buffer
	resultWriter := bufio.NewWriter(&resultBuffer)
	err = jpeg.Encode(resultWriter, dst, &jpeg.Options{Quality: 70})
	if err != nil {
		return nil, fmt.Errorf("failed to encode an image: %s", err)
	}
	err = resultWriter.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to encode an image: %s", err)
	}

	return resultBuffer.Bytes(), nil
}

func loadFonts() (map[string]*truetype.Font, error) {
	fonts := make(map[string]*truetype.Font, len(fontNames))
	for _, fontName := range fontNames {
		fontBytes, err := bindata.Asset(fontName)
		if err != nil {
			return nil, err
		}
		font, err := freetype.ParseFont(fontBytes)
		if err != nil {
			return nil, err
		}
		fonts[fontName] = font
	}
	return fonts, nil
}

func (fm *FaceMaker) getContext(fontName, colorName string, size float64) (*freetype.Context, draw.Image, error) {
	font, ok := fm.Fonts[fontName]
	if !ok {
		return nil, nil, fmt.Errorf("font %q not found", fontName)
	}

	bounds := fm.Original.Bounds()

	dst := image.NewRGBA(bounds)
	draw.Draw(dst, dst.Bounds(), fm.Original, image.ZP, draw.Src)
	draw.Draw(
		dst,
		image.Rect(50, 150, bounds.Max.X-50, bounds.Max.Y/3+50),
		image.NewUniform(color.RGBA{R: 100, G: 100, B: 100, A: 1}),
		image.ZP,
		draw.Src,
	)

	c := freetype.NewContext()

	c.SetFont(font)
	c.SetDPI(181)
	c.SetFontSize(size)

	c.SetSrc(image.White)
	c.SetDst(dst)
	c.SetClip(bounds)

	return c, dst, nil
}
