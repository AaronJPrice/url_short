# URL Short
Simple URL shortening HTTP service

* Port defined in main.go
* To shorten URL: 
  * POST URL in body to "/"
  * Response will contain short-string in body
* To get original URL: GET "/{short-string from the response to the POST request}"
