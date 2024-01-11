# AI-service
Simple API integration service that serves as a gateway to use other API services.
## Circle CI/CD status
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/RSO-project-Prepih/AI-service/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/RSO-project-Prepih/AI-service/tree/main)

## GO status 
[![Go Report Card](https://goreportcard.com/badge/github.com/RSO-project-Prepih/AI-service)](https://goreportcard.com/report/github.com/RSO-project-Prepih/AI-service)

## Description
Example of using the end points of the AI service:
1. To use the famous places api you need to send a GET request to the following endpoint:
``` 
http://localhost:8080/famous-places
```
2. To use the color enhancement api you need to send a POST request to the following endpoint:
```
http://localhost:8080/enhance-color?userID=<your-user-id>&imageID=<your-image-id>
```
3. To get all the images form the database that where edit you need to send a GET request to the following endpoint:
```
http://localhost:8080/image-processing
```

## Swagger openapi documentation
To see the swagger documentation you need to run the application and go to the following endpoint:
```
http://localhost:8080/openapi/index.html
```
