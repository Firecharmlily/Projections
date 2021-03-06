/*
{-
    - Author: Liana Villafuerte, lvillafuerte2018@my.fit.edu
    - Author: Matthew Craven, mcraven2015@my.fit.edu
    - Course: CSE 4250, Fall 2019
    - Project: Proj1, Projection Please
    - Language implementations:
      - go version go1.7.4 linux/amd64
      - go version go1.12.9 windows/amd64
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

    usageStr := "Usage: program_name input_file output_file [projection_type] [std_latitude, if Lambert projection]"

    if (len(os.Args) < 3 || len(os.Args) > 5) { //outputs error if improperly called on
        fmt.Printf(usageStr) //prints message
        os.Exit(1) //leaves program
    }

    /*
    // It looks like the presence rather than the spelling of the third arg is what counts.
    if (len(os.Args) > 3 && os.Args[3] != "Lambert") { //if wording is not exact
        fmt.Printf("Unknown projection type")
        os.Exit(2) //exits program since error
    }
    */

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

        // This was by far the most useful formula on the
        // wikipedia page "Cylindrical equal-area projection"

        // For most formulas on the pages consulted, it was easier
        // to reason about the maximum and minimum values of the
        // various coordinates than to explicitly find the parameters.

        aspectRatio = math.Pi * math.Pow(math.Cos(standLat), 2)
    }


    infile, err := os.Open(filename) //learned code from site on how to turn colored images in grayscaled

    if err != nil { //if there was no file to open it would print the error
        log.Printf("failed opening %s: %s", filename, err)
        panic(err.Error())
    }
    defer infile.Close()

    imgSrc, _, err := image.Decode(infile) //decodes the image file

    if err != nil { //if it can't decode the file it stops and prints an error
        panic(err.Error())
    }

    bounds := imgSrc.Bounds() //grabs the bounds of the image
    sourceWidth, sourceHeight := bounds.Max.X, bounds.Max.Y //finds image width and height

    var width int
    height := int(sourceHeight) //height always remains the same in both programs

    if (isLambert) { //if the statement of it being lambert is true, it changes the width
        width = int(float64(height) * aspectRatio + 0.5)
    } else {
        width = sourceWidth //mollweide has the same width as the input image
    }

    image := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

    if(isLambert) {
        widthRatio := float64(sourceWidth)/float64(width)
        for y := 0; y < height; y++ {
            // + 0.5 is used to take the center of a pixel
            // flooring is used without rounding for that reason
            latitude := math.Asin(((2/float64(height))*(float64(y) + 0.5)) - 1)
            spy := (latitude/math.Pi + 0.5)*float64(sourceHeight)
            sourcePixelY := int(spy) //rounds to an int to make implementing image.Set easier

            for x := 0; x < width; x++ { //pulls the needed pixel from the source file and places it at (x,y)
                spx := (float64(x) + 0.5) * widthRatio
                sourcePixelX := int(spx)
                image.Set(x, y, imgSrc.At(sourcePixelX, sourcePixelY))
            }
        }
    } else { //defaults to Mollweide Projection
        half_width, half_height := float64(width)/2, float64(height)/2 //calculate center
        bottomx := half_width*half_width //to calculate x height of ellipse
        bottomy := half_height*half_height //to calculate y width of ellipse

        for y := 0; y < height; y++ {
            dy := float64(y) + 0.5 - half_height
            scalary := dy*dy/bottomy

            // This algorithm, using an auxilliary angle theta to
            // invert the projection map, was adapted from the
            // Wikipedia article on the "Mollweide projection"
            theta := math.Asin(dy / half_height)
            latitude := math.Asin((2*theta + math.Sin(2*theta))/math.Pi)
            spy := (latitude/math.Pi + 0.5)*float64(sourceHeight)
            sourcePixelY := int(spy) //turned to int for Set to work

            //helps set and scale the image into the ellipse for mollweide
            xOffset := 0.5 * float64(sourceWidth)
            xScale := float64(sourceWidth) / (float64(width) * math.Cos(theta)) 

            for x := 0; x < width; x++ {
                dx := float64(x) + 0.5 - half_width
                scalarx := dx*dx/bottomx

                if(scalarx + scalary >= 1) {
                    White := color.Gray{uint8(255)}
                    image.Set(x, y, White)
                } else {
                    spx := xOffset + dx * xScale //formula to calculate where to grab the pixel from
                    sourcePixelX := int(spx)
                    image.Set(x, y, imgSrc.At(sourcePixelX, sourcePixelY))
                }
            }
        }

    }


    // Encode the elipse image to the new file
    //learned also from the color to black/white image blog
    //link: https://www.golangprograms.com/how-to-convert-colorful-png-image-to-gray-scale.html
    newFileName := os.Args[2]
    newfile, err := os.Create(newFileName)
    if err != nil {
        log.Printf("failed creating %s: %s", newFileName, err)
        panic(err.Error())
    }
    defer newfile.Close()
    png.Encode(newfile,image)
}
