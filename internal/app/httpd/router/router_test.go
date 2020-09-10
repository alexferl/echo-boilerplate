package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	name := "MyRoute"
	r := &Router{}
	r.Routes = []Route{{Name: name}, {Name: "OtherRoute"}}

	assert.Equal(t, name, r.FindRouteByName(name).Name)
	assert.Nil(t, r.FindRouteByName("does_not_exists"))
}
