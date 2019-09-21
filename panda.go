package main

import (
	"fmt"
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

	dprFactor := 1.0

	if dprHeader := r.Header.Get("DPR"); len(dprHeader) > 0 {
		if dprFactor, err = strconv.ParseFloat(dprHeader, 64); err != nil {
			dprFactor = 1.0
		}
	}

	width = int(float64(width) * dprFactor)
	height = int(float64(height) * dprFactor)

	key := fmt.Sprintf("cache/%v/%v.jpg", width, height)

	var buf []byte
	if checkCache(width, height) {
		buf, err = downloadPanda(key)

		if err != nil {
			respondWithError(err, w, http.StatusBadRequest)
			return
		}

		w.Header().Set("X-Cache-Info", "HIT")
	} else {
		pandaNum := rand.Int() % len(pandas)
		pandaFile, err := downloadPanda(pandas[pandaNum])
		if err != nil {
			respondWithError(err, w, http.StatusInternalServerError)
			return
		}

		img := bimg.NewImage(pandaFile)
		img.SmartCrop(width, height)
		buf = img.Image()

		w.Header().Set("X-Cache-Info", "MISS")

		go func() {
			err := uploadPanda(key, buf)
			if err != nil {
				log.Printf("Error uploading cached image: %v", err)
			}
		}()
	}

	writeResponse(buf, w)
}

func writeResponse(buf []byte, w http.ResponseWriter) (err error) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=15552000")
	w.Header().Set("Vary", "DPR, Accept")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(buf)
	return
}

func checkCache(width, height int) bool {
	key := fmt.Sprintf("cache/%v/%v.jpg", width, height)

	info, err := getPandaInfo(key)

	if err != nil || info == nil {
		return false
	}
	return true
}
