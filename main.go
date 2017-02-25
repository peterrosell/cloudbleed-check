package main

import "fmt"
import (
	"bufio"
	"os"
	"io"
	"net/http"
)
import "strings"
import (
	"regexp"
	"sort"
	"net/url"
)

const affectedSiteUrl = "https://raw.githubusercontent.com/pirate/sites-using-cloudflare/master/sorted_unique_cf.txt"

var domainRegex = regexp.MustCompile(`.co\...$`)
var verbose = false
var debug = false

func main() {

	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "v") {
			verbose = true
		}
		if strings.Contains(os.Args[1], "d") {
			debug = true
		}
	}

	responseBody, err := openSiteStream(affectedSiteUrl)
	if err != nil {
		fmt.Printf("Error open url with affected sites. %s", err)
		return
	}
	defer responseBody.Close()

	CheckForMatchingSites(responseBody, os.Stdin, foundMatchOnAffectedSite)
}

func openSiteStream(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Got bad status code: %d", response.StatusCode)
	}
	return response.Body, nil
}

func GetSitesToCheck(reader io.Reader) []string {
	scanner := bufio.NewScanner(reader)

	sites := []string{}

	for scanner.Scan() {
		urlAsString := scanner.Text()
		if urlAsString != "" {
			//fmt.Printf("Got input urlAsString: %s", urlAsString)
			//site := ExtractSite(urlAsString)
			siteUrl, err := url.Parse(urlAsString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse URL: %s\n", urlAsString)
				continue
			}
			var site string
			if siteUrl.Host != "" {
				site = extractMainDomain(siteUrl.Host)
			} else {
				site = urlAsString
			}
			if debug {
				fmt.Fprintf(os.Stderr, "Line: '%s' results in site to check: '%s'\n", urlAsString, site)
			} else if verbose {
				fmt.Fprintf(os.Stderr, "Site to check: '%s'\n", site)
			}
			sites = append(sites, site)
		}
	}
	return sites
}

func sortedSitesToCheck(reader io.Reader) []string {
	sitesToCheck := GetSitesToCheck(reader)
	sort.Strings(sitesToCheck)
	return sitesToCheck
}

func extractMainDomain(host string) string {
	expectedNrOfDotsInDomain := getNrOfDotsInDomain(host)
	nrOfDotsInDomain := 0
	endIndex := len(host)
	startIndex := 0
	for i := endIndex - 1; i >= 0; {
		if host[i] == '.' {
			if nrOfDotsInDomain == expectedNrOfDotsInDomain {
				startIndex = i + 1
				break
			}
			nrOfDotsInDomain++
		} else if host[i] == ':' {
			endIndex = i
		}
		i--
	}
	return host[startIndex:endIndex]
}
func getNrOfDotsInDomain(host string) int {
	if domainRegex.MatchString(host) {
		return 2
	}
	return 1
}

func CheckForMatchingSites(affectedSitesReader io.Reader, sitesToCheckReader io.Reader, matchHandler func(string)) {
	siteIndex := 0
	var getNextAffectedSite bool
	affectedSitesScanner := bufio.NewScanner(affectedSitesReader)

	sitesToCheck := sortedSitesToCheck(sitesToCheckReader)
	fmt.Fprintf(os.Stderr, "Will search the list of sites that might be affected. Your list contains %d sites\n", len(sitesToCheck))

	matchedSite := ""
	loggedMySite := ""

	for ; affectedSitesScanner.Scan() && siteIndex < len(sitesToCheck); {
		affectedSite := affectedSitesScanner.Text()
		if affectedSite == "" {
			continue
		}
		getNextAffectedSite = false

		for ; !getNextAffectedSite && siteIndex < len(sitesToCheck); {
			mySite := sitesToCheck[siteIndex]
			if verbose {
				if loggedMySite != mySite {
					loggedMySite = mySite
					fmt.Fprintf(os.Stderr, "Checking: '%s'\n", loggedMySite)
				}
			}
			if ( mySite == "") {
				siteIndex++
				continue
			}
			if debug {
				fmt.Fprintf(os.Stderr, "Checking your site '%s' against '%s' \n", mySite, affectedSite)
			}
			switch strings.Compare(affectedSite, mySite) {
			case 0:
				if matchedSite != mySite {
					matchedSite = mySite
					matchHandler(mySite)
				}
				siteIndex++
			case 1:
				siteIndex++
			case -1:
				getNextAffectedSite = true
			}
		}
	}
}

func foundMatchOnAffectedSite(site string) {
	fmt.Println(site)
}


