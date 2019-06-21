package helpers

import "flag"

//StartupParams contains information about parameters provided during application start
type StartupParams struct {
	IP   *string
	Port *string
}

//MakeStringPointer provides easy way to obtain inline string pointer
//Its sole purpose is to provide cleaner code during filling structs
func MakeStringPointer(defaultValue string) *string {
	return &defaultValue
}

//NewParams defines and parses command line arguments returning them in Object manner
func NewParams() *StartupParams {
	sp := &StartupParams{
		IP:   MakeStringPointer(""),
		Port: MakeStringPointer(""),
	}

	flag.StringVar(sp.IP, "addr", "0.0.0.0", "Server orchestrator listen ip address")
	flag.StringVar(sp.Port, "port", "8080", "Server orchestrator listen ip port")

	flag.Parse()

	return sp
}
