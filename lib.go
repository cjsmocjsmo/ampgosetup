package ampgosetup

import (
	"os"
	"log"
	"fmt"
	"time"
	"context"
	"strings"
	"strconv"
	"io/ioutil"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"github.com/bogem/id3v2"
	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Tagmap exported
type Tagmap struct {
	Dirpath     string `bson:"dirpath"`
	Filename    string `bson:"filename"` 
	Extension   string `bson:"extension"`
	FileID      string `bson:"fileID"`
	Filesize    string `bson:"filesize"`
	Artist      string `bson:"artist"`
	ArtistID    string `bson:"artistID"`
	Album       string `bson:"album"`
	AlbumID     string `bson:"albumID"`
	Title       string `bson:"title"`
	Genre       string `bson:"genre"`
	TitlePage   string `bson:"titlepage"`
	PicID       string `bson:"picID"`
	PicDB       string `bson:"picDB"` 
	PicPath     string `bson:"picPath"`
	PicHttpAddr string `bson:"picHttpAddr"`
	Idx         string `bson:"idx"`
	HttpAddr    string `bson:"httpaddr"`
	Duration	string `bson:"duration"`

	ArtStart    string `bson:"artstart"`
	AlbStart    string `bson:"albstart"`
	TitStart    string `bson:"titstart"`
}

type ArtVieW2 struct {
	Artist   string              `bson:"artist"`
	ArtistID string              `bson:"artistID"`
	Albums   []map[string]string `bson:"albums"`
	AlbCount string              `bson:"albcount"`
	Page     string              `bson:"page"`
	Index    string              `bson:"idx"`
}

type AlbVieW2 struct {
	Artist      string              `bson:"artist"`
	ArtistID    string              `bson:"artistID"`
	Album       string              `bson:"album"`
	AlbumID     string              `bson:"albumID"`
	Songs       []map[string]string `bson:"songs"`
	AlbumPage   string              `bson:"albumpage"`
	NumSongs    string              `bson:"numsongs"`
	PicPath     string              `bson:"picPath"`
	Idx         string              `bson:"idx"`
	PicHttpAddr string              `bson:"picHttpAddr"`
}

type Imageinfomap struct {
	Dirpath       string `bson:"dirpath"`
	Filename      string `bson:"filename"`
	Imagesize     string `bson:"imagesize"`
	ImageHttpAddr string `bson:"imagehttpaddr`
	Index         string `bson:"index"`
	IType         string `bson:"itype"`
	Page          string `bson:"page"`
}

func startLogging() string {
	logtxtfile := os.Getenv("AMPGO_LIB_LOG_PATH")
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(logtxtfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Logging started")
	return "Logging started"
}

var lStart = startLogging()

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
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "AmpgoDistinct: MongoDB connection has failed")
	collection := client.Database(db).Collection(coll)
	DD1, err2 := collection.Distinct(context.TODO(), fieldd, filter, opts)
	CheckError(err2, "AmpgoDistinct: MongoDB distinct album has failed")
	var DAlbum1 []string
	for _, DD := range DD1 {
		zoo := fmt.Sprintf("%s", DD)
		DAlbum1 = append(DAlbum1, zoo)
	}
	return DAlbum1
}

func AmpgoFindOne(db string, coll string, filtertype string, filterstring string) map[string]string {
	filter := bson.M{filtertype: filterstring}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "AmpgoFindOne: MongoDB connection has failed")
	collection := client.Database(db).Collection(coll)
	var results map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&results)
	if err != nil {
		log.Println("AmpgoFindOne: find one has fucked up")
		log.Fatal(err)
	}
	return results
}

func AmpgoFind(dbb string, collb string, filtertype string, filterstring string) []map[string]string {
	filter := bson.D{}
	if (filtertype == "None" && filterstring == "None") {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{filtertype, filterstring}}
	}
	// filter := bson.D{{filtertype, filterstring}}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "AmpgoFind: MongoDB connection has failed")
	coll := client.Database(dbb).Collection(collb)
	cur, err := coll.Find(context.TODO(), filter)
	CheckError(err, "AmpgoFind: ArtPipeline find has failed")
	var results []map[string]string //all albums for artist to include double entries
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Println("AmpgoFind: cur.All has fucked up")
		log.Fatal(err)
	}
	return results
}

