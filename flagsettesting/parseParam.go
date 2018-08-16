package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type param struct {
	Ip string `yaml:"ip"`
	Port string `yaml:"port"`
}

var TEMPLATE = `ip: 127.0.0.1
port: 3344
`

func main() {
	//myparam := new(param)
	//fs := flag.NewFlagSet("mytest",flag.ExitOnError)
	//r_ip := fs.String("ip",myparam.ip,"remote ip")
	//fs.String("port",myparam.port,"remote port")
	//fs.Parse(os.Args[1:])
	//fmt.Println(*r_ip)
	fBytes,_ := ioutil.ReadFile("./test.yml")

	cfg := &param{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.
	//*cfg = DefaultConfig
	//m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(fBytes, cfg)
	if err != nil {

	}
	fmt.Println(cfg.Ip,":",cfg.Port)
}

var DefaultConfig = param{
	Ip:"123",
	Port:"456",
}


