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

package main

import (
	"os"
	// "io"
	"fmt"
	"log"
	"path"
	"sync"
	"time"
	// "strings"
	//"crypto/sha512"
	//"encoding/hex"
	// "gopkg.in/mgo.v2"
	"github.com/globalsign/mgo"
	"path/filepath"
	"runtime"
	// "strconv"
)

//Set Constants
const (
	OffSet = 10
)

//DBcon exported
func DBcon() *mgo.Session {
	s, err := mgo.Dial(os.Getenv("AMP_AMPDB_ADDR"))
	if err != nil {
		log.Println("Session creation dial error")
		log.Println(err)
	}
	log.Println("Session Connection to db established")
	return s
}

//CheckError exported
func CheckError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

// //GMainAll exported
// func GMainAll() (Main2SL []map[string]string) {
// 	sesCopy := DBcon()
// 	defer sesCopy.Close()
// 	MAINc := sesCopy.DB("tempdb2").C("titleoffset")
// 	MAINc.Find(nil).All(&Main2SL)
// 	return
// }

//GMAll exported
func GMAll() (Main2SL []map[string]string) {
	sesC := DBcon()
	defer sesC.Close()
	MAINc := sesC.DB("tempdb2").C("titleoffset")
	MAINc.Find(nil).All(&Main2SL)
	return
}

func visit(pAth string, f os.FileInfo, err error) error {
	// println("this is path from visit \n")
	ext := path.Ext(pAth)
	if ext == ".jpg" {
		fmt.Println(pAth)
		// UnknownJpg(pAth)
	} else if ext == ".mp3" {
		fmt.Println("fuck")
		TaGmap(pAth)
	} else {
		fmt.Println("WTF are you? You must be a Dir")
		fmt.Println(pAth)
	}
	return nil
}

//SetUp is exported to main
func main() {
	ti := time.Now()
	fmt.Println(ti)
	runtime.GOMAXPROCS(runtime.NumCPU())

	filepath.Walk(os.Getenv("AMPGO_MEDIA_PATH"), visit)

	dalb := GDistAlbum()
	var wg1 sync.WaitGroup
	for _, alb := range dalb {
		wg1.Add(1)
		go func(alb string) {
			InsAlbumID(alb)
			wg1.Done()
		}(alb)
		wg1.Wait()
	}

	dart := GDistArtist()
	var wg2 sync.WaitGroup
	for _, art := range dart {
		wg2.Add(1)
		go func(art string) {
			InsArtistID(art)
			wg2.Done()
		}(art)
		wg2.Wait()
	}

	TitleOffset()

	AllObj := GMAll()
	var wg3 sync.WaitGroup
	for _, blob := range AllObj {
		wg3.Add(1)
		go func(blob map[string]string) {
			UpdateMainDB(blob)
			wg3.Done()
		}(blob)
		wg3.Wait()
	}

	fmt.Println("creating and inserting thumbnails is complete")
	fmt.Println("Inserting album and artists ids is complete")

	//AggArtist
	DistArtist := GDistArtist2()
	fmt.Printf("\n\n\n this is DistArtist %s \n\n\n", DistArtist)
	var wg5 sync.WaitGroup
	artIdx := 0
	for _, DArtt := range DistArtist {
		wg5.Add(1)
		artIdx++
		go func(DArtt map[string]string, artIdx int) {
			GAI := GArtInfo2(DArtt)
			APL := ArtPipeline(DArtt)
			AlbID := AddAlbumID(APL)
			InsArtIPipe2(GAI, AlbID, artIdx)
			wg5.Done()
		}(DArtt, artIdx)
		wg5.Wait()
	}
	fmt.Println("AggArtists is complete")

	// ArtistOffset()
	// fmt.Println("ArtistOffset is complete")

	// // //AggAlbum
	// fmt.Println("AggAlbum has started")
	// DistAlbum3 := GDistAlbum3()

	// for i, DAlb := range DistAlbum3 {
	// 	B64Image := GetSImage(DAlb)
	// 	GAI := GAlbInfo(DAlb)
	// 	APL := AlbPipeline(DAlb)
	// 	ATID := AddTitleID(APL)

	// 	Artist := GAI["artist"]
	// 	ArtistID := GAI["artistID"]
	// 	Album := GAI["album"]
	// 	AlbumID := GAI["albumID"]

	// 	nss := len(APL) //titlez
	// 	numson := strconv.Itoa(nss)

	// 	var AView AlbvieW
	// 	AView.Artist = Artist
	// 	AView.ArtistID = ArtistID
	// 	AView.Album = Album
	// 	AView.AlbumID = AlbumID
	// 	AView.Songs = ATID
	// 	AView.Page = 0
	// 	AView.NumSongs = numson
	// 	AView.HSImage = B64Image
	// 	AView.Idx = i
	// 	InsAlbViewID(AView)
	// }

	// AlbumOffset()
	// fmt.Println("AlbumOffset is complete")

	// FileSizeF()
	// fmt.Println("Stats are complete")

	// // // RPics()
	// // // fmt.Println("RPics is complete")

	// // TextSearchIndexes()
	// // fmt.Println("text search indexes complete")

	// DropTempDBs()

	t2 := time.Now().Sub(ti)
	fmt.Println(t2)
	fmt.Println("THE END")
	// fmt.Printf("\n This is noartlist \n %v \n", NoArtList)
}
