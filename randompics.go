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

import (
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"
	"github.com/globalsign/mgo/bson"
	mathrand "math/rand"
	"path/filepath"
)





// func getDBCount() int {
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	ALBVc := sesC.DB("albview").C("albview")
// 	count, _ := ALBVc.Count()
// 	return count
// }

// func getRandNumb(dbc int) (myRand int) {
// 	mathrand.Seed(time.Now().UnixNano())
// 	min := 1
// 	max := dbc
// 	myRand = mathrand.Intn(max-min+1) + min
// 	return
// }

// func createRandNumList(dbc int) (randlist []int) {
// 	ranpicnum := os.Getenv("AMPGO_NUM_RAND_PICS")
// 	rpn, _ := strconv.Atoi(ranpicnum)
// 	for i := 0; i < rpn; i++ {
// 		r := getRandNumb(dbc)
// 		randlist = append(randlist, r)
// 	}
// 	return
// }

// func getAlbIDList(randlist []int) (randpicidlist []string) {
// 	for _, r := range randlist {
// 		sesC := DBcon()
// 		defer sesC.Close()
// 		RANDc := sesC.DB("goampgo").C("albidx")
// 		b1 := bson.M{"idxNum": r}
// 		var randresult map[string]string = make(map[string]string)
// 		RANDc.Find(b1).One(&randresult)
// 		randpicidlist = append(randpicidlist, randresult["albID"])
// 	}
// 	return
// }

// func getPicList(albidlist []string) (piclist []map[string]string) {
// 	for _, a := range albidlist {
// 		sesC := DBcon()
// 		defer sesC.Close()
// 		ALBVc := sesC.DB("albview").C("albview")
// 		b1 := bson.M{"albumID": a}
// 		var albinfo map[string]string = make(map[string]string)
// 		err := ALBVc.Find(b1).One(&albinfo)
// 		if err != nil {
// 			fmt.Println("CheckThumbDB has fucked up")
// 		}
// 		piclist = append(piclist, albinfo)
// 	}
// 	return
// }

//RanNewList exported
// func RanNewList(d []int, a []map[string]string) []map[string]string {
// 	var NewList []map[string]string
// 	for _, v := range d {
// 		NewList = append(NewList, a[v])
// 	}
// 	return NewList
// }
type Imageinfomap struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Dirpath   string `bson:"dirpath"`
	Filename  string `bson:"filename"`
	Imagesize string `bson:"imagesize"`
	ImageHttpAddr string `bson:"imagehttpaddr`
	Index string `bson:"index"`
	IType string `bson:"itype"`
}
//RanPics exported
func CreateRandomPicsDB() (BulkImages []Imageinfomap) {
	// =/root/static/
	thumb_path := os.Getenv("AMPGO_THUMB_PATH")
	thumb_glob_path := thumb_path + "/*.jpg"
	thumb_glob, err := filepath.Glob(thumb_glob_path)
	if err != nil {
		fmt.Println("CheckThumbDB has fucked up")
	}
	for i, v := range thumb_glob {
		var iim Imageinfomap = create_image_info_map(i, v)
		BulkImages := append(BulkImages, iim)
		return BulkImages
	}
	fmt.Println(BulkImages)
	return 
}

func create_image_info_map(i int, afile string) (ImageInfoMap Imageinfomap) {
	itype := "None"
	if strings.Contains(afile, "thumb") {
		itype = "thumb"
	} else {
		itype = "original"
	}
	dir, filename := filepath.Split(afile)
	image_size := get_image_size(afile)
	image_http_path := create_image_http_addr(afile)
	ii := i + 1
	idx := strconv.Itoa(ii)
	ImageInfoMap.Dirpath = dir
	ImageInfoMap.Filename = filename
	ImageInfoMap.Imagesize = image_size
	ImageInfoMap.ImageHttpAddr = image_http_path
	ImageInfoMap.Index = idx
	ImageInfoMap.IType = itype
	ses := DBcon()
	defer ses.Close()
	imageinfo := ses.DB("coverart").C("coverart")
	imageinfo.Insert(ImageInfoMap)
	return 
}

func get_image_size(apath string) string {
	fi, err := os.Stat(apath)
	if err != nil {
		fmt.Println(err)
	}
	size := fi.Size()
	newsize := int(size)
	
	return strconv.Itoa(newsize)
}

func create_image_http_addr(aimage string) string {
	newpath := aimage[5:]
	httppath := "https://192.168.0.91:9090" + newpath
	return httppath
}

func pagonate_coverart(alist []Imageinfomap) {
	fmt.Println("STARTING PAGINATION \n")
	mathrand.Seed(time.Now().UnixNano())
	mathrand.Shuffle(len(alist), func(i, j int) { alist[i], alist[j] = alist[j], alist[i] })
	var outslice []Imageinfomap
	count := 0
	for _, v := range alist {
		fmt.Println(v)
		// count ++
		// if count == 5 {
		// 	outslice = append(outslice, v)
		// 	fmt.Println(outslice)
		// 	sesCopy := DBcon()
		// 	defer sesCopy.Close()
		// 	RPICc := sesCopy.DB("coverart").C("rppages")
		// 	RPICc.Insert(outslice)
		// 	count = 0
		// 	outslice = nil
		// } else if count < 5 {
		// 	outslice = append(outslice, v)
		// } else {
		// 	fmt.Println("end of loop")
		// }
	}
}

// func shuffleList(num []int) []int {
// 	dest := make([]int, len(num))
// 	perm := mathrand.Perm(len(num))
// 	for i, v :=  range perm {
// 		dest[v] = num[i]
// 	}
// 	return dest
// }
