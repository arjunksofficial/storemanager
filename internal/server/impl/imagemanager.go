package impl

import (
	"context"
	"image"
	"net/http"
)

type ImageManagerIF interface {
	DownloadImageAndCalculatePerimeter(context.Context, string) (int, error)
}

type ImageManager struct{}

func NewImageManager() ImageManagerIF {
	return &ImageManager{}
}

func (iM *ImageManager) DownloadImageAndCalculatePerimeter(ctx context.Context, imageUrl string) (perimeter int, err error) {
	firstTime := true
	retry := 0
	var res *http.Response
	if firstTime || (err != nil && retry < 15) {
		retry++
		res, err = http.Get(imageUrl)
	}
	if err != nil {
		return 0, err
	}
	im, _, err := image.DecodeConfig(res.Body)
	if err != nil {
		return 0, err
	}
	return 2 * (im.Height + im.Width), nil
}
