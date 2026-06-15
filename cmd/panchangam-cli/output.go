package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	ppb "github.com/naren-m/panchangam/proto"
	"gopkg.in/yaml.v3"
)

func outputTable(resp *ppb.GetPanchangamResponse, req *ppb.GetPanchangamRequest) error {
	data := resp.PanchangamData

	fmt.Printf("\nPanchangam Data\n")
	fmt.Println("---------------")
	fmt.Printf("Date: %s\n", data.Date)
	fmt.Printf("Location: lat %.4f, lon %.4f\n", req.Latitude, req.Longitude)
	if req.Timezone != "" {
		fmt.Printf("Timezone: %s\n", req.Timezone)
	}

	fmt.Printf("\nSunrise/Sunset Times (UTC):\n")
	fmt.Println("---------------------------")
	fmt.Printf("Sunrise: %s\n", data.SunriseTime)
	fmt.Printf("Sunset:  %s\n", data.SunsetTime)

	sunrise, err := time.Parse("15:04:05", data.SunriseTime)
	if err == nil {
		sunset, err := time.Parse("15:04:05", data.SunsetTime)
		if err == nil {
			dayLength := sunset.Sub(sunrise)
			if dayLength < 0 {
				dayLength += 24 * time.Hour
			}
			fmt.Printf("Day Length: %v\n", dayLength)
		}
	}

	fmt.Printf("\nTraditional Panchangam:\n")
	fmt.Printf("- Tithi: %s\n", data.Tithi)
	fmt.Printf("- Nakshatra: %s\n", data.Nakshatra)
	fmt.Printf("- Yoga: %s\n", data.Yoga)
	fmt.Printf("- Karana: %s\n", data.Karana)

	if len(data.Events) > 0 {
		fmt.Printf("\nEvents:\n")
		for _, event := range data.Events {
			fmt.Printf("- %s at %s", event.Name, event.Time)
			if event.EventType != "" {
				fmt.Printf(" (%s)", event.EventType)
			}
			fmt.Println()
		}
	}

	return nil
}

func outputJSON(resp *ppb.GetPanchangamResponse) error {
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputYAML(resp *ppb.GetPanchangamResponse) error {
	data := resp.PanchangamData
	fmt.Printf("panchangam_data:\n")
	fmt.Printf("  date: %s\n", data.Date)
	fmt.Printf("  sunrise_time: %s\n", data.SunriseTime)
	fmt.Printf("  sunset_time: %s\n", data.SunsetTime)
	fmt.Printf("  tithi: %s\n", data.Tithi)
	fmt.Printf("  nakshatra: %s\n", data.Nakshatra)
	fmt.Printf("  yoga: %s\n", data.Yoga)
	fmt.Printf("  karana: %s\n", data.Karana)

	if len(data.Events) > 0 {
		fmt.Printf("  events:\n")
		for _, event := range data.Events {
			fmt.Printf("    - name: %s\n", event.Name)
			fmt.Printf("      time: %s\n", event.Time)
			if event.EventType != "" {
				fmt.Printf("      type: %s\n", event.EventType)
			}
		}
	}

	return nil
}

func outputTithiTable(tithi *astronomy.TithiInfo, detailed bool) error {
	fmt.Printf("Tithi Information\n")
	fmt.Println("-----------------")
	fmt.Printf("Number: %d\n", tithi.Number)
	fmt.Printf("Name: %s\n", tithi.Name)
	fmt.Printf("Type: %s\n", tithi.Type)
	fmt.Printf("Paksha: %s\n", map[bool]string{true: "Shukla (Waxing)", false: "Krishna (Waning)"}[tithi.IsShukla])

	if detailed {
		fmt.Printf("Start Time: %s\n", tithi.StartTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("End Time: %s\n", tithi.EndTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("Duration: %.2f hours\n", tithi.Duration)
		fmt.Printf("Moon-Sun Difference: %.2f degrees\n", tithi.MoonSunDiff)
	}
	return nil
}

func outputTithiJSON(tithi *astronomy.TithiInfo) error {
	data, err := json.MarshalIndent(tithi, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputTithiYAML(tithi *astronomy.TithiInfo) error {
	data, err := yaml.Marshal(tithi)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}

func outputSunTimesTable(sunTimes *astronomy.SunTimes, location astronomy.Location, date time.Time, detailed bool) error {
	fmt.Printf("Sun Times Information\n")
	fmt.Println("---------------------")
	fmt.Printf("Date: %s\n", date.Format("2006-01-02"))
	fmt.Printf("Location: lat %.4f, lon %.4f\n", location.Latitude, location.Longitude)
	fmt.Printf("Timezone: %s\n", date.Location().String())
	fmt.Printf("Sunrise: %s\n", sunTimes.Sunrise.Format("15:04:05"))
	fmt.Printf("Sunset: %s\n", sunTimes.Sunset.Format("15:04:05"))

	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	if dayLength < 0 {
		dayLength += 24 * time.Hour
	}
	fmt.Printf("Day Length: %v\n", dayLength)

	if detailed {
		solarNoon := sunTimes.Sunrise.Add(dayLength / 2)
		fmt.Printf("Solar Noon: %s\n", solarNoon.Format("15:04:05"))
		fmt.Printf("Night Length: %v\n", 24*time.Hour-dayLength)
	}
	return nil
}

func outputSunTimesJSON(sunTimes *astronomy.SunTimes, location astronomy.Location, date time.Time) error {
	data := map[string]interface{}{
		"date":       date.Format("2006-01-02"),
		"location":   map[string]float64{"latitude": location.Latitude, "longitude": location.Longitude},
		"timezone":   date.Location().String(),
		"sunrise":    sunTimes.Sunrise.Format("15:04:05"),
		"sunset":     sunTimes.Sunset.Format("15:04:05"),
		"day_length": sunTimes.Sunset.Sub(sunTimes.Sunrise).String(),
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

func outputSunTimesYAML(sunTimes *astronomy.SunTimes, location astronomy.Location, date time.Time) error {
	data := map[string]interface{}{
		"date":       date.Format("2006-01-02"),
		"location":   map[string]float64{"latitude": location.Latitude, "longitude": location.Longitude},
		"timezone":   date.Location().String(),
		"sunrise":    sunTimes.Sunrise.Format("15:04:05"),
		"sunset":     sunTimes.Sunset.Format("15:04:05"),
		"day_length": sunTimes.Sunset.Sub(sunTimes.Sunrise).String(),
	}
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Print(string(yamlData))
	return nil
}
