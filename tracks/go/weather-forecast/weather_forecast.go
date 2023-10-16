// Package weather contains functionality to retrieve weather data.
package weather

// CurrentCondition is a string representation of the current weather.
var CurrentCondition string

// CurrentLocation is the city name used in weather forecasting.
var CurrentLocation string

// Forecast returns the weather forecast for a city.
func Forecast(city, condition string) string {
	CurrentLocation, CurrentCondition = city, condition
	return CurrentLocation + " - current weather condition: " + CurrentCondition
}
