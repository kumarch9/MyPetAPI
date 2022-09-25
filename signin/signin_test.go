package signin

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminLogin(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdminLogin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserLogin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
