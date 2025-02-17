package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestFieldNamesLowerSnakeCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without fields",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{},
				},
			},
		},
		{
			name: "no failures for proto with valid field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "song_name",
							},
							&parser.Field{
								FieldName: "singer",
							},
							&parser.MapField{
								MapName: "song_name2",
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "song_name3",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "failures for proto with invalid field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "song_Name",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   210,
										Line:     14,
										Column:   30,
									},
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName",
										Meta: meta.Meta{
											Pos: meta.Position{
												Filename: "example.proto",
												Offset:   300,
												Line:     21,
												Column:   45,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					`Field name "song_Name" must be underscore_separated_names like "song_name"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   210,
						Line:     14,
						Column:   30,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					`Field name "MapFieldName" must be underscore_separated_names like "map_field_name"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     21,
						Column:   45,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					`Field name "OneofFieldName" must be underscore_separated_names like "oneof_field_name"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewFieldNamesLowerSnakeCaseRule(false)

			got, err := rule.Apply(test.inputProto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got %v, but want %v", got, test.wantFailures)
			}
		})
	}
}

func TestFieldNamesLowerSnakeCaseRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "lower_snake_case.proto",
			wantFilename:  "lower_snake_case.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "lower_snake_case.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewFieldNamesLowerSnakeCaseRule(true)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
