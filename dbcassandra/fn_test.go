package dbcassandra

import (
	"fmt"
	"testing"
)

func TestBuildWhereQuery(t *testing.T) {
	values := map[string]interface{}{"Name": "Bia", "Age": 45, "Maried": true}
	str := BuildWhereQuery(values)
	fmt.Println(">  ", str)
}
