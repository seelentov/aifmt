package service

import (
	"os"
	"strings"
	"testing"
)

func TestFormatCode(t *testing.T) {
	lg := `packge main
	
	funnc main()
	{
		retun "Hello " + ", " + "World" + "!"	
	}`

	token := os.Getenv("API_KEY")

	fmtd, upds, err := FormatCode(lg, "go", "deepseek/deepseek-chat:free", token, nil)

	if err != nil {
		t.Error(err)
	}

	for _, cs := range []string{
		"func main()",
		"package",
		"Hello",
		"World",
		",",
		"!",
	} {
		if !strings.Contains(fmtd, cs) {
			t.Errorf("Expected %s but not", cs)
			t.Log(fmtd)
			t.Log(upds)
			return
		}
	}
}