func AmpgoInsertOne(db string, coll string, ablob map[string]string) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "AmpgoInsertOne: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "AmpgoInsertOne has failed")
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
	if err != nil {
		log.Println(err)
		log.Println(apath)
		return "None", "None", "None", "None", "None"
	}
	defer tag.Close()
	artist := tag.Artist()
	album := tag.Album()
	title := tag.Title()
	genre := tag.Genre()
	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	newdumpOutFile2 := ""
	newdumpOutFileThumb := ""
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if !ok {
			log.Fatal("DumpArtToFile: Couldn't assert picture frame")
		}
		dumpOutFile2 := os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + ".jpg"
		newdumpOutFile2 = strings.Replace(dumpOutFile2, " ", "_", -1)
		dumpOutFileThumb := os.Getenv("AMPGO_THUMB_PATH") + tag.Artist() + "_-_" + tag.Album() + "_thumb.jpg"
		newdumpOutFileThumb = strings.Replace(dumpOutFileThumb, " ", "_", -1)
		g, err := os.Create(newdumpOutFile2)
		defer g.Close()
		CheckError(err, "DumpArtToFile: Unable to create newdumpOutFile2")
		n3, err := g.Write(pic.Picture)
		CheckError(err, "DumpArtToFile: newdumpOutfile2 Write has fucked up")
		fmt.Println(n3, "DumpArtToFile: bytes written successfully")
	}
	outfile22 := resizeImage(newdumpOutFile2, newdumpOutFileThumb)
	return artist, album, title, genre, outfile22
}

func TaGmap(apath string, apage int, idx int) (TaGmaP Tagmap) {
	artist, album, title, genre, picpath := DumpArtToFile(apath)
	if artist != "None" && album != "None" && title != "None" {
		log.Println(apath)
		page := strconv.Itoa(apage)
		index := strconv.Itoa(idx)
		uuid, _ := UUID()
		pichttpaddr := os.Getenv("AMPGO_SERVER_ADDRESS") + ":" + os.Getenv("AMPGO_SERVER_PORT") + picpath[5:]
		fname, size := getFileInfo(apath)
		httpaddr := os.Getenv("AMPGO_SERVER_ADDRESS") + ":" + os.Getenv("AMPGO_SERVER_PORT") + apath[5:]
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
		TaGmaP.PicHttpAddr = pichttpaddr
		TaGmaP.Idx = index
		TaGmaP.HttpAddr = httpaddr
		TaGmaP.Duration = "None"
		TaGmaP.ArtStart = "None"
		TaGmaP.AlbStart = "None"
		TaGmaP.TitStart = "None"
		client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
		CheckError(err, "TaGmap: Connections has failed")
		defer Close(client, ctx, cancel)
		_, err2 := InsertOne(client, ctx, "tempdb1", "meta1", &TaGmaP)
		CheckError(err2, "TaGmap: Tempdb1 insertion has failed")
		return
	} else {
		os.Remove(apath)
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////////////

func InsAlbumID(alb string) {
	uuid, _ := UUID()
	Albid := map[string]string{"album" : alb, "albumID": uuid}
	AmpgoInsertOne("tempdb2", "albumid", Albid)
	return
}

func startLibLogging() string {
	var logtxtfile string = os.Getenv("AMPGO_LIB_LOG_PATH")
	file, err := os.OpenFile(logtxtfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	return "Logging started"
}

func GetPicForAlbum(alb string) map[string]string {
	// startLibLogging()
	log.Printf("GetPicForAlbum: %s this is alb", alb)
	filter := bson.M{"album": alb}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "GetPicForAlbum: MongoDB connection has failed")
	collection := client.Database("maindb").Collection("maindb")
	var albuminfo Tagmap
	err = collection.FindOne(context.Background(), filter).Decode(&albuminfo)
	if err != nil { log.Fatal(err) }
	log.Printf("GetPicForAlbum: %s this is album", alb)
	log.Printf("GetPicForAlbum: %s this is AlbumID", albuminfo.AlbumID)
	log.Printf("GetPicForAlbum: %s this is PicHttpAddr", albuminfo.PicHttpAddr)

	var albinfo map[string]string = make(map[string]string)
	albinfo["Album"] = alb
	albinfo["AlbumID"] = albuminfo.AlbumID
	albinfo["PicPath"] = albuminfo.PicHttpAddr
	AmpgoInsertOne("tempdb2", "artidpic", albinfo)
	fmt.Println(albinfo)
	log.Println(albinfo)
	return albinfo
}

func InsArtistID(art string) {
	uuid, _ := UUID()
	Artid := map[string]string{"artist" : art, "artistID" : uuid}
	AmpgoInsertOne("tempdb2", "artistid", Artid)
	return
}

func GetTitleOffsetAll() (Main2SL []map[string]string) {
	filter := bson.D{}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "GetTitleOffsetAll: MongoDB connection has failed")
	collection := client.Database("tempdb1").Collection("meta1")
	cur, err := collection.Find(context.Background(), filter)
	if err != nil { log.Fatal(err) }
	if err = cur.All(context.Background(), &Main2SL); err != nil {
		log.Println("GetTitleOffsetAll: cur.All has failed")
		log.Fatal(err)
	}
	return
}

func gArtistInfo(Art string) map[string]string {
	filter := bson.M{"artist": Art}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "gArtistInfo: MongoDB connection has failed")
	collection := client.Database("tempdb2").Collection("artistid")
	var ArtInfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&ArtInfo)
	if err != nil { 
		log.Println("gArtistInfo: has failed")
		log.Fatal(err)
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
		log.Println("gAlbumInfo: has failed")
		log.Fatal(err)
	}
	return AlbInfo
}

