package meta

import "fmt"

var (
	Version     = "0.0.0"
	CommitHash  = ""
	CompileDate = ""
	Vendor      = "cynalytica"
	VendorMail  = "maintainer@cynalytica.com"
	Name        = "doc-tools"
	Usage       = fmt.Sprintf("%s %s v%s", Vendor, Name, Version)
)
