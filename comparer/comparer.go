package comparer

//Comparer Compares Different objects
type Comparer interface {

	//CompareString returns a bool as true if both the string matches
	CompareString(string, string) bool

	//CompareMapString returns a bool as true if both the map matches
	CompareMapString(map[string]string, map[string]string) bool
}
