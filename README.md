# Welcome to myHttp tool!

This tool makes http requests and prints the address of the request along with the MD5 hash of the response. There is a concurrency limit option that if not specified will default to a preset number. 
Please see the below explanation of the approach. Happy reading!

# Approach

The orchestration of the tool is in the **main function**.
Here we read the command arguments from the command line. 
Check if the `-parallel` option is specified and construct the array of urls accordingly.
### Pass the url array to the **MakeRequest**
Thread management is implemented here. Two channels are created: a buffered one(with the concurrency limit) for sending the requests and a non-buffered one to store the results.
For each of the urls we start a go routine to process the url and store the result in the results slice. 
- This calls the **GetEncodedResponse** this will make the get request and read the body and finally get the md5 hax of it.
### Format the urls
Important here is that before we pass the url to the get function we format it. Because in the requirements we allow urls like `google.com`, those will fail to pass the validation and will not be pinged. For that reason we match those kind of urls and add the protocol in front of it.

# Tests

Unit tests are in `myhttp_test.go` for each of the functions. I am sure there is a more beautiful way of supplying the test data to the test. But I didn't manage to have time to go through it in go.

# Improve
Error handling
