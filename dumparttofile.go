package ampgosetup

import (
	// "context"
	// "crypto/rand"
	// "encoding/hex"
	// "encoding/json"
	"fmt"
	"github.com/bogem/id3v2"
	// "github.com/disintegration/imaging"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	"os"
	"path/filepath"
	// "strconv"
	"strings"
	// "time"
)

func folderjpg_check(apath string) Fjpg {
	fjpg := "folder.jpg"
	dir, _ := filepath.Split(apath)
	testfile := dir + fjpg
	_, error := os.Stat(testfile)
	if os.IsNotExist(error) {
		var pic Fjpg
		pic.exists = false
		pic.path = testfile
		return pic
	} else {
		var pic Fjpg
		pic.exists = true
		pic.path = testfile
		return pic
	}
}

func DumpArtToFile(apath string) (string, string, string, string, string) {
	folderjpgcheck := folderjpg_check(apath)
	tag, err := id3v2.Open(apath, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Println(err)
		fmt.Println(apath)
		return "None", "None", "None", "None", "None"
	}
	defer tag.Close()
	if folderjpgcheck.exists {
		artist := tag.Artist()
		album := tag.Album()
		title := tag.Title()
		genre := tag.Genre()
		albumart := folderjpgcheck.path
		return artist, album, title, genre, albumart
	} else {
		artist := tag.Artist()
		album := tag.Album()
		title := tag.Title()
		genre := tag.Genre()
		pictures := tag.GetFrames(tag.CommonID("Attached picture"))
		newdumpOutFile2 := ""
		newdumpOutFileThumb := ""
		for _, f := range pictures {
			pic, ok := f.(id3v2.PictureFrame)
			if !ok {
				fmt.Println("DumpArtToFile: Couldn't assert picture frame")
			}
			dumpOutFile2 := os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + ".jpg"
			newdumpOutFile2 = strings.Replace(dumpOutFile2, " ", "_", -1)
			dumpOutFileThumb := os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + "_thumb.jpg"
			newdumpOutFileThumb = strings.Replace(dumpOutFileThumb, " ", "_", -1)
			g, err := os.Create(newdumpOutFile2)
			CheckError(err, "DumpArtToFile: Unable to create newdumpOutFile2")
			defer g.Close()

			n3, err := g.Write(pic.Picture)
			CheckError(err, "DumpArtToFile: newdumpOutfile2 Write has fucked up")
			fmt.Println(n3, "DumpArtToFile: bytes written successfully")
		}
		outfile22 := resizeImage(newdumpOutFile2, newdumpOutFileThumb)
		return artist, album, title, genre, outfile22
	}
}
