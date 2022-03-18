/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
			yaml.Unmarshal(yamlFile, vals)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		return
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
			defer writer.Flush()
			err = tpl(path, vals, writer)
			if err != nil {
				log.Println(err)
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
