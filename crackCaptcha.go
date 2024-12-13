package main

import (
	"GoPassCAPTCHA/request"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func crackerCaptcha(c *request.Client) error {
	// 1. 获取验证码token
	rtk := getRTK(c)
	if rtk == "" {
		panic("rtk is null")
	}

	// 2. 获取验证码参数
	captchaData, err := getCaptchaData(c, rtk)
	if err != nil {
		return err
	}

	// 3. 下载验证码图片
	imtk := captchaData.Imtk
	bgURL := captchaData.Si
	targetURL := captchaData.Mi
	bg, target, err := fetchCaptcha(c, bgURL, targetURL, imtk)
	if err != nil {
		return err
	}

	// 4. 识别并生成参数
	movement, err := generateMovement(bg, target)
	if err != nil {
		return err
	}

	// 5. 打请求过验证码
	extend := getExtend(UserAgent)
	return fetchVerify(c, rtk, movement, extend)
}

// 获取验证码token
func getRTK(c *request.Client) string {
	params := map[string]string{
		"type":       "resource",
		"instanceId": "zfcaptchaLogin",
		"name":       "zfdun_captcha.js",
	}

	resp, err := c.Request().
		EnableTrace().
		SetQueryParams(params).Get(CaptchaURL)
	if err != nil {
		panic(err)
	}

	rtk := extractRTK(resp.String())
	return rtk
}

// 获取验证码参数
func getCaptchaData(c *request.Client, rtk string) (data *CaptchaData, err error) {
	params := map[string]string{
		"type":       "refresh",
		"rtk":        rtk,
		"time":       getTimestampStr(),
		"instanceId": "zfcaptchaLogin",
	}

	_, err = c.Request().
		EnableTrace().
		SetQueryParams(params).
		SetResult(&data).
		Get(CaptchaURL)
	if err != nil {
		return nil, err
	}
	return data, err
}

func getExtend(userAgent string) string {
	data := map[string]string{
		"appName":    "Netscape",
		"appVersion": userAgent,
		"userAgent":  userAgent,
	}
	bytes, _ := json.Marshal(data)
	//	base64
	return base64.StdEncoding.EncodeToString(bytes)

}

// 下载图片
func fetchImage(c *request.Client, url, imtk string) ([]byte, error) {
	params := map[string]string{
		"type":       "image",
		"id":         url,
		"imtk":       imtk,
		"t":          getTimestampStr(),
		"instanceId": "zfcaptchaLogin",
	}

	resp, err := c.Request().
		EnableTrace().
		SetQueryParams(params).
		Get(CaptchaURL)

	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

// 下载验证码图片
func fetchCaptcha(c *request.Client, bgURL string, targetURL string, imtk string) (bg []byte, target []byte, err error) {
	bg, err = fetchImage(c, bgURL, imtk)
	if err != nil {
		return nil, nil, err
	}

	target, err = fetchImage(c, targetURL, imtk)
	if err != nil {
		return nil, nil, err
	}
	return bg, target, nil
}

type move struct {
	X int   `json:"x"`
	Y int   `json:"y"`
	T int64 `json:"t"`
}

// 生成滑块轨迹参数
func generateMovement(bg []byte, target []byte) (string, error) {
	const steps = 10
	var movement []move
	start := getTimestamp()

	// 滑块识别
	result, err := SlideMatch(bg, target)
	if err != nil {
		return "", err
	}
	xEnd := result.Target[0]
	step := xEnd / steps

	// 生成轨迹
	for xMove := 0; xMove < xEnd; xMove += step {
		point := move{
			X: 50 + xMove,
			Y: 50,
			T: start,
		}
		movement = append(movement, point)
		start += 250
	}
	bytes, err := json.Marshal(movement)
	if err != nil {
		return "", err
	}

	// base64
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func fetchVerify(c *request.Client, rtk, movement, extend string) error {
	resp, err := c.Request().
		SetFormData(map[string]string{
			"type":       "verify",
			"rtk":        rtk,
			"time":       getTimestampStr(),
			"mt":         movement,
			"instanceId": "zfcaptchaLogin",
			"extend":     extend,
		}).
		Post(CaptchaURL)
	fmt.Println(resp.String())
	return err
}
