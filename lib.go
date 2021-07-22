package ampgosetup

import (
	"os"
	"log"
	"fmt"
	"time"
	// "sort"
	"context"
	"strings"
	"strconv"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"github.com/bogem/id3v2"
	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

type ArtVieW2 struct {
	Artist   string              `bson:"artist"`
	ArtistID string              `bson:"artistID"`
	Albums   []map[string]string `bson:"albums"`
	Page     string              `bson:"page"`
	Index    string              `bson:"idx"`
}

type AlbVieW2 struct {
	Artist    string              `bson:"artist"`
	ArtistID  string              `bson:"artistID"`
	Album     string              `bson:"album"`
	AlbumID   string              `bson:"albumID"`
	Songs     []map[string]string `bson:"songs"`
	AlbumPage string              `bson:"albumpage"`
	NumSongs  string              `bson:"numsongs"`
	PicPath   string              `bson:"picPath"`
	Idx       string              `bson:"idx"`
}

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return client, ctx, cancel, err
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
    collection := client.Database(dataBase).Collection(col)
    result, err := collection.InsertOne(ctx, doc)
    return result, err
}

func Query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func AmpgoDistinct(db string, coll string, fieldd string) []string {
	filter := bson.D{}
	opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database(db).Collection(coll)
	DD1, err2 := collection.Distinct(context.TODO(), fieldd, filter, opts)
	CheckError(err2, "MongoDB distinct album has failed")
	var DAlbum1 []string
	for _, DD := range DD1 {
		zoo := fmt.Sprintf("%s", DD)
		DAlbum1 = append(DAlbum1, zoo)
	}
	return DAlbum1
}

func AmpgoFindOne(db string, coll string, filtertype string, filterstring string) map[string]string {
	filter := bson.M{filtertype: filterstring}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database(db).Collection(coll)
	var results map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&results)
	if err != nil { log.Fatal(err) }
	return results
}

func AmpgoFind(dbb string, collb string, filtertype string, filterstring string) []map[string]string {
	filter := bson.D{{filtertype, filterstring}}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	coll := client.Database(dbb).Collection(collb)
	cur, err := coll.Find(context.TODO(), filter)
	CheckError(err, "ArtPipeline find has failed")
	var results []map[string]string //all albums for artist to include double entries
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	return results
}

func AmpgoInsertOne(db string, coll string, ablob map[string]string) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "albumID insertion has failed")
}

//////////////////////////////////////////////////////////////////////////

func getFileInfo(apath string) (filename string, size string) {
	ltn, err := os.Open(apath)
	CheckError(err, "getFileInfo: file open has fucked up")
	defer ltn.Close()
	ltnInfo, _ := ltn.Stat()
	filename = ltnInfo.Name()
	size = strconv.FormatInt(ltnInfo.Size(), 10)
	return
}

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
		CheckError(err, "Unable to create newdumpOutFile2")
		n3, err := g.Write(pic.Picture)
		CheckError(err, "newdumpOutfile2 Write has fucked up")
		fmt.Println(n3, "bytes written successfully")
	}
	outfile22 := resizeImage(newdumpOutFile2, newdumpOutFileThumb)
	return artist, album, title, genre, outfile22
}

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

/////////////////////////////////////////////////////////////////////////////////////////////

func InsAlbumID(alb string) {
	uuid, _ := UUID()
	Albid := map[string]string{"album" : alb, "albumID": uuid}
	AmpgoInsertOne("tempdb2", "albumid", Albid)
	return
}

func InsArtistID(art string) {
	uuid, _ := UUID()
	Artid := map[string]string{"artist" : art, "artistID" : uuid}
	AmpgoInsertOne("tempdb2", "artistid", Artid)
	return
}

func GetTitleOffsetAll() (Main2SL []map[string]string) {
	filter := bson.D{}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb1").Collection("meta1")
	cur, err := collection.Find(context.Background(), filter)
	if err != nil { log.Fatal(err) }
	if err = cur.All(context.Background(), &Main2SL); err != nil {
		log.Fatal(err)
	}
	return
}

func gArtistInfo(Art string) map[string]string {
	filter := bson.M{"artist": Art}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("artistid")
	var ArtInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&ArtInfo)
	if err != nil { log.Fatal(err) }
	return ArtInfo
}

func gAlbumInfo(Alb string) map[string]string {
	filter := bson.M{"album": Alb}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("albumid")
	var AlbInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&AlbInfo)
	if err != nil { log.Fatal(err) }
	return AlbInfo
}

