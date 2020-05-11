package tor
import(
	"path/filepath"
	"fmt"
	"errors"
	"github.com/chocolatesofts/cloudfoundry/apt"
	"github.com/chocolatesofts/cloudfoundry/util"
	"github.com/cloudfoundry/libbuildpack"
)
var Sockports=make(map[int]string)
var Contports=make(map[int]string)
var Configfiles=make(map[string]string)
var Scripts=make(map[string]string)
const(
	header="#Begin Setup : %s\n"
	footer="#End Setup : %s\n"
	scriptfile="torprepare.sh"
)
func PrepareScript(confs []ConfigPort,fname string,depdir string) error{
	if(Scripts[fname]!=""){
		return fmt.Errorf("Config : %s already used",fname)
	}
	script:=fmt.Sprintf(header,fname)
	for i:=0;i<len(confs);i++{
		config:=confs[i]
		err:=writeLog(&config,"sock",fname)
		if(err!=nil){
			return err
		}
		err=writeLog(&config,"control ",fname)
		if(err!=nil){
			return err
		}		
		if(config.Config!=""){
			cfile:=GetConfFile(config.Config,depdir)
			Configfiles[config.Config]=cfile
			config.Config=cfile
		}
		script=script+"\n"+config.getCommand()
	}
	script=script+"\n"+fmt.Sprintf(footer,fname)
	Scripts[fname]=script
	return nil
}
func WriteTorConfigs(stager *apt.Stager) error{
	for infile,outfile:= range Configfiles{
		oinfile:=infile
		infile=filepath.Join((*stager).BuildDir(),infile)
		exists,err:=libbuildpack.FileExists(infile)
		if(err!=nil){
			return err
		}
		isf,err:= util.IsFile(infile)
		if(err!=nil){
			return err
		}
		if(!exists || !isf){
			return fmt.Errorf("Config : %s file not found",oinfile)
		}

		exists,err=libbuildpack.FileExists(outfile)
		if(err!=nil){
			return err
		}
		if(exists){
			return fmt.Errorf("Config : %s already written to %s",oinfile,outfile)
		}
		err=libbuildpack.CopyFile(infile,outfile)
		if(err!=nil){
			return err
		}
	}
	return nil
}

func WritePrepareScript(stager *apt.Stager) error{
	totalscript:=""
	for _,script:= range Scripts{
		totalscript=totalscript+script
	}
	return (*stager).WriteProfileD(scriptfile,totalscript)
}
func WriteExportScript(){
	//TODO : Implement Export option
}
func GetConfFile(conffile string,depdir string)string{
	return filepath.Join(depdir,apt.Aptname,Packagename,"configs",conffile)
}
func writeLog(conf *ConfigPort,typ string,fname string) error{
	var vmap map[int]string
	var port int
	if(typ=="sock"){
		vmap=Sockports
		port=conf.SockPort
	} else{
		vmap=Contports
		port=conf.ControlPort
	}
	usef:=vmap[port]
	if(usef!=""){
		errorf:=fmt.Sprintf("in config :%s\n%sport %v already used",fname,typ,port)
		if(usef!=fname){
			errorf=errorf+fmt.Sprintf(" in config file %s",usef)
		}
		return errors.New(errorf)
	}
	vmap[port]=fname
	return nil
}