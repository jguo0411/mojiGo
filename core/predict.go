package core

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"mojiGo/utils/file"
	img "mojiGo/utils/image"
	"os"
	"time"
)

var Net = gocv.ReadNetFromONNX("core/weight/mlt_25k.onnx")

func init() {

	if Net.Empty() {
		fmt.Println("Error reading network model")
		return
	}
	Net.SetPreferableBackend(gocv.NetBackendCUDA)
	//Net.SetPreferableBackend(gocv.NetBackendOpenCV)
	//Net.SetPreferableBackend(gocv.NetBackendType(gocv.NetBackendDefault))
	Net.SetPreferableTarget(gocv.NetTargetCUDA)
}

func fit(img0 *gocv.Mat, net *gocv.Net) {

	blob := gocv.BlobFromImage(*img0, 1, image.Pt(768, 768), gocv.NewScalar(0, 0, 0, 0), false, false)
	net.SetInput(blob, "input_0")
	_ = net.Forward("output_0")
}

func valueClip(img *gocv.Mat, ty gocv.MatType, min, max int) {
	//size:=img.Size()
	for h := 0; h < 384; h++ {
		for w := 0; w < 384; w++ {
			if img.GetFloatAt(h, w) < float32(min) {
				img.SetFloatAt(h, w, float32(min))
			} else if img.GetFloatAt(h, w) > float32(max) {
				img.SetFloatAt(h, w, float32(max))
			}
		}
	}
}

func getDetBoxes(sT, sL gocv.Mat, hyp *file.Hyp) {
	txtScoreComb := gocv.NewMat()
	if hyp.AddLink == 1 {
		gocv.Threshold(sT, &sT, float32(hyp.LowText), 1, 0)
		gocv.Threshold(sL, &sL, float32(hyp.LinkThreshold), 1, 0)
		gocv.Add(sT, sL, &txtScoreComb)
		valueClip(&txtScoreComb, gocv.MatTypeCV8U, 0, 1)
	} else {
		gocv.Threshold(sT, &txtScoreComb, float32(hyp.LowText), 1, 0)
	}

	txtScoreComb.ConvertTo(&txtScoreComb, gocv.MatTypeCV8U)
	labels := gocv.NewMat()
	stats := gocv.NewMat()
	centroids := gocv.NewMat()
	nLables := gocv.ConnectedComponentsWithStatsWithParams(txtScoreComb, &labels, &stats, &centroids, 4, gocv.MatTypeCV32S, 1)

	fmt.Println("label: ", labels.Size())
	fmt.Println("stats: ", stats.Size())
	fmt.Println("centroids: ", centroids.Size())

	fmt.Println(stats.Type())

	areaMin := int32(hyp.LowAreaThreshold)
	areaMax := int32(hyp.HighAreaThreshold)

	wMax := int32(hyp.HighWidthThreshold)
	hMax := int32(hyp.HighHeightThreshold)

	var area int32
	var width int32
	var height int32
	// stats : (x, y, width, height, area)
	for k := 1; k < nLables; k++ {
		area = stats.GetIntAt(k, 4)
		if area < areaMin || area > areaMax {
			continue
		}
		if hyp.AddLink==0{
			width = stats.GetIntAt(k, 2)
			height = stats.GetIntAt(k, 3)
			if width > wMax || height > hMax{
				continue
		}
		//if np.max(textmap[labels == k]) < hyp.TextThreshold{
		//
		//}
		}
	}

	sT.MultiplyUChar(255)
	//st_temp := gocv.NewMat()
	sT.ConvertTo(&sT, gocv.MatTypeCV8UC1)
	temp := gocv.NewMat()
	gocv.ApplyColorMap(sT, &temp, gocv.ColormapJet)
	//img.ShowImg(temp)
}

func outputFile(path string, img gocv.Mat) {
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	clos := img.Cols()
	rows := img.Rows()
	a := gocv.Split(img)
	for i := 0; i < len(a); i++ {
		fmt.Fprintf(file, "==== ch %d ====\n", i)
		for h := 0; h < clos; h++ {
			for w := 0; w < rows; w++ {
				fmt.Fprintf(file, "%.5f ", a[i].GetFloatAt(h, w))
			}
			fmt.Fprintln(file)
		}
	}
}


func netFit(img0 *gocv.Mat, net *gocv.Net, hyp *file.Hyp) {
	rgbimg := gocv.NewMat()

	gocv.CvtColor(*img0, &rgbimg, gocv.ColorBGRToRGB)

	blob := gocv.BlobFromImage(rgbimg, 1, image.Pt(768, 768), gocv.NewScalar(0, 0, 0, 0), false, false)

	net.SetInput(blob, "input_0")
	prob := net.Forward("output_0")

	a, _ := prob.DataPtrFloat32()
	size := []int{384, 384}
	scoreText := gocv.NewMatWithSizes(size, gocv.MatTypeCV32FC1)
	scoreLink := gocv.NewMatWithSizes(size, gocv.MatTypeCV32FC1)

	// h,w,c => idx:= len(ch)*w+ch+(len(ch)*len(w)*h)
	for h := 0; h < 384; h++ {
		for w := 0; w < 384; w++ {
			scoreText.SetFloatAt(h, w, a[2*w+(2*h*384)])
			scoreLink.SetFloatAt(h, w, a[2*w+1+(2*h*384)])
		}
		//fmt.Println()
	}
	getDetBoxes(scoreText, scoreLink, file.HYP)
}

func PredictProgress(hyp *file.Hyp) {
	txtFilePath := hyp.ImgPath
	imgOrg := img.LoadImage(txtFilePath)
	imgBG := imgOrg.MkBg()
	fmt.Println(imgBG.Size()) // TODO d
	offsetXYs, cutImages := img.CutIntoSmall(imgOrg, hyp)

	start0 := time.Now() // 获取当前时间
	for _, img0 := range cutImages {
		start := time.Now() // 获取当前时间
		fit(&img0, &Net)
		//netFit(&img0, &Net, file.HYP)
		elapsed := time.Since(start)
		fmt.Println("cost:", elapsed)
	}
	elapsed0 := time.Since(start0)
	fmt.Println("total cost:", elapsed0)
	fmt.Println("Keep offsetYX: ", len(offsetXYs))
}

//net := cv.ReadNetFromONNX("core/weight/mlt_25k.onnx")
//net.SetPreferableBackend(cv.NetBackendDefault)
//fmt.Println(net)
