package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/microlib/simple"
)

type Component struct {
	ID          string        `json:"id"`
	Component   string        `json:"component"`
	Tab         string        `json:"tab"`
	Name        string        `json:"name"`
	Reference   string        `json:"reference,omitempty"`
	X           int           `json:"x"`
	Y           int           `json:"y"`
	Connections Connection    `json:"connections,omitempty"`
	Options     OptionDetails `json:"options,omitempty"`
}

type Connection struct {
	Num0 []struct {
		Index string `json:"index"`
		ID    string `json:"id"`
	} `json:"0"`
}

type OptionDetails struct {
	Code       string        `json:"code,omitempty"`
	Outputs    int           `json:"outputs,omitempty"`
	Filename   string        `json:"filename,omitempty"`
	Append     bool          `json:"append,omitempty"`
	Delimiter  string        `json:"delimiter,omitempty"`
	Template   string        `json:"template,omitempty"`
	Layout     bool          `json:"layout,omitempty"`
	Name       string        `json:"name,omitempty"`
	Arg        []interface{} `json:"arg,omitempty"`
	Convert    string        `json:"convert,omitempty"`
	Timeout    int           `json:"timeout,omitempty"`
	Type       string        `json:"type,omitempty"`
	Repository bool          `json:"repository,omitempty"`
	Enabled    bool          `json:"enabled,omitempty"`
	Parser     string        `json:"parser,omitempty"`
	URL        string        `json:"url,omitempty"`
	Method     string        `json:"method,omitempty"`
	Stringify  string        `json:"stringify,omitempty"`
	Props      []string      `json:"props,omitempty"`
	ID         bool          `json:"id,omitempty"`
	Fn         string        `json:"fn,omitempty"`
}

type Flow struct {
	Tabs []struct {
		Name   string `json:"name"`
		Linker string `json:"linker"`
		ID     string `json:"id"`
		Index  int    `json:"index"`
	} `json:"tabs"`
	Components []Component `json:"components,omitempty"`
	Disabledio struct {
		Input  []interface{} `json:"input"`
		Output []interface{} `json:"output"`
	} `json:"disabledio"`
	State struct {
		Text  string `json:"text"`
		Color string `json:"color"`
	} `json:"state"`
	Color     string    `json:"color"`
	Notes     string    `json:"notes"`
	Variables string    `json:"variables"`
	Panel     string    `json:"panel"`
	URL       string    `json:"url"`
	Created   time.Time `json:"created"`
}

type FileDetails struct {
	Output  string `json:"output"`
	Content string `json:"content"`
}

type HtmlSchema struct {
	Title                     string `json:"title"`
	TitleDescription          string `json:"titleDescription"`
	Headline                  string `json:"headline"`
	SubHeadline               string `json:"subheadline"`
	StatementATitle           string `json:"statementATitle"`
	StatementA                string `json:"statementA"`
	StatementBTitle           string `json:"statementBTitle"`
	StatementB                string `json:"statementB"`
	ProofATitle               string `json:"proofATitle"`
	ProofA                    string `json:"proofA"`
	ProofB                    string `json:"proofB"`
	ProofBTitle               string `json:"proofBTitle"`
	ProofC                    string `json:"proofC"`
	ProofCTitle               string `json:"proofCTitle"`
	ContradictionHandlerTitle string `json:"contradictionHandlerTitle"`
	ContradictionHandler      string `json:"contradictionHandler"`
	ButtonA                   string `json:"buttonA"`
	ButtonB                   string `json:"buttonB"`
	VideoA                    string `json:"videoA"`
	VideoATitle               string `json:"videoATitle"`
	VideoADescription         string `json:"videoADescription"`
	VideoB                    string `json:"videoB"`
	VideoBTitle               string `json:"videoBTitle"`
	VideoBDescription         string `json:"videoBDescription"`
	Audio                     string `json:"audio"`
	AudioTitle                string `json:"audioTitle"`
	AudioDescription          string `json:"audioDescription"`
	DataA                     string `json:"dataA"`
	DataATitle                string `json:"dataATitle"`
	DataAUrl                  string `json:"dataAUrl"`
	DataB                     string `json:"dataB"`
	DataBTitle                string `json:"dataBTitle"`
	DataBUrl                  string `json:"dataBUrl"`
	DataC                     string `json:"dataC"`
	DataCTitle                string `json:"dataCTitle"`
	DataCUrl                  string `json:"dataCUrl"`
	DataD                     string `json:"dataD"`
	DataDTitle                string `json:"dataDTitle"`
	DataDUrl                  string `json:"dataDUrl"`
	Quote                     string `json:"quote"`
	SummaryTitle              string `json:"summaryTitle"`
	Summary                   string `json:"summary"`
	Address                   string `json:"address"`
	Pricing                   string `json:"pricing"`
	PricingDescription        string `json:"pricingDescription"`
	PlanA                     string `json:"planA"`
	PlanADescription          string `json:"planADescription"`
	PlanADetails              string `json:"planADetails"`
	PlanAUrl                  string `json:"planAUrl"`
	PlanABtn                  string `json:"planABtn"`
	AS                        string `json:"AS"`
	AP                        string `json:"AP"`
	AM                        string `json:"AM"`
	PlanB                     string `json:"planB"`
	PlanBDescription          string `json:"planBDescription"`
	PlanBDetails              string `json:"planBDetails"`
	PlanBUrl                  string `json:"planBUrl"`
	PlanBBtn                  string `json:"planBBtn"`
	BS                        string `json:"BS"`
	BP                        string `json:"BP"`
	BM                        string `json:"BM"`
	About                     string `json:"about"`
	AboutDescription          string `json:"aboutdescription"`
	Contact                   string `json:"contact"`
	Phone                     string `json:"phone"`
	CTA                       string `json:"cta"`
	CTAUrl                    string `json:"ctaurl"`
	LinkAUrl                  string `json:"linkaurl"`
	LinkBUrl                  string `json:"linkburl"`
	LinkCUrl                  string `json:"linkcurl"`
	LinkABtn                  string `json:"linkabtn"`
	LinkBBtn                  string `json:"linkbbtn"`
	LinkCBtn                  string `json:"linkcbtn"`
	Pagename                  string `json:"pagename"`
	Pagetype                  string `json:"pagetype"`
}

