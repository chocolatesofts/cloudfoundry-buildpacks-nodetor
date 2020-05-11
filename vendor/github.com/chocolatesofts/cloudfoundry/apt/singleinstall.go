package apt

func SingleInstall(s *Supplier,pkg string,typ string) error {
	installDir := InstallDir(s,pkg)
	err,options,doptions,archiveDir :=AptSetup(s,installDir)
	if(err!=nil){
	 	return err
	 }

	err = AptUpdate(s,options)
	if(err!=nil){
	 	return err
	 }

	err = DownloadRepoPackage(s,doptions,pkg)
	if(err!=nil){
	 	return err
	 }

	err = InstallPackages(s,archiveDir,installDir)
	if(err!=nil){
	 	return err
	 }

	err = LinkPackages(s,installDir)
	if(err!=nil){
	 	return err
	 }
	return nil
}