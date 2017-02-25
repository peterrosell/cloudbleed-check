package main

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"os"
	"bufio"
)

func TestURLParser(t *testing.T) {
	parseData, err := os.Open("test-data/url-parser/input.txt")
	require.Nil(t, err)
	defer parseData.Close()

	expectedData, err := os.Open("test-data/url-parser/expected.txt")
	require.Nil(t, err)
	defer expectedData.Close()

	sitesToCheck := GetSitesToCheck(parseData)

	scanner := bufio.NewScanner(expectedData)

	for _, site := range sitesToCheck {
		require.True(t, scanner.Scan())
		require.Equal(t, scanner.Text(), site)
	}
}

var collectedSites []string

func TestMatching(t *testing.T) {
	yourSites, err := os.Open("test-data/matching-test/input-your-sites.txt")
	require.Nil(t, err)
	defer yourSites.Close()

	affectedSites, err := os.Open("test-data/matching-test/input-affected-sites.txt")
	require.Nil(t, err)
	defer affectedSites.Close()

	expectedData, err := os.Open("test-data/matching-test/expected.txt")
	require.Nil(t, err)
	defer expectedData.Close()

	CheckForMatchingSites(affectedSites, yourSites, collectSites)

	scanner := bufio.NewScanner(expectedData)

	for _, site := range collectedSites {
		require.True(t, scanner.Scan())
		expected := scanner.Text()
		require.Equal(t, expected, site)
	}
	for scanner.Scan() {
		text := scanner.Text()
		assert.Equal(t,"", text, "Missing '%s' in matched result.", text)
	}

}

func collectSites(site string) {
	collectedSites = append(collectedSites,site)
}