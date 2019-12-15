package model

//App struct for getting all api
type App struct {
	Version string
	API     map[string][]API
}

//API struct for getting different api scenarios
type API struct {
	Description string
	Endpoint    string
	Scenarios   []Scenario
}

//Scenario consists of Request and the response for them
type Scenario struct {
	Request  Request
	Response Response
}

//Request contains the request params
type Request struct {
	Header  map[string][]string
	Query   map[string][]string
	Payload Payload
}

//Response contains the response components
type Response struct {
	Header     map[string][]string
	Payload    Payload
	StatusCode int
	Delay      int
}

//Payload contains data about payload
type Payload struct {
	Type string
	Data string
}

//New returns an instance of app
func New() App {
	return App{
		API: make(map[string][]API),
	}
}
