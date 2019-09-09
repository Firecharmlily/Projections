/*
{-
    - Author: Liana Villafuerte, lvillafuerte2018@my.fit.edu
    - Author: Matthew Craven, mcraven2015@my.fit.edu
    - Course: CSE 4250, Fall 2019
    - Project: Proj1, Projection Please
    - Language implementation: go version go1.10.4 linux/amd64
    -
}
*/

package main

import (
    "image"
    "image/color"
    "image/png"
    "log"
    "os"
    "strconv"
    "math"
    "fmt"
)


func main() {
    filename := os.Args[1] //learned from website blog on how to take inputs from command line
    isLambert := len(os.Args) >= 4 && os.Args[3] == "Lambert"

    if (len(os.Args) < 3) {
        fmt.Printf("Usage: program_name input_file output_file [projection_type] [std_latitude, if Lambert projection]")
        os.Exit(1)
    }
    if (len(os.Args) > 3 && os.Args[3] != "Lambert") { //if wording is not exact
        fmt.Printf("Unknown projection type")
        os.Exit(2) //exits program since error
    }

    var standLat float64
    var aspectRatio float64

    if (isLambert) { //In this case, need to determine standard latitude
        standLat = 0.0 //set standard latitude as default 0.0
        if(len(os.Args) == 5){ //see if there is another degree point

            stand, err := strconv.ParseFloat(os.Args[4], 64) //turns string to float

            if (err != nil) { //makes sure float is less then 50.0
                panic(err.Error()) //print error if wrong
            } else if (stand > 50.0 || stand < 0.0) {
                fmt.Printf("Error in Standard Latitude.  Out of Bounds")
                os.Exit(2) //exits program since error
            }
            standLat = stand * math.Pi / 180 //radians are better
        }
        aspectRatio = math.Pi * math.Pow(math.Cos(standLat), 2)
    }


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

    bounds := imgSrc.Bounds()
    sourceWidth, sourceHeight := bounds.Max.X, bounds.Max.Y //finds image width and height
    var width int
    height := int(sourceHeight)

    if (isLambert) {
        width = int(float64(height) * aspectRatio + 0.5)
    } else {
        width = sourceWidth
    }

    image := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

    if(isLambert) {
        for y := 0; y < height; y++ {
            // + 0.5 is used to round or to take the center of a pixel
            latitude := math.Asin(((2/float64(height))*(float64(y) + 0.5)) - 1)
            spy := (((latitude/math.Pi + 0.5)*float64(sourceHeight)) + 0.5)
            sourcePixelY := int(spy) //turned to int for Set to work
            for x := 0; x < width; x++ {
                spx := (float64(x) + 0.5) * float64(sourceWidth)/float64(width) + 0.5
                sourcePixelX := int(spx)
                //k = (float64(height) * math.Cos(standLat)) / 2
                image.Set(x, y, imgSrc.At(sourcePixelX, sourcePixelY))
            }
        }
    } else { //defaults to Mollweide Projection

        h, k := (width/2), (height/2) //calculate center
        a, b := (width - h), (height - k) //calculate radius

        for x := 0; x < width; x++ {
            for y := 0; y < height; y++ {

                topx := float64((x - h) * (x - h)) //to calculate x height of ellipse
                bottomx := float64(a*a) //to calculate x height of ellipse

                topy := float64((y - k) * (y - k)) //to calculate y width of ellipse
                bottomy := float64(b*b) //to calculate y width of ellipse

                scalarx := float64(topx/bottomx)
                scalary := float64(topy/bottomy)

                //sourcePixelX := float64(x) * scalarx
                //sourcePixelY := float64(y) * scalary

                if((scalarx + scalary) > 1){
                    White := color.Gray{uint8(255)}
                    image.Set(x, y, White)
                    //image.Set(int(sourcePixelX), int(sourcePixelY), imgSrc.At(x, y))
                }
                if((scalarx + scalary) <= 1){
                    //sourcePixelX, spy := someFunction(x,y,width,height,proj)
                    //ellipse.Set(x, y, imgSrc.At(sourcePixelX, spy))
                    image.Set(x, y, imgSrc.At(x,y))
                }
            }
        }

    }


	// Encode the elipse image to the new file
    newFileName := os.Args[2]
    newfile, err := os.Create(newFileName)
    if err != nil {
        log.Printf("failed creating %s: %s", newfile, err)
        panic(err.Error())
    }
    defer newfile.Close()
    png.Encode(newfile,image)
}
