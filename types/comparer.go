package types

//Comparer Compares Different objects
type Comparer interface {

	//String returns a bool as true if both the string matches
	String(string, string) bool

	//MapString returns a bool as true if both the map matches
	MapString(map[string]string, map[string]string) bool

	//MapStringArr returns a bool as true if both the map matches
	MapStringArr(map[string][]string, map[string][]string) bool
}
