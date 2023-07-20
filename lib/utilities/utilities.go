package utils

import (
	"fmt"
	"os"
)

func WriteTemplatesToFile(templatedirpath string,templatename string, jsontemplate []byte) {
	// Write the templates to file
	_,err := os.Stat(templatedirpath)
	if err != nil {
		fmt.Println("")
		os.MkdirAll(templatedirpath, 0755)
	} 
	
	templatepath := templatedirpath + "/" + templatename
	_,err2 := os.Stat(templatepath)
	if err2 != nil {
		fmt.Println("")
		_,err := os.Create(templatepath)
		if err != nil {
			fmt.Println("")
		}
	}
	err3 := os.WriteFile(templatepath, jsontemplate, 0755)
	if err3 != nil {
		fmt.Println("")
	}
}
