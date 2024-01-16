package translation

import (
	"fmt"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"testing"
)

func TestGrade(t *testing.T) {
	myconfig := config.New("config_temp")
	d := db.New(myconfig)
	tr := New(d)
	if err := tr.Grade(); err != nil {
		fmt.Println(err)
	}
}
