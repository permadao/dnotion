package guild

import (
	"fmt"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"testing"
	"time"
)

func TestGrade(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	c = config.New("config_temp2")
	d = db.New(c)
	translator, err := d.GetTranslator("d8c270f68a8f44aaa6b24e17c927df2b", nil)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(translator)
	start := time.Now().AddDate(0, 0, -28).Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	err = g.GenGrade("e8d79c55c0394cba83664f3e5737b0bd", "d8c270f68a8f44aaa6b24e17c927df2b", start, end)
	fmt.Println(err)
}
