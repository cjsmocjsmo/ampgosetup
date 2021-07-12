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
	"github.com/globalsign/mgo"
	"path/filepath"
	"runtime"
	"strconv"
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
		fmt.Println("FOOUND JPG")
		fmt.Println(pAth)
		// UnknownJpg(pAth)
	} else if ext == ".mp3" {
		fmt.Println("fuck yea mp3")
		TaGmap(pAth)
	} else {
		fmt.Println("WTF are you? You must be a Dir")
		fmt.Println(pAth)
	}
	return nil
}

func SetUpCheck() {
	fileinfo, err := os.Stat("setup.txt")
    if os.IsNotExist(err) {
		Setup()
        // log.Fatal("File does not exist.")
		// panic(err)
    }
    log.Println(fileinfo)
}

//SetUp is exported to main
func Setup() {
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

	ArtistOffset()
	fmt.Println("ArtistOffset is complete")

	//AggAlbum
	fmt.Println("AggAlbum has started")
	DistAlbum3 := GDistAlbum3()
	var wg6 sync.WaitGroup
	albIdx := 0
	for _, DAlb := range DistAlbum3 {
		wg6.Add(1)
		albIdx++
		go func(DAlb map[string]string, albIdx int) {
			artist, artistID, album, albumID, picPath, page, idx := GAlbInfo(DAlb)
			
			APL := AlbPipeline(DAlb)
			nss := len(APL)
			songcount := strconv.Itoa(nss)

			ATID := AddTitleID(APL)
			InsAlbViewID(artist, artistID, album, albumID, picPath, songcount, ATID, page, idx)
			wg6.Done()
		}(DAlb, albIdx)
		wg6.Wait()
	}

	AlbumOffset()

	bulklist := CreateRandomPicsDB()
	fmt.Println(bulklist)
	pagonate_coverart(bulklist)

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

	fmt.Println("AlbumOffset is complete")
	t2 := time.Now().Sub(ti)
	fmt.Println(t2)
	fmt.Println("THE END")
}
