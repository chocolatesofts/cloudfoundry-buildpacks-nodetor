package tor
import (
	"strconv"
	"fmt"
	"strings"
	"errors"
	"github.com/cloudfoundry/libbuildpack"
)

const configfilename="tor*.yml"
type ConfigYML struct{
	Ports	[]struct{
		SockPort    string `yaml:"sockport"`
		ControlPort string `yaml:"controlport"`
		Pass 		string `yaml:"pass"`
		Config 		string `yaml:"config"`
	}	`yaml:"defaultports"`

	ControlPass	string    `yaml:"defaultpass"`
}

type ConfigPort struct{
	SockPort		int
	ControlPort		int
	Pass 			string
	Config 			string
	Datadir 		string
}
func determinePort(Port string) []int {
	if(strings.Contains(Port,"..")){
		ports:=strings.Split(Port,"..")
		if(len(ports)!=2){
			panic(fmt.Errorf("invalid port range syntax. Try : start..end"))
		}
		port1,err:= strconv.Atoi(ports[0])
		if(err!=nil || port1<0 || port1>65535){
			panic(fmt.Errorf("invalid starting port : '%s'",ports[0]))
		}
		_,err= strconv.Atoi(ports[1])
		if(err!=nil){
			panic(fmt.Errorf("invalid ending port : '%s'",ports[1]))
		}
		if ((len(ports[0])<len(ports[1]))|| (len(ports[1])==0)){
			if(len(ports[1])==0){
				panic("port number missing. Try : start..end")
			}
			panic(fmt.Errorf("invalid port range syntax. Try : %s%s..%s",strings.Repeat("0",len(ports[1])-len(ports[0])),ports[0],ports[1]))
		}
		com:=string(([]rune(ports[0]))[0:len(ports[0])-len(ports[1])])	
		ports1:=com+ports[1]
		port2,err:= strconv.Atoi(ports1)
		if(err!=nil || port2<0 || port2>65535){
			panic(fmt.Errorf("invalid ending port : '%s'",ports[1]))
		}
		if(port1>port2){
			panic(fmt.Errorf("starting port is greater than ending port. Try : %s..%s",ports[1],ports[0])) 
		}
		if(port1==port2){
			panic(fmt.Errorf("port range not a range. Try single port syntax"))
		}
		portarr:= make([]int,port2-port1+1)
		for i:=port1;i<port1+len(portarr);i++{
			portarr[i-port1]=i
		}
		return portarr
	} else {
		port,err:= strconv.Atoi(Port)
		if(err!=nil){
			panic(fmt.Errorf("invalid port : %s",port))
		}
		return []int{port}
	}

}
func ParseConfig(file string) (configs []ConfigPort,err error){
	cml:= &ConfigYML{}
	libbuildpack.NewYAML().Load(file,cml)
	pnum:=0
	ptype:=""
	defer func(){
		r:=recover()
		if(r!=nil){
			errormsg:=fmt.Sprintf("in defaultports config\nin port config index: %v\n",pnum)
			if(ptype!=""){
				errormsg=errormsg+fmt.Sprintf("in %s : ",ptype)
			}
			errormsg=fmt.Sprint(errormsg,r)
			err=errors.New(errormsg)
			configs=nil
		}
	}()
	if(cml.Ports==nil){
		panic(fmt.Errorf("no default port defined"))
	}
	portconfigs:=[]ConfigPort{}
	for i:=0;i<len(cml.Ports);i++{
		pnum=i+1
		portset:=cml.Ports[i]
		if(portset.SockPort==""){
			panic(fmt.Errorf("no sockport specified"))
		}
		ptype= "sockport"
		sport:=determinePort(portset.SockPort)
		ptype= ""
		cport:=[]int{}
		if(portset.ControlPort==""){
			if((len(sport)==1 && (portset.Config==""))|| (len(sport)>1)){
				panic(fmt.Errorf("insufficient control ports"))
			}
			fmt.Printf("in port config index: %v\nexplicit control port not set. Not checking for port conflict\n",pnum)
		} else{
			ptype= "controlport"
			cport=determinePort(portset.ControlPort)
			ptype= ""
			if(len(cport)!=len(sport)){
				panic(fmt.Errorf("sock and control port length mismatch"))
			}
		}
		cpass:=portset.Pass
		conf:=portset.Config
		if(conf=="" && cpass==""){
			if(cml.ControlPass==""){
				panic(fmt.Errorf("no matching control port password. Try : setting 'defaultpass'"))
			} else{
				cpass=cml.ControlPass
			}
		}
		for num:=0;num<len(sport);num++{
			tcport:=-1
			if(len(cport)>num){
				tcport=cport[num]
			}	
			portconfigs=append(portconfigs,ConfigPort{sport[num],tcport,cpass,conf,""})
		}
	}
	configs = portconfigs
	err = nil
	return
}
func (conp *ConfigPort) getCommand() string{
	basecomm:= "mkdir -p "+conp.Datadir
	if(conp.Pass!=""){
		basecomm=basecomm+"\ntorpass=$(tor --hash-password \""+conp.Pass+"\")"
	}
	basecomm=basecomm+"\n#nohup tor "
	if(conp.Config!=""){
		basecomm=basecomm+"-f "+conp.Config+" "
	}
	basecomm=basecomm+"SOCKSPort "+strconv.Itoa(conp.SockPort)+" CONTROLPort "+strconv.Itoa(conp.ControlPort)
	if(conp.Pass!=""){
		basecomm=basecomm+" HashedControlPassword $torpass"
	}
	basecomm=basecomm+" DATADirectory "+conp.Datadir+" &"
	return basecomm
}

