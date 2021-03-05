package main

// MASSIVE inspiration from https://github.com/macvk/dnsleaktest/blob/master/dnsleaktest.go

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var noColorOption bool
var elaspsedOption bool

func init() {
	flag.BoolVar(&noColorOption, "n", false, "Disable color output")
	flag.BoolVar(&elaspsedOption, "e", false, "Show elapsed time for prefetch")
	flag.Parse()
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: dnsleak-go -nc\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type block struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	ASN         string `json:"asn"`
	Type        string `json:"type"`
}

func getRandomNumber(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func getID() int {
	start := time.Now()
	wg := sync.WaitGroup{}
	subDomain := getRandomNumber(1000000, 9999999)
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			url := fmt.Sprintf("https://%d.%d.bash.ws", i, subDomain)
			http.Get(url)
			wg.Done()
		}(i)
	}
	wg.Wait()
	if elaspsedOption {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("\nGathering results took %s.\n\n", yellow(time.Since(start)))
	}
	return subDomain
}

func getResult(id int) ([]block, error) {
	var data []block
	url := fmt.Sprintf("https://bash.ws/dnsleak/test/%d?json", id)
	r, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return data, err
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

func main() {
	color.NoColor = noColorOption
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fail := color.New(color.FgWhite).Add(color.BgRed).SprintFunc()
	pass := color.New(color.FgBlack).Add(color.BgGreen).SprintFunc()
	s := spinner.New(spinner.CharSets[39], 250*time.Millisecond)
	s.Prefix = "Gathering results from https://bash.ws/dnsleak "
	s.Start()
	id := getID()
	s.Stop()
	result, err := getResult(id)
	if err != nil {
		fmt.Println(err)
	}
	dns := 1
	for _, b := range result {
		switch b.Type {
		case "ip":
			fmt.Printf("Your IP Address: %-15s (%s, %s)\n", green(b.IP), b.CountryName, b.ASN)
			fmt.Println()
		case "dns":
			fmt.Printf("DNS [%2d]: %15s (%s, %s)\n", dns, yellow(b.IP), b.CountryName, b.ASN)
			dns++
		case "conclusion":
			fmt.Println()
			if strings.Contains(b.IP, "may be") {
				fmt.Println(fail(b.IP))
			} else {
				fmt.Println(pass(b.IP))
			}
		}
	}
}
