package skeleton

//App: struct for getting all api
type App struct {
	Post Api
	Get  Api
}

type Api struct {
	Description string
	Endpoint    string
	Scenarios   []Scenario
}

type Scenario struct {
	Request  Request
	Response Response
}

type Request struct {
	Header  map[string]string
	Query   map[string]string
	Payload Payload
}

type Response struct {
	Header     map[string]string
	Data       string
	Type       string
	StatusCode int
}

type Payload struct {
	Type string
	Data string
}
