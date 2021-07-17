// LICENSE = GNU General Public License, version 2 (GPLv2)
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
	"time"
	"context"
	"log"
	"strconv"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
// 	"github.com/globalsign/mgo/bson"
)

// const (
// 	// OffSet := os.Getenv("AMPGO_OFFSET")
// 	Offset = 35
// )

// // GDistAlbum3 exported
// func GDistAlbum3() (DAlbum []map[string]string) {
// 	sess := DBcon()
// 	defer sess.Close()
// 	MAINc := sess.DB("maindb").C("maindb")
// 	var dlist []string
// 	MAINc.Find(nil).Distinct("album", &dlist)
// 	for _, d := range dlist {
// 		DMainc := sess.DB("maindb").C("maindb")
// 		b1 := bson.M{"album": d}
// 		b2 := bson.M{}
// 		var Boo map[string]string = make(map[string]string)
// 		DMainc.Find(b1).Select(b2).One(&Boo)
// 		DAlbum = append(DAlbum, Boo)
// 	}
// 	return
// }

// func GAlbInfo(DAlb map[string]string) (string, string, string, string, string, string) {
// 	var AlbInfo2 map[string]string = make(map[string]string)
// 	sess := DBcon()
// 	defer sess.Close()
// 	AMPc := sess.DB("maindb").C("maindb")
// 	AMPc.Find(bson.M{"album": DAlb["album"]}).One(&AlbInfo2)
// 	// return AlbInfo2["artist"], AlbInfo2["artistID"], AlbInfo2["album"], AlbInfo2["albumID"], AlbInfo2["picPath"], AlbInfo2["page"], AlbInfo2["idx"]
// 	return AlbInfo2["artist"], AlbInfo2["artistID"], AlbInfo2["album"], AlbInfo2["albumID"], AlbInfo2["picPath"], AlbInfo2["idx"] //, AlbInfo2["mainpage"]
// }

// //  exported
// type p2 struct {
// 	Titlez []string
// }

// // AlbPipeline exported
// func AlbPipeline(DAlb map[string]string) []string {
// 	var P2 []p2
// 	sess := DBcon()
// 	defer sess.Close()
// 	AMPc := sess.DB("maindb").C("maindb")
// 	pipeline2 := AMPc.Pipe([]bson.M{
// 		{"$match": bson.M{"album": DAlb["album"]}},
// 		{"$group": bson.M{"_id": "title", "titlez": bson.M{"$addToSet": "$title"}}},
// 		{"$project": bson.M{"titlez": 1}},
// 	}).Iter()
// 	err := pipeline2.All(&P2)
// 	CheckError(err, "\n AlbPipeline: Agg Album pipeline2 fucked up")
// 	// fmt.Printf("this is P2 %s", P2)
// 	return P2[0].Titlez
// }

// // AddTitleID exported
// func AddTitleID(titlez []string) []map[string]string {
// 	var TAAID []map[string]string
// 	sess := DBcon()
// 	defer sess.Close()
// 	AMP2c := sess.DB("maindb").C("maindb")
// 	for _, boo := range titlez {
// 		var TAid map[string]string = make(map[string]string)
// 		AMP2c.Find(bson.M{"title": boo}).Select(bson.M{"title": 1, "fileID": 1, "_id": 0}).One(&TAid)
// 		TAAID = append(TAAID, TAid)
// 	}
// 	return TAAID
// }

// AlbvieW exported
type AlbvieW struct {
	Artist    string              `bson:"artist"`
	ArtistID  string              `bson:"artistID"`
	Album     string              `bson:"album"`
	AlbumID   string              `bson:"albumID"`
	Songs     []map[string]string `bson:"songs"`
	AlbumPage string                 `bson:"albumpage"`
	NumSongs  string              `bson:"numsongs"`
	PicPath   string              `bson:"picPath"`
	Idx       string                 `bson:"idx"`
}

