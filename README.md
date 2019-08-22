![Project Logo](files/SecShare.gif)

# Concurrency-Decentralized-Network
It is a decentralized network to share files among peer on the network

## About
 This is a project designed to understand concurrency through The Go Programming Language as a fragment of course CS2433 (Principles Of Programming Language - 2) offered by Dr. Saurabh Joshi for Spring'19.

### Prerequisites
 To run the project, any version of a binary distribution of The Go Programming Language is necessary.
 
### Installing
 Use the provided link to get the go enviroment
 ```
 https://golang.org/dl/
 ```
 
## Running the Tests
 To run every test case in the project execute :
 ```
 go test ./...
 ```
 To run an individual test case in the project execute :
 ```
 go test <TestFileName>.go
 ```
 
### These test cases are required in order to test whether the designed functions opearte precisely or not.

## Deployment
 * To use this project
 ```
 go get -u bitbucket.org/deep_shanky/secshare-decentralisednetwork
 ```
 * To launch the sever, initiate :
 ```
 go run server.go
 ```
 * To generate clients, inititate :
 ```
 go run client.go <client name> <client port> <server IP>
 ```
 * To achieve message passing, file sharing, downloading :
 ```
 follow the options guided in the client terminal
 ```
 
## External packages used : package ui
```
import "github.com/andlabs/ui"
```
```
@ github : https://github.com/andlabs/ui
```

## Authors
* Deeptanshu Sankhwar
* Surgan Jandial
* Jatin Chauhan
* Happy Mahto

## Acknowledgments
 A huge proportion of inspiritation behind this project was Dr Joshi's unconditional efforts in educating and expertise in his area
