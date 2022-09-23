package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"text/template"
)

const providerTemplate = `# Generated from consul-peering-setup
{{ range .Providers }}
provider "consul" {
  address = "{{ .Address }}"
  alias = "{{ .Name }}"
}
{{end}}`

const acceptorTemplate = `# Generated from consul-peering-setup
{{ range .Peerings }}
resource "consul_peering_token" "{{ .Name }}-to-{{ .PeerName }}" {
  provider  = consul.{{ .ProviderName }}
  peer_name = "{{ .PeerName }}"
  {{if .Partition}}partition = "{{ .Partition }}"{{end}}
}
{{end}}`

const dialerTemplate = `# Generated from consul-peering-setup
{{ range .Peerings }}
resource "consul_peering" "{{ .Name }}-to-{{ .PeerName }}" {
  provider  = consul.{{ .ProviderName }}
  peer_name = "{{ .PeerName }}"
  peering_token = consul_peering_token.{{ .PeerName }}-to-{{ .Name }}.peering_token
  {{if .Partition}}partition = "{{ .Partition }}"{{end}}
}
{{end}}`

func (config Config) makeProviders() {
	pi := config.ToProviderInput()

	t, err := template.New("providers").Parse(providerTemplate)
	if err != nil {
		log.Fatal("error parsing provider template: %s", err)
	}

	file, err := os.Create("providers.tf")
	if err != nil {
		log.Fatal("error making tf providers: %s", err)
	}

	err = t.Execute(file, pi)
	if err != nil {
		log.Fatalf("error constructing providers: %s", err)
	}
}

func (config Config) makeAcceptors() {
	ai := config.ToAcceptorInput()
	t, err := template.New("acceptors").Parse(acceptorTemplate)
	if err != nil {
		log.Fatal("error parsing acceptor template: %s", err)
	}

	file, err := os.Create("acceptors.tf")
	if err != nil {
		log.Fatal("error making tf aceptors: %s", err)
	}

	err = t.Execute(file, ai)
	if err != nil {
		log.Fatalf("error constructing acceptors: %s", err)
	}
}

func (config Config) makeDialers() {
	ai := config.ToDialerInput()
	t, err := template.New("dialers").Parse(dialerTemplate)
	if err != nil {
		log.Fatal("error parsing dialer template: %s", err)
	}

	file, err := os.Create("dialers.tf")
	if err != nil {
		log.Fatal("error making tf dialers: %s", err)
	}

	err = t.Execute(file, ai)
	if err != nil {
		log.Fatalf("error constructing dialers: %s", err)
	}
}

type ProviderInput struct {
	Providers []ProviderData
}

type ProviderData struct {
	Address string
	Name    string
}

type Config []ConfigDatum

type ConfigDatum struct {
	PeerName  string `json:"peer_name"`
	Address   string `json:"address"`
	Partition string `json:"partition"`
}

func (c ConfigDatum) ProviderName() string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	return re.ReplaceAllString(c.Address, "")
}

type PeeringInput struct {
	Peerings []Peering
}

type Peering struct {
	Name         string
	PeerName     string
	Partition    string
	ProviderName string
}

func (c Config) ToProviderInput() ProviderInput {
	var pi ProviderInput
	aliases := make(map[string]struct{})
	for _, clusterData := range c {
		providerName := clusterData.ProviderName()
		if _, ok := aliases[providerName]; !ok {
			pi.Providers = append(pi.Providers, ProviderData{
				Address: clusterData.Address,
				Name:    providerName,
			})
			aliases[providerName] = struct{}{}
		}
	}

	return pi
}

func (c Config) ToAcceptorInput() PeeringInput {
	var ai PeeringInput
	for i := 0; i < len(c)/2+1; i++ {
		acceptor := c[i]
		for j := i + 1; j < len(c); j++ {
			dialer := c[j]
			ai.Peerings = append(ai.Peerings, Peering{
				Name:         acceptor.PeerName,
				PeerName:     dialer.PeerName,
				ProviderName: acceptor.ProviderName(),
				Partition:    acceptor.Partition,
			})
		}
	}

	return ai
}

func (c Config) ToDialerInput() PeeringInput {
	var ai PeeringInput
	for i := 0; i < len(c)/2+1; i++ {
		acceptor := c[i]
		for j := i + 1; j < len(c); j++ {
			dialer := c[j]
			ai.Peerings = append(ai.Peerings, Peering{
				Name:         dialer.PeerName,
				PeerName:     acceptor.PeerName,
				ProviderName: dialer.ProviderName(),
				Partition:    dialer.Partition,
			})
		}
	}

	return ai
}

func (c Config) validateConfig() error {
	// TODO verify thingsl like peer name uniqueness.
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("unexpected arguments")
	}

	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}
	defer file.Close()

	configReader, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	var config Config
	err = json.Unmarshal([]byte(configReader), &config)
	if err != nil {
		log.Fatalf("error parsing file as json: %s", err)
	}

	err = config.validateConfig()
	if err != nil {
		log.Fatalf("config validation error: %s", err)
	}

	if err != nil {
		log.Fatalf("unexpected error reading config file: %s", err)
	}

	config.makeProviders()
	config.makeAcceptors()
	config.makeDialers()
}
