# Proxy Example
A Simple Rest Proxy Server written in golang.

you can find example curl requests in curl_commands.txt

Request:
```Example Request
{
"method": "GET",
"url": "http://google.com",
"headers": {
"Authentication": "Basic bG9naW46cGFzc3dvcmQ=",
....
  }
}
  ```
Response:
```Example Request
{
"id": <generated unique id>,
"status": <HTTP status of 3rd-party service response>,
"headers": {
<headers array from 3rd-party service response>
},
"length": <content length of 3rd-party service response>
}
  ```
  
Example curl command to run proxy:  
```Example curl command

curl -X GET \
  -H "Content-type: application/json" \
  -H "Accept: application/json" \
  -d '{"method":"GET HTTP 1.1", "url": "http://www.yahoo.com"}' \
  "http://localhost:8080"
  
  ```
