package main

import (
    "crypto/md5"
    "fmt"
	"net/http"
	"io/ioutil"
    "os"
    "strconv"
    "errors"
    "log"
    "regexp"
    "net/url"
    "strings"
)

var DEFAULT_LIMIT = 10

func main() {
    // read arguments from the command line
    // slices the command name from it
    argsWithoutCommand := os.Args[1:]
    
    var urls []string
    concurrencyLimit := 0
    urlStartingIndex := 0
    
    // if there is a "-parallel" option specified it should be the first argument
    // followed by the value of the limit
    // the rest should be the urls
    if (argsWithoutCommand[0] == "-parallel") {  
        
        func() {
            s, err := strconv.Atoi(argsWithoutCommand[1]);
            if(err == nil) {
                urlStartingIndex = 2
                concurrencyLimit = s
                
            } else {
               log.Fatal("Expected a number after -parallel option")
            }
        }()
    } 
    
    // depending if there was a -parallel specified we slice urls from the arguments
    urls = argsWithoutCommand[urlStartingIndex:]
            
    results := MakeRequests(urls, concurrencyLimit)
    
    // each element of the result slice will contain either the string with addrress and md5 and nil for the error
    // or an empty string and the error
    // we print on the console accordingly
    for _, result := range results {
        if result.err == nil{
           fmt.Println(result.res) 
        } else {
           fmt.Println(result.err)
        }
    }
}

// a struct to hold the result from each request 
type result struct {
	res   string
	err   error
} 

// MakeRequests sends requests in parallel to the concurrency level
func MakeRequests(urls []string, limit int) []result {

    var concurrencyLimit int
    
    if (limit == 0) {
        concurrencyLimit = DEFAULT_LIMIT
    } else {
        concurrencyLimit = limit
    }
    
	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	// this channel will not block and collect the http request results
	resultsChan := make(chan *result)

	// we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()
    
	// loop through every url we send a request in goroutine 
    // and keep the result in an instance of a result struct
	for _, url := range urls {
        
		// we create a goroutine to get the responce for each url
		go func(httpUrl string) {
            
			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			//  along with any error that might have occoured
            res, err := GetEncodedResponse(httpUrl)
            result := &result{res, err}

			// now we can send the result struct through the resultsChan
			resultsChan <- result

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

        }(url)
	}

	// make a slice to hold the results we're expecting
	var results []result

	// for next result in the result Channel resultsChan
	// append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we got the results for all the urls we break
		if len(results) == len(urls) {
			break
		}
	}

	// now we're done we return the results
	return results
}

// will do the actual request to the given url
// returns a string with the address and the md5 code and the error in case there is one
func GetEncodedResponse(url string) (string, error) {
    
    // first in case the url by the user was incomplete(google.com) 
    // for the request to not fail we format it
    formattedUrl := FormatUrl(url)
    
    // make the request
    resp, err := http.Get(formattedUrl)
	if err != nil {
        return "", errors.New("Failed to request the url: " + formattedUrl)
	}

    // read the body of the request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to read the responce from url: " + formattedUrl)
	}

    // encode the body with MD5
	data := []byte(string(body))
    formattedResponce := fmt.Sprintf("%x", md5.Sum(data))
    
    // construct the output string 
    resultString := formattedUrl + " " + formattedResponce
    return resultString, nil
}                   

// because we need to support the case of incomplete urls
// we format it to have the protocol name so that it can be requested
func FormatUrl(address string) string {
    // a regexp to match incompleate url like reddit.com/r/notfunny
    urlWithoutProtocol := "(([a-zA-Z])+(.)([a-zA-Z])+(/([a-zA-Z])+)*){1}"
    formattedUrl := strings.TrimSpace(address)
    httpProtocol := "http://"
    
     _, err := url.ParseRequestURI(address)
    if err != nil {
        // if the url is of the regex match then add protocol name from the front
        match, _ := regexp.MatchString(urlWithoutProtocol, address)
        if  match == true {
            formattedUrl = httpProtocol + address
        }
    }
    
    return formattedUrl
}