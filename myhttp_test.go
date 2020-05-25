package main

import (
    "testing"
    "errors"
    "strings"
)

func TestFormatUrl(t *testing.T) {
    validUrl := FormatUrl("http://adjust.com")
    formattedIncompleteUrl := FormatUrl("adjust.com")
    if validUrl != formattedIncompleteUrl {
       t.Errorf("Formatting was not incorrect, got: %s, want: %s.", formattedIncompleteUrl, validUrl)
    }
}

func TestRequestAndGetEncodedResponse(t *testing.T) {
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
    if invalidUrlError.Error() != expectedInvalidUrlError.Error() {
        t.Errorf("Failed to throw an error, got: %s, want: %s.", invalidUrlError, expectedInvalidUrlError)
    }
}

func TestMakeRequests(t *testing.T) {
        
    requestError := errors.New("Failed to request the url: https://rrrrreddit.com/r/notfunny")
    
    validUrls := []string{"http://google.com", "http://reddit.com/r/notfunny"}
    urlsWithError := []string{"https://rrrrreddit.com/r/notfunny"}
    
    validUrlWithParallelResults := MakeRequests(validUrls, 2)
    validUrlWithoutParallelResults := MakeRequests(validUrls, 0)
    invalidUrlResults := MakeRequests(urlsWithError, 0)
        
    for i, result := range validUrlWithParallelResults {
        if ((strings.HasPrefix(result.res, validUrls[i]) != true)  || (result.err != nil)) {
			t.Errorf("[valid urls with limit] Request results were incorrect, urls: %v, results: %v.", validUrls, validUrlWithParallelResults)
        }
    } 
    
    for i, result := range validUrlWithoutParallelResults {
        if ((strings.HasPrefix(result.res, validUrls[i]) != true) || (result.err != nil)) {
			t.Errorf("[valid urls without limit] Request results were incorrect, urls: %v, results: %v.", validUrls, validUrlWithParallelResults)
        }
    }
    
    for _, result := range invalidUrlResults {
        if (result.err.Error() != requestError.Error()) {
			t.Errorf("[invalid urls without limit] Request results were incorrect, got: %s, want: %s.", result.err.Error(), requestError.Error())
        }
    }
}