package processing

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/evilcry/pastemon/db"

	"github.com/evilcry/pastemon/configs"
)

func savePaste(conf *configs.Config, key, content string) {
	if len(content) > conf.MaxSize {
		return
	}

	log.Println("[+] Saving paste: ", key)

	if conf.SaveFile {
		err := ioutil.WriteFile(path.Join(conf.PasteDir, key), []byte(content), 0644)
		if err != nil {
			log.Println("[!] Error while storing paste: ", err)
		}
	}

	if conf.Save == false {
		return
	}

	conf.Ds.Write("pastes", key, []byte(content))
}

func processRegexes(conf *configs.Config, key, content string) {
	save := false
	for i, _ := range conf.Regexes {
		r := conf.Regexes[i]

		switch r.Match {
		case "all":
			items := r.Compiled.FindAllString(content, -1)

			if items != nil {
				save = true
			}

			for k := range items {
				rKey := fmt.Sprintf("%s-%s-%d", r.Prefix, key, k)
				conf.Ds.Write("regexes", rKey, []byte(items[k]))
			}
		case "one":
			match := r.Compiled.FindString(content)
			rKey := fmt.Sprintf("%s-%s", r.Prefix, key)

			if match != "" {
				save = true
				conf.Ds.Write("regexes", rKey, []byte(match))
			}
		default:
		}
	}

	if save {
		savePaste(conf, key, content)
	}
}

func processKeywords(conf *configs.Config, key, content string) {
	save := false
	for i, _ := range conf.Keywords {
		kwd := conf.Keywords[i]
		kwdKey := fmt.Sprintf("%s-%s", kwd.Prefix, key)

		if strings.Contains(strings.ToLower(content), strings.ToLower(kwd.Keyword)) {
			save = true
			conf.Ds.Write("keywords", kwdKey, []byte(key))
		}
	}

	if save {
		savePaste(conf, key, content)
	}
}

// ProcessContent func
// pastebin analysis
func ProcessContent(conf *configs.Config, key, content string) {
	conf.Ds = db.GetStorageConnection(conf)
	defer conf.Ds.Close()

	processRegexes(conf, key, content)
	processKeywords(conf, key, content)
}
