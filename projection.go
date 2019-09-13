/*
{-
    - Author: Matthew Craven, mcraven2015@my.fit.edu
    - Author: Liana Villafuerte, lvillafuerte2018@my.fit.edu
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
    isLambert := len(os.Args) >= 4 && os.Args[3] == "Lambert" //checks input to see if Lambert is present and correctly spelled

    if (len(os.Args) < 3) { //if input is less then needed to run
        fmt.Printf("Usage: program_name input_file output_file [projection_type] [std_latitude, if Lambert projection]")
        os.Exit(1) //exits program since error
    }
    if (len(os.Args) > 3 && os.Args[3] != "Lambert") { //if wording is not exact
        fmt.Printf("Unknown projection type")
        os.Exit(2) //exits program since error
    }

    var standLat float64 //sets variable for standard latitude
    var aspectRatio float64 //sets variable for aspect ratio needed for Lambert formulae

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

    if err != nil { //if file is not present or an error occurred
        log.Printf("failed opening %s: %s", filename, err)
        panic(err.Error())
    }
    defer infile.Close()

    imgSrc, _, err := image.Decode(infile) //also learned from color to grayscale site
    if err != nil { //if decoding file caused an error
        panic(err.Error())
    }

    bounds := imgSrc.Bounds() //sets a variable that points to the (x, y) pixel bounds
    sourceWidth, sourceHeight := bounds.Max.X, bounds.Max.Y //finds image width and height

    var width int //set variable for when width may need to be changed
    height := int(sourceHeight) //always set height as source height since that stays consistent

    if (isLambert) {
        width = int(float64(height) * aspectRatio + 0.5) //if lambert prepare to change width
    } else {
        width = sourceWidth //if mollweide keep source width
    }

    image := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}}) //new colored image is set to be coded

    if(isLambert) {
        for y := 0; y < height; y++ {
            // + 0.5 is used to round or to take the center of a pixel
            latitude := math.Asin(((2/float64(height))*(float64(y) + 0.5)) - 1) //calculates latitude for y

            spy := (latitude/math.Pi + 0.5)*float64(sourceHeight) //calculates y coordinate to pull from
            sourcePixelY := int(spy + 0.5) //math was found and translated from Lambert wikipedia

            for x := 0; x < width; x++ {

                spx := (float64(x) + 0.5) * float64(sourceWidth)/float64(width) //calculates x coordinate to pull from
                sourcePixelX := int(spx + 0.5) //mth was found and translated from Lambert wikipedia

                image.Set(x, y, imgSrc.At(sourcePixelX, sourcePixelY)) //sets image pixel loaction as the pulled pixel from input
            }
        }
    } else { //defaults to Mollweide Projection

        half_width, half_height := float64(width)/2, float64(height)/2 //calculate center

        for y := 0; y < height; y++ {
             // + 0.5 is used to round or to take the center of a pixel
            theta := math.Asin((2/float64(height)*(float64(y) + 0.5)) - 1) //calulates theta for formulae
            latitude := math.Asin((2*theta + math.Sin(2*theta))/math.Pi) //calculates latitude of y for placement

            spy := (latitude/math.Pi + 0.5)*float64(sourceHeight) //calculates y pixel to pull from
            sourcePixelY := int(spy + 0.5) //math was found and translated from mollweide wikipedia

            for x := 0; x < width; x++ {

                 // + 0.5 is used to round or to take the center of a pixel
                dx := float64(x) + 0.5 - half_width
                dy := float64(y) + 0.5 - half_height

                bottomx := float64(half_width*half_width) //to calculate x height of ellipse
                bottomy := float64(half_height*half_height) //to calculate y width of ellipse

                scalarx := float64(dx*dx/bottomx) //helps scale the ellipse radius
                scalary := float64(dy*dy/bottomy) //helps scale the ellipse radius

                if((scalarx + scalary) > 1) {
                    //sets out of bounds as white
                    White := color.Gray{uint8(255)} //found on how to use image/Color library
                    image.Set(x, y, White)
                } else {
                     // + 0.5 is used to round or to take the center of a pixel
                    spx := float64(sourceWidth) * (0.5 + dx / float64(width) / math.Cos(theta)) //calculates x pixel to pull from
                    sourcePixelX := int(spx + 0.5) //math was found and translated from mollweide wikipedia

                    image.Set(x, y, imgSrc.At(sourcePixelX, sourcePixelY)) //sets image pixel loaction as the pulled pixel from input
                }
            }
        }

    }


	// Encode the elipse image to the new file
    newFileName := os.Args[2] //name of the newfile from input command line
    newfile, err := os.Create(newFileName) //creates new file
    if err != nil { //
        log.Printf("failed creating %s: %s", newfile, err)
        panic(err.Error())
    }
    defer newfile.Close()
    png.Encode(newfile,image)
}
