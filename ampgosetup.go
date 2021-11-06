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
	// "context"
	"strconv"
	"path/filepath"
)

var OFFSET string = os.Getenv("AMPGO_OFFSET")
var OffSet int = convertSTR(OFFSET)

func convertSTR(astring string) int {
	Ofset, err := strconv.Atoi(astring)
	CheckError(err, "strconv has failed")
	return Ofset
}

//CheckError exported
func CheckError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		log.Println(msg)
		log.Println(err)
		panic(err)
	}
}

func durationVisit(pAth string, f os.FileInfo, err error) error {
	ext := path.Ext(pAth)
	if ext == ".mp3info" {
		InsertDurationInfo(pAth)
	} else {
		fmt.Println("WTF are you? You must be a Dir")
		fmt.Println(pAth)
	}
	return nil
}

var titlepage int = 0
var ii int = 0
func visit(pAth string, f os.FileInfo, err error) error {
	log.Println(pAth)
	if ii < OffSet {
		ii++
		titlepage = 1
	} else if ii % OffSet == 0 {
		ii++
		titlepage++
	} else {
		fmt.Println("I'm Not A Page")
		ii++
		titlepage = titlepage + 0
	}
	ext := path.Ext(pAth)
	if ext == ".jpg" {
		fmt.Println("FOOUND JPG")
	} else if ext == ".mp3" {
		TaGmap(pAth, titlepage, ii)
	} else {
		fmt.Println("WTF are you? You must be a Dir")
		fmt.Println(pAth)
	}
	log.Println(pAth)
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
	logtxtfile := os.Getenv("AMPGO_SETUP_LOG_PATH")
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

	log.Println("starting duration walk \n")
	filepath.Walk(os.Getenv("AMPGO_MEDIA_PATH"), durationVisit)
	log.Println("duration walk is complete \n")

	log.Println("starting walk \n")
	filepath.Walk(os.Getenv("AMPGO_MEDIA_PATH"), visit)
	log.Println("walk is complete \n")

	log.Println("starting GetDistAlbumMeta1 \n")
	dalb := AmpgoDistinct("tempdb1", "meta1", "album")
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
	dart := AmpgoDistinct("tempdb1", "meta1", "artist")
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

	log.Println("starting GetTitleOffSetAll")
	AllObj := GetTitleOffsetAll()
	log.Println("GetTitleOffSetAll is complete \n")

	log.Println("starting UpdateMainDB")
	var wg3 sync.WaitGroup
	for _, blob := range AllObj {
		log.Println(blob)
		wg3.Add(1)
		go func(blob map[string]string) {
			UpdateMainDB(blob)
			wg3.Done()
		}(blob)
		wg3.Wait()
	}
	log.Println("UpdateMainDB is complete \n")

	log.Println("starting ArtistFirst \n")
	var wg99a sync.WaitGroup
	for _, art := range dart {
		wg99a.Add(1)
		go func(art string) {
			ArtistFirst(art)
			wg99a.Done()
		}(art)
		wg99a.Wait()
	}
	log.Println("ArtistFirst is complete \n")

	log.Println("starting AlbumFirst \n")
	var wg99 sync.WaitGroup
	for _, alb := range dalb {
		wg99.Add(1)
		go func(alb string) {
			AlbumFirst(alb)
			wg99.Done()
		}(alb)
		wg99.Wait()
	}
	log.Println("AlbumFirst is complete \n")
	
	log.Println("starting GetPicForAlbum \n")
	var wg133 sync.WaitGroup
	for _, alb := range dalb {
		wg133.Add(1)
		go func(alb string) {
			zoo := GetPicForAlbum(alb)
			fmt.Println(zoo)
			wg133.Done()
		}(alb)
		wg133.Wait()
	}
	log.Println("GetPicForAlbum is complete \n")

	// //AggArtist
	log.Println("starting UpdateMainDB")
	DistArtist := GDistArtist2()
	log.Println("GDistArtist2 is complete \n")

	log.Println("starting GArtInfo2")
	var wg5 sync.WaitGroup
	var artpage int = 0
	for artIdx, DArtt := range DistArtist {
		if artIdx < OffSet {
			artpage = 1
		} else if artIdx % OffSet == 0 {
			artpage++
		} else {
			artpage = artpage + 0
		}
		
		APL := ArtPipline(DArtt, artpage, artIdx)
		
		wg5.Add(1)
		go func(APL ArtVieW2) {
			InsArtPipeline(APL)
			wg5.Done()
		}(APL)
		wg5.Wait()
	}
	fmt.Println("AggArtists is complete")
	log.Println("AggArtists is complete")
	// // ArtistOffSet()w11
	// // fmt.Println("ArtistOffSet is complete")

	// //AggAlbum
	// fmt.Println("AggAlbum has started")

	log.Println("Starting GDistAlbum3")
	DistAlbum := GDistAlbum()

	var wg6 sync.WaitGroup
	var albpage int = 0
	for albIdx, DAlb := range DistAlbum {
		wg6.Add(1)
		if albIdx < OffSet {
			albpage = 1
		} else if albIdx % OffSet == 0 {
			albpage++
		} else {
			albpage = albpage + 0
		}
		APLX := AlbPipeline(DAlb, albpage, albIdx)
		go func(APLX AlbVieW2) {
			InsAlbViewID(APLX)
			wg6.Done()
		}(APLX)
		wg6.Wait()
	}
	CreateRandomPicsDB()

	CreateRandomPlaylistDB()
	
	var lines = []string{
		"Go",
		"is",
		"the",
		"best",
		"programming",
		"language",
		"in",
		"the",
		"world",
	}

	f, err := os.Create("setup.txt")
    if err != nil {
        log.Fatal(err)
    }
    // remember to close the file
    defer f.Close()

    for _, line := range lines {
        _, err := f.WriteString(line + "\n")
        if err != nil {
            log.Fatal(err)
        }
    }

	// fmt.Println("AlbumOffSet is complete")
	t2 := time.Now().Sub(ti)
	fmt.Println(t2)
	fmt.Println("THE END")


// func Update() {
// 	logtxtfile := os.Getenv("AMPGO_SETUP_LOG_PATH")
// 	// If the file doesn't exist, create it or append to the file
// 	file, err := os.OpenFile(logtxtfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.SetOutput(file)
// 	log.Println("Logging started")

	// ti = time.Now()
	// fmt.Println(ti)
	// log.Println(ti)
	// runtime.GOMAXPROCS(runtime.NumCPU())

}


