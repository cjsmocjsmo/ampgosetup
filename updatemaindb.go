package ampgosetup

import (
	"context"
	// "crypto/rand"
	// "encoding/hex"
	// "encoding/json"
	"fmt"
	// "github.com/bogem/id3v2"
	// "github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"
)

func gArtistInfo(Art string) map[string]string {
	filter := bson.M{"artist": Art}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "gArtistInfo: MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("artistid")
	var ArtInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&ArtInfo)
	if err != nil {
		fmt.Println("gArtistInfo: has failed")
		fmt.Println(err)
	}
	return ArtInfo
}

func gAlbumInfo(Alb string) map[string]string {
	filter := bson.M{"album": Alb}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "gAlbumInfo: MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("albumid")
	var AlbInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&AlbInfo)
	if err != nil {
		fmt.Println("gAlbumInfo: has failed")
		fmt.Println(err)
	}
	return AlbInfo
}

func gDurationInfo(filename string) map[string]string {
	fmt.Println(filename)
	filter := bson.M{"filename": filename}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("durdb").Collection("durdb")
	var durinfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&durinfo)
	if err != nil {
		fmt.Println("gDuration has failed")
		fmt.Println(err)
	}
	fmt.Println(durinfo)
	return durinfo
}

// func StartsWith(astring string) string {
// 	if len(astring) > 3 {
// 		if astring[3:] == "The" || astring[3:] == "the" {
// 			return strings.ToUpper(astring[4:5])
// 		} else {
// 			return strings.ToUpper(astring[:1])
// 		}
// 	} else {
// 		return strings.ToUpper(astring[:1])
// 	}
// }

// func UpdateMainDB(m2 map[string]string) (Doko JsonMP3) {
// 	fmt.Println(m2["Filename"])

// 	artID := gArtistInfo(m2["Tags_artist"])
// 	fmt.Println(artID)
// 	albID := gAlbumInfo(m2["album"])
// 	fmt.Println("this is albID")
// 	fmt.Println(albID)
// 	fullpath := m2["dirpath"] + "/" + m2["filename"]
// 	fmt.Println(fullpath)
// 	duration := gDurationInfo(fullpath)
// 	fmt.Println("this is duration")
// 	fmt.Println(duration)


// 	Doko.Dirpath = m2["dirpath"]
// 	Doko.Filename = m2["filename"]
// 	Doko.Ext = m2["Ext"]
// 	Doko.File_id = m2["File_id"]
// 	Doko.File_Size = m2["File_Size"]
// 	Doko.Artist = m2["artist"]
// 	Doko.ArtistID = artID["artistID"]
// 	Doko.Album = m2["album"]
// 	Doko.AlbumID = albID["albumID"]
// 	Doko.Title = m2["title"]
// 	Doko.Genre = m2["genre"]
// 	Doko.PicID = m2["picID"]
// 	Doko.PicDB = "thumbnails"
// 	Doko.TitlePage = m2["titlepage"]
// 	Doko.Idx = m2["idx"]
// 	Doko.PicPath = m2["picPath"]
// 	Doko.PicHttpAddr = m2["picHttpAddr"]
// 	Doko.HttpAddr = m2["httpaddr"]
// 	Doko.Duration = duration["duration"]
// 	// Doko.ArtStart = StartsWith(m2["artist"])
// 	Doko.ArtStart = strings.ToUpper(m2["artist"][:1])
// 	Doko.AlbStart = strings.ToUpper(m2["album"][:1])
// 	Doko.TitStart = strings.ToUpper(m2["title"][:1])
// 	Doko.Howl = m2["howl"]
// 	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
// 	CheckError(err, "UpdateMainDB: Connections has failed")
// 	defer Close(client, ctx, cancel)
// 	_, err2 := InsertOne(client, ctx, "maindb", "maindb", &Doko)
// 	CheckError(err2, "UpdateMainDB: maindb insertion has failed")
// 	return
// }
