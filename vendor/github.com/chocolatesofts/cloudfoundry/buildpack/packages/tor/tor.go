package tor
import(
	"github.com/chocolatesofts/cloudfoundry/apt"
	"github.com/chocolatesofts/cloudfoundry/util"	
	"github.com/cloudfoundry/libbuildpack"
	"path/filepath"
)
const (
	Packagename="tor"
)
type Supplier struct{
	AptSupplier			*apt.Supplier
	Logger			    *libbuildpack.Logger
	Config				string
}

func InstallTor(s *Supplier) error {
	s.Logger.Info("Installing Tor.....")
	err:=apt.SingleInstall(s.AptSupplier,Packagename,"repo")
	if(err!=nil){
		return err
	}
	s.Logger.Info("Tor installed!!!")
	cfile:=configfilename
	if(s.Config!=""){
		cfile=s.Config
	}

	cfiles,err:=filepath.Glob(filepath.Join(s.AptSupplier.Stager.BuildDir(),cfile))
	if(err!=nil){
		return err
	}
	for _,cfile=range cfiles {
		isf,err:= util.IsFile(cfile)
		if(err!=nil){
			return err
		}
		if(!isf){
			continue
		}
		confs,err:=ParseConfig(cfile)
		if(err!=nil){
			return err
		}
		PrepareScript(confs,cfile,s.AptSupplier.Stager.DepDir())
	}
	err= WriteTorConfigs(&(s.AptSupplier.Stager))
	if(err!=nil){
		return err
	}

	err= WritePrepareScript(&(s.AptSupplier.Stager))
	if(err!=nil){
		return err
	}

//	torscript:=`export TOR_PORT_1=58`
//	s.AptSupplier.Stager.WriteProfileD("tor.sh",torscript)
	return nil
}