package apt
import(
	"fmt"
)
func DownloadRepoPackage(s *Supplier,options []string,pkg string) error{
	aptArgs :=append(options,pkg)
	out, err := s.Command.Output("/", "apt-get", aptArgs...)
	if err != nil {
		return fmt.Errorf("failed apt-get install %s\n\n%s", out, err)
	}
	return nil;
}
