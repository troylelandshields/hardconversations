package bird

type Bird struct {
	Species  string //`hardc-instruction:"What is the species of the bird, or just say 'bird' if unknown?"`
	Behavior string `hardc-instruction:"must be one of: Singing, Flying, Sitting, Eating"`
	IsWild   bool   `hardc-instruction:"indicates if the bird is known to be wild or not, answer must be one of: true, false"`
}
