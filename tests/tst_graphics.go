// This test file is a giant soup that had one goal: act as a testing ground
// to get various graphical routines working in X. Namely:
//
// Blending one image on top of another (with alpha).
// Painting an arbitrary image into a window.
// Drawing text on to an image.
//
// All of this is done here successfully. Most of this file will be split up
// into nicer pieces in my window manager.
package main

import (
    "fmt"
    // "image" 
    "image/color"
    // "image/draw" 
    "time"

    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/ewmh"
    "github.com/BurntSushi/xgbutil/xgraphics"
)

var X *xgbutil.XUtil
var Xerr error

func Recovery() {
    if r := recover(); r != nil {
        fmt.Println("ERROR:", r)
        // os.Exit(1) 
    }
}

var fontFile string = "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf"

func main() {
    defer Recovery()

    X, Xerr = xgbutil.Dial("")
    if Xerr != nil {
        panic(Xerr)
    }

    active, _ := ewmh.ActiveWindowGet(X)
    icons, _ := ewmh.WmIconGet(X, active)
    fmt.Printf("Active window's (%x) icon data: (length: %v)\n", 
               active, len(icons))
    for _, icon := range icons {
        fmt.Printf("\t(%d, %d)", icon.Width, icon.Height)
        fmt.Printf(" :: %d == %d\n", icon.Width * icon.Height, len(icon.Data))
    }

    work := icons[2]
    fmt.Printf("Working with (%d, %d)\n", work.Width, work.Height)


    img, mask := xgraphics.EwmhIconToImage(work)

    dest := xgraphics.BlendBg(img, mask, 70, color.RGBA{0, 0, 255, 255})

    // Let's try to write some text...
    xgraphics.DrawText(dest, 5, 5, color.RGBA{255, 255, 255, 255}, 10,
                       fontFile, "Hello, world!")

    tw, th, err := xgraphics.TextExtents(fontFile, 11, "Hiya")
    fmt.Println(tw, th, err)

    win := xgraphics.CreateImageWindow(X, dest, 3940, 400)
    X.Conn().MapWindow(win)

    time.Sleep(20 * time.Second)
}

