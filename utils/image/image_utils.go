package image

import (
	cv "gocv.io/x/gocv"
	"image"
	"math"
	"mojiGo/utils/file"
)

type OrgImage struct {
	h, w, c int
	m       cv.Mat
}

func ShowImg(img cv.Mat) {
	window := cv.NewWindow("test")
	for {
		window.IMShow(img)
		window.WaitKey(0)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func (o *OrgImage) MkBg() cv.Mat {
	bg := cv.NewMatWithSize(o.h, o.w, cv.MatTypeCV8UC3)
	return bg
}

func LoadImage(imgFile string) *OrgImage {
	img := cv.IMRead(imgFile, cv.IMReadColor)
	if img.Empty() {
		//TODO failed to load image
	}
	return &OrgImage{
		h: img.Rows(),
		w: img.Cols(),
		c: img.Channels(),
		m: img,
	}
}

//getRange return image content boundary.gocv 0.27.0
func getRange(ref cv.Mat, orgImage *OrgImage) image.Rectangle {
	refCntsPointer := cv.FindContours(ref, cv.RetrievalExternal, cv.ChainApproxSimple)
	defer refCntsPointer.Close()

	var mRetc image.Rectangle
	if refCntsPointer.Size() < 2 {
		mRetc.Min.X = 0
		mRetc.Min.Y = 0
		mRetc.Max.X = orgImage.w
		mRetc.Max.Y = orgImage.h
		return mRetc
	}
	tmp0 := cv.BoundingRect(refCntsPointer.At(0))
	mRetc.Min.X = tmp0.Min.X
	mRetc.Min.Y = tmp0.Min.X
	mRetc.Max.X = tmp0.Max.X
	mRetc.Max.Y = tmp0.Max.Y

	for i := 1; i < refCntsPointer.Size(); i++ {
		tmp := cv.BoundingRect(refCntsPointer.At(i))
		mRetc.Min.X = int(math.Min(float64(mRetc.Min.X), float64(tmp.Min.X)))
		mRetc.Min.Y = int(math.Min(float64(mRetc.Min.Y), float64(tmp.Min.Y)))
		mRetc.Max.X = int(math.Max(float64(mRetc.Max.X), float64(tmp.Max.X)))
		mRetc.Max.Y = int(math.Max(float64(mRetc.Max.Y), float64(tmp.Max.Y)))
	}
	return mRetc
}

func rangePre(orgImage *OrgImage) cv.Mat {
	ref := cv.NewMat()
	cv.CvtColor(orgImage.m, &ref, cv.ColorBGRToGray)
	cv.Threshold(ref, &ref, 0, 255, cv.ThresholdOtsu)

	//cv.IMWrite("/home/dac/Desktop/txt.png", ref)

	// invert
	tmpMat := ref.Clone()
	defer tmpMat.Close()
	tmpMat.AddUChar(255)
	cv.Subtract(tmpMat, ref, &ref)

	return ref
}

//func NormalizeMeanVariance(img *cv.Mat, mean, variance [3]float32){
func NormalizeMeanVariance(img *cv.Mat) (dst cv.Mat) {

	//mean := [3]float32{0.485*255, 0.456*255, 0.406*255}
	//variance := [3]float32{0.229*255, 0.224*255, 0.225*255}
	mean := [3]float32{123.675, 116.28, 103.53}
	variance := [3]float32{58.395, 57.12, 57.375}

	dst = cv.NewMat()

	a := cv.Split(*img)
	for i := 0; i < 3; i++ {
		// m(x,y)=saturate_cast<rType>(α(∗this)(x,y)+β)
		a[i].ConvertToWithParams(&a[i], cv.MatTypeCV32FC1, 1.0/variance[i], (0.0-mean[i])/variance[i])
		//defer a[i].Close()
	}
	cv.Merge(a, &dst)
	return
}

func CutIntoSmall(orgImg *OrgImage, hyp *file.Hyp) ([]image.Point, []cv.Mat) {

	//var corpBox image.Rectangle
	baseSize := hyp.CutSize - hyp.Pad
	// inverted 2value picture
	ref := rangePre(orgImg)
	mRetc := getRange(ref, orgImg)

	offsetX0 := mRetc.Min.X
	offsetY0 := mRetc.Min.Y

	cropImg := orgImg.m.Region(mRetc)

	w, h := cropImg.Cols(), cropImg.Rows()
	h1 := int(math.Ceil(float64(h)/float64(baseSize)) * float64(baseSize))
	w1 := int(math.Ceil(float64(w)/float64(baseSize)) * float64(baseSize))

	cropImgRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
	cropBG := cv.NewMatWithSize(h1, w1, cv.MatTypeCV8UC3)

	cropImaBG := cropBG.Region(cropImgRect)
	cropImg.CopyTo(&cropImaBG)
	_ = cropImg.Close()
	_ = cropImaBG.Close()

	var offsetXYs []image.Point
	var cutImages []cv.Mat

	for r := 0; r < (h1 / baseSize); r++ {
		for l := 0; l < (w1 / baseSize); l++ {
			cropRectRL := image.Rectangle{
				Min: image.Point{
					X: l * baseSize,
					Y: r * baseSize,
				},
				Max: image.Point{
					X: (l + 1) * baseSize,
					Y: (r + 1) * baseSize,
				},
			}
			offsetXY := image.Point{
				X: cropRectRL.Min.X + offsetX0,
				Y: cropRectRL.Min.Y + offsetY0,
			}
			offsetXYs = append(offsetXYs, offsetXY)

			cropRLImg := cropBG.Region(cropRectRL)
			//savename := fmt.Sprintf("/home/dac/Desktop/192.168.0.5/%d_%dgo.png", r, l)
			//cv.IMWrite(savename, cropRLImg)
			croped := NormalizeMeanVariance(&cropRLImg)

			_ = cropRLImg.Close()
			cutImages = append(cutImages, croped.Clone())
			_ = croped.Close()
		}
	}
	return offsetXYs, cutImages
}