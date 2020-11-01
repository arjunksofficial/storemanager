package impltest

import (
	"context"
	"fmt"
	"storemanager/internal/server/impl"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadImageAndCalculatePerimeter(t *testing.T) {
	iM := impl.NewImageManager()
	testCases := []struct {
		desc string
		url  string
	}{
		{
			desc: "success",
			url:  "https://www.gstatic.com/webp/gallery/2.jpg",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			perimeter, err := iM.DownloadImageAndCalculatePerimeter(ctx, tC.url)
			assert.NoError(t, err)
			fmt.Println(perimeter)
		})
	}
}
