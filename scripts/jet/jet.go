package main

import (
	"slices"

	"github.com/go-jet/jet/v2/generator/metadata"
	postgresGen "github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	skipTables = []string{
		"schema_migrations",
	}
)

const dbFolderName = "db"

func main() {
	err := postgresGen.Generate(
		"./gen",
		postgresGen.DBConnection{ // nolint
			Host:       "localhost",
			Port:       5432,
			User:       "dev",
			Password:   "dev",
			DBName:     "dev",
			SchemaName: "public",
			SslMode:    "disable",
		},
		template.Default(postgres.Dialect).
			UseSchema(func(schemaMetaData metadata.Schema) template.Schema {
				return template.DefaultSchema(schemaMetaData).
					UseModel(template.DefaultModel().
						UseTable(func(table metadata.Table) template.TableModel {
							if slices.Contains(skipTables, table.Name) {
								return template.TableModel{Skip: true} //nolint: exhaustruct
							}
							return template.DefaultTableModel(table).
								UseField(func(columnMetaData metadata.Column) template.TableModelField {
									defaultTableModelField := template.DefaultTableModelField(
										columnMetaData,
									)
									switch defaultTableModelField.Type.Name {
									case "int32":
										defaultTableModelField.Type.Name = "int64"
									case "*int32":
										defaultTableModelField.Type.Name = "*int64"
									}

									return defaultTableModelField
								})
						}),
					)
			}),
	)

	if err != nil {
		panic(err)
	}
}
