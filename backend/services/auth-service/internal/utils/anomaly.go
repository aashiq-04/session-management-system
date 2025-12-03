package utils

import (
	"math"
	"time"
	"fmt"
)

// AnomalyDetection represents detected security anomalies
type AnomalyDetection struct {
	ImpossibleTravel bool
	NewCountry       bool
	NewDevice        bool
	Severity         string
	Description      string
	Details          map[string]interface{}
}

// DetectAnomalies checks for security anomalies in login attempt
func DetectAnomalies(
	currentLocation Location,
	previousLocation *Location,
	previousLoginTime *time.Time,
	isNewDevice bool,
) *AnomalyDetection {
	anomaly := &AnomalyDetection{
		Details: make(map[string]interface{}),
	}

	// Check for new device
	if isNewDevice {
		anomaly.NewDevice = true
		anomaly.Severity = "medium"
		anomaly.Description = "Login from new unrecognized device"
		return anomaly
	}

	// If no previous location data, cannot detect travel anomaly
	if previousLocation == nil || previousLoginTime == nil {
		return nil
	}

	// Check for new country
	if previousLocation.Country != "" && currentLocation.Country != "" &&
		previousLocation.Country != currentLocation.Country {
		anomaly.NewCountry = true
	}

	// Calculate distance between locations
	distance := calculateDistance(
		previousLocation.Latitude, previousLocation.Longitude,
		currentLocation.Latitude, currentLocation.Longitude,
	)

	// Calculate time difference
	timeDiff := time.Since(*previousLoginTime)
	hoursDiff := timeDiff.Hours()

	// Store details
	anomaly.Details["distance_km"] = distance
	anomaly.Details["time_hours"] = hoursDiff
	anomaly.Details["previous_location"] = previousLocation.City + ", " + previousLocation.Country
	anomaly.Details["current_location"] = currentLocation.City + ", " + currentLocation.Country

	// Detect impossible travel
	// Average commercial flight speed: 800 km/h
	// Add 2 hours buffer for airport time
	maxPossibleSpeed := 800.0 // km/h
	minTimeNeeded := (distance / maxPossibleSpeed)

	if distance > 100 && hoursDiff < minTimeNeeded {
		// Impossible travel detected
		anomaly.ImpossibleTravel = true
		anomaly.Severity = "critical"
		anomaly.Description = "Impossible travel detected: Login from " +
			previousLocation.City + ", " + previousLocation.Country +
			" then " + currentLocation.City + ", " + currentLocation.Country +
			" within " + formatDuration(timeDiff) +
			" - physically impossible to travel " + formatDistance(distance)
		return anomaly
	}

	// Check for new country (less severe than impossible travel)
	if anomaly.NewCountry {
		anomaly.Severity = "high"
		anomaly.Description = "Login from new country detected: " +
			currentLocation.Country + " (previous: " + previousLocation.Country + ")"
		return anomaly
	}

	// No anomalies detected
	return nil
}

// Location represents geographic location
type Location struct {
	Country   string
	City      string
	Latitude  float64
	Longitude float64
}

// calculateDistance calculates distance between two coordinates (Haversine formula)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // Earth's radius in kilometers

	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// formatDuration formats duration to human readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	}
	return fmt.Sprintf("%d minutes", minutes)
}

// formatDistance formats distance to human readable string
func formatDistance(km float64) string {
	if km > 1000 {
		return fmt.Sprintf("%.0f km", km)
	}
	return fmt.Sprintf("%.1f km", km)
}

