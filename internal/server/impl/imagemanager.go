package impl

import (
	"image"
	_ "image/jpeg"
	"net/http"

	"github.com/pkg/errors"
)

type ImageManagerIF interface {
	DownloadImageAndCalculatePerimeter(string) (int, error)
}

type ImageManager struct{}

func NewImageManager() ImageManagerIF {
	return &ImageManager{}
}

func (iM *ImageManager) DownloadImageAndCalculatePerimeter(imageUrl string) (perimeter int, err error) {
	firstTime := true
	retry := 0
	var res *http.Response
	req, err := http.NewRequest(http.MethodGet, imageUrl, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "error creating request; url: %s", imageUrl)
	}
	if firstTime || (err != nil && retry < 15) {
		retry++
		res, err = http.DefaultClient.Do(req)
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
