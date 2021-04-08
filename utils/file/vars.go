package file

import (
	"encoding/json"
	"fmt"
	"log"
	"mojiGo/utils/sys"
	"os"
	"strconv"
	"strings"
)

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func (e *Err) SaveErrTxt(h *Hyp) {
	txtFile := h.TxtFilePath
	if Exist(txtFile) {
		f, err := os.OpenFile(txtFile, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		wStr := fmt.Sprintf("\r\n\r\n### ERROR ###\r\n%s\r\n", e.Msg)
		if _, err = f.WriteString(wStr); err != nil {
			panic(err)
		}
		_ = f.Close()
		sys.SignalWriter(sys.ERROR)
		log.Println(e.Msg)

		newName := strings.Replace(txtFile, ".txt", "_ERROR.txt", -1)
		_ = os.Rename(txtFile, newName)
	}
}

func NewError(code int, msg string) *Err {
	return &Err{
		Code: code,
		Msg:  msg,
	}
}

type Hyp struct {
	CutSize, Pad uint64

	TxtFilePath   string
	ImgPath       string
	TextThreshold float64
	LowText       float64
	LinkThreshold float64
	AddLink       uint64

	Connectivity        uint64
	LowAreaThreshold    uint64
	HighAreaThreshold   uint64
	HighWidthThreshold  uint64
	HighHeightThreshold uint64

	//hidden parameters
	Poly       bool
	MagRatio   float64
	CanvasSize uint64
}

func (h *Hyp) Init() {
	h.CutSize, h.Pad = 768, 30
	h.TxtFilePath = ""

	h.ImgPath = ""
	h.TextThreshold = 0.55
	h.LowText = 0.6
	h.LinkThreshold = 0.35
	h.AddLink = 0

	h.Connectivity = 4
	h.LowAreaThreshold = 10
	h.HighAreaThreshold = 10000
	h.HighWidthThreshold = 100
	h.HighHeightThreshold = 100

	//hidden parameters
	h.Poly = false
	h.MagRatio = 1.5
	h.CanvasSize = 768
}

func (h *Hyp) LoadCfg(txtPath string) error {
	h.Init()
	h.TxtFilePath = txtPath
	tempMap, err := ReadConfig(txtPath)
	if err != nil {
		err.SaveErrTxt(h)
	}
	imgPath, err := win2mnt(tempMap["parameters"]["img_path"])
	if err != nil {
		err.SaveErrTxt(h)
		return nil
	}
	h.ImgPath = imgPath
	h.TextThreshold, _ = strconv.ParseFloat(tempMap["parameters"]["text_threshold"], 8)
	h.LowText, _ = strconv.ParseFloat(tempMap["parameters"]["low_text"], 8)
	h.LinkThreshold, _ = strconv.ParseFloat(tempMap["parameters"]["link_threshold"], 8)
	h.AddLink, _ = strconv.ParseUint(tempMap["parameters"]["link_threshold"], 10, 4)
	h.Connectivity, _ = strconv.ParseUint(tempMap["parameters"]["connectivity"], 10, 8)
	h.LowAreaThreshold, _ = strconv.ParseUint(tempMap["parameters"]["low_area_threshold"], 10, 8)
	h.HighAreaThreshold, _ = strconv.ParseUint(tempMap["parameters"]["high_area_threshold"], 10, 8)
	h.HighWidthThreshold, _ = strconv.ParseUint(tempMap["parameters"]["high_width_threshold"], 10, 8)
	h.HighHeightThreshold, _ = strconv.ParseUint(tempMap["parameters"]["high_height_threshold"], 10, 8)

	return nil
}
