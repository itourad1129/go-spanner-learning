package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"go-spanner-learning/domain"
	"log"
	"os"
	"reflect"
	"strings"
)

func NewSpannerClient() (*spanner.Client, error) {
	ctx := context.Background()
	database := ""
	if domain.ServerEnv == domain.SERVER_ENV_LOCAL {
		database = "projects/project-local/instances/ins-local/databases/db-local"
		emulatorHost := os.Getenv("SPANNER_EMULATOR_HOST")
		if emulatorHost == "" {
			log.Fatalf("Failed to ENV -> SPANNER_EMULATOR_HOST == NULL")
			return nil, nil
		}
	} else if domain.ServerEnv == domain.SERVER_ENV_TEST_1 {
		database = "projects/project-test1/instances/ins-test1/databases/db-test1"
	} else if domain.ServerEnv == domain.SERVER_ENV_PROD {
		database = "projects/pj/instances/ins/databases/db"
	}
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatalf("Failed to create Spanner client: %v", err)
		return nil, err
	}
	return client, err
}

func GetSpannerColumns(obj interface{}) (string, map[string]string, error) {
	val := reflect.TypeOf(obj)
	columns := make(map[string]string)
	columnNames := []string{}

	// 再帰的にフィールドを処理するヘルパー関数
	var processField func(reflect.Type)
	processField = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("spanner")
			// 埋め込まれた構造体の場合は再帰的にフィールドを探索
			if field.Anonymous {
				processField(field.Type)
			} else if tag != "" {
				columns[field.Name] = tag
				columnNames = append(columnNames, tag)
			}
		}
	}

	processField(val)
	return strings.Join(columnNames, ", "), columns, nil
}
