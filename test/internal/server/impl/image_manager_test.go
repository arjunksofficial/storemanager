package impltest

import (
	"storemanager/internal/server/impl"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadImageAndCalculatePerimeter(t *testing.T) {
	iM := impl.NewImageManager()
	testCases := []struct {
		desc          string
		url           string
		wantErr       bool
		wantPerimeter int
	}{
		{
			desc:          "success",
			url:           "https://www.gstatic.com/webp/gallery/2.jpg",
			wantErr:       false,
			wantPerimeter: 1908,
		},
		{
			desc:          "fail invalid url",
			url:           "https://abc",
			wantErr:       true,
			wantPerimeter: 0,
		},
		{
			desc:          "fail unknown format",
			url:           "https://www.gstatic.com/meet/meet_logo_dark_2020q4_8955caafa87e403c96e24e8aa63f2433.svg",
			wantErr:       true,
			wantPerimeter: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			perimeter, err := iM.DownloadImageAndCalculatePerimeter(tC.url)
			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tC.wantPerimeter, perimeter)
		})
	}
}
