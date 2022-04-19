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
func SongFirst() string {

	aAll := AmpgoFind("maindb", "maindb", "titstart", "A")
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "SongFirst: Connections has failed")
	defer Close(client, ctx, cancel)
	for _, a := range aAll {
		_, err = InsertOne(client, ctx, "songalpha", "A", a)
		CheckError(err, "SongFirst: a insertion has failed")
	}
	aa := len(aAll)

	bAll := AmpgoFind("maindb", "maindb", "titstart", "B")
	for _, b := range bAll {
		_, err = InsertOne(client, ctx, "songalpha", "B", b)
		CheckError(err, "SongFirst: b insertion has failed")
	}
	bb := len(bAll)

	cAll := AmpgoFind("maindb", "maindb", "titstart", "C")
	for _, c := range cAll {
		_, err = InsertOne(client, ctx, "songalpha", "C", c)
		CheckError(err, "SongFirst: c insertion has failed")
	}
	cc := len(cAll)

	dAll := AmpgoFind("maindb", "maindb", "titstart", "D")
	for _, d := range dAll {
		_, err = InsertOne(client, ctx, "songalpha", "D", d)
		CheckError(err, "SongFirst: d insertion has failed")
	}
	dd := len(dAll)

	eAll := AmpgoFind("maindb", "maindb", "titstart", "E")
	for _, e := range eAll {
		_, err = InsertOne(client, ctx, "songalpha", "E", e)
		CheckError(err, "SongFirst: e insertion has failed")
	}
	ee := len(eAll)

	fAll := AmpgoFind("maindb", "maindb", "titstart", "F")
	for _, f := range fAll {
		_, err = InsertOne(client, ctx, "songalpha", "F", f)
		CheckError(err, "SongFirst: f insertion has failed")
	}
	ff := len(fAll)

	gAll := AmpgoFind("maindb", "maindb", "titstart", "G")
	gg := len(gAll)
	for _, g := range gAll {
		_, err = InsertOne(client, ctx, "songalpha", "G", g)
		CheckError(err, "SongFirst: g insertion has failed")
	}

	hAll := AmpgoFind("maindb", "maindb", "titstart", "H")
	hh := len(hAll)
	for _, h := range hAll {
		_, err = InsertOne(client, ctx, "songalpha", "H", h)
		CheckError(err, "SongFirst: h insertion has failed")
	}

	iAll := AmpgoFind("maindb", "maindb", "titstart", "I")
	ii := len(iAll)
	for _, i := range iAll {
		_, err = InsertOne(client, ctx, "songalpha", "I", i)
		CheckError(err, "SongFirst: i insertion has failed")
	}

	jAll := AmpgoFind("maindb", "maindb", "titstart", "J")
	jj := len(jAll)
	for _, j := range jAll {
		_, err = InsertOne(client, ctx, "songalpha", "J", j)
		CheckError(err, "SongFirst: j insertion has failed")
	}

	kAll := AmpgoFind("maindb", "maindb", "titstart", "K")
	kk := len(kAll)
	for _, k := range kAll {
		_, err = InsertOne(client, ctx, "songalpha", "K", k)
		CheckError(err, "SongFirst: k insertion has failed")
	}

	lAll := AmpgoFind("maindb", "maindb", "titstart", "L")
	ll := len(lAll)
	for _, l := range lAll {
		_, err = InsertOne(client, ctx, "songalpha", "L", l)
		CheckError(err, "SongFirst: l insertion has failed")
	}

	mAll := AmpgoFind("maindb", "maindb", "titstart", "M")
	mm := len(mAll)
	for _, m := range mAll {
		_, err = InsertOne(client, ctx, "songalpha", "M", m)
		CheckError(err, "SongFirst: m insertion has failed")
	}

	nAll := AmpgoFind("maindb", "maindb", "titstart", "N")
	nn := len(nAll)
	for _, n := range nAll {
		_, err = InsertOne(client, ctx, "songalpha", "N", n)
		CheckError(err, "SongFirst: n insertion has failed")
	}

	oAll := AmpgoFind("maindb", "maindb", "titstart", "O")
	oo := len(oAll)
	for _, o := range oAll {
		_, err = InsertOne(client, ctx, "songalpha", "O", o)
		CheckError(err, "SongFirst: o insertion has failed")
	}

	pAll := AmpgoFind("maindb", "maindb", "titstart", "P")
	pp := len(pAll)
	for _, p := range pAll {
		_, err = InsertOne(client, ctx, "songalpha", "P", p)
		CheckError(err, "SongFirst: p insertion has failed")
	}

	qAll := AmpgoFind("maindb", "maindb", "titstart", "Q")
	qq := len(qAll)
	for _, q := range qAll {
		_, err = InsertOne(client, ctx, "songalpha", "Q", q)
		CheckError(err, "SongFirst: q insertion has failed")
	}

	rAll := AmpgoFind("maindb", "maindb", "titstart", "R")
	rr := len(rAll)
	for _, r := range rAll {
		_, err = InsertOne(client, ctx, "songalpha", "R", r)
		CheckError(err, "SongFirst: r insertion has failed")
	}

	sAll := AmpgoFind("maindb", "maindb", "titstart", "S")
	ss := len(sAll)
	for _, s := range sAll {
		_, err = InsertOne(client, ctx, "songalpha", "S", s)
		CheckError(err, "SongFirst: s insertion has failed")
	}

	tAll := AmpgoFind("maindb", "maindb", "titstart", "T")
	tt := len(tAll)
	for _, t := range tAll {
		_, err = InsertOne(client, ctx, "songalpha", "T", t)
		CheckError(err, "SongFirst: t insertion has failed")
	}

	uAll := AmpgoFind("maindb", "maindb", "titstart", "U")
	uu := len(uAll)
	for _, u := range uAll {
		_, err = InsertOne(client, ctx, "songalpha", "U", u)
		CheckError(err, "SongFirst: u insertion has failed")
	}

	vAll := AmpgoFind("maindb", "maindb", "titstart", "V")
	vv := len(vAll)
	for _, v := range vAll {
		_, err = InsertOne(client, ctx, "songalpha", "V", v)
		CheckError(err, "SongFirst: v insertion has failed")
	}

	wAll := AmpgoFind("maindb", "maindb", "titstart", "W")
	ww := len(wAll)
	for _, w := range wAll {
		_, err = InsertOne(client, ctx, "songalpha", "W", w)
		CheckError(err, "SongFirst: w insertion has failed")
	}

	xAll := AmpgoFind("maindb", "maindb", "titstart", "X")
	xx := len(xAll)
	for _, x := range xAll {
		_, err = InsertOne(client, ctx, "songalpha", "X", x)
		CheckError(err, "SongFirst: x insertion has failed")
	}

	yAll := AmpgoFind("maindb", "maindb", "titstart", "Y")
	yy := len(yAll)
	for _, y := range yAll {
		_, err = InsertOne(client, ctx, "songalpha", "Y", y)
		CheckError(err, "SongFirst: y insertion has failed")
	}

	zAll := AmpgoFind("maindb", "maindb", "titstart", "Z")
	zz := len(zAll)
	for _, z := range zAll {
		_, err = InsertOne(client, ctx, "songalpha", "Z", z)
		CheckError(err, "SongFirst: z insertion has failed")
	}

	t1 := aa + bb + cc + dd + ee + ff + gg + hh + ii + jj + kk + ll + mm
	t2 := nn + oo + pp + qq + rr + ss + tt + uu + vv + ww + xx + yy + zz
	tot := t1 + t2
	total := strconv.Itoa(tot)
	var total2 map[string]string = map[string]string{"total": total}

	client, ctx, cancel, err = Connect("mongodb://db:27017/ampgo")
	CheckError(err, "SongFirst: Connections has failed")
	defer Close(client, ctx, cancel)

	_, err = InsertOne(client, ctx, "songtotal", "total", &total2)
	CheckError(err, "SongFirst: z insertion has failed")

	return "Complete"
}
