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

Start Database Container
docker-compose up -d
<details> <summary><strong>View <code>docker-compose.yml</code></strong></summary>
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
</details>
Install Dependencies
go mod download
Configure Osquery (macOS)
Open System Preferences → Security & Privacy → Privacy
Enable Full Disk Access for:
Your terminal application (Terminal/iTerm2)
osqueryi (usually located at /usr/local/bin/osqueryi)
Usage

Start the Service
# Ensure database is running
docker-compose up -d

# Start the Go service
go run main.go
# Endpoint	Description
  GET /latest_data	JSON API Response
  GET /latest_data_table	HTML Dashboard
  Example JSON Response:
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

# API Documentation

  GET /latest_data
  Returns the most recent system snapshot in JSON format.
  Response Schema:
  {
    "os_name": "string",
    "os_version": "string",
    "osquery_version": "string",
    "installed_apps": ["string"],
    "created_at": "ISO8601 timestamp"
  }
s is a minimal viable product (MVP) for demonstration purposes.
