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
	"os"
	"fmt"
	"log"
	"path"
	"sync"
	"time"
	"runtime"
	"context"
	// "strconv"
	"path/filepath"
	// "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

//Set Constants
const (
	OffSet = 35
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return client, ctx, cancel, err
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
    collection := client.Database(dataBase).Collection(col)
    result, err := collection.InsertOne(ctx, doc)
    return result, err
}

// func Query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
// 	collection := client.Database(dataBase).Collection(col)
// 	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
// 	return
// }


//CheckError exported
func CheckError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

var titlepage int = 0
var i int = 0
func visit(pAth string, f os.FileInfo, err error) error {
	if i < OffSet {
		i++
		titlepage = 1
	} else if i % OffSet == 0 {
		i++
		titlepage++
	} else {
		fmt.Println("I'm Not A Page")
		i++
		titlepage = titlepage + 0
	}
	ext := path.Ext(pAth)
	if ext == ".jpg" {
		fmt.Println("FOOUND JPG")
	} else if ext == ".mp3" {
		TaGmap(pAth, titlepage, i)
	} else {
		fmt.Println("WTF are you? You must be a Dir")
		fmt.Println(pAth)
	}
	return nil
}

func SetUpCheck() {
	Setup()
	// fileinfo, err := os.Stat("setup.txt")
    // if os.IsNotExist(err) {
	// 	Setup()
    // }
    // log.Println(fileinfo)
}

//SetUp is exported to main
func Setup() {
	logtxtfile := os.Getenv("AMPGO_LOG_PATH")
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(logtxtfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Logging started")
	




	ti := time.Now()
	fmt.Println(ti)
	log.Println(ti)
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("starting walk \n")
	filepath.Walk(os.Getenv("AMPGO_MEDIA_PATH"), visit)
	log.Println("walk is complete \n")

	log.Println("starting GetDistAlbumMeta1 \n")
	dalb := GetDistAlbumMeta1()
	fmt.Println(dalb)
	log.Println(dalb)
	log.Println("GetDistAlbumMeta1 is complete \n")


	log.Println("starting InsAlbumID \n")
	var wg1 sync.WaitGroup
	for _, alb := range dalb {
		wg1.Add(1)
		go func(alb string) {
			InsAlbumID(alb)
			wg1.Done()
		}(alb)
		wg1.Wait()
	}
	log.Println("InsAlbumID is complete \n")

	log.Println("starting GDistArtist")
	dart := GDistArtist()
	fmt.Println(dart)
	log.Println(dart)
	log.Println("GDistArtist is complete \n")

	log.Println("starting InsArtistID")
	var wg2 sync.WaitGroup
	for _, art := range dart {
		wg2.Add(1)
		go func(art string) {
			InsArtistID(art)
			wg2.Done()
		}(art)
		wg2.Wait()
	}
	log.Println("InsArtistID is complete \n")

	log.Println("starting GetTitleOffsetAll")
	AllObj := GetTitleOffsetAll()
	log.Println("GetTitleOffsetAll is complete \n")

	log.Println("starting UpdateMainDB")
	var wg3 sync.WaitGroup
	for _, blob := range AllObj {
		wg3.Add(1)
		go func(blob map[string]string) {
			UpdateMainDB(blob)
			wg3.Done()
		}(blob)
		wg3.Wait()
	}
	log.Println("UpdateMainDB is complete \n")

	// fmt.Println("creating and inserting thumbnails is complete")
	// fmt.Println("Inserting album and artists ids is complete")

	// //AggArtist
	log.Println("starting UpdateMainDB")
	DistArtist := GDistArtist2()
	for _, v := range DistArtist {
		fmt.Printf("%v this is GDistArtist2", v)
	}
	log.Println("GDistArtist2 is complete \n")
	// var wg5 sync.WaitGroup
	// var artpage int = 0
	// for artIdx, DArtt := range DistArtist {
	// 	if artIdx < OffSet {
	// 		artpage = 1
	// 	} else if artIdx % OffSet == 0 {
	// 		artpage++
	// 	} else {
	// 		artpage = artpage + 0
	// 	}
	// 	wg5.Add(1)
	// 	go func(DArtt map[string]string, artIdx int, artpage int) {
	// 		GAI := GArtInfo2(DArtt)
	// 		for _, g := range GAI {
	// 			fmt.Println("%v THIS IS GGGGGGGGG\n\n\n", g)
	// 		}
	// // // 		APL := ArtPipeline(DArtt)
	// // // 		AlbID := AddAlbumID(APL)
	// // // 		// aartIdX := strconv.Itoa(artIdx)
	// // // 		// aartpage := strconv.Itoa(artpage)
	// // // 		InsArtIPipe2(GAI, AlbID, artIdx, artpage)
	// 		wg5.Done()
	// 	}(DArtt, artIdx, artpage)
	// 	wg5.Wait()
	// }
	// fmt.Println("AggArtists is complete")

	// // ArtistOffset()w11
	// // fmt.Println("ArtistOffset is complete")

	// //AggAlbum
	// fmt.Println("AggAlbum has started")


	// DistAlbum3 := GDistAlbum3()
	// for _, v := range DistAlbum3 {
	// 	fmt.Printf("%v this is DistAlbum3", v)
	// }

	// var wg6 sync.WaitGroup
	// var albpage int
	// for albIdx, DAlb := range DistAlbum3 {
	// 	wg6.Add(1)
	// 	if albIdx < OffSet {
	// 		albpage = 1
	// 	} else if albIdx % OffSet == 0 {
	// 		albpage++
	// 	} else {
	// 		albpage = albpage + 0
	// 	}
	// 	fmt.Println("\n THIS IS ALBPAGE")
	// 	fmt.Println(albpage)
	// 	fmt.Println("\n THIS IS ALBIDX")
	// 	fmt.Println(albIdx)

	// 	go func(DAlb map[string]string, albIdx int, albpage int) {
			
	// 		APL := AlbPipeline(DAlb)
	// 		songcount := len(APL)
	// 		ATID := AddTitleID(APL)
	// 		// songcount := strconv.Itoa(nss)
	// 		// aidx, _ := strconv.Atoi(idx)
	// 		artist, artistID, album, albumID, picPath, _ := GAlbInfo(DAlb)
	// 		InsAlbViewID(artist, artistID, album, albumID, picPath, songcount, ATID, albpage, albIdx)
	// 		wg6.Done()
	// 	}(DAlb, albIdx, albpage)
	// 	wg6.Wait()
	// }

	// // AlbumOffset()

	// var bulklist []Imageinfomap = CreateRandomPicsDB()
	// fmt.Println(bulklist)

	// var lines = []string{
	// 	"Go",
	// 	"is",
	// 	"the",
	// 	"best",
	// 	"programming",
	// 	"language",
	// 	"in",
	// 	"the",
	// 	"world",
	// }

	// f, err := os.Create("setup.txt")
    // if err != nil {
    //     log.Fatal(err)
    // }
    // // remember to close the file
    // defer f.Close()

    // for _, line := range lines {
    //     _, err := f.WriteString(line + "\n")
    //     if err != nil {
    //         log.Fatal(err)
    //     }
    // }

	// fmt.Println("AlbumOffset is complete")
	t2 := time.Now().Sub(ti)
	fmt.Println(t2)
	fmt.Println("THE END")
}


