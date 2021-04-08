package image

import (
	"fmt"
	cv "gocv.io/x/gocv"
	"image"
	"mojiGo/utils/file"
)

type OrgImage struct {
	h, w, c int
	m       cv.Mat
}

func (o *OrgImage) MkBg() cv.Mat {
	bg := cv.NewMatWithSize(o.h, o.w, cv.MatTypeCV8UC3)
	return bg
}

func LoadImage(imgFile string) *OrgImage {
	img := cv.IMRead(imgFile, cv.IMReadColor)
	return &OrgImage{
		h: img.Rows(),
		w: img.Cols(),
		c: img.Channels(),
		m: img,
	}
}

func CutIntoSmall(orImg *OrgImage, hyp *file.Hyp) {

	var boundingBoxes []image.Rectangle
	//var corpBox image.Rectangle

	//baseSize := hyp.CutSize - hyp.Pad
	ref := cv.NewMat()
	cv.CvtColor(orImg.m, &ref, cv.ColorBGRToGray)
	cv.Threshold(ref, &ref, 0, 255, cv.ThresholdOtsu)

	//cv.IMWrite("/home/dac/Desktop/txt.png", ref)

	refCnts := cv.FindContours(ref, cv.RetrievalExternal, cv.ChainApproxSimple)
	//net := cv.ReadNetFromONNX("core/weight/mlt_25k.onnx")
	//net.SetPreferableBackend(cv.NetBackendDefault)
	//fmt.Println(net)

	fmt.Println(len(refCnts), len(refCnts[0]))

	for i := 1; i < len(refCnts); i++ {
		tmp := cv.BoundingRect(refCnts[i])
		boundingBoxes = append(boundingBoxes, tmp)
	}

	for _, v := range boundingBoxes {
		fmt.Println(v.Min.X, v.Min.Y, v.Max.X, v.Max.Y)
	}
}
