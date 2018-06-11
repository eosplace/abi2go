package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dimiro1/banner"
	"github.com/eoscanada/eos-go"
	"github.com/mattn/go-colorable"
)

func saveFile(filename, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func toGoName(v string) string {
	ret := strings.Title(v)

	return ret
}

func convertType(typ string) string {
	switch {
	case typ == "time_point_sec":
		return "string"
	case typ == "checksum256":
		return "string"
	default:
		return typ
	}
}

func writeStructs(abi eos.ABI, pack, prefix string) error {
	filename := prefix + "_structs.go"
	os.Remove(filename)

	content := "// Code generated by abi2go.\n"
	content += "// DO NOT EDIT!\n"
	content += "package " + pack + "\n"

	for _, st := range abi.Structs {
		content += "type " + toGoName(st.Name) + " struct {\n"
		if len(st.Base) > 0 {
			content += st.Base + "\n"
		}
		for _, field := range st.Fields {
			content += toGoName(field.Name) + " " + convertType(field.Type) + " `json:\"" + field.Name + "\"`\n"
		}
		content += "}\n\n"
	}

	return saveFile(filename, content)
}

func writeActions(abi eos.ABI, pack, prefix string) error {
	filename := prefix + "_actions.go"

	content := "// Code generated by abi2go.\n"
	content += "// DO NOT EDIT!\n"
	content += "package " + pack + "\n"
	content += "import (\n  eos \"github.com/eoscanada/eos-go\"\n)\n"

	for _, ac := range abi.Actions {
		content += "func Send" + toGoName(string(ac.Name)) + "(input " + toGoName(ac.Type) + ", account, permission string) error {\n"
		content += "			action := &eos.Action{\n"
		content += "				Account: eos.AccountName(account),\n"
		content += "				Name:    eos.ActionName(\"" + string(ac.Name) + "\"),\n"
		content += "				Authorization: []eos.PermissionLevel{\n"
		content += "					{Actor: eos.AccountName(account), Permission: eos.PermissionName(permission)},\n"
		content += "				},\n"
		content += "				ActionData: eos.NewActionData(input),\n"
		content += "			}\n\n"
		content += "			api := getAPI()\n"
		content += "			api.Debug = true\n"
		content += "			if _, err := api.SignPushActions(action); err != nil {\n"
		content += "				return err\n"
		content += "			}\n\n"
		content += "			return nil\n"
		content += "		}\n"

	}

	return saveFile(filename, content)
}

func writeQueryTables(abi eos.ABI, pack, prefix string) error {
	filename := prefix + "_tables.go"
	content := "package " + pack + "\n"

	return saveFile(filename, content)
}

func writeOutput(abi eos.ABI, pack, prefix string) error {
	err := writeStructs(abi, pack, prefix)
	if err != nil {
		return err
	}

	writeActions(abi, pack, prefix)
	if err != nil {
		return err
	}

	writeQueryTables(abi, pack, prefix)
	if err != nil {
		return err
	}
	return nil
}

const textBanner = `
{{ .AnsiColor.Yellow }}
1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111
1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111
1188888881111888888811118888888111188888811111881111111111888881111118888811111888888811
1188111111111881118811118811111111188111881111881111111118811188131188111881111881111111
1188111111111881118811118111111111188111881111881111111118811188111188111111111881111111
1188888811111881118811118888888111188888811111881111111118888888111188111111111888888111
1188111111111881118811111111188111188111111111881111111118811188111188111111111881111111
1188111111111881118811111111188111188111111111881111111118811188111188111881111881111111
1188888881111888888811118888888111188111111111888888811118811188111118888811111888888811
1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111
1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111
{{ .AnsiColor.Default }}{{ .AnsiBackground.Default }}
ABI2GO
GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
`

func main() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(textBanner))

	var input, prefix, pack string
	flag.StringVar(&input, "input", "", "abi input filename")
	flag.StringVar(&prefix, "prefix", "", "go source output filename prefix")
	flag.StringVar(&pack, "package", "main", "go source package")
	flag.Parse()

	if len(input) == 0 {
		log.Fatal("input is mandatory")
	}

	if len(prefix) == 0 {
		log.Fatal("prefix is mandatory")
	}

	content, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	var abi eos.ABI

	err = json.Unmarshal(content, &abi)
	if err != nil {
		panic(err)
	}

	err = writeOutput(abi, pack, prefix)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "fmt")
	log.Println("Running command go fmt...")
	err = cmd.Run()
	log.Println("DONE!! ")

}
