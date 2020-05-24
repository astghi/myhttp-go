package main

import (
    "testing"
//    "fmt"
    "errors"
)

func TestFormatUrl(t *testing.T) {
    validUrl := FormatUrl("http://adjust.com")
    formattedIncompleteUrl := FormatUrl("adjust.com")
    if validUrl != formattedIncompleteUrl {
       t.Errorf("Formatting was not incorrect, got: %s, want: %s.", formattedIncompleteUrl, validUrl)
    }
}

func TestGetEncodedResponse(t *testing.T) {
    validUrl := "http://adjust.com"
    incompleteUrl := "adjust.com"
    invalidUrl := "http://bellakio.com"
    
    result1, _ := GetEncodedResponse(validUrl)
    result2, _ := GetEncodedResponse(incompleteUrl)
    
    if  result1 != result2 {
       t.Errorf("Formatting was incorrect, got: %s, want: %s.", result1, result2)
    }
    
    _, invalidUrlError := GetEncodedResponse(invalidUrl)
    
    expectedInvalidUrlError := errors.New("Failed to request the url: " + invalidUrl)
    if invalidUrlError != expectedInvalidUrlError {
        t.Errorf("Failed to throw an error, got: %s, want: %s.", invalidUrlError, expectedInvalidUrlError)
    }
}

func TestMakeRequests(t *testing.T) {
    testData := []struct {
		urls[] string
		concurrencyLimit int
		expectedResult[] string
	}{
		{
            ["http://google.com", "http://reddit.com/r/notfunny"], 
            2, 
            [
                ("http://google.com c64ad8d1822553328b0cca7c6154ceae", nil)
                ("http://reddit.com/r/notfunny 5a088861dbb7bd828803388d086ec52d", nil)
            ]
        },
        
        {
            ["http://google.com", "http://reddit.com/r/notfunny"], 
            0, 
            [
                ("http://google.com c64ad8d1822553328b0cca7c6154ceae", nil)
                ("http://reddit.com/r/notfunny 5a088861dbb7bd828803388d086ec52d", nil)
            ]
        
        },
        
        {
            ["http://google.com", "https://rrrrreddit.com/r/notfunny"], 
            3, 
            [
                ("http://google.com c64ad8d1822553328b0cca7c6154ceae", nil)
                ("", errors.New("Failed to read the responce from url: https://rrrrreddit.com/r/notfunny"))
            ]
        
        },
	}
    
	for _, data := range testData {
        actualResult = MakeRequests(data.urls, data.concurrencyLimit)
        
		if actualResult != expectedResult {
			t.Errorf("Request results were incorrect, got: %v, want: %v.", actualResult, data.expectedResult)
		}
	}
}