package main

import (
    "image"
    "image/color"
    "image/png"
    "log"
    "os"
)


func main() {
    filename := os.Args[1] //learned from website blog on how to take inputs from command line
    infile, err := os.Open(filename) //learned code from site on how to turn colored images in grayscaled

    if err != nil {
        log.Printf("failed opening %s: %s", filename, err)
        panic(err.Error())
    }
    defer infile.Close()

    imgSrc, _, err := image.Decode(infile)
    if err != nil {
        panic(err.Error())
    }

    // Create a mollweide projection
    bounds := imgSrc.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y
    h, k := (width/2), (height/2)
    a, b := (width - h), (height - k)

    ellipse := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {

            //if(os.Args[3] == nil){
            if (true) { //temporary
                topx := float64((x - h) * (x - h))
                bottomx := float64(a*a)
                topy := float64((y - k) * (y - k))
                bottomy := float64(b*b)

                if((float64(topx/bottomx) + float64(topy/bottomy)) > 1){
                    White := color.Gray{uint8(255)}
                    ellipse.Set(x, y, White)
                }
                if((float64(topx/bottomx) + float64(topy/bottomy)) <= 1){
                    //sourcePixelX, spy := someFunction(x,y,width,height,proj)
                    //ellipse.Set(x, y, imgSrc.At(sourcePixelX, spy))
                }
            }
            /*
            lambert code along with parameters
            */
        }
	}

	// Encode the grayscale image to the new file
    newFileName := os.Args[2]
    newfile, err := os.Create(newFileName)
    if err != nil {
        log.Printf("failed creating %s: %s", newfile, err)
        panic(err.Error())
    }
    defer newfile.Close()
    png.Encode(newfile,ellipse)
}
