package internal

type WeatherPayload struct {
	City        string `structs:"City"`
	Temperature string `structs:"Temperature"`
	Humidity    string `structs:"Humidity"`
	WindSpeed   string `structs:"WindSpeed"`
	Condition   string `structs:"Condition"`
}
