///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
// LICENSE: GNU General Public License, version 2 (GPLv2)
// Copyright 2016, Charlie J. Smotherman
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License v2
// as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

package ampgosetup

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"strconv"
// 	"path/filepath"
// 	"github.com/globalsign/mgo/bson"
// )

// type Imageinfomap struct {
// 	ID bson.ObjectId `bson:"_id,omitempty"`
// 	Dirpath   string `bson:"dirpath"`
// 	Filename  string `bson:"filename"`
// 	Imagesize string `bson:"imagesize"`
// 	ImageHttpAddr string `bson:"imagehttpaddr`
// 	Index string `bson:"index"`
// 	IType string `bson:"itype"`
// 	Page string `bson:"page"`
// }

// //RanPics exported
// func CreateRandomPicsDB() []Imageinfomap {
// 	thumb_path := os.Getenv("AMPGO_THUMB_PATH")
// 	thumb_glob_path := thumb_path + "/*.jpg"
// 	thumb_glob, err := filepath.Glob(thumb_glob_path)
// 	if err != nil {
// 		fmt.Println("CheckThumbDB has fucked up")
// 	}
// 	var BulkImages []Imageinfomap
// 	var page int
// 	for i, v := range thumb_glob {
// 		if i < 5 {
// 			page = 1
// 		} else if i % 5 == 0 {
// 			page += 1
// 		} else {
// 			page = page
// 		}
// 		var iim Imageinfomap = create_image_info_map(i, v, page)
// 		BulkImages = append(BulkImages, iim)
// 	}
// 	return BulkImages
// }

// func create_image_info_map(i int, afile string, page int) Imageinfomap {
// 	itype := get_type(afile)
// 	dir, filename := filepath.Split(afile)
// 	image_size := get_image_size(afile)
// 	image_http_path := create_image_http_addr(afile)
// 	ii := i + 1
// 	idx := strconv.Itoa(ii)
// 	pgx := strconv.Itoa(page)
// 	var ImageInfoMap Imageinfomap
// 	ImageInfoMap.Dirpath = dir
// 	ImageInfoMap.Filename = filename
// 	ImageInfoMap.Imagesize = image_size
// 	ImageInfoMap.ImageHttpAddr = image_http_path
// 	ImageInfoMap.Index = idx
// 	ImageInfoMap.IType = itype
// 	ImageInfoMap.Page = pgx
// 	ses := DBcon()
// 	defer ses.Close()
// 	imageinfo := ses.DB("coverart").C("coverart")
// 	imageinfo.Insert(ImageInfoMap)
// 	return ImageInfoMap
// }

// func get_type(afile string) string {
// 	if strings.Contains(afile, "thumb") {
// 		return "thumb"
// 	} else {
// 		return "original"
// 	}
// }

// func get_image_size(apath string) string {
// 	fi, err := os.Stat(apath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	size := fi.Size()
// 	newsize := int(size)
// 	return strconv.Itoa(newsize)
// }

// func create_image_http_addr(aimage string) string {
// 	httppath := "https://192.168.0.91:9090" + aimage[5:]
// 	return httppath
// }
