package servicehandlers

import (
	"net/http"
	"log"
	"reflect"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/jsonschema"
	"time"
	"encoding/csv"
	"bytes"
	"strings"
	"os"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/weissleb/peloton-tableau-connector/service/extractors"
	"github.com/weissleb/peloton-tableau-connector/config"
	"github.com/weissleb/peloton-tableau-connector/service/clients"
)

var (
	httpClient clients.HttpClientInterface
)

func init() {
	httpClient = &http.Client{}
}

var pelotonClient *extractors.PelotonClient

type dummy struct {
	stringField string
	intField    int32
	boolField   bool
	floatField  float32
	dateField   time.Time
}

type TableSchemas struct {
	Tables []TableSchema `json:"tables"`
}

type TableSchema struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Columns     []TableColumn `json:"columns"`
}

type TableColumn struct {
	Name   string `json:"name"`
	GoType string `json:"goType"`
}

var tableTypes = map[string]interface{}{
	"dummy":    dummy{},
	"workouts": extractors.Workout{},
}

// helpers
func simpleSchemaForTable(tableName string, tableType interface{}) TableSchema {
	schema := TableSchema{
		Name: tableName,
	}
	t := reflect.TypeOf(tableType)
	var cols []TableColumn
	for i := 0; i < t.NumField(); i++ {
		cols = append(cols, TableColumn{
			Name:   t.Field(i).Name,
			GoType: t.Field(i).Type.String(),
		})
	}
	schema.Columns = cols
	return schema
}

// methods serving endpoints
func GetSimpleCyclingSchemas(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)
	tableSchemas := &TableSchemas{Tables: []TableSchema{}}

	tableParams, ok := r.URL.Query()["tables"]
	allTables := false
	if !ok || len(tableParams) < 1 {
		allTables = true
	}

	if allTables {
		for k, v := range tableTypes {
			tableSchemas.Tables = append(tableSchemas.Tables, simpleSchemaForTable(k, v))
		}
	} else {
		for _, table := range tableParams {
			tableSchemas.Tables = append(tableSchemas.Tables, simpleSchemaForTable(table, tableTypes[table]))
		}
	}

	tableSchemaJson, err := json.Marshal(tableSchemas)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting schemas: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(tableSchemaJson)
}

// NOT USED OR FULLY IMPLEMENTED
// using jsonschema to write a full schema rather than the simple one above
type JsonSchemas struct {
	Tables []*jsonschema.Schema `json:"tables"`
}

func jsonSchemaForTable(tableName string, tableType interface{}) *jsonschema.Schema {
	return jsonschema.Reflect(tableType)
}

func GetJsonCyclingSchemas(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// return data for requested table
func GetCyclingDataJson(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	table, ok := vars["table"]
	if !ok {
		msg := "{\"error\": \"path parameter is missing, expecting table name in path (/cycling/data/{table})\"}"
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
		log.Println(msg)
		return
	}

	if config.LogLevel == "DEBUG" {
		log.Printf("DEBUG authorization header in workouts request = %s", r.Header.Get("Authorization"))
	}
	pelotonClient, authError, err := getClient(r, w)
	if len(authError) > 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(authError))
		log.Println(authError)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var data [][]string
	if strings.ToLower(table) == "workouts" {
		workouts, err := extractors.ExtractCyclingWorkouts(*pelotonClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting data for workouts: %v", err.Error()), http.StatusInternalServerError)
		}

		data = workouts.GetAsRecords(true)
	} else {
		msg := fmt.Sprintf("{\"error\": \"table '%s' not found\"}", table)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
		log.Println(msg)
		return
	}

	header := data[0]
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	buf := bytes.Buffer{}
	jsonRecords := make([]string, len(data)-1)
	for i, datum := range data[1:] {
		jsonRecord := make([]string, len(datum))
		for i2, s := range datum {
			jsonRecord[i2] = fmt.Sprintf(`"%s": "%s"`, header[i2], s)
		}
		jsonRecords[i] = fmt.Sprintf("{%s}", strings.Join(jsonRecord, ","))
	}
	buf.WriteString(fmt.Sprintf(`{"data":[%s]}`, strings.Join(jsonRecords, ",")))
	buf.WriteTo(w)
}

func getClient(r *http.Request, w http.ResponseWriter) (*extractors.PelotonClient, string, error) {
	var (
		pelotonClient *extractors.PelotonClient
		err error
	)

	if config.RequireAuth {
		userToken := r.Header.Get("Authorization")
		userToken = strings.TrimPrefix(userToken, "Bearer ")
		if config.LogLevel == "DEBUG" {
		log.Printf("DEBUG userToken bytes = %s", userToken)
		}
		tokenBytes, err := base64.StdEncoding.DecodeString(userToken)
		if err != nil {
			msg := "{\"error\": \"could not decode user_token\"}"
			return nil, msg, err
		}
		userToken = string(tokenBytes)
		if config.LogLevel == "DEBUG" {
			log.Print("DEBUG userToken string: " + userToken)
		}
		tmp := strings.SplitN(userToken, ":", 2)
		user := tmp[0]
		token := tmp[1]

		pelotonClient, err = extractors.ExistingPelotonClient(httpClient, user, token)
	} else {
		u, p := os.Getenv("PELO_USER"), os.Getenv("PELO_PASS")
		pelotonClient, err = extractors.NewPelotonClient(httpClient, u, p)
	}
	return pelotonClient, "", err
}

// * * * NOT FULLY IMPLEMENTED * * *
// writes data to response as CSV
// currently not using.
func GetCyclingDataCsv(w http.ResponseWriter, r *http.Request) {
	pelotonClient, err := extractors.NewPelotonClient(httpClient, os.Getenv("PELO_USER"), os.Getenv("PELO_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	workouts, err := extractors.ExtractCyclingWorkouts(*pelotonClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting data for workouts: %v", err.Error()), http.StatusInternalServerError)
	}

	data := workouts.GetAsRecords(true)
	w.Header().Set("Content-Type", "text/csv")
	buf := bytes.Buffer{}
	writer := csv.NewWriter(&buf)
	for _, datum := range data {
		err := writer.Write(datum)
		if err != nil {
			http.Error(w, fmt.Sprintf("error sending csv: %v", err.Error()), http.StatusInternalServerError)
			return
		}
	}
	writer.Flush()
	w.Write(buf.Bytes())
}

func GetCyclingDataSummary(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	pelotonClient, authError, err := getClient(r, w)
	if len(authError) > 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(authError))
		log.Println(authError)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	workoutsSummary, err := extractors.GetCyclingWorkoutsSummary(*pelotonClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting data for workout summary: %v", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, err := json.Marshal(workoutsSummary)
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshalling data for workout summary: %v", err.Error()), http.StatusInternalServerError)
	}
	w.Write(data)
}