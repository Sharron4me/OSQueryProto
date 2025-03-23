package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
	"runtime"
	"html/template" 
	_ "github.com/go-sql-driver/mysql"
)

type OSVersionData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type OsqueryInfoData struct {
	Version string `json:"version"`
}

type InstalledApp struct {
	Name string `json:"name"`
}

type TemplateData struct {
	OSName         string
	OSVersion      string
	OsqueryVersion string
	InstalledApps  []string
	Timestamp      string
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/osquery_db?parseTime=true&charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS system_snapshots (
		id INT AUTO_INCREMENT PRIMARY KEY,
		os_name VARCHAR(255),
		os_version VARCHAR(255),
		osquery_version VARCHAR(255),
		installed_apps JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	collectAndStoreData(db)

	http.HandleFunc("/latest_data", latestDataHandler(db))
	http.HandleFunc("/latest_data_table", latestDataTableHandler(db)) // Add this line
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func collectAndStoreData(db *sql.DB) {
	osVersion, err := queryOSVersion()
	if err != nil {
		log.Printf("Error querying OS version: %v", err)
		return
	}

	osqueryVersion, err := queryOsqueryVersion()
	if err != nil {
		log.Printf("Error querying osquery version: %v", err)
		return
	}

	apps, err := queryInstalledApps()
	if err != nil {
		log.Printf("Error querying installed apps: %v", err)
		return
	}

	appsJSON, err := json.Marshal(apps)
	if err != nil {
		log.Printf("Error marshaling apps: %v", err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO system_snapshots (os_name, os_version, osquery_version, installed_apps) VALUES (?, ?, ?, ?)",
		osVersion.Name, osVersion.Version, osqueryVersion, appsJSON,
	)
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	log.Println("Data stored successfully")
}

func queryOSVersion() (*OSVersionData, error) {
	cmd := exec.Command("osqueryi", "--json", "SELECT name, version FROM os_version;")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("osqueryi failed: %v", err)
	}

	var versions []OSVersionData
	if err := json.Unmarshal(output, &versions); err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no OS data")
	}

	return &versions[0], nil
}

func queryOsqueryVersion() (string, error) {
	cmd := exec.Command("osqueryi", "--json", "SELECT version FROM osquery_info;")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("osqueryi failed: %v", err)
	}

	var info []OsqueryInfoData
	if err := json.Unmarshal(output, &info); err != nil {
		return "", fmt.Errorf("parse error: %v", err)
	}

	if len(info) == 0 {
		return "", fmt.Errorf("no osquery data")
	}

	return info[0].Version, nil
}

func queryInstalledApps() ([]InstalledApp, error) {
	var query string
    switch runtime.GOOS {
    case "darwin":
        query = "SELECT name FROM apps;" 
    case "windows":
        query = "SELECT name FROM programs;"
    default:
        return nil, fmt.Errorf("unsupported OS")
    }

    cmd := exec.Command("osqueryi", "--json", query)	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("osqueryi failed: %v", err)
	}

	var apps []InstalledApp
	if err := json.Unmarshal(output, &apps); err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	return apps, nil
}

func latestDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			osName         string
			osVersion      string
			osqueryVersion string
			appsJSON       []byte
			createdAt      time.Time
		)

		err := db.QueryRow(`
			SELECT 
				os_name, 
				os_version, 
				osquery_version, 
				installed_apps, 
				CONVERT_TZ(created_at, @@session.time_zone, '+00:00') 
			FROM system_snapshots
			ORDER BY created_at DESC
			LIMIT 1
		`).Scan(&osName, &osVersion, &osqueryVersion, &appsJSON, &createdAt)
        if err != nil {
            log.Printf("Database query error: %v", err) // <-- Add this line
            if err == sql.ErrNoRows {
                http.Error(w, "No data found", http.StatusNotFound)
                return
            }
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

		var apps []InstalledApp
		if err := json.Unmarshal(appsJSON, &apps); err != nil {
			http.Error(w, "Error parsing apps", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"os_name":         osName,
			"os_version":      osVersion,
			"osquery_version": osqueryVersion,
			"installed_apps":  apps,
			"created_at":      createdAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func latestDataTableHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			osName         string
			osVersion      string
			osqueryVersion string
			appsJSON       []byte
			createdAt      time.Time
		)

		err := db.QueryRow(`
			SELECT os_name, os_version, osquery_version, installed_apps, created_at
			FROM system_snapshots
			ORDER BY created_at DESC
			LIMIT 1
		`).Scan(&osName, &osVersion, &osqueryVersion, &appsJSON, &createdAt)
		
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No data found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Parse installed apps
		var apps []InstalledApp
		if err := json.Unmarshal(appsJSON, &apps); err != nil {
			http.Error(w, "Error parsing apps", http.StatusInternalServerError)
			return
		}

		// Convert to simple string slice
		appNames := make([]string, len(apps))
		for i, app := range apps {
			appNames[i] = app.Name
		}

		// Create template data
		data := TemplateData{
			OSName:         osName,
			OSVersion:      osVersion,
			OsqueryVersion: osqueryVersion,
			InstalledApps:  appNames,
			Timestamp:      createdAt.Local().Format("Jan 02, 2006 15:04:05 MST"),
		}

		// Render template
		tmpl, err := template.New("table").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, "Error creating template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>System Snapshot</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #f2f2f2; }
        tr:nth-child(even) { background-color: #f9f9f9; }
    </style>
</head>
<body>
    <h2>Latest System Snapshot</h2>
    <table>
        <tr>
            <th>OS Name</th>
            <td>{{.OSName}}</td>
        </tr>
        <tr>
            <th>OS Version</th>
            <td>{{.OSVersion}}</td>
        </tr>
        <tr>
            <th>Osquery Version</th>
            <td>{{.OsqueryVersion}}</td>
        </tr>
        <tr>
            <th>Timestamp</th>
            <td>{{.Timestamp}}</td>
        </tr>
        <tr>
            <th>Installed Apps ({{len .InstalledApps}})</th>
            <td>
                <ul>
                    {{range .InstalledApps}}
                    <li>{{.}}</li>
                    {{end}}
                </ul>
            </td>
        </tr>
    </table>
</body>
</html>
`