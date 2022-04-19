package ampgosetup

// "context"
// // "crypto/rand"
// // "encoding/hex"
// // "encoding/json"
// "fmt"
// // "github.com/bogem/id3v2"
// // "github.com/disintegration/imaging"
// "go.mongodb.org/mongo-driver/bson"
// // "go.mongodb.org/mongo-driver/mongo"
// // "go.mongodb.org/mongo-driver/mongo/options"
// // "io/ioutil"
// // "log"
// // "os"
// // "path/filepath"
// // "strconv"
// "strings"
// "time"

func ArtistFirst(astring string) string {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "ArtistFirst: Connections has failed")
	defer Close(client, ctx, cancel)

	char := StartsWith(astring)
	switch {
	case char == "A":
		var item map[string]string = map[string]string{"artist": astring}
		_, erra := InsertOne(client, ctx, "artistalpha", "A", item)
		CheckError(erra, "ArtistFirst: A insertion has failed")
		return "A Created"

	case char == "B":
		var item map[string]string = map[string]string{"artist": astring}
		_, errb := InsertOne(client, ctx, "artistalpha", "B", item)
		CheckError(errb, "ArtistFirst: B insertion has failed")
		return "B Created"

	case char == "C":
		var item map[string]string = map[string]string{"artist": astring}
		_, errc := InsertOne(client, ctx, "artistalpha", "C", item)
		CheckError(errc, "ArtistFirst: C insertion has failed")
		return "C Created"

	case char == "D":
		var item map[string]string = map[string]string{"artist": astring}
		_, errd := InsertOne(client, ctx, "artistalpha", "D", item)
		CheckError(errd, "ArtistFirst: D insertion has failed")
		return "D Created"

	case char == "E":
		var item map[string]string = map[string]string{"artist": astring}
		_, erre := InsertOne(client, ctx, "artistalpha", "E", item)
		CheckError(erre, "ArtistFirst: E insertion has failed")
		return "E Created"

	case char == "F":
		var item map[string]string = map[string]string{"artist": astring}
		_, errf := InsertOne(client, ctx, "artistalpha", "F", item)
		CheckError(errf, "ArtistFirst: F insertion has failed")
		return "F Created"

	case char == "G":
		var item map[string]string = map[string]string{"artist": astring}
		_, errg := InsertOne(client, ctx, "artistalpha", "G", item)
		CheckError(errg, "ArtistFirst: G insertion has failed")
		return "G Created"

	case char == "H":
		var item map[string]string = map[string]string{"artist": astring}
		_, errh := InsertOne(client, ctx, "artistalpha", "H", item)
		CheckError(errh, "ArtistFirst: H insertion has failed")
		return "H Created"

	case char == "I":
		var item map[string]string = map[string]string{"artist": astring}
		_, erri := InsertOne(client, ctx, "artistalpha", "I", item)
		CheckError(erri, "ArtistFirst: I insertion has failed")
		return "I Created"

	case char == "J":
		var item map[string]string = map[string]string{"artist": astring}
		_, errj := InsertOne(client, ctx, "artistalpha", "J", item)
		CheckError(errj, "ArtistFirst: J insertion has failed")
		return "J Created"

	case char == "K":
		var item map[string]string = map[string]string{"artist": astring}
		_, errk := InsertOne(client, ctx, "artistalpha", "K", item)
		CheckError(errk, "ArtistFirst: K insertion has failed")
		return "K Created"

	case char == "L":
		var item map[string]string = map[string]string{"artist": astring}
		_, errl := InsertOne(client, ctx, "artistalpha", "L", item)
		CheckError(errl, "ArtistFirst: L insertion has failed")
		return "L Created"

	case char == "M":
		var item map[string]string = map[string]string{"artist": astring}
		_, errm := InsertOne(client, ctx, "artistalpha", "M", item)
		CheckError(errm, "ArtistFirst: M insertion has failed")
		return "M Created"

	case char == "N":
		var item map[string]string = map[string]string{"artist": astring}
		_, errn := InsertOne(client, ctx, "artistalpha", "N", item)
		CheckError(errn, "ArtistFirst: N insertion has failed")
		return "N Created"

	case char == "O":
		var item map[string]string = map[string]string{"artist": astring}
		_, erro := InsertOne(client, ctx, "artistalpha", "O", item)
		CheckError(erro, "ArtistFirst: O insertion has failed")
		return "O Created"

	case char == "P":
		var item map[string]string = map[string]string{"artist": astring}
		_, errp := InsertOne(client, ctx, "artistalpha", "P", item)
		CheckError(errp, "ArtistFirst: P insertion has failed")
		return "P Created"

	case char == "Q":
		var item map[string]string = map[string]string{"artist": astring}
		_, errq := InsertOne(client, ctx, "artistalpha", "Q", item)
		CheckError(errq, "ArtistFirst: Q insertion has failed")
		return "Q Created"

	case char == "R":
		var item map[string]string = map[string]string{"artist": astring}
		_, errr := InsertOne(client, ctx, "artistalpha", "R", item)
		CheckError(errr, "ArtistFirst: R insertion has failed")
		return "R Created"

	case char == "S":
		var item map[string]string = map[string]string{"artist": astring}
		_, errs := InsertOne(client, ctx, "artistalpha", "S", item)
		CheckError(errs, "ArtistFirst: S insertion has failed")
		return "S Created"

	case char == "T":
		var item map[string]string = map[string]string{"artist": astring}
		_, errt := InsertOne(client, ctx, "artistalpha", "T", item)
		CheckError(errt, "ArtistFirst: T insertion has failed")
		return "T Created"

	case char == "U":
		var item map[string]string = map[string]string{"artist": astring}
		_, erru := InsertOne(client, ctx, "artistalpha", "U", item)
		CheckError(erru, "ArtistFirst: U insertion has failed")
		return "U Created"

	case char == "V":
		var item map[string]string = map[string]string{"artist": astring}
		_, errv := InsertOne(client, ctx, "artistalpha", "V", item)
		CheckError(errv, "ArtistFirst: V insertion has failed")
		return "V Created"

	case char == "W":
		var item map[string]string = map[string]string{"artist": astring}
		_, errw := InsertOne(client, ctx, "artistalpha", "W", item)
		CheckError(errw, "ArtistFirst: W insertion has failed")
		return "W Created"

	case char == "X":
		var item map[string]string = map[string]string{"artist": astring}
		_, errx := InsertOne(client, ctx, "artistalpha", "X", item)
		CheckError(errx, "ArtistFirst: X insertion has failed")
		return "X Created"

	case char == "Z":
		var item map[string]string = map[string]string{"artist": astring}
		_, errz := InsertOne(client, ctx, "artistalpha", "Z", item)
		CheckError(errz, "ArtistFirst: Z insertion has failed")
		return "Z Created"
	}
	return "None"
}
