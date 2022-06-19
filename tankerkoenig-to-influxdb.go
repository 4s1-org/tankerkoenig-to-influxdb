package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type ConfigurationStation struct {
	Id     string `json:"id"`
	Brand  string `json:"brand"`
	City   string `json:"city"`
	Street string `json:"street"`
}

type ConfigurationInfluxDB struct {
	ServerUrl   string `json:"serverUrl"`
	Token       string `json:"token"`
	Bucket      string `json:"bucket"`
	Org         string `json:"org"`
	Measurement string `json:"measurement"`
}

type Configuration struct {
	Stations []ConfigurationStation `json:"stations"`
	InfluxDB ConfigurationInfluxDB  `json:"influxDB"`
}

type Row struct {
	Timestamp    time.Time
	Id           string
	Diesel       float64
	E5           float64
	E10          float64
	DieselChange string
	E5Change     string
	E10Change    string
	Station      ConfigurationStation
}

func main() {
	// Program parameters
	var configFileName string
	flag.StringVar(&configFileName, "c", "", "configuration file")
	flag.Parse()
	sourceFiles := flag.Args()

	// Load configuration
	config, err := loadConfigurationFile(configFileName)
	if err != nil {
		panic(err)
	}

	// create new client with default option for server url authenticate by token
	client := influxdb2.NewClientWithOptions(
		config.InfluxDB.ServerUrl,
		config.InfluxDB.Token,
		influxdb2.DefaultOptions().SetBatchSize(20))
	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteAPI(config.InfluxDB.Org, config.InfluxDB.Bucket)

	for _, filename := range sourceFiles {
		fmt.Println(filename)
		doIt(config, filename, writeAPI)
	}

	writeAPI.Flush()
	client.Close()
}

func loadConfigurationFile(configFileName string) (*Configuration, error) {
	_, err := os.Stat(configFileName)
	if errors.Is(err, os.ErrNotExist) {
		// Config doesn't exists
		return nil, errors.New("Configuration file not found")
	}
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		panic(err)
	}
	configuration := Configuration{}
	err = json.Unmarshal([]byte(file), &configuration)
	if err != nil {
		panic(err)
	}

	return &configuration, nil
}

func doIt(config *Configuration, filename string, writeAPI api.WriteAPI) {
	srcFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	csvReader := csv.NewReader(srcFile)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		stationId := row[1]

		for _, station := range config.Stations {
			if station.Id == stationId {
				detail := Row{
					Timestamp:    convertDate(row[0]),
					Id:           row[1],
					Diesel:       convertCurrency(row[2]),
					E5:           convertCurrency(row[3]),
					E10:          convertCurrency(row[4]),
					DieselChange: row[5],
					E5Change:     row[6],
					E10Change:    row[7],
					Station:      station,
				}

				point := influxdb2.NewPointWithMeasurement(config.InfluxDB.Measurement)

				// 0=keine Änderung, 1=Änderung, 2=Entfernt, 3=Neu

				point.
					AddField("Diesel [EUR]", detail.Diesel).
					AddField("E5 [EUR]", detail.E5).
					AddField("E10 [EUR]", detail.E10).
					AddTag("Marke", detail.Station.Brand).
					AddTag("Stadt", detail.Station.City).
					AddTag("Strasse", detail.Station.Street).
					SetTime(detail.Timestamp)
				writeAPI.WritePoint(point)
			}
		}
	}
}

func convertDate(value string) time.Time {
	timestamp, err := time.Parse(time.RFC3339, value[0:10]+"T"+value[11:22]+":00")
	if err != nil {
		log.Fatal(err)
	}
	return timestamp
}

func convertCurrency(value string) float64 {
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	return result
}
