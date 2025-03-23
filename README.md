# Osquery Data Collector Service

**Go Version**: 1.20+  
**License**: MIT  
**Docker Required**

A lightweight Go service that collects system information using osquery, stores it in MySQL, and exposes API endpoints for data access.

## Features

- **System Data Collection via osquery**:
  - Operating System Version
  - Osquery Version Information
  - Installed Applications List
- **Dockerized MySQL Database Storage**
- **Dual API Endpoints**:
  - RESTful JSON API (`/latest_data`)
  - HTML Dashboard (`/latest_data_table`)
- **Automatic Data Collection on Service Startup**

## Prerequisites
- Go 1.20+
- osquery (macOS/Windows)

## Start Database Container
``` bash
docker-compose up -d
```
- docker-compose.yml

``` bash
  version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: fleetmvp
      MYSQL_USER: user
      MYSQL_PASSWORD: password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
volumes:
  mysql_data:
```
## Install Dependencies
  - go mod download
  - Configure Osquery (macOS)
  - Open System Preferences → Security & Privacy → Privacy
  - Enable Full Disk Access for:
  Your terminal application (Terminal/iTerm2)
  osqueryi (usually located at /usr/local/bin/osqueryi)

# Starting The service
  ## Ensure database is running
  docker-compose up -d

  ## Login to Docker container
  
  ``` bash
      docker exec -it <container_name> bash
  ```
  ## Login to SQL CLI
  ``` bash
      mysql -u <username> -p -h 127.0.0.1 osquery_db
      #ENTER PASSWORD HERE
  ```  
  ## Create a database 
  ``` SQL
     CREATE DATABASE osquery_db 
  ```
  ## Start the Go service
  ```bash
    go run main.go
  ```
# API Documentation

  ## GET /latest_data
  
  Returns the most recent system snapshot in JSON format.
  
  Response Schema:
  ``` bash
  {
    "os_name": "string",
    "os_version": "string",
    "osquery_version": "string",
    "installed_apps": ["string"],
    "created_at": "ISO8601 timestamp"
  }
  ```
  Usage:
  ``` bash
      curl http://localhost:8080/latest_data                              
  ```

  Example JSON Response:
  ``` JSON
    {
      "os_name": "macOS",
      "os_version": "14.5",
      "osquery_version": "5.10.2",
      "installed_apps": [
        "Google Chrome",
        "Visual Studio Code",
        "Docker Desktop"
      ],
      "created_at": "2025-03-24T09:15:32-07:00"
    }
  ```
  ## GET /latest_data_table
  Returns the same data as above in a tabular format. 