func main() {
	var utm_campaign string
	var utm_content string
	var utm_affiliate string
	var utm_medium string
	var base_url string
	var html []byte
	var flow Flow
	var files map[string]string
	var filesTo map[string]string
	var htmlschema HtmlSchema
	var DIR string = ""

	logger := &simple.Logger{Level: "info"}

	if len(os.Args) == 1 {
		logger.Error(fmt.Sprintf("Command line args are missing"))
		os.Exit(-1)
	}
	logger.Level = os.Args[2]
	logger.Info(fmt.Sprintf("Command line args %s %d", os.Args, len(os.Args)))

	if len(os.Args) == 4 {
		DIR = "../html-templates/"
	}

	file, _ := ioutil.ReadFile(DIR + os.Args[1])
	// update our schema
	err := json.Unmarshal([]byte(file), &flow)
	if err != nil {
		logger.Error(fmt.Sprintf("Converting designer flow json %v", err))
	}

	var designer []Component
	for comp, _ := range flow.Components {
		if flow.Components[comp].Component != "comment" {
			designer = append(designer, flow.Components[comp])
		}
		if flow.Components[comp].Reference == "utm_campaign" {
			utm_campaign = flow.Components[comp].Name
		}
		if flow.Components[comp].Reference == "utm_content" {
			utm_content = flow.Components[comp].Name
		}
		if flow.Components[comp].Reference == "utm_medium" {
			utm_medium = flow.Components[comp].Name
		}
		if flow.Components[comp].Reference == "affiliate" {
			utm_affiliate = flow.Components[comp].Name
		}
		if flow.Components[comp].Reference == "base_url" {
			base_url = flow.Components[comp].Name
		}
	}

	size := len(designer)
	urlLinks := make([][]string, size)
	for i := range urlLinks {
		urlLinks[i] = make([]string, 5)
	}

	for from, _ := range designer {
		if len(designer[from].Connections.Num0) > 0 {
			for links, _ := range designer[from].Connections.Num0 {
				//logger.Info(fmt.Sprintf("LMZ Num Links %s %d", designer[from].Reference, links))
				for to, _ := range designer {
					if designer[to].ID == designer[from].Connections.Num0[links].ID {
						logger.Trace(fmt.Sprintf("Template dump %s %s", designer[from].Options.Template, designer[from].Reference))
						err := json.Unmarshal([]byte(designer[from].Options.Template), &files)
						if err != nil {
							logger.Error(fmt.Sprintf("Converting embedded file [from] data json %v %s", err, designer[from].Name+" "+designer[from].Reference))
							break
						}
						err = json.Unmarshal([]byte(designer[to].Options.Template), &filesTo)
						if err != nil {
							logger.Error(fmt.Sprintf("Converting embedded file [to] data json %v", err))
							break
						}
						logger.Debug(fmt.Sprintf("Files %v %v", files, filesTo))
						_, err = os.Stat(DIR + designer[from].Name + "/" + files["content"])
						if err != nil {
							logger.Error(fmt.Sprintf("No json content file found %v", err))
						} else {
							// ok we can open the files now
							html, _ = ioutil.ReadFile(DIR + designer[from].Name + "/template.html")
							content, _ := ioutil.ReadFile(DIR + designer[from].Name + "/" + files["content"])
							err := json.Unmarshal(content, &htmlschema)
							if err != nil {
								logger.Error(fmt.Sprintf("Umarshalling data %v\n", err))
								break
							}
							// add in links now
							if files["pagetype"] == "origin" {
								urlLinks[from][links] = base_url + designer[to].Name + "/" + filesTo["output"] +
									"?utm_campaign=" + utm_campaign +
									"&utm_source=" + strings.ToLower(designer[from].Reference) +
									"&utm_content=" + utm_content +
									"&utm_affiliate=" + utm_affiliate +
									"&utm_medium=" + utm_medium +
									"&pagename=" + files["pagename"] +
									"&pagetype=" + files["pagetype"]
							} else {
								urlLinks[from][links] = "javascript:injectParams('" + base_url + designer[to].Name + "/" + filesTo["output"] +
									"?utm_campaign=" + utm_campaign +
									"&utm_source=" + strings.ToLower(designer[from].Reference) +
									"&utm_content=" + utm_content +
									"&utm_affiliate=" + utm_affiliate +
									"&utm_medium=" + utm_medium +
									"&pagename=" + files["pagename"] +
									"&pagetype=" + files["pagetype"] + "');"
							}
						}
					}
				}
			}
		} else {
			err := json.Unmarshal([]byte(designer[from].Options.Template), &files)
			if err != nil {
				logger.Error(fmt.Sprintf("Converting embedded file [from] data json %v", err))
				break
			}
			_, err = os.Stat(DIR + designer[from].Name + "/" + files["content"])
			if err != nil {
				logger.Error(fmt.Sprintf("No json content file found %v", err))
			} else {
				// ok we can open the files now
				html, _ = ioutil.ReadFile(DIR + designer[from].Name + "/template.html")
				content, _ := ioutil.ReadFile(DIR + designer[from].Name + "/" + files["content"])
				err := json.Unmarshal(content, &htmlschema)
				if err != nil {
					logger.Error(fmt.Sprintf("Umarshalling data %v\n", err))
					break
				}
			}
		}

		var data bytes.Buffer
		if (htmlschema != HtmlSchema{}) {
			tmpl, err := template.New("transform").Parse(string(html))
			if err != nil {
				logger.Error(fmt.Sprintf("Creating transform for %s %v\n", designer[from].Reference, err))
				break
			}
			// we add in our pagename and pagetype variables
			if filesTo["pagename"] == "" || filesTo["pagetype"] == "" {
				logger.Error(fmt.Sprintf("Please ensure pagename and pagetype variables are included in the page [%s]", designer[from].Reference+" "+designer[from].Name))
				break
			} else {
				htmlschema.Pagename = files["pagename"]
				htmlschema.Pagetype = files["pagetype"]
				for y, _ := range urlLinks[from] {
					switch y {
					case 0:
						htmlschema.CTAUrl = urlLinks[from][0]
					case 1:
						htmlschema.DataAUrl = urlLinks[from][1]
					case 2:
						htmlschema.DataBUrl = urlLinks[from][2]
					case 3:
						htmlschema.DataCUrl = urlLinks[from][3]
					case 4:
						htmlschema.DataDUrl = urlLinks[from][4]
					}
				}
			}
			err = tmpl.Execute(&data, htmlschema)
			if err != nil {
				logger.Error(fmt.Sprintf("Executing transform for %s %v\n", designer[from].Reference, err))
				break
			}
			err = ioutil.WriteFile(DIR+designer[from].Name+"/"+files["output"], data.Bytes(), 0755)
			if err != nil {
				logger.Error(fmt.Sprintf("Writing file %v\n", err))
				break
			} else {
				logger.Info(fmt.Sprintf("Succesfully saved file %s\n", DIR+designer[from].Name+"/"+files["output"]))
			}
		}
	}

	os.Exit(0)
}
