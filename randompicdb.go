package ampgosetup

import (
	// "context"
	// // "crypto/rand"
	// // "encoding/hex"
	// // "encoding/json"
	// "fmt"
	// // "github.com/bogem/id3v2"
	// // "github.com/disintegration/imaging"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	// "time"
)

func get_type(afile string) string {
	if strings.Contains(afile, "thumb") {
		return "thumb"
	} else {
		return "original"
	}
}

func get_image_size(apath string) string {
	fi, err := os.Stat(apath)
	CheckError(err, "get_image_size: os.stat has failed")
	size := fi.Size()
	newsize := int(size)
	return strconv.Itoa(newsize)
}

func create_image_http_addr(aimage string) string {
	return os.Getenv("AMPGO_SERVER_ADDRESS") + ":" + os.Getenv("AMPGO_SERVER_PORT") + aimage[5:]
}

func create_image_info_map(i int, afile string, page int) Imageinfomap {
	itype := get_type(afile)
	dir, filename := filepath.Split(afile)
	image_size := get_image_size(afile)
	image_http_path := create_image_http_addr(afile)
	ii := i + 1
	var ImageInfoMap Imageinfomap
	ImageInfoMap.Dirpath = dir
	ImageInfoMap.Filename = filename
	ImageInfoMap.Imagesize = image_size
	ImageInfoMap.ImageHttpAddr = image_http_path
	ImageInfoMap.Index = strconv.Itoa(ii)
	ImageInfoMap.IType = itype
	ImageInfoMap.Page = strconv.Itoa(page)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "create_image_info_map: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "coverart", "coverart", ImageInfoMap)
	CheckError(err2, "create_image_info_map: coverart insertion has failed")
	return ImageInfoMap
}

func CreateFolderJpgImageInfoMap(afile string) {
	itype := get_type(afile)
	dir, filename := filepath.Split(afile)
	image_size := get_image_size(afile)
	image_http_path := create_image_http_addr(afile)
	var ImageInfoMap Imageinfomap
	ImageInfoMap.Dirpath = dir
	ImageInfoMap.Filename = filename
	ImageInfoMap.Imagesize = image_size
	ImageInfoMap.ImageHttpAddr = image_http_path
	ImageInfoMap.IType = itype
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "create_image_info_map: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "foldercoverart", "foldercoverart", ImageInfoMap)
	CheckError(err2, "create_image_info_map: coverart insertion has failed")
}

func CreateRandomPicsDB() []Imageinfomap {
	thumb_path := os.Getenv("AMPGO_THUMB_PATH")
	thumb_glob_path := thumb_path + "/*.jpg"
	thumb_glob, err := filepath.Glob(thumb_glob_path)
	CheckError(err, "CreateRandomPicsDB: CheckThumbDB has fucked up")
	var BulkImages []Imageinfomap
	var page int
	for i, v := range thumb_glob {
		if i < 5 {
			page = 1
		} else if i%5 == 0 {
			page++
		} else {
			page = page + 0
		}
		var iim Imageinfomap = create_image_info_map(i, v, page)
		BulkImages = append(BulkImages, iim)
	}
	return BulkImages
}
