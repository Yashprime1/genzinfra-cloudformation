package utils

import (
	"fmt"
	"os"
)

func WriteTemplatesToFile(templatedirpath string,templatename string, jsontemplate []byte) {
	// Write the templates to file
	_,err := os.Stat(templatedirpath)
	if err != nil {
		fmt.Println(err)
		os.MkdirAll(templatedirpath, 0755)
	} 
	
	templatepath := templatedirpath + "/" + templatename
	_,err2 := os.Stat(templatepath)
	if err2 != nil {
		fmt.Println(err2)
		j,err := os.Create(templatepath)
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println(j)
		}
	}
	err3 := os.WriteFile(templatepath, jsontemplate, 0755)
	if err3 != nil {
		fmt.Println(err3)
	}
}