func gDurationInfo(filename string) map[string]string {
	log.Println(filename)
	filter := bson.M{"filename": filename}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("durdb").Collection("durdb")
	var durinfo map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&durinfo)
	if err != nil { log.Fatal(err) }
	log.Println(durinfo)
	return durinfo
}

func startsWith(astring string) string {
	if (len(astring) > 3) {
		if (astring[4:] == "The " || astring[4:] == "the ") {
			return strings.ToUpper(astring[4:5])
		}
	} else {
		return strings.ToUpper(astring[:1])
	}
	return "None"
}

func UpdateMainDB(m2 map[string]string) (Doko Tagmap) {
	log.Println(m2["filename"])
	artID := gArtistInfo(m2["artist"])
	log.Println(artID)
	albID := gAlbumInfo(m2["album"])
	log.Println(albID)
	fullpath := m2["dirpath"] + "/" + m2["filename"]
	log.Println(fullpath)
	duration := gDurationInfo(fullpath)
	log.Println(duration)
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
	Doko.PicHttpAddr = m2["picHttpAddr"]
	Doko.HttpAddr = m2["httpaddr"]
	Doko.Duration = duration["duration"]
	Doko.ArtStart = startsWith(m2["artist"])
	Doko.AlbStart = strings.ToUpper(m2["album"][:1])
	Doko.TitStart = strings.ToUpper(m2["title"][:1])
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "UpdateMainDB: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "maindb", "maindb", &Doko)
	CheckError(err2, "UpdateMainDB: maindb insertion has failed")
	
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
	albc := len(results2)
	albcount := strconv.Itoa(albc)
	MyArView.Artist = artmap["artist"]
	MyArView.ArtistID = artmap["artistID"]
	MyArView.Albums = results2
	MyArView.AlbCount = albcount
	MyArView.Page = strconv.Itoa(page)
	MyArView.Index = strconv.Itoa(idx)
	return
}

func InsArtPipeline(AV1 ArtVieW2) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsArtPipeline: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "artistview", "artistview", &AV1)
	CheckError(err2, "InsArtPipeline: artistview insertion has failed")
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
		just_songID_list = append(just_songID_list, song["fileID"])
	}

	//remove double songID entries
	unique_items := unique(just_songID_list)
	for _, uitem := range unique_items {
		songINFO := AmpgoFindOne("maindb", "maindb", "fileID", uitem)
		final_songlist = append(final_songlist, songINFO)
	}
	return final_songlist
}

// // // AlbPipeline exported
func AlbPipeline(DAlb map[string]string, page int, idx int) (MyAlbview AlbVieW2) {
	dirtysonglist := AmpgoFind("maindb","maindb", "albumID", DAlb["albumID"])
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

// //InsAlbViewID exported
func InsAlbViewID(MyAlbview AlbVieW2) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsAlbViewID: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "albumview", "albumview", &MyAlbview)
	CheckError(err2, "InsAlbViewID: AmpgoInsertOne has failed")
	return
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

//RanPics exported
func CreateRandomPicsDB() []Imageinfomap {
	thumb_path := os.Getenv("AMPGO_THUMB_PATH")
	thumb_glob_path := thumb_path + "/*.jpg"
	thumb_glob, err := filepath.Glob(thumb_glob_path)
	CheckError(err, "CreateRandomPicsDB: CheckThumbDB has fucked up")
	var BulkImages []Imageinfomap
	var page int
	for i, v := range thumb_glob {
		if i < 5 {
			page = 1
		} else if i % 5 == 0 {
			page++
		} else {
			page = page + 0
		}
		var iim Imageinfomap = create_image_info_map(i, v, page)
		BulkImages = append(BulkImages, iim)
	}
	return BulkImages
}

