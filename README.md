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
