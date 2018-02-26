package main

import (
	"net/http"
	"encoding/json"
	"strconv"
	"os"
	"bytes"
	"log/syslog"
);


type myData struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type serverConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var config serverConfig


func PrepareRequest(writer http.ResponseWriter, request *http.Request) {

	var smg myData;

	json.NewDecoder(request.Body).Decode(&smg)

	logger, _ := syslog.NewLogger(syslog.LOG_DEBUG, 6)
	str, _ := json.Marshal(smg)
	logger.Println(string(str))

}

func prepareHttp()  {
	http.HandleFunc("/", PrepareRequest);
}

func parseConfig()  {
	file, error := os.Open("listener-config.json");

	if error != nil {
		panic(error)
	}

	json.NewDecoder(file).Decode(&config)

}

func formAddress() (string, error)  {
	var addr bytes.Buffer
	var err error

	if config.Host == "" {
		panic("Config host not found!\n")
	}

	if strconv.Itoa(config.Port) == "" {
		panic("Config port not found!\n")
	}

	addr.WriteString(config.Host)
	addr.WriteString(":");
	addr.WriteString(strconv.Itoa(config.Port))

	var address = addr.String()
	return address, err
}

/**
Start server
 */
func main() {

	parseConfig()

	var address, err = formAddress();
	if err != nil {
		panic(err)
	}

	prepareHttp()

	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err);
	}

}
