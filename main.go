package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/urfave/cli/v3"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "otel_config_validator",
		Usage: "Validate a configuration file against the OpenTelemetry Configuration Schema",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				OnlyOnce: true,
				Usage:    "where to output the configuration as json after validation",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 1 {
				log.Fatalf("Must pass a configuration filename")
			} else {
				config_file_path := cmd.Args().Get(0)
				json_config := validate_configuration(config_file_path)
				if o := cmd.String("output"); o != "" {
					json_to_file(json_config, o)
				}
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func validate_configuration(config_file string) interface{} {
	schema_files, err := filepath.Glob("schema/*.json")
	if err != nil {
		log.Fatalf("can't find schema")
	}

	c := jsonschema.NewCompiler()

	for _, file := range schema_files {
		schema_url, err := url.JoinPath("https://opentelemetry.io/otelconfig/", filepath.Base(file))
		schema, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		if err := c.AddResource(schema_url, bytes.NewReader(schema)); err != nil {
			log.Fatal(err)
		}
	}

	schema, err := c.Compile("schema/opentelemetry_configuration.json")
	if err != nil {
		log.Fatalf("%#v", err)
	}

	v := decodeFile(config_file)

	if err = schema.Validate(v); err != nil {
		if ve, ok := err.(*jsonschema.ValidationError); ok {
			b, _ := json.MarshalIndent(ve.DetailedOutput(), "", "  ")

			fmt.Println(string(b))
		} else {
			log.Fatalf("%#v", err)
		}
	} else {
		fmt.Println("Valid OpenTelemetry Configuration!")
	}

	return v
}

func decodeFile(file string) interface{} {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	ext := filepath.Ext(file)
	if ext == ".yaml" || ext == ".yml" {
		return decodeYAML(file)
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatalf("Invalid json file %s: %#v", file, err)
	}

	return v
}

func decodeYAML(file string) interface{} {
	var v interface{}

	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read configuration file %s: %v", file, err)
	}

	reader := bytes.NewReader(body)
	dec := yaml.NewDecoder(reader)

	if err := dec.Decode(&v); err != nil {
		log.Fatalf("Invalid yaml file %s: %v", file, err)
	}

	return v
}

func json_to_file(j interface{}, out_file string) {
	json_string, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		log.Fatalf("Unable to convert to json: %v", err)
	}

	err = os.WriteFile(out_file, json_string, 0644)
	if err != nil {
		log.Fatalf("Unable to write output file: %v", err)
	}
}
