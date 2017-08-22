package game

type Word struct {
	Value  string
	Length int
}

func NewWord(str string) *Word {
	word := &Word{}
	word.Value = str
	word.Length = len(str)
	return word
}
