package octoyaml

/*

Test
	DSL -> yaml
	YAML -> JSON

	001-test.js
	001-test.yaml
	001-test.json
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unicode"

	"github.com/tdewolff/minify"
	minjson "github.com/tdewolff/minify/json"
)

const (
	testDir  = "octodns/octoyaml/parse_tests"
	errorDir = "octodns/octoyaml/error_tests"
)

func init() {
	os.Chdir("../..") // go up a directory so we helpers.js is in a consistent place.
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func TestYamlWrite(t *testing.T) {

	// Read a .JS and make sure we can generate the expected YAML.

	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {

		// Run JS -> conf

		// run all js files that start with a number. Skip others.
		if filepath.Ext(f.Name()) != ".js" || !unicode.IsNumber(rune(f.Name()[0])) {
			continue
		}

		m := minify.New()
		m.AddFunc("json", minjson.Minify)

		t.Run(f.Name(), func(t *testing.T) {
			fname := filepath.Join(testDir, f.Name())
			fmt.Printf("Filename: %v\n", fname)
			content, err := ioutil.ReadFile(fname)
			if err != nil {
				t.Fatal(err)
			}
			conf, err := ExecuteJavascript(string(content), true)
			if err != nil {
				t.Fatal(err)
			}
			basename := f.Name()[:len(f.Name())-3] // Remove ".js"

			// Run conf -> YAML

			actualYAML := bytes.NewBuffer([]byte{})
			dom := conf.FindDomain("example.tld")
			if dom == nil {
				panic(fmt.Sprintf("FILE %s does not mention domain '%s'", f.Name(), "example.tld"))
			}

			err = WriteYaml(
				actualYAML, conf.FindDomain("example.tld").Records, "example.tld")

			// Read expected YAML

			expectedFile := filepath.Join(testDir, basename+".yaml")
			expectedData, err := ioutil.ReadFile(expectedFile)
			if err != nil {
				t.Fatal(err)
			}
			expectedYAML := expectedData

			// Compare YAML and expectedData

			if string(expectedYAML) != actualYAML.String() {
				t.Error("Expected and actual YAML don't match")
				t.Log("Expected:", string(expectedYAML))
				t.Log("Actual  :", actualYAML.String())
			}

		})
	}
}

func TestYamlRead(t *testing.T) {

	// Read a .YAML and make sure it matches the RecordConfig (.JSON).

	minifyFlag := true

	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		// run all yaml files that start with a number. Skip others.
		if filepath.Ext(f.Name()) != ".yaml" || !unicode.IsNumber(rune(f.Name()[0])) {
			continue
		}
		basename := f.Name()[:len(f.Name())-5] // remove ".yaml"

		m := minify.New()
		m.AddFunc("json", minjson.Minify)

		t.Run(f.Name(), func(t *testing.T) {

			// Parse YAML

			content, err := ioutil.ReadFile(filepath.Join(testDir, f.Name()))
			if err != nil {
				if os.IsNotExist(err) {
					content = nil
				} else {
					t.Fatal(err)
				}
			}
			recs, err := ReadYaml(bytes.NewBufferString(string(content)), "example.tld")
			if err != nil {
				t.Fatal(err)
			}

			// YAML -> JSON

			actualJSON, err := json.MarshalIndent(recs, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if minifyFlag {
				actualJSON, err = m.Bytes("json", actualJSON)
			}
			if err != nil {
				t.Fatal(err)
			}

			// Read expected JSON
			expectedFile := filepath.Join(testDir, basename+".json")
			expectedData, err := ioutil.ReadFile(expectedFile)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("SKIPPING")
					t.Log("Skipping (no .json)")
					return
				}
				t.Fatal(err)
			}

			var expectedJSON []byte
			if minifyFlag {
				expectedJSON, err = m.Bytes("json", expectedData)
			} else {
				expectedJSON = expectedData
			}
			if err != nil {
				t.Fatal(err)
			}

			if b, err := JSONBytesEqual(expectedJSON, actualJSON); (!b) || (err != nil) {
				t.Error("Expected and actual json don't match")
				t.Log("Expected:", string(expectedJSON))
				t.Log("Actual  :", string(actualJSON))
			}
		})
	}
}
