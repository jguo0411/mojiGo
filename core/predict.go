package core

import (
	"fmt"
	"mojiGo/utils/file"
	"mojiGo/utils/image"
)


func PredictProgress(hyp *file.Hyp)  {
	txtFilePath := hyp.ImgPath
	imgOrg:= image.LoadImage(txtFilePath)
	imgBG := imgOrg.MkBg()
	fmt.Println(imgBG)
	image.CutIntoSmall(imgOrg, hyp)

	//TODO read img fail
}
