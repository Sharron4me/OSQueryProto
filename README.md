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
      - "3307:3306"
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
  ``` bash
    docker-compose up -d
  ```
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

  # FlowChart
  ![flowchart](https://mermaid.ink/img/pako:eNp1lA1v2jAQhv_KyZOmVmrVjLUdRFMnGqC05Tt00uagyk0csGbs1Ha2Mcp_n2PSEqYRKVGce-7jvXO8RrFMKPLRXJFsAdNWJMBezSMcoZT4KTnNOFlBaIgyEZodw-npFVzjW8EMI5z9odC6hkAKQWPDpJht3a8dFqzXEdqZafIlQpvNlggK4uUb1S_QwoGixFCYkidO9awKDOQLtHFbKangPbR_M1Oat8-Wy9PBgeTcZoBwpQ1dQosYUnIdR9zgYQhfqdK7EreGLh7q55yq1X-tt1anNsTGTqCZZXov940tqGvvW4feFVqbnMPYRmNUQ5jHMdU6zXlF9d1O9T0eKVkQNkRopKLABPRX4bg3q7JFA3q4J-fgmrBXwb3L3MduONCdTkcQUvWTqs9P6sqve3WvxPsOHOAJnTPbIAVtkWSSCfOqaOCAIT6LELeT0OYxsS2Eo7twODiO0NkeNvoHezTF3OCoO-33KvD2ObTyRs5tjDvUxAvoOc_qjMbOPnG7pUjbkbmo7pXJrmuh1WByJWBCdSaFprMqUjRris-9cxhIsw1T2kOX4gEHnFGxv4e0Wdnqm5Ayzv1350Gzc-GdxJJL5f9aMEOrVLuk0vSidlE7RN2_UY26dzBWv6Qace3T00FqXFKe17is1w9RDyVV-9C4TD9WKXSCllQtCUvsP74ufCJkFnRJI-Tb14SoHxGKxMZyJDcyXIkY-Ubl9AQpmc8XyB4BXNtVntlZ0xYj9qBYvn3NiPgu5et68xcriFKY?type=png)
