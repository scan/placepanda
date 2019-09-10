package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/h2non/bimg.v1"
)

var pandas [][]byte

func init() {
	rand.Seed(time.Now().Unix())

	baseDir, err := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "pandas"))
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	pandas = make([][]byte, len(files))
	for i := range files {
		pandas[i], err = bimg.Read(filepath.Join(baseDir, files[i].Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func respondWithError(err error, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(status)
	log.Printf("%v\n", err)
	w.Write([]byte(err.Error()))
}

func pandaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		respondWithError(err, w, http.StatusBadRequest)
		return
	}
	if width < 0 {
		width = 0
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		respondWithError(err, w, http.StatusBadRequest)
		return
	}
	if height < 0 {
		height = 0
	}

	pandaNum := rand.Int() % len(pandas)
	pandaFile := pandas[pandaNum]

	img := bimg.NewImage(pandaFile)
	img.SmartCrop(width, height)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=15552000")
	w.WriteHeader(http.StatusOK)
	w.Write(img.Image())
}