func UpdateMainDB(m2 map[string]string) (Doko Tagmap) {
	artID := gArtistInfo(m2["artist"])
	albID := gAlbumInfo(m2["album"])
	Doko.Dirpath = m2["dirpath"]
	Doko.Filename = m2["filename"]
	Doko.Extension = m2["extension"]
	Doko.FileID = m2["fileID"]
	Doko.Filesize = m2["filesize"]
	Doko.Artist = m2["artist"]
	Doko.ArtistID = artID["artistID"]
	Doko.Album = m2["album"]
	Doko.AlbumID = albID["albumID"]
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
	return
}

func GDistArtist2() (dArtAll []map[string]string) {
	dArtist := AmpgoDistinct("maindb", "maindb", "artist")
	for _, art := range dArtist {
		dArt := AmpgoFindOne("maindb", "maindb", "artist", art)
		dArtAll = append(dArtAll, dArt)
	}
	return dArtAll
}

func unique(arr []string) []string {
    occured := map[string]bool{}
    result := []string{}
    for e := range arr {
        if occured[arr[e]] != true {
            occured[arr[e]] = true
            result = append(result, arr[e])
        }
    }
    return result
}

func create_just_albumID_list(alist []map[string]string) (just_albumID_list []string) {
	for _, albID := range alist {
		just_albumID_list = append(just_albumID_list, albID["albumID"])
	}
	return
}

func get_albums_for_artist(fullalblist []map[string]string) (final_alblist []map[string]string) {
	just_albumID_list := create_just_albumID_list(fullalblist)
	//remove double albumid entries
	unique_items := unique(just_albumID_list)
	for _, uitem := range unique_items {
		albINFO := AmpgoFindOne("maindb", "maindb", "albumID", uitem)
		final_alblist = append(final_alblist, albINFO)
	}
	return 
}

func ArtPipline(artmap map[string]string, page int, idx int) (MyArView ArtVieW2) {
	dirtyalblist := AmpgoFind("maindb","maindb", "artistID", artmap["artistID"]) //[]map[string]string
	results2 := get_albums_for_artist(dirtyalblist)
	MyArView.Artist = artmap["artist"]
	MyArView.ArtistID = artmap["artistID"]
	MyArView.Albums = results2
	MyArView.Page = strconv.Itoa(page)
	MyArView.Index = strconv.Itoa(idx)
	return
}

func InsArtPipeline(AV1 ArtVieW2) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "artistview", "artistview", &AV1)
	CheckError(err2, "artistview insertion has failed")
}

func GDistAlbum() (DAlbAll []map[string]string) {
	DAlbumID := AmpgoDistinct("maindb", "maindb", "albumID")
	for _, albID := range DAlbumID {
		DAlb := AmpgoFindOne("maindb", "maindb", "albumID", albID)
		DAlbAll = append(DAlbAll, DAlb)
	}
	return
}

func get_songs_for_album(fullsonglist []map[string]string) (final_songlist []map[string]string) {
	//a list of just albumid's
	var just_songID_list []string
	for _, song := range fullsonglist {
		songID := song["artisID"]
		just_songID_list = append(just_songID_list, songID)
	}

	//remove double songID entries
	unique_items := unique(just_songID_list)
	for _, uitem := range unique_items {
		songINFO := AmpgoFindOne("maindb", "maindb", "songID", uitem)
		final_songlist = append(final_songlist, songINFO)
	}
	return final_songlist
}

// // // AlbPipeline exported
func AlbPipeline(DAlb map[string]string, page int, idx int) (MyAlbview AlbVieW2) {
	_artist := DAlb["artist"]
	_artistID := DAlb["artistID"]
	_album := DAlb["album"]
	_albumID := DAlb["albumID"]
	_picPath := DAlb["picpath"]
	dirtysonglist := AmpgoFind("maindb","maindb", "albumID", _albumID)
	results := get_songs_for_album(dirtysonglist)
	songcount := len(results)
	MyAlbview.Artist = _artist
	MyAlbview.ArtistID = _artistID
	MyAlbview.Album = _album
	MyAlbview.AlbumID = _albumID
	MyAlbview.NumSongs = strconv.Itoa(songcount)
	MyAlbview.PicPath = _picPath
	MyAlbview.Songs = results
	MyAlbview.AlbumPage = strconv.Itoa(page)
	MyAlbview.Idx = strconv.Itoa(idx)
	return 
}

// //InsAlbViewID exported
func InsAlbViewID(MyAlbview AlbVieW2) {
	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	CheckError(err, "Connections has failed")
    defer Close(client, ctx, cancel)
	_, err = InsertOne(client, ctx, "albumview", "albumview", &MyAlbview)
	CheckError(err, "albumview insertion has fucked up")
	return
}