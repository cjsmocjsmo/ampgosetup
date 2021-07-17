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
	"time"
	"context"
	"strings"
	"strconv"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/globalsign/mgo/bson"
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
		return os.Getenv("AMPGO_NO_ART_PIC_PATH")
	}
	sjImage := imaging.Resize(pic, 200, 0, imaging.Lanczos)
	err = imaging.Save(sjImage, outfile)
	CheckError(err, "resizeImage: image save has fucked up")
	return outfile
}

// //DumpArtToFile is exported
func DumpArtToFile(apath string) (string, string, string, string, string) {
	tag, err := id3v2.Open(apath, id3v2.Options{Parse: true})
	artist := tag.Artist()
	album := tag.Album()
	title := tag.Title()
	genre := tag.Genre()
	CheckError(err, "Error while opening mp3 file")
	defer tag.Close()
	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	newdumpOutFile2 := ""
	newdumpOutFileThumb := ""
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if !ok {
			log.Fatal("Couldn't assert picture frame")
		}
		dumpOutFile2 := os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + ".jpg"
		newdumpOutFile2 = strings.Replace(dumpOutFile2, " ", "_", -1)
		dumpOutFileThumb := os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + "_thumb.jpg"
		newdumpOutFileThumb = strings.Replace(dumpOutFileThumb, " ", "_", -1)
		g, err := os.Create(newdumpOutFile2)
		defer g.Close()
		if err != nil {
			fmt.Println("Unable to create newdumpOutFile2")
			fmt.Println(err)
		}
		n3, err := g.Write(pic.Picture)
		CheckError(err, "newdumpOutfile2 Write has fucked up")
		fmt.Println(n3, "bytes written successfully")
	}
	outfile22 := resizeImage(newdumpOutFile2, newdumpOutFileThumb)
	return artist, album, title, genre, outfile22
}

// Tagmap exported
type Tagmap struct {
	Dirpath    string `bson:"dirpath"`
	Filename   string `bson:"filename"` 
	Extension  string `bson:"extension"`
	FileID     string `bson:"fileID"`
	Filesize   string `bson:"filesize"`
	Artist     string `bson:"artist"`
	ArtistID   string `bson:"artistID"`
	Album      string `bson:"album"`
	AlbumID    string `bson:"albumID"`
	Title      string `bson:"title"`
	Genre      string `bson:"genre"`
	TitlePage  string `bson:"titlepage"`
	PicID      string `bson:"picID"`
	PicDB      string `bson:"picDB"` 
	PicPath    string `bson:"picPath"`
	Idx        string    `bson:"idx"`
}

// TAgMap exported
func TaGmap(apath string, apage int, idx int) (TaGmaP Tagmap) {
	page := strconv.Itoa(apage)
	index := strconv.Itoa(idx)
	uuid, _ := UUID()
	artist, album, title, genre, picpath := DumpArtToFile(apath)
	fname, size := getFileInfo(apath)
	TaGmaP.Dirpath = filepath.Dir(apath)
	TaGmaP.Filename = fname
	TaGmaP.Extension = filepath.Ext(apath)
	TaGmaP.FileID = uuid
	TaGmaP.Filesize = size
	TaGmaP.Artist = artist
	TaGmaP.ArtistID = "None"
	TaGmaP.Album = album
	TaGmaP.AlbumID = "None"
	TaGmaP.Title = title
	TaGmaP.Genre = genre
	TaGmaP.TitlePage = page
	TaGmaP.PicID = uuid
	TaGmaP.PicDB = "None"
	TaGmaP.PicPath = picpath
	TaGmaP.Idx = index
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "tempdb1", "meta1", &TaGmaP)
	CheckError(err2, "Tempdb1 insertion has failed")
	return
}









func GetDistAlbumMeta1() []string {
	filter := bson.D{{}}
	opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb1").Collection("meta1")

	DD1, err2 := collection.Distinct(context.TODO(), "album", filter, opts)
	CheckError(err2, "MongoDB distinct album has failed")
	var DAlbum1 []string
	for _, DD := range DD1 {
		zoo := fmt.Sprintf("%s", DD)
		DAlbum1 = append(DAlbum1, zoo)
	}
	return DAlbum1


// // // 	sess := DBcon()
// // // 	defer sess.Close()
// // // 	MAINc := sess.DB("tempdb1").C("meta1")
// // // 	MAINc.Find(nil).Distinct("album", &DAlbum)
	
}

// InsAlbumID exported
func InsAlbumID(alb string) {
	uuid, _ := UUID()
	var Albid interface{}
	Albid = bson.D{
		{"album", alb},
		{"albumID", uuid},
	}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "tempdb2", "albumid", Albid)
	CheckError(err2, "albumID insertion has failed")
	return
// 	sess := DBcon()
// 	defer sess.Close()
// 	TAlbIc := sess.DB("tempdb2").C("albumid")
// 	DALBI := map[string]string{"album": alb, "albumID": uuid}
// 	TAlbIc.Insert(&DALBI)
}

func GDistArtist() []string {
	filter := bson.D{{}}
	opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb1").Collection("meta1")
	DArtist, err2 := collection.Distinct(context.TODO(), "artist", filter, opts)
	CheckError(err2, "MongoDB distinct album has failed")
	var DArtist1 []string
	for _, DA := range DArtist {
		zooo := fmt.Sprintf("%s", DA)
		DArtist1 = append(DArtist1, zooo)
	}
	return DArtist1
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	MAINc := sesC.DB("tempdb1").C("meta1")
// 	MAINc.Find(nil).Distinct("artist", &DArtist)
// 	return
}

//InsArtistID exported
func InsArtistID(art string) {
	uuid, _ := UUID()
	var Artid interface{}
	Artid = bson.D{
		{"artist", art},
		{"artistID", uuid},
	}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "tempdb2", "artistid", Artid)
	CheckError(err2, "artistID insertion has failed")
	return
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	TArtIc := sesC.DB("tempdb2").C("artistid")
// 	uuid, _ := UUID()
// 	DARTI := map[string]string{"artist": art, "artistID": uuid}
// 	TArtIc.Insert(&DARTI)
}

//GMAll exported
// func GMAll() (Main2SL []map[string]string) {
	func GetTitleOffsetAll() (Main2SL []map[string]string) {
		// filter := bson.D{{}}
		// opts := options.Distinct().SetMaxTime(2 * time.Second)
		client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
		defer Close(client, ctx, cancel)
		CheckError(err, "MongoDB connection has failed")
		collection := client.Database("tempdb1").Collection("meta1")
		cur, err := collection.Find(context.Background(), bson.D{})
		if err != nil { log.Fatal(err) }
	
	
		if err = cur.All(context.Background(), &Main2SL); err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n\n\n This is Main2SL \n")
		fmt.Println(Main2SL)
		// sesC := DBcon()
		// defer sesC.Close()
		// MAINc := sesC.DB("tempdb2").C("titleoffset")
		// MAINc.Find(nil).All(&Main2SL)
		return
	}