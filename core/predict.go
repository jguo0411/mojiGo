package core

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"mojiGo/utils/file"
	img "mojiGo/utils/image"
)

func fit(img0 gocv.Mat, net *gocv.Net) {

	blob := gocv.BlobFromImage(img0, 1.0, image.Pt(768, 768), gocv.NewScalar(0, 0, 0, 0), false, false)
	net.SetInput(blob, "input_0")
	prob := net.Forward("output_0")

	fmt.Println("Rows:", prob.Size())
}

func PredictProgress(hyp *file.Hyp) {
	txtFilePath := hyp.ImgPath
	imgOrg := img.LoadImage(txtFilePath)
	imgBG := imgOrg.MkBg()
	fmt.Println(imgBG) // TODO d
	offsetXYs, cutImages := img.CutIntoSmall(imgOrg, hyp)

	net := gocv.ReadNetFromONNX("core/weight/mlt_25k.onnx")
	if net.Empty() {
		fmt.Println("Error reading network model")
		return
	}

	net.SetPreferableBackend(gocv.NetBackendCUDA)
	//net.SetPreferableBackend(gocv.NetBackendType(gocv.NetBackendDefault))
	net.SetPreferableTarget(gocv.NetTargetCUDA)
	fmt.Println(net)

	for _, img0 := range cutImages {
		fit(img0, &net)
	}

	fmt.Println(len(offsetXYs), len(cutImages))
}

//net := cv.ReadNetFromONNX("core/weight/mlt_25k.onnx")
//net.SetPreferableBackend(cv.NetBackendDefault)
//fmt.Println(net)
