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
	"strings"
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

type ConfigurationInfluxDb struct {
	ServerUrl   string `json:"serverUrl"`
	Token       string `json:"token"`
	Bucket      string `json:"bucket"`
	Org         string `json:"org"`
	Measurement string `json:"measurement"`
}

type Configuration struct {
	Stations []ConfigurationStation `json:"stations"`
	InfluxDb ConfigurationInfluxDb  `json:"influxDb"`
}

type Row struct {
	Timestamp    int64
	Id           string
	Diesel       int // in Cent
	E5           int // in Cent
	E10          int // in Cent
	DieselChange bool
	E5Change     bool
	E10Change    bool
	Station      ConfigurationStation
}

func main() {
	fmt.Println("--- Hello World ---")

	// Program parameters
	var configFileName string
	flag.StringVar(&configFileName, "c", "", "configuration file")
	flag.Parse()
	sourceFiles := flag.Args()
	fmt.Printf("Tail: %+q\n", sourceFiles)

	fmt.Println(configFileName)

	// Load configuration
	config, err := loadConfigurationFile(configFileName)
	if err != nil {
		panic(err)
	}

	// create new client with default option for server url authenticate by token
	client := influxdb2.NewClient(config.InfluxDb.ServerUrl, config.InfluxDb.Token)
	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(config.InfluxDb.Org, config.InfluxDb.Bucket)
	defer client.Close()

	for _, filename := range sourceFiles {
		fmt.Println(filename)
		doIt(config, filename, writeAPI)
	}
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

func doIt(config *Configuration, filename string, writeAPI api.WriteAPIBlocking) {
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
					DieselChange: row[5] == "1",
					E5Change:     row[6] == "1",
					E10Change:    row[7] == "1",
					Station:      station,
				}
				fmt.Println(detail.Timestamp)

				p = influxdb2.NewPointWithMeasurement(config.InfluxDb.Measurement).
					AddTag("unit", "temperature").
					AddField("avg", 23.2).
					AddField("max", 45).
					SetTime(time.Now())
				writeAPI.WritePoint(context.Background(), p)
			}
		}
	}
}

func convertDate(value string) int64 {
	timestamp, err := time.Parse(time.RFC3339, value[0:10]+"T"+value[11:22]+":00")
	if err != nil {
		log.Fatal(err)
	}
	return timestamp.UTC().Unix()
}

func convertCurrency(value string) int {
	centStr := strings.Replace(value, ".", "", 1)
	cent, err := strconv.Atoi(centStr)
	if err != nil {
		panic(err)
	}
	return cent
}

// func parseFile(date time.Time, c *Configuration) {
// 	year := fmt.Sprintf("%04d", date.Year())
// 	month := fmt.Sprintf("%02d", date.Month())
// 	day := fmt.Sprintf("%02d", date.Day())

// 	fmt.Printf("%s-%s-%s", year, month, day)
// 	fmt.Println()

// 	srcPath := filepath.Join(c.TankerkoenigDataFolder, "prices", year, month)
// 	srcFilename := fmt.Sprintf("%s-%s-%s-prices.csv", year, month, day)
// 	srcFile2 := filepath.Join(srcPath, srcFilename)

// 	if _, err := os.Stat(srcFile2); errors.Is(err, os.ErrNotExist) {
// 		fmt.Println("Sourcefile \"" + srcFile2 + "\" does not exists.")
// 		return
// 	}

// 	destPath := filepath.Join(c.CsvDataFolder, year, month)
// 	destFilename := fmt.Sprintf("%s-%s-%s.csv", year, month, day)
// 	destFilenameInfluxDb := fmt.Sprintf("%s-%s-%s-influxdb.txt", year, month, day)
// 	destFile2 := filepath.Join(destPath, destFilename)
// 	destFile2InfluxDb := filepath.Join(destPath, destFilenameInfluxDb)

// 	srcFile, err := os.Open(srcFile2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer srcFile.Close()

// 	err2 := os.MkdirAll(destPath, os.ModePerm)
// 	if err2 != nil {
// 		log.Fatal(err2)
// 	}

