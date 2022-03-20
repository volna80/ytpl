package main

import (
	"bufio"
	"flag"
	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var workDir = flag.String("input", ".", "a directory with yaml files which should be rendered")
var outputDir = flag.String("output", "."+string(os.PathSeparator)+"output", "a folder with the result yaml files")

func init() {
	//log.SetFlags(log.Lshortfile)
	flag.Parse()
}

func main() {

	if *flag.Bool("help", false, "-help - print this helm message") {
		flag.PrintDefaults()
		return
	}

	os.MkdirAll(*outputDir, os.ModePerm)

	var vals = make(map[string]interface{})

	//load variables
	err := filepath.Walk(*workDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		matched, err := regexp.Match("_.*.yaml", []byte(info.Name()))
		if err != nil {
			return err
		}
		if matched {
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(yamlFile, vals)
			if err != nil {
				log.Fatal("Invalid yaml: "+path+" ", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	//execute templates
	err = filepath.Walk(*workDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		matched, err := regexp.Match("_.*.yaml", []byte(info.Name()))
		if err != nil {
			return err
		}

		if !matched {
			file, err := create(*outputDir + string(os.PathSeparator) + path)
			if err != nil {
				log.Fatal(err)
			}

			writer := bufio.NewWriter(file)
			err = tpl(path, vals, writer)
			if err != nil {
				log.Fatal(err)
			}
			writer.Flush()

			//validate the result
			yamlFile, err := ioutil.ReadFile(file.Name())
			if err != nil {
				log.Println(err)
				return err
			}

			err = yaml.Unmarshal(yamlFile, make(map[string]interface{}))
			if err != nil {
				log.Fatal("Invalid yaml: "+file.Name()+" ", err)
			}

		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

}

func tpl(t string, vals map[string]interface{}, out io.Writer) error {
	file, err := ioutil.ReadFile(t)
	if err != nil {
		return err
	}
	tt, err := template.New("_").Funcs(sprig.FuncMap()).Parse(string(file))
	if err != nil {
		return err
	}
	return tt.Execute(out, vals)
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
