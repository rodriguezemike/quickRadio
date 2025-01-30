package radioErrors

import "log"

func ErrorLog(err error) {
	if err != nil {
		log.Println(err)
	}
}
func ErrorFail(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
