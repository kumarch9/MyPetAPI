package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPetRoute(t *testing.T) {
	type args struct {
		en *gin.Engine
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			PetRoute(tt.args.en)
		})

	}
}
