package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/h2non/bimg.v1"
)

func respondWithError(err error, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(status)
	log.Printf("%v\n", err)
	w.Write([]byte(err.Error()))
}

func pandaHandler(w http.ResponseWriter, r *http.Request) {
	pandas, err := listOfPandas()
	if err != nil {
		respondWithError(err, w, http.StatusInternalServerError)
		return
	}

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
	pandaFile, err := downloadPanda(pandas[pandaNum])
	if err != nil {
		respondWithError(err, w, http.StatusInternalServerError)
		return
	}

	img := bimg.NewImage(pandaFile)
	img.SmartCrop(width, height)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=15552000")
	w.WriteHeader(http.StatusOK)
	w.Write(img.Image())
}
