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
    
    image := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
    
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {

            if( (len(os.Args) >= 4) && (os.Args[3] == "Lambert")) { //see if Lambert is present and correctly spelled

                var standLat float64
                standLat = 0.0 //set standard latitude as default 0.0

                if(len(os.Args) == 5){ //see if there is another degree point
                    standLat, err := strconv.ParseFloat(os.Args[4], 64) //turns string to float
                    if ((err != nil) || standLat > 50.0) { //makes sure float is less then 50.0
                        panic(err.Error()) //print error if wrong
                    }
                }

                var spy float64

                //x is longitude - prime meridian aka x - 0
                spy = float64(y) - math.Sin(standLat) //y = sin(latitude) 
                sourcePixelY := int(spy) //turned to int for Set to work
                image.Set(x, sourcePixelY, imgSrc.At(x, y))

            } else if ((len(os.Args) != 3) && os.Args[3] != "Lambert"){ //if wording is not exact

                fmt.Printf("Error in projection type")
                os.Exit(2) //exits program since error

            } else { //defaults to Mollweide Projection

                h, k := (width/2), (height/2)
                a, b := (width - h), (height - k)
                topx := float64((x - h) * (x - h))
                bottomx := float64(a*a)
                topy := float64((y - k) * (y - k))
                bottomy := float64(b*b)

                if((float64(topx/bottomx) + float64(topy/bottomy)) > 1){
                    White := color.Gray{uint8(255)}
                    image.Set(x, y, White)
                }
                if((float64(topx/bottomx) + float64(topy/bottomy)) <= 1){
                    //sourcePixelX, spy := someFunction(x,y,width,height,proj)
                    //ellipse.Set(x, y, imgSrc.At(sourcePixelX, spy))
                    image.Set(x, y, imgSrc.At(x, y))
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