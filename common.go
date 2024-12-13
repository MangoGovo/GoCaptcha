package main

const (
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

	LoginURL   = "http://www.gdjw.zjut.edu.cn/jwglxt/xtgl/login_slogin.html"
	CaptchaURL = "http://www.gdjw.zjut.edu.cn/jwglxt/zfcaptchaLogin"
)

type CaptchaData struct {
	Msg    string `json:"msg"`
	T      int64  `json:"t"`
	Si     string `json:"si"`
	Imtk   string `json:"imtk"`
	Mi     string `json:"mi"`
	Vs     string `json:"vs"`
	Status string `json:"status"`
}
