package main

import (
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	"golang.org/x/image/bmp"
	"golang.org/x/sys/windows/registry"
)

type WallpaperStyle uint

func (wps WallpaperStyle) String() string {
	return wallpaperStyles[wps]
}

const (
	Fill    WallpaperStyle = iota // 填充
	Fit                           // 适应
	Stretch                       // 拉伸
	Tile                          // 平铺
	Center                        // 居中
	Cross                         // 跨区

)

var wallpaperStyles = map[WallpaperStyle]string{
	0: "填充",
	1: "适应",
	2: "拉伸",
	3: "平铺",
	4: "居中",
	5: "跨区"}

var (
	bgFile       string
	bgStyle      int
	sFile        string
	waitTime     int
	activeScreen bool
	passwd       bool
)

var (
	regist registry.Key
)

// 设置本地图片为桌面壁纸
func setWallpaper(imgFile string) {
	var err error
	regist, err = registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.ALL_ACCESS)
	checkErr(err)
	defer regist.Close()
	style := WallpaperStyle(2)
	setDesktopWallpaper(imgFile, style)
	Println(fmt.Sprintf("Set wallpaper file and style --> %s, %s", imgFile, style))
}

func checkErr(err error) {
	if err != nil {
		Panicln(err)
	}
}

// Check that the OS is Vista or later (Vista is v6.0).
func checkVersion() bool {
	version := GetVersion()
	major := version & 0xFF
	if major < 6 {
		return false
	}
	return true
}

// jpg转换为bmp
func ConvertedWallpaper(bgfile string) string {
	file, err := os.Open(bgfile)
	checkErr(err)
	defer file.Close()

	img, err := jpeg.Decode(file) //解码
	checkErr(err)

	bmpPath := os.Getenv("USERPROFILE") + `\Local Settings\Application Data\Microsoft\Wallpaper1.bmp`
	bmpfile, err := os.Create(bmpPath)
	checkErr(err)
	defer bmpfile.Close()

	err = bmp.Encode(bmpfile, img)
	checkErr(err)
	return bmpPath
}

func setDesktopWallpaper(bgFile string, style WallpaperStyle) error {
	ext := filepath.Ext(bgFile)
	// vista 以下的系统需要转换jpg为bmp（xp、2003）
	if !checkVersion() && ext != ".bmp" {
		setRegistString("ConvertedWallpaper", bgFile)
		bgFile = ConvertedWallpaper(bgFile)
	}

	// 设置桌面背景
	setRegistString("WallPaper", bgFile)

	/* 设置壁纸风格和展开方式
	   在Control Panel\Desktop中的两个键值将被设置
	   TileWallpaper
	    0: 图片不被平铺
	    1: 被平铺
	   WallpaperStyle
	    0:  0表示图片居中，1表示平铺
	    2:  拉伸填充整个屏幕
	    6:  拉伸适应屏幕并保持高度比
	    10: 图片被调整大小裁剪适应屏幕保持纵横比
	    22: 跨区
	*/
	var bgTileWallpaper, bgWallpaperStyle string
	bgTileWallpaper = "0"
	switch style {
	case Fill: // (Windows 7 or later)
		bgWallpaperStyle = "10"
	case Fit: // (Windows 7 or later)
		bgWallpaperStyle = "6"
	case Stretch:
		bgWallpaperStyle = "2"
	case Tile:
		bgTileWallpaper = "1"
		bgWallpaperStyle = "0"
	case Center:
		bgWallpaperStyle = "0"
	case Cross: // win10 or later
		bgWallpaperStyle = "22"
	}

	setRegistString("WallpaperStyle", bgWallpaperStyle)
	setRegistString("TileWallpaper", bgTileWallpaper)

	ok := SystemParametersInfo(SPI_SETDESKWALLPAPER, FALSE, nil, SPIF_UPDATEINIFILE|SPIF_SENDWININICHANGE)
	if !ok {
		return errors.New("Desktop background settings fail")
	}
	return nil
}

func setRegistString(name, value string) {
	oldvalue, _, err := regist.GetStringValue(name)
	checkErr(err)
	if oldvalue != value {
		err = regist.SetStringValue(name, value)
		checkErr(err)
	}
}
