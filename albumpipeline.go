package ampgosetup

import (
	// "context"
	// // "crypto/rand"
	// // "encoding/hex"
	// // "encoding/json"
	// "fmt"
	// // "github.com/bogem/id3v2"
	// // "github.com/disintegration/imaging"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	// "os"
	// "path/filepath"
	"strconv"
	// "strings"
	// "time"
)

func GDistAlbum() (DAlbAll []map[string]string) {
	DAlbumID := AmpgoDistinct("maindb", "maindb", "albumID")
	for _, albID := range DAlbumID {
		DAlb := AmpgoFindOne("maindb", "maindb", "albumID", albID)
		DAlbAll = append(DAlbAll, DAlb)
	}
	return
}

func AlbPipeline(DAlb map[string]string, page int, idx int) (MyAlbview AlbVieW2) {
	dirtysonglist := AmpgoFind("maindb", "maindb", "albumID", DAlb["albumID"])
	results := get_songs_for_album(dirtysonglist)
	songcount := len(results)
	MyAlbview.Artist = DAlb["artist"]
	MyAlbview.ArtistID = DAlb["artistID"]
	MyAlbview.Album = DAlb["album"]
	MyAlbview.AlbumID = DAlb["albumID"]
	MyAlbview.NumSongs = strconv.Itoa(songcount)
	MyAlbview.PicPath = DAlb["picPath"]
	MyAlbview.Songs = results
	MyAlbview.AlbumPage = strconv.Itoa(page)
	MyAlbview.Idx = strconv.Itoa(idx)
	MyAlbview.PicHttpAddr = DAlb["picHttpAddr"]
	return
}

func InsAlbViewID(MyAlbview AlbVieW2) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsAlbViewID: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "albumview", "albumview", &MyAlbview)
	CheckError(err2, "InsAlbViewID: AmpgoInsertOne has failed")
}
