# DO_Aug2025

## Overview
Simple API for imaginary cats with missions and targets

## How to run

Clone this repo:

    git clone https://github.com/vaal12-go/DO_Aug2025.git

Install dependencies:

    go mod download

Run server:

    go run main.go

Can use optional -port parameter to specify port on which server will run
Custom log will appear in custom.log file in the root directory

## Stack
* Golang 1.23.6
* Gin
* SQLite DB (it will create itself on the first run)
* Endpoints listed in the CatApp_PostmanCollection.json 