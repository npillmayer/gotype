package pngcanvas

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

/* Output routine to produce PNG format files.
 */
type PNGOutputRoutine struct {
	filenamePattern string
	serial          int
}

/* Create a new PNG output routine.
 */
func NewPNGOutputRoutine() *PNGOutputRoutine {
	return &PNGOutputRoutine{
		filenamePattern: "%d",
		serial:          1,
	}
}

/* Ship out an image to file. The filename is chosen from the image name.
 * If the image name is empty, this routine will use a filename pattern
 * and a serial number for each image to output.
 */
func (pngoutp *PNGOutputRoutine) Shipout(name string, image image.Image) {
	if name == "" {
		name = fmt.Sprintf(pngoutp.filenamePattern, pngoutp.serial)
		pngoutp.serial++
	}
	name = name + ".png"
	f, err := os.Create(name)
	if err != nil {
		panic(fmt.Sprintf("cannot open output file: %s", name))
	}
	defer f.Close()
	err = png.Encode(f, image)
}
