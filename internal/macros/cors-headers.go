package macros

var CorsHeaders = map[string]string{
	"Content-Type":                     "application/json",
	"Access-Control-Allow-Origin":      "*", // Or use "http://localhost:3000"
	"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
	"Access-Control-Allow-Headers":     "X-Forwaded-For, Content-Type, Authorization",
	"Access-Control-Allow-Credentials": "true",
}
