package controllers

import (
	"github.com/gin-gonic/gin"
	"gocaptcha/internal/exceptions"
	"gocaptcha/internal/utils/response"
	"gocv.io/x/gocv"
	"image"
)

func Captcha(c *gin.Context) {
	bgImage, err := c.FormFile("background")
	if err != nil {
		response.AbortWithException(c, exceptions.ParamsError, err)
		return
	}

	targetImage, err := c.FormFile("target")
	if err != nil {
		response.AbortWithException(c, exceptions.ParamsError, err)
		return
	}

	bgFile, _ := bgImage.Open()
	targetFile, _ := targetImage.Open()
	var bgBytes, targetBytes []byte

	_, _ = bgFile.Read(bgBytes)
	_, _ = targetFile.Read(targetBytes)

	result, err := slideMatch(bgBytes, targetBytes)
	response.JsonSuccess(c, result)
}

// slideResult 滑块匹配结果
type slideResult struct {
	TargetX int
	TargetY int
	Target  []int
}

// slideMatch 滑块匹配
func slideMatch(backgroundBytes []byte, targetBytes []byte) (*slideResult, error) {
	// 解码目标图像
	targetMat, err := gocv.IMDecode(targetBytes, gocv.IMReadAnyColor)
	if err != nil {
		return nil, err
	}
	defer targetMat.Close()

	// 解码背景图像
	backgroundMat, err := gocv.IMDecode(backgroundBytes, gocv.IMReadAnyColor)
	if err != nil {
		return nil, err
	}
	defer backgroundMat.Close()

	// 对图像应用 Canny 边缘检测
	targetCanny := gocv.NewMat()
	defer targetCanny.Close()
	gocv.Canny(targetMat, &targetCanny, 100, 200)

	backgroundCanny := gocv.NewMat()
	defer backgroundCanny.Close()
	gocv.Canny(backgroundMat, &backgroundCanny, 100, 200)

	// 将图像转换为 RGB
	targetRGB := gocv.NewMat()
	defer targetRGB.Close()
	gocv.CvtColor(targetCanny, &targetRGB, gocv.ColorGrayToBGR)

	backgroundRGB := gocv.NewMat()
	defer backgroundRGB.Close()
	gocv.CvtColor(backgroundCanny, &backgroundRGB, gocv.ColorGrayToBGR)

	// 模板匹配
	mask := gocv.NewMat()
	defer mask.Close()
	result := gocv.NewMat()
	defer result.Close()
	gocv.MatchTemplate(backgroundRGB, targetRGB, &result, gocv.TmCcoeffNormed, mask)

	// 获取匹配结果
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	if maxVal == 0 {
		return nil, err
	}

	// 计算匹配框的右下角坐标
	h, w := targetRGB.Rows(), targetRGB.Cols()
	bottomRight := image.Point{X: maxLoc.X + w, Y: maxLoc.Y + h}

	return &slideResult{
		TargetX: 0,
		TargetY: 0,
		Target:  []int{maxLoc.X, maxLoc.Y, bottomRight.X, bottomRight.Y},
	}, nil
}
