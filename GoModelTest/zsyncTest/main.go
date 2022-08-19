package main

import zsync "github.com/AppImageCrafters/libzsync-go"

func main() {
	sync, _ := zsync.NewZSync("https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage.zsync")
	sync.RemoteFileUrl = "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage"
}
