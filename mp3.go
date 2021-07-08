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
	"log"
	"fmt"
	"strconv"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"gopkg.in/mgo.v2/bson"
	"github.com/bogem/id3v2"
	"github.com/disintegration/imaging"
)

func getFileInfo(apath string) (filename string, size string) {
	ltn, err := os.Open(apath)
	CheckError(err, "getFileInfo: file open has fucked up")
	defer ltn.Close()
	ltnInfo, _ := ltn.Stat()
	filename = ltnInfo.Name()
	size = strconv.FormatInt(ltnInfo.Size(), 10)
	return
}

// UUID exported
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = 0x80
	uuid[4] = 0x40
	boo := hex.EncodeToString(uuid)
	return boo, nil
}

func resizeImage(infile string, outfile string) string {
	pic, err := imaging.Open(infile)
	if err != nil {
		fmt.Println(infile)
		fmt.Println("this is file Open error noartthumb")
		print(err)
	}
	sjImage := imaging.Resize(pic, 200, 0, imaging.Lanczos)
	err = imaging.Save(sjImage, outfile)
	CheckError(err, "resizeImage: image save has fucked up")
	return outfile
}

//DumpArtToFile is exported
func DumpArtToFile(apath string) (string, string, string, string, string) {
	tag, err := id3v2.Open(apath, id3v2.Options{Parse: true})
	artist := tag.Artist()
	album := tag.Album()
	title := tag.Title()
	genre := tag.Genre()
	CheckError(err, "Error while opening mp3 file")
	defer tag.Close()
	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	dumpOutFile2 := ""
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if !ok {
			log.Fatal("Couldn't assert picture frame")
		}
		dumpOutFile2 = os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + ".jpg"
		fmt.Println("\n\n this is apath")
		fmt.Println(apath)
		fmt.Println(dumpOutFile2)
		fmt.Println("this is dumpOutfile2 \n\n")

		g, err := os.Create(dumpOutFile2)
		defer g.Close()
		if err != nil {
			fmt.Println("Unable to create dumpOutFile2")
			// fmt.Println(f)
			fmt.Println(err)
		}
		n3, err := g.Write(pic.Picture)
		CheckError(err, "dumpOutfile2 Write has fucked up")
		fmt.Println(n3, "bytes written successfully")
	}
	outfile22 := resizeImage(dumpOutFile2, dumpOutFile2)
	return artist, album, title, genre, outfile22
}
// Tagmap exported
type Tagmap struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Dirpath   string `bson:"dirpath"`
	Filename  string `bson:"filename"`
	Extension string `bson:"extension"`
	FileID    string `bson:"fileID"`
	Filesize  string `bson:"filesize"`
	Artist    string `bson:"artist"`
	ArtistID  string `bson:"artistID"`
	Album     string `bson:"album"`
	AlbumID   string `bson:"albumID"`
	Title     string `bson:"title"`
	Genre     string `bson:"genre"`
	Page      string `bson:"page"`
	PicID     string `bson:"picID"`
	PicDB     string `bson:"picDB"`
	PicPath    string `bson:"picPath"`
	Idx       string `bson:"idx"`
}

// TAgMap exported
func TaGmap(apath string) (TAGmap Tagmap) {
	uuid, _ := UUID()
	artist, album, title, genre, picpath := DumpArtToFile(apath)
	fname, size := getFileInfo(apath)
	TAGmap.Dirpath = filepath.Dir(apath)
	TAGmap.Filename = fname
	TAGmap.Extension = filepath.Ext(apath)
	TAGmap.FileID = uuid
	TAGmap.Filesize = size
	TAGmap.Artist = artist
	TAGmap.ArtistID = "None"
	TAGmap.Album = album
	TAGmap.AlbumID = "None"
	TAGmap.Title = title
	TAGmap.Genre = genre
	TAGmap.Page = "None"
	TAGmap.PicID = uuid
	TAGmap.PicDB = "None"
	TAGmap.PicPath = picpath
	TAGmap.Idx = "None"
	ses := DBcon()
	defer ses.Close()
	tagz := ses.DB("tempdb1").C("meta1")
	tagz.Insert(TAGmap)
	return TAGmap
}

func GDistAlbum() (DAlbum []string) {
	sess := DBcon()
	defer sess.Close()
	MAINc := sess.DB("tempdb1").C("meta1")
	MAINc.Find(nil).Distinct("album", &DAlbum)
	return
}

// InsAlbumID exported
func InsAlbumID(alb string) {
	uuid, _ := UUID()
	sess := DBcon()
	defer sess.Close()
	TAlbIc := sess.DB("tempdb2").C("albumid")
	DALBI := map[string]string{"album": alb, "albumID": uuid}
	TAlbIc.Insert(&DALBI)
}

func GDistArtist() (DArtist []string) {
	sesC := DBcon()
	defer sesC.Close()
	MAINc := sesC.DB("tempdb1").C("meta1")
	MAINc.Find(nil).Distinct("artist", &DArtist)
	return
}

//InsArtistID exported
func InsArtistID(art string) {
	sesC := DBcon()
	defer sesC.Close()
	TArtIc := sesC.DB("tempdb2").C("artistid")
	uuid, _ := UUID()
	DARTI := map[string]string{"artist": art, "artistID": uuid}
	TArtIc.Insert(&DARTI)
}