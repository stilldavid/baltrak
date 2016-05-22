package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type sentence struct {
	Rssi      float64 `json:"rssi"`
	Count     int     `json:"count"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Alt       float64 `json:"alt"`
	Spd       float64 `json:"spd"`
	Tmpi      float64 `json:"tmpint"`
	Tmpo      float64 `json:"tmpext"`
	Press     float64 `json:"press"`
	Volts     float64 `json:"volts"`
	Chase_lat float64 `json:"chase_lat"`
	Chase_lng float64 `json:"chase_lng"`
}

// -21,2,40.037207,-105.263775,1633.70,0.55,26.83,21.50,pascal,4.10
func parseSentence(message []byte) sentence {

	split := strings.Split(string(message), ",")

	if len(split) != 12 {
		log.Println("wrong number of params.")
		return sentence{}
	}

	rssi, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		log.Fatal("can't parse rssi")
	}
	count, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatal("can't parse count")
	}
	lat, err := strconv.ParseFloat(split[2], 64)
	if err != nil {
		log.Fatal("can't parse lat")
	}
	lng, err := strconv.ParseFloat(split[3], 64)
	if err != nil {
		log.Fatal("can't parse lng")
	}
	alt, err := strconv.ParseFloat(split[4], 64)
	if err != nil {
		log.Fatal("can't parse alt")
	}
	spd, err := strconv.ParseFloat(split[5], 64)
	if err != nil {
		log.Fatal("can't parse speed")
	}
	itmp, err := strconv.ParseFloat(split[6], 64)
	if err != nil {
		log.Fatal("can't parse internal temp")
	}
	etmp, err := strconv.ParseFloat(strings.TrimSpace(split[7]), 64)
	if err != nil {
		log.Fatal("can't parse out temp")
	}
	press, err := strconv.ParseFloat(strings.TrimSpace(split[8]), 64)
	if err != nil {
		log.Fatal("can't parse pressure")
	}
	volts, err := strconv.ParseFloat(strings.TrimSpace(split[9]), 64)
	if err != nil {
		log.Fatal("can't parse voltage")
	}
	chase_lat, err := strconv.ParseFloat(strings.TrimSpace(split[10]), 64)
	if err != nil {
		log.Fatal("can't parse voltage")
	}
	chase_lng, err := strconv.ParseFloat(strings.TrimSpace(split[11]), 64)
	if err != nil {
		log.Fatal("can't parse voltage")
	}

	ret := sentence{rssi, count, lat, lng, alt, spd, itmp, etmp, press, volts, chase_lat, chase_lng}

	return ret
}

func writeToFile(towrite []byte) {
	t := time.Now()
	filename := t.Format("2006-1-_2.csv")

	var f *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
	}

	defer f.Close()

	if _, err := f.Write(towrite); err != nil {
		panic(err)
	}
}

func readFromFile() ([]sentence, error) {
	t := time.Now()
	filename := t.Format("2006-1-_2.csv")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []sentence{}, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return []sentence{}, err
	}
	defer f.Close()

	var sentences []sentence

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sentences = append(sentences, parseSentence([]byte(scanner.Text())))
	}

	return sentences, nil
}