func create_image_info_map(i int, afile string, page int) Imageinfomap {
	itype := get_type(afile)
	dir, filename := filepath.Split(afile)
	image_size := get_image_size(afile)
	image_http_path := create_image_http_addr(afile)
	ii := i + 1
	var ImageInfoMap Imageinfomap
	ImageInfoMap.Dirpath = dir
	ImageInfoMap.Filename = filename
	ImageInfoMap.Imagesize = image_size
	ImageInfoMap.ImageHttpAddr = image_http_path
	ImageInfoMap.Index = strconv.Itoa(ii)
	ImageInfoMap.IType = itype
	ImageInfoMap.Page = strconv.Itoa(page)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "create_image_info_map: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "coverart", "coverart", ImageInfoMap)
	CheckError(err2, "create_image_info_map: coverart insertion has failed")
	return ImageInfoMap
}

func get_type(afile string) string {
	if strings.Contains(afile, "thumb") {
		return "thumb"
	} else {
		return "original"
	}
}

func get_image_size(apath string) string {
	fi, err := os.Stat(apath)
	CheckError(err, "get_image_size: os.stat has failed")
	size := fi.Size()
	newsize := int(size)
	return strconv.Itoa(newsize)
}

func create_image_http_addr(aimage string) string {
	return os.Getenv("AMPGO_SERVER_ADDRESS") + ":" + os.Getenv("AMPGO_SERVER_PORT") + aimage[5:]
}

type randDb struct {
	PlayListName string `bson:"playlistname"`
	PlayListID string `bson:"playlistID"`
	PlayListCount string `bson:"playlistcount"`
	Playlist []string `bson:"playlist"`
}

func CreateRandomPlaylistDB() string {
	var ranDBInfo randDb
	var emptylist []string
	uuid, _ := UUID()
	ranDBInfo.PlayListName = "EmptyRandomPlaylist"
	ranDBInfo.PlayListID = uuid
	ranDBInfo.PlayListCount = "0"
	ranDBInfo.Playlist = emptylist

	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "randplaylists", "randplaylists", ranDBInfo)
	CheckError(err2, "CreateRandomPlaylistDB: randplaylists insertion has failed")
	return "Created"
}

func ReadDurationFile(apath string) map[string]string {
	data, err := ioutil.ReadFile(apath)
	CheckError(err, "ReadDurationFile: mp3info read has failed")
	var mp3info map[string]string
	err2 := json.Unmarshal(data, &mp3info)
	CheckError(err2, "ReadDurationFile: json unmarshal has failed")
	return mp3info
}

func InsertDurationInfo(apath string) (string) {
	mp3 := ReadDurationFile(apath)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "InsertDurationInfo: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "durdb", "durdb", mp3)
	CheckError(err2, "InsertDurationInfo: durdb insertion has failed")
	return "durdb Created"
	
}

