package util
import (
	"os"
)
func IsFile(path string) (bool,error) {
	info,err:= os.Stat(path)
	if(err!=nil){
		return false,err
	}
	return info.Mode().IsRegular(),nil
}