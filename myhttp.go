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
)

var DEFAULT_LIMIT = 10

func main() {
    // read arguments from the command line
    // slices the command name from it
    argsWithoutCommand := os.Args[1:]
    
    var urls []string
    concurrencyLimit := 0
    urlStartingIndex := 0
    
    // go through the arguments slice to look for the concurrencyLimit
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
    
    urls = argsWithoutCommand[urlStartingIndex:]
        
    results := MakeRequests(urls, concurrencyLimit)
    
    for _, result := range results {
        if result.err == nil{
           fmt.Println(result.res) 
        } else {
           fmt.Println(result.err)
        }
    }
}

// a struct to hold the result from each request including an index
// which will be used for sorting the results after they come in
type result struct {
	res   string
	err   error
} 

// MakeRequests sends requests in parallel but only up to a certain
// limit, and furthermore it's only parallel up to the amount of CPUs but
// is always concurrent up to the concurrency limit
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

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()
    
	// keen an index and loop through every url we will send a request to
	for _, url := range urls {
        url := url
		// url in a closure
		go func(httpUrl string) {
            
			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			// along with the index so we can sort them later along with
			// any error that might have occoured
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

func GetEncodedResponse(url string) (string, error) {
    formattedUrl := FormatUrl(url)
    
    resp, err := http.Get(formattedUrl)
	if err != nil {
        return "", errors.New("Failed to request the url: " + formattedUrl)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to read the responce from url: " + formattedUrl)
	}

	data := []byte(string(body))
    formattedResponce := fmt.Sprintf("%x", md5.Sum(data))
    
    resultString := formattedUrl + " " + formattedResponce
    return resultString, nil
}

func FormatUrl(address string) string {
    urlWithoutProtocol := "(([a-zA-Z])+(.)([a-zA-Z])+(/([a-zA-Z])+)*){1}"
    formattedUrl := address
    httpProtocol := "http://"
    
    _, err := url.ParseRequestURI(address)
    if err != nil {
        match, _ := regexp.MatchString(urlWithoutProtocol, address)
        if  match == true {
            formattedUrl = httpProtocol + address
        }
    }
    
    return formattedUrl
}