// //InsAlbViewID exported
func InsAlbViewID(artist string, artistID string, album string, albumID string, picPath string, songcount int, ATID []map[string]string, albpage int, idx int) {
// 	// MyAlbview.Artist = artist
// 	// MyAlbview.ArtistID = artistID
// 	// MyAlbview.Album = album
// 	// MyAlbview.AlbumID = albumID
// 	// MyAlbview.NumSongs = strconv.Itoa(songcount)
// 	// MyAlbview.PicPath = picPath
// 	// MyAlbview.Songs = ATID
// 	// MyAlbview.AlbumPage = strconv.Itoa(albpage)
// 	// MyAlbview.Idx = strconv.Itoa(idx)
	numsongs := strconv.Itoa(songcount)
	albpagE := strconv.Itoa(albpage)
	index := strconv.Itoa(idx)
	AlbView := bson.D{
		{"Artist", artist},
		{"ArtistID", artistID},
		{"Album", album},
		{"AlbumID", albumID},
		{"NumSongs", numsongs},
		{"PicPath", picPath},
		{"Songs", ATID},
		{"AlbumPage", albpagE},
		{"Idx", index},
	}

	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	CheckError(err, "Connections has failed")
    defer Close(client, ctx, cancel)

	insertOneResult, err := InsertOne(client, ctx, "tempdb2", "titleoffset", AlbView)
	CheckError(err, "titleoffset insertion has fucked up")
	fmt.Println(insertOneResult)

// 	// sess := DBcon()
// 	// defer sess.Close()
// 	// AVc := sess.DB("albview").C("albview")
// 	// AVc.Insert(MyAlbview)
	return
}

// // GAlbVCount exported
func GAlbVCount() int64 {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
		
	coll := client.Database("tempdb1").Collection("meta1")

    defer Close(client, ctx, cancel)
	opts := options.Count().SetMaxTime(2 * time.Second)

	count, err := coll.CountDocuments(
		context.TODO(),
		bson.D{{}},
		opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("the alb count is %v documents", count)
	return count

// 	// sess := DBcon()
// 	// defer sess.Close()
// 	// ALBc := sess.DB("albview").C("albview")
// 	// err := ALBc.Find(nil).All(&AlbV)
// 	// CheckError(err, "GALBVCount: albumcount has fucked up")
// 	// return
}

// //AlbumOffset exported
// // func AlbumOffset() {
// // 	sess := DBcon()
// // 	defer sess.Close()
// // 	ALBcc := sess.DB("albview").C("albview")
// // 	ALBview := GAlbVCount()
	
// // 	fmt.Printf("THIS IS ALBview FOR ALBUMVIEW %v", ALBview[0].Idx)
// // 	var page1 int = 1
// // 	for i, alb := range ALBview {
// // 		if i < Offset {
// // 			var BOO AlbvieW
// // 			BOO.Artist = alb.Artist
// // 			BOO.ArtistID = alb.ArtistID
// // 			BOO.Album = alb.Album
// // 			BOO.AlbumID = alb.AlbumID
// // 			BOO.Songs = alb.Songs
// // 			BOO.Page = strconv.Itoa(page1)
// // 			BOO.NumSongs = alb.NumSongs
// // 			BOO.PicPath = alb.PicPath
// // 			BOO.Idx = alb.Idx
// // 			ALBcc.Update(bson.M{"ArtistID": alb.ArtistID}, BOO)
// // 			ALBcc.Update(bson.M{"Page": alb.Page}, BOO)
// // 		} else if i % Offset == 0 {
// // 			page1++
// // 			var MOO AlbvieW
// // 			MOO.Artist = alb.Artist
// // 			MOO.ArtistID = alb.ArtistID
// // 			MOO.Album = alb.Album
// // 			MOO.AlbumID = alb.AlbumID
// // 			MOO.Songs = alb.Songs
// // 			MOO.Page = strconv.Itoa(page1)
// // 			MOO.NumSongs = alb.NumSongs
// // 			MOO.PicPath = alb.PicPath
// // 			MOO.Idx = alb.Idx
// // 			ALBcc.Update(bson.M{"AlbumID": alb.AlbumID, "Page": alb.Page}, MOO)
// // 		} else {
// // 			var MOO AlbvieW
// // 			MOO.Artist = alb.Artist
// // 			MOO.ArtistID = alb.ArtistID
// // 			MOO.Album = alb.Album
// // 			MOO.AlbumID = alb.AlbumID
// // 			MOO.Songs = alb.Songs
// // 			MOO.Page = strconv.Itoa(page1)
// // 			MOO.NumSongs = alb.NumSongs
// // 			MOO.PicPath = alb.PicPath
// // 			MOO.Idx = alb.Idx
// // 			ALBcc.Update(bson.M{"AlbumID": alb.AlbumID, "Page": alb.Page}, MOO)
// // 		}
// // 	}
// // }