func ArtistFirst(astring string) string {
	// client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	// CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
	// defer Close(client, ctx, cancel)
	char := startsWith(astring)
	switch {
		case strings.Contains(char, "A"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, erra := InsertOne(client, ctx, "artistalpha", "A", item)
			CheckError(erra, "ArtistFirst: A insertion has failed")
			return "A Created"

		case strings.Contains(char, "B"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errb := InsertOne(client, ctx, "artistalpha", "B", item)
			CheckError(errb, "ArtistFirst: B insertion has failed")
			return "B Created"

		case strings.Contains(char, "C"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errc := InsertOne(client, ctx, "artistalpha", "C", item)
			CheckError(errc, "ArtistFirst: C insertion has failed")
			return "C Created"

		case strings.Contains(char, "D"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errd := InsertOne(client, ctx, "artistalpha", "D", item)
			CheckError(errd, "ArtistFirst: D insertion has failed")
			return "D Created"

		case strings.Contains(char, "E"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, erre := InsertOne(client, ctx, "artistalpha", "E", item)
			CheckError(erre, "ArtistFirst: E insertion has failed")
			return "E Created"
			
		case strings.Contains(char, "F"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errf := InsertOne(client, ctx, "artistalpha", "F", item)
			CheckError(errf, "ArtistFirst: F insertion has failed")
			return "F Created"

		case strings.Contains(char, "G"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errg := InsertOne(client, ctx, "artistalpha", "G", item)
			CheckError(errg, "ArtistFirst: G insertion has failed")
			return "G Created"

		case strings.Contains(char, "H"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errh := InsertOne(client, ctx, "artistalpha", "H", item)
			CheckError(errh, "ArtistFirst: H insertion has failed")
			return "H Created"

		case strings.Contains(char, "I"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, erri := InsertOne(client, ctx, "artistalpha", "I", item)
			CheckError(erri, "ArtistFirst: I insertion has failed")
			return "I Created"

		case strings.Contains(char, "J"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errj := InsertOne(client, ctx, "artistalpha", "J", item)
			CheckError(errj, "ArtistFirst: J insertion has failed")
			return "J Created"

		case strings.Contains(char, "K"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errk := InsertOne(client, ctx, "artistalpha", "K", item)
			CheckError(errk, "ArtistFirst: K insertion has failed")
			return "K Created"

		case strings.Contains(char, "L"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errl := InsertOne(client, ctx, "artistalpha", "L", item)
			CheckError(errl, "ArtistFirst: L insertion has failed")
			return "L Created"

		case strings.Contains(char, "M"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errm := InsertOne(client, ctx, "artistalpha", "M", item)
			CheckError(errm, "ArtistFirst: M insertion has failed")
			return "M Created"

		case strings.Contains(char, "N"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errn := InsertOne(client, ctx, "artistalpha", "N", item)
			CheckError(errn, "ArtistFirst: N insertion has failed")
			return "N Created"

		case strings.Contains(char, "O"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, erro := InsertOne(client, ctx, "artistalpha", "O", item)
			CheckError(erro, "ArtistFirst: O insertion has failed")
			return "O Created"

		case strings.Contains(char, "P"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errp := InsertOne(client, ctx, "artistalpha", "P", item)
			CheckError(errp, "ArtistFirst: P insertion has failed")
			return "P Created"

		case strings.Contains(char, "Q"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errq := InsertOne(client, ctx, "artistalpha", "Q", item)
			CheckError(errq, "ArtistFirst: Q insertion has failed")
			return "Q Created"

		case strings.Contains(char, "R"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errr := InsertOne(client, ctx, "artistalpha", "R", item)
			CheckError(errr, "ArtistFirst: R insertion has failed")
			return "R Created"

		case strings.Contains(char, "S"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errs := InsertOne(client, ctx, "artistalpha", "S", item)
			CheckError(errs, "ArtistFirst: S insertion has failed")
			return "S Created"

		case strings.Contains(char, "T"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errt := InsertOne(client, ctx, "artistalpha", "T", item)
			CheckError(errt, "ArtistFirst: T insertion has failed")
			return "T Created"

		case strings.Contains(char, "U"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, erru := InsertOne(client, ctx, "artistalpha", "U", item)
			CheckError(erru, "ArtistFirst: U insertion has failed")
			return "U Created"

		case strings.Contains(char, "V"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errv := InsertOne(client, ctx, "artistalpha", "V", item)
			CheckError(errv, "ArtistFirst: V insertion has failed")
			return "V Created"

		case strings.Contains(char, "W"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errw := InsertOne(client, ctx, "artistalpha", "W", item)
			CheckError(errw, "ArtistFirst: W insertion has failed")
			return "W Created"

		case strings.Contains(char, "X"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errx := InsertOne(client, ctx, "artistalpha", "X", item)
			CheckError(errx, "ArtistFirst: X insertion has failed")
			return "X Created"

		case strings.Contains(char, "Z"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"artist": astring}
			_, errz := InsertOne(client, ctx, "artistalpha", "Z", item)
			CheckError(errz, "ArtistFirst: Z insertion has failed")
			return "Z Created"
	}
	return "None"
}

func AlbumFirst(astring string) string {
	// client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	// CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
	// defer Close(client, ctx, cancel)
	char := startsWith(astring)
	switch {
		case strings.Contains(char, "A"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, erra := InsertOne(client, ctx, "albumalpha", "A", item)
			CheckError(erra, "AlbumFirst: A insertion has failed")
			return "A Created"

		case strings.Contains(char, "B"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errb := InsertOne(client, ctx, "albumalpha", "B", item)
			CheckError(errb, "AlbumFirst: B insertion has failed")
			return "B Created"

		case strings.Contains(char, "C"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errc := InsertOne(client, ctx, "albumalpha", "C", item)
			CheckError(errc, "AlbumFirst: C insertion has failed")
			return "C Created"

		case strings.Contains(char, "D"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errd := InsertOne(client, ctx, "albumalpha", "D", item)
			CheckError(errd, "AlbumFirst: D insertion has failed")
			return "D Created"

		case strings.Contains(char, "E"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, erre := InsertOne(client, ctx, "albumalpha", "E", item)
			CheckError(erre, "AlbumFirst: E insertion has failed")
			return "E Created"

		case strings.Contains(char, "F"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errf := InsertOne(client, ctx, "albumalpha", "F", item)
			CheckError(errf, "AlbumFirst: F insertion has failed")
			return "F Created"

		case strings.Contains(char, "G"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errg := InsertOne(client, ctx, "albumalpha", "G", item)
			CheckError(errg, "AlbumFirst: G insertion has failed")
			return "G Created"

		case strings.Contains(char, "H"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errh := InsertOne(client, ctx, "albumalpha", "H", item)
			CheckError(errh, "AlbumFirst: H insertion has failed")
			return "H Created"

		case strings.Contains(char, "I"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, erri := InsertOne(client, ctx, "albumalpha", "I", item)
			CheckError(erri, "AlbumFirst: I insertion has failed")
			return "I Created"

		case strings.Contains(char, "J"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errj := InsertOne(client, ctx, "albumalpha", "J", item)
			CheckError(errj, "AlbumFirst: J insertion has failed")
			return "J Created"

		case strings.Contains(char, "K"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errk := InsertOne(client, ctx, "albumalpha", "K", item)
			CheckError(errk, "AlbumFirst: K insertion has failed")
			return "K Created"

		case strings.Contains(char, "L"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errl := InsertOne(client, ctx, "albumalpha", "L", item)
			CheckError(errl, "AlbumFirst: L insertion has failed")
			return "L Created"

		case strings.Contains(char, "M"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errm := InsertOne(client, ctx, "albumalpha", "M", item)
			CheckError(errm, "AlbumFirst: M insertion has failed")
			return "M Created"

		case strings.Contains(char, "N"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errn := InsertOne(client, ctx, "albumalpha", "N", item)
			CheckError(errn, "AlbumFirst: N insertion has failed")
			return "N Created"

		case strings.Contains(char, "O"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, erro := InsertOne(client, ctx, "albumalpha", "O", item)
			CheckError(erro, "AlbumFirst: O insertion has failed")
			return "O Created"

		case strings.Contains(char, "P"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errp := InsertOne(client, ctx, "albumalpha", "P", item)
			CheckError(errp, "AlbumFirst: P insertion has failed")
			return "P Created"

		case strings.Contains(char, "Q"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errq := InsertOne(client, ctx, "albumalpha", "Q", item)
			CheckError(errq, "AlbumFirst: Q insertion has failed")
			return "Q Created"

		case strings.Contains(char, "R"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errr := InsertOne(client, ctx, "albumalpha", "R", item)
			CheckError(errr, "AlbumFirst: R insertion has failed")
			return "R Created"

		case strings.Contains(char, "S"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errs := InsertOne(client, ctx, "albumalpha", "S", item)
			CheckError(errs, "AlbumFirst: S insertion has failed")
			return "S Created"

		case strings.Contains(char, "T"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errt := InsertOne(client, ctx, "albumalpha", "T", item)
			CheckError(errt, "AlbumFirst: T insertion has failed")
			return "T Created"

		case strings.Contains(char, "U"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, erru := InsertOne(client, ctx, "albumalpha", "U", item)
			CheckError(erru, "AlbumFirst: U insertion has failed")
			return "U Created"

		case strings.Contains(char, "V"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errv := InsertOne(client, ctx, "albumalpha", "V", item)
			CheckError(errv, "AlbumFirst: V insertion has failed")
			return "V Created"

		case strings.Contains(char, "W"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errw := InsertOne(client, ctx, "albumalpha", "W", item)
			CheckError(errw, "AlbumFirst: W insertion has failed")
			return "W Created"

		case strings.Contains(char, "X"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)

			var item map[string]string = map[string]string{"album": astring}
			_, errx := InsertOne(client, ctx, "albumalpha", "X", item)
			CheckError(errx, "AlbumFirst: X insertion has failed")
			return "X Created"

		case strings.Contains(char, "Z"):
			client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
			CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
			defer Close(client, ctx, cancel)
			
			var item map[string]string = map[string]string{"album": astring}
			_, errz := InsertOne(client, ctx, "albumalpha", "Z", item)
			CheckError(errz, "AlbumFirst: Z insertion has failed")
			return "Z Created"
	}
	return "None"
}