// 	destFile, err := os.Create(destFile2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer destFile.Close()

// 	destFileInfluxDb, err := os.Create(destFile2InfluxDb)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer destFileInfluxDb.Close()

// 	csvReader := csv.NewReader(srcFile)
// 	csvWriter := csv.NewWriter(destFile)
// 	defer csvWriter.Flush()
// 	txtWriter := bufio.NewWriter(destFileInfluxDb)
// 	defer txtWriter.Flush()
// 	firstRow := true

// 	for {
// 		row, err := csvReader.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if firstRow {
// 			header := []string{"date", "brand", "city", "street", "fueltype", "price"}
// 			csvWriter.Write(header)
// 			firstRow = false
// 			continue
// 		}

// 		timestampParsed, err := time.Parse(time.RFC3339, row[0][0:10]+"T"+row[0][11:22]+":00")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		timestampString := timestampParsed.Format(time.RFC3339)
// 		timestampUnix := timestampParsed.UTC().Unix()
// 		uuid := row[1]

// 		// ToDo: Das geht besser in Go
// 		for i := 0; i < len(c.Stations); i++ {
// 			station := c.Stations[i]
// 			if uuid == station.Id {
// 				parseLine(csvWriter, txtWriter, row, timestampString, timestampUnix, station)
// 			}
// 		}
// 	}

// 	// ToDo: Add Brennpaste
// }

// func parseLine(csvWriter *csv.Writer, txtWriter *bufio.Writer, row []string, timestampString string, timestampUnix int64, station Station) {
// 	diesel := row[2]
// 	e5 := row[3]
// 	e10 := row[4]
// 	dieselchange := row[5]
// 	e5change := row[6]
// 	e10change := row[7]

// 	details := []string{timestampString, station.Brand, station.City, station.Street}

// 	// 0=keine Änderung, 1=Änderung, 2=Entfernt, 3=Neu
// 	// ToDo -1 und 2 beachten

// 	// Hint: Manchmal ist bei einer Änderung (1) der Preis -0.001.
// 	// Keine Ahnung warum, aber die Preise werden ignoriert.

// 	if dieselchange == "1" && diesel[0] != '-' {
// 		csvWriter.Write(append(details, "Diesel", diesel))
// 		// ToDo Das geht besser
// 		txtWriter.WriteString("kraftstoffpreise,marke=" + strings.Replace(station.Brand, " ", "\\ ", -1) + ",ort=" + strings.Replace(station.City, " ", "\\ ", -1) + ",strasse=" + strings.Replace(station.Street, " ", "\\ ", -1) + ",")
// 		txtWriter.WriteString("sorte=Diesel preis=" + diesel + " " + strconv.FormatInt(timestampUnix, 10) + "\n")
// 	}
// 	if e5change == "1" && e5[0] != '-' {
// 		csvWriter.Write(append(details, "E5", e5))
// 		// ToDo Das geht besser
// 		txtWriter.WriteString("kraftstoffpreise,marke=" + strings.Replace(station.Brand, " ", "\\ ", -1) + ",ort=" + strings.Replace(station.City, " ", "\\ ", -1) + ",strasse=" + strings.Replace(station.Street, " ", "\\ ", -1) + ",")
// 		txtWriter.WriteString("sorte=E5 preis=" + e5 + " " + strconv.FormatInt(timestampUnix, 10) + "\n")
// 	}
// 	if e10change == "1" && e10[0] != '-' {
// 		csvWriter.Write(append(details, "E10", e10))
// 		// ToDo Das geht besser
// 		txtWriter.WriteString("kraftstoffpreise,marke=" + strings.Replace(station.Brand, " ", "\\ ", -1) + ",ort=" + strings.Replace(station.City, " ", "\\ ", -1) + ",strasse=" + strings.Replace(station.Street, " ", "\\ ", -1) + ",")
// 		txtWriter.WriteString("sorte=E10 preis=" + e10 + " " + strconv.FormatInt(timestampUnix, 10) + "\n")
// 	}
// }
