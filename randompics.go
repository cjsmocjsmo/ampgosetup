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
	"strconv"
	// "time"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// mathrand "math/rand"
	"path/filepath"
)

//use index in albums db

//create a count db.  Index and albumID
//get db count
//get 5 random numbers from 0-dbcount
//use random numbers to search count db for album Index
//return art

// func distAlbIDList() (Alblist []string) {
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	ALBVc := sesC.DB("maindb").C("maindb")
// 	err := ALBVc.Find(nil).Distinct("albumID", &Alblist)
// 	if err != nil {
// 		fmt.Println("random art fucked up")
// 		fmt.Println(err)
// 	}
// 	return
// }

//CreateIndexAlbumIDDB exported
// func CreateIndexAlbumIDDB() {
// 	var NumPics []map[string]string
// 	albidlist := distAlbIDList()
// 	for i, ab := range albidlist {
// 		ii := strconv.Itoa(i)
// 		ace := map[string]string{"idxNum": ii, "albID": ab}
// 		NumPics = append(NumPics, ace)
// 	}
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	ALBVc := sesC.DB("goampgo").C("albidx")
// 	ALBVc.Insert(&NumPics)
// 	return
// }

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
}
//RanPics exported
func CreateRandomPicsDB() (ImageInfoMap Imageinfomap) {
	// =/root/static/
	thumb_path := os.Getenv("AMPGO_THUMB_PATH")
	thumb_glob_path := thumb_path + "/*.jpg"
	thumb_glob, err := filepath.Glob(thumb_glob_path)
	if err != nil {
		fmt.Println("CheckThumbDB has fucked up")
	}
	
	for _, v := range thumb_glob {
		dir, filename := filepath.Split(v)
		image_size := get_image_size(v)
		image_http_path := create_image_http_addr(v)

		fmt.Printf("This is crap %d, %d, %d", dir, filename, image_size)

		ImageInfoMap.Dirpath = dir
		ImageInfoMap.Filename = filename
		ImageInfoMap.Imagesize = image_size
		ImageInfoMap.imagehttpaddr = image_http_path

		fmt.Printf("this is ImageInfoMap %d", ImageInfoMap)

		// ses := DBcon()
		// defer ses.Close()
		// tagz := ses.DB("coverart").C("meta1")
		// tagz.Insert(TAGmap)
		// return TAGmap


	}


	// CreateIndexAlbumIDDB()
	// ofse := os.Getenv("AMPGO_OFFSET")
	// offset, _ := strconv.Atoi(ofse)


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
	httppath := "https://192.168.0.91:9090/" + newpath
	return httppath
}







// func chunckit(nl []map[string]string, ) {
// 	var outslice []string
// 	count := 0
// 	for _, v := range nl {
// 		count ++
// 		if count == 5 {
// 			outslice = append(outslice, v["albumid"])
// 			sesCopy := DBcon()
// 			defer sesCopy.Close()
// 			RPICc := sesCopy.DB("goampgo").C("randompics")
// 			RPICc.Insert(outslice)
// 			count = 0
// 			outslice = nil
// 		} else if count < 5 {
// 			outslice = append(outslice, v["albumid"])
// 		} else {
// 			fmt.Println("end of loop")
// 		}
// 	}
// }

// func shuffleList(num []int) []int {
// 	dest := make([]int, len(num))
// 	perm := mathrand.Perm(len(num))
// 	for i, v :=  range perm {
// 		dest[v] = num[i]
// 	}
// 	return dest
// }
