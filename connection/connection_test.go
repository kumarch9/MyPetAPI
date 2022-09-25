package connection

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestConnectionDb(t *testing.T) {
	tests := []struct {
		name string
		want *mongo.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConnectionDb(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectionDb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDBCollection(t *testing.T) {
	type args struct {
		client         *mongo.Client
		databaseName   string
		collectionName string
	}
	tests := []struct {
		name string
		args args
		want *mongo.Collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDBCollection(tt.args.client, tt.args.databaseName, tt.args.collectionName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDBCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
