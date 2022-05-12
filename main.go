package main

import (
	"context"
	_ "embed"
	"go/format"
	"log"
	"os"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

//go:embed pkg/bundles/generator.arraiz
var generator []byte

var (
	nameFlagValue    string
	packageFlagValue string
	outFlagValue     string
)

func main() {
	app := &cli.App{
		Name:      "Adobo",
		Usage:     "Generate types in Go from a jsonschema",
		UsageText: "adobo [flags...] filepath",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "Name for your type. If not provided, it will generate one from title annotation. If title is not provided, it uses the default value",
				Value:       "Schema",
				Destination: &nameFlagValue,
			},
			&cli.StringFlag{
				Name:        "package",
				Aliases:     []string{"p"},
				Usage:       "The package name for the generated code",
				Value:       "schema",
				Destination: &packageFlagValue,
			},
			&cli.StringFlag{
				Name:        "out",
				Aliases:     []string{"o"},
				Usage:       "Specify location of file output. If not specified, it will output to stdout.",
				Destination: &outFlagValue,
			},
		},

		Action: generate,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generate(c *cli.Context) error {
	config := createConfig(packageFlagValue, nameFlagValue)
	filename := c.Args().First()
	return generateCode(afero.NewOsFs(), config, rel.NewString([]rune(filename)).(rel.String), outFlagValue)
}

func generateCode(fs afero.Fs, config rel.Tuple, schemaPath rel.String, outputPath string) error {
	generatorFn, err := syntax.EvaluateBundle(generator)
	if err != nil {
		// impossible, unless the arrai script is screwed, which is possible
		panic(err)
	}

	// this is just a way to call the function generate_from_json_path.arrai
	// with the argument config and filename.
	code, err := rel.NewCallExprCurry(
		generatorFn.Source(),
		generatorFn, config, schemaPath,
	).Eval(arraictx.InitRunCtx(context.Background()), rel.EmptyScope)

	if err != nil {
		// TODO: make this better, arrai errors are scary
		return err
	}

	return outputValue(fs, code.String(), outputPath)
}

func outputValue(fs afero.Fs, content, path string) error {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return nil
	}
	if path == "" {
		_, err := os.Stdout.Write(formatted)
		return err
	}

	f, err := fs.Create(path)
	if err != nil {
		if err == afero.ErrFileExists {
			f, err = fs.Open(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()

	_, err = f.Write(formatted)
	return err
}

func createConfig(packageName, name string) rel.Tuple {
	// create tuple config
	// (
	//		package: packageFlagValue,
	//		name:    nameFlagValue,
	// )
	return rel.NewTuple(
		rel.NewAttr("package", rel.NewString([]rune(packageName))),
		rel.NewAttr("name", rel.NewString([]rune(name))),
	)
}
