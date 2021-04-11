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
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

//GArtistInfo exported
func GArtistInfo(Art string) string {
	sesCopy := DBcon()
	defer sesCopy.Close()
	DARTc := sesCopy.DB("tempdb2").C("artistid")
	var ArtInfo map[string]string = make(map[string]string)
	DARTc.Find(bson.M{"artist": Art}).One(&ArtInfo)
	return ArtInfo["artistid"]
}

//GAlbumInfo exported
func GAlbumInfo(Alb string) string {
	sesCopy := DBcon()
	defer sesCopy.Close()
	DALBc := sesCopy.DB("tempdb2").C("albumid")
	var AlbInfo map[string]string = make(map[string]string)
	DALBc.Find(bson.M{"album": Alb}).One(&AlbInfo)
	return AlbInfo["albumid"]
}

//NoArtList Exported
var NoArtList []string

//CheckThumbDB ex
func CheckThumbDB(alb string) string {
	sses := DBcon()
	defer sses.Close()
	Thumbdb := sses.DB("thumbnails").C("meta")
	b1 := bson.M{"album": alb}
	b2 := bson.M{}
	var Gofthumb map[string]string = make(map[string]string)
	err := Thumbdb.Find(b1).Select(b2).One(&Gofthumb)
	if err != nil {
		fmt.Println("CheckThumbDB has fucked up")
		NoArtList = append(NoArtList, alb)
		return "None"
	}
	return Gofthumb["imgID"]
}

//CheckUnknownDB exported
// func CheckUnknownDB(dirp string) string {
// 	fmt.Printf("this is dirp from unknowdb %s", dirp)
// 	sess := DBcon()
// 	defer sess.Close()
// 	Unknowndb := sess.DB("unknownjpg").C("meta")
// 	b1 := bson.M{"dirpath": dirp}
// 	b2 := bson.M{}
// 	var Gofunknown map[string]string = make(map[string]string)
// 	err := Unknowndb.Find(b1).Select(b2).One(&Gofunknown)
// 	if err != nil {
// 		fmt.Println("CheckUnknownDB has fucked up")
// 	}
// 	fmt.Printf("\n\n\n THIS IS GOFUNKNOWN %s \n\n\n", &Gofunknown)
// 	return Gofunknown["imgID"]
// }

//UpdateMainDB exported
func UpdateMainDB(m2 map[string]string) (Doko Tagmap) {
	artID := GArtistInfo(m2["artist"])
	albID := GAlbumInfo(m2["album"])
	Doko.Dirpath = m2["dirpath"]
	Doko.Filename = m2["filename"]
	Doko.FileID = m2["fileID"]
	Doko.Extension = m2["extension"]
	Doko.Filesize = m2["filesize"]
	Doko.Artist = m2["artist"]
	Doko.ArtistID = artID
	Doko.Album = m2["album"]
	Doko.AlbumID = albID
	Doko.Title = m2["title"]
	Doko.Genre = m2["genre"]
	imid := CheckThumbDB(m2["album"])
	if imid != "None" {
		Doko.PicID = imid
		Doko.PicDB = "thumbnails"
		// Doko.PicCol = "meta"
	}
	Doko.Page = m2["page"]
	Doko.Idx = m2["idx"]

	sesC := DBcon()
	defer sesC.Close()
	DOKOc := sesC.DB("maindb").C("maindb")
	DOKOc.Insert(Doko)
	return
}

//TitleOffset exported
//Tagmap
func TitleOffset() (Tmap Tagmap) {
	fmt.Println("n\n\n STARTING TITLEOFFSET")
	var IPS []map[string]string
	sesCopy := DBcon()
	defer sesCopy.Close()
	MAINc := sesCopy.DB("tempdb1").C("meta1")
	MAINc.Find(nil).All(&IPS)
	var count int = 1
	var Pa int = 1
	for _, A := range IPS {
		count++
		if count <= OffSet {
			PPage := strconv.Itoa(Pa)
			Tmap.Page = PPage
		} else {
			count = 1
			Pa++
			PPage := strconv.Itoa(Pa)
			Tmap.Page = PPage
		}
		var idx string = strconv.Itoa(count)
		Tmap.Filename = A["filename"]
		Tmap.FileID = A["fileID"]
		Tmap.Extension = A["extension"]
		Tmap.Filesize = A["filesize"]
		Tmap.Artist = A["artist"]
		Tmap.ArtistID = A["artistID"]
		Tmap.Album = A["album"]
		Tmap.AlbumID = A["albumID"]
		Tmap.Title = A["title"]
		Tmap.Genre = A["genre"]
		Tmap.PicID = A["picID"]
		Tmap.Dirpath = A["dirpath"]
		Tmap.PicDB = A["picDB"]
		// Tmap.PicCol = A["picCol"]
		Tmap.Idx = idx
		TmpDBc := sesCopy.DB("tempdb2").C("titleoffset")
		TmpDBc.Insert(Tmap)
	}
	return
}
