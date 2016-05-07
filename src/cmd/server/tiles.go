package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func downloadTile(z, x, y int) error {
	filename := fmt.Sprintf("%d.png", y)
	path := fmt.Sprintf("%d/%d", z, x)

	err := os.MkdirAll("tiles/"+path, 0777)
	if err != nil {
		return err
	}

	output, err := os.Create("tiles/" + path + "/" + filename)
	if err != nil {
		return err
	}
	defer output.Close()

	url := "http://b.tile.openstreetmap.org/" + path + "/" + filename
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	fmt.Println("downloaded tile: ", path, filename)

	return nil
}

func tileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/tiles/"):]

	if _, err := os.Stat("tiles/" + path); err == nil {
		http.ServeFile(w, r, "tiles/"+path)
		return
	}

	split := strings.Split(string(path), "/")

	z, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Println("can't parse z")
		http.NotFound(w, r)
		return
	}
	x, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("can't parse x")
		http.NotFound(w, r)
		return
	}
	y, err := strconv.Atoi(strings.Trim(split[2], ".png"))
	if err != nil {
		fmt.Println("can't parse y", strings.Trim(split[3], ".png"))
		http.NotFound(w, r)
		return
	}

	err = downloadTile(z, x, y)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "tiles/"+path)
}
