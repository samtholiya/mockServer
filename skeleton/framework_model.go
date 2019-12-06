package skeleton

//App struct for getting all api
type App struct {
	Post []API
	Get  []API
}

//API struct for getting different api scenarios
type API struct {
	Description string
	Endpoint    string
	Scenarios   []Scenario
}

//Scenario
type Scenario struct {
	Request  Request
	Response Response
}

type Request struct {
	Header  map[string][]string
	Query   map[string][]string
	Payload Payload
}

type Response struct {
	Header     map[string][]string
	Payload    Payload
	StatusCode int
}

type Payload struct {
	Type string
	Data string
}
