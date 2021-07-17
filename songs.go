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
	// "fmt"
	"log"
	"context"
	// "strconv"
	"go.mongodb.org/mongo-driver/bson"
)

//GArtistInfo exported
func gArtistInfo(Art string) string {
	filter := bson.M{"artist": Art}
	// opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("artistid")
	var ArtInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&ArtInfo)
	if err != nil { log.Fatal(err) }
	return ArtInfo["artistid"]

	// sesCopy := DBcon()
	// defer sesCopy.Close()
	// DARTc := sesCopy.DB("tempdb2").C("artistid")
	// var ArtInfo map[string]string = make(map[string]string)
	// DARTc.Find(bson.M{"artist": Art}).One(&ArtInfo)
	// return ArtInfo["artistID"]
}

// //GAlbumInfo exported
func gAlbumInfo(Alb string) string {
	filter := bson.M{"album": Alb}
	// opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("albumid")
	var AlbInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&AlbInfo)
	if err != nil { log.Fatal(err) }
	return AlbInfo["albumid"]
// 	sesCopy := DBcon()
// 	defer sesCopy.Close()
// 	DALBc := sesCopy.DB("tempdb2").C("albumid")
// 	var AlbInfo map[string]string = make(map[string]string)
// 	DALBc.Find(bson.M{"album": Alb}).One(&AlbInfo)
// 	return AlbInfo["albumID"]
}

// //UpdateMainDB exported
func UpdateMainDB(m2 map[string]string) (Doko Tagmap) {
	artID := gArtistInfo(m2["artist"])
	albID := gAlbumInfo(m2["album"])
	Doko.Dirpath = m2["dirpath"]
	Doko.Filename = m2["filename"]
	Doko.Extension = m2["extension"]
	Doko.FileID = m2["fileID"]
	Doko.Filesize = m2["filesize"]
	Doko.Artist = m2["artist"]
	Doko.ArtistID = artID
	Doko.Album = m2["album"]
	Doko.AlbumID = albID
	Doko.Title = m2["title"]
	Doko.Genre = m2["genre"]
	Doko.PicID = m2["picID"]
	Doko.PicDB = "thumbnails"
	Doko.TitlePage = m2["titlepage"]
	Doko.Idx = m2["idx"]
	Doko.PicPath = m2["picPath"]

	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "maindb", "maindb", &Doko)
	CheckError(err2, "maindb insertion has failed")


	// sesC := DBcon()
	// defer sesC.Close()
	// DOKOc := sesC.DB("maindb").C("maindb")
	// DOKOc.Insert(Doko)
	return
}

// //TitleOffset exported
// //Tagmap
// func TitleOffset() (Tmap Tagmap) {
// 	fmt.Println("n\n\n STARTING TITLEOFFSET")
// 	var IPS []map[string]string
// 	sesCopy := DBcon()
// 	defer sesCopy.Close()
// 	MAINc := sesCopy.DB("tempdb1").C("meta1")
// 	MAINc.Find(nil).All(&IPS)
// 	var count int = 1
// 	var Pa int = 1
// 	for _, A := range IPS {
// 		count++
// 		if count <= OffSet {
// 			PPage := strconv.Itoa(Pa)
// 			Tmap.TitlePage = PPage
// 		} else {
// 			count = 1
// 			Pa++
// 			PPage := strconv.Itoa(Pa)
// 			Tmap.TitlePage = PPage
// 		}
// 		var idx string = strconv.Itoa(count)
// 		Tmap.Filename = A["filename"]
// 		Tmap.Extension = A["extension"]
// 		Tmap.FileID = A["fileID"]
// 		Tmap.Filesize = A["filesize"]
// 		Tmap.Artist = A["artist"]
// 		Tmap.ArtistID = A["artistID"]
// 		Tmap.Album = A["album"]
// 		Tmap.AlbumID = A["albumID"]
// 		Tmap.Title = A["title"]
// 		Tmap.Genre = A["genre"]
// 		Tmap.PicID = A["picID"]
// 		Tmap.PicPath = A["picPath"]
// 		Tmap.Dirpath = A["dirpath"]
// 		Tmap.PicDB = A["picDB"]
// 		Tmap.Idx = idx
// 		TmpDBc := sesCopy.DB("tempdb2").C("titleoffset")
// 		TmpDBc.Insert(Tmap)
// 	}
// 	return
// }
