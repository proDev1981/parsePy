package main

import "io/ioutil"
import "fmt"
import "strings"
import "os"
import "os/exec"
import "flag"


func main(){
name := "./gui"

flag.Parse()
if flag.Arg(0) != ""{
	name = flag.Arg(0)
}

html := getHtml("./src/index.html")
css := getCss("./src/style.css")
build(name,html,css)

}

func getHtml(path string)string{

	file ,_ := ioutil.ReadFile(path)

	ele := strings.Split(string(file),"\n")
	for index,item := range ele{
		if strings.Contains(item ,"<win"){
			ele[index] = strings.ReplaceAll(ele[index],"=",":")		
			ele[index] = strings.Replace(ele[index],"<win ","win({",1)
			ele[index] = strings.Replace(ele[index],">","},",1)
			ele[index] = strings.ReplaceAll(ele[index]," ",",")		
		}
		if strings.Contains(item,"<box"){
			ele[index] = strings.ReplaceAll(ele[index],"=",":")		
			ele[index] = strings.Replace(ele[index],"<box ","box({",1)
			ele[index] = strings.Replace(ele[index],">","},[",1)
			ele[index] = strings.ReplaceAll(ele[index]," ",",")
			ele[index] = strings.ReplaceAll(ele[index],",,","  ")
		}
		if strings.Contains(item,"<entry"){
			ele[index] = strings.ReplaceAll(ele[index],"=",":")		
			ele[index] = strings.Replace(ele[index],"<entry ","entry({",1)
			ele[index] = strings.Replace(ele[index],">","}",1)
			ele[index] = strings.ReplaceAll(ele[index]," ",",")
			ele[index] = strings.Replace(ele[index],"</entry>","),",1)
			ele[index] = strings.ReplaceAll(ele[index],",,","  ")
		}
		if strings.Contains(item,"<label"){	
			open := strings.Index(item,">")
			header := item[:open+1]
			inner := item[open+1:]
			header = strings.ReplaceAll(header,"=",":")		
			header = strings.Replace(header,"<label ","label({",1)
			header = strings.Replace(header,"<label","label({",1)
			header = strings.ReplaceAll(header," ",",")
			inner = `,"`+ inner
			inner = strings.Replace(inner,"<",`"<`,1)
			ele[index] = fmt.Sprint(header,inner)
			ele[index] = strings.Replace(ele[index],">","}",1)
			ele[index] = strings.Replace(ele[index],"</label>","),",1)
			ele[index] = strings.ReplaceAll(ele[index],",,","  ")
		}
		if strings.Contains(item,"<button"){
			open := strings.Index(item,">")
			header := item[:open+1]
			inner := item[open+1:]
			header = strings.ReplaceAll(header,"=",":")				
			header = strings.Replace(header,"<button ","button({",1)
			header = strings.Replace(header,"<button","button({",1)
			header = strings.ReplaceAll(header," ",",")
			inner = `,"`+ inner
			inner = strings.Replace(inner,"<",`"<`,1)
			ele[index] = fmt.Sprint(header,inner)
			ele[index] = strings.Replace(ele[index],">","}",1)
			ele[index] = strings.Replace(ele[index],"</button>","),",1)
			ele[index] = strings.ReplaceAll(ele[index],",,","  ")
		}
		if strings.Contains(item,"</box>"){
			ele[index] = strings.Replace(ele[index],"</box>","]),",1)
		}
		if strings.Contains(item,"</win>"){
			ele[index] = strings.Replace(ele[index],"</win>",")",1)
		}
	
	}
	return strings.Join(ele,"\n")

}
func getCss(path string)string{
	file ,_ := ioutil.ReadFile(path)
	return string(file)
}
func build(name , html , css string){

	tags := []string{"win","box","entry","button","label"}


	file , err := ioutil.ReadFile("/usr/lib/gtk-gui/.lib/base.lib")

	if err != nil{
		fmt.Println(err)
		return
	}
	base := string(file)
	for _,tag := range tags{
		if strings.Contains(html,tag){
			file , err = ioutil.ReadFile("/usr/lib/gtk-gui/.lib/"+tag+".lib")
			if err != nil{
				fmt.Println(err)
				return
			}
			base = strings.Replace(base,"##methods##","\n##methods##\n\n"+string(file),1)
		}
		
	}
	base = strings.Replace(base,"##html##",html,1)
	file , err = ioutil.ReadFile("/usr/lib/gtk-gui/.lib/css.lib")
	if err != nil{
		fmt.Println(err)
		return
	}
	css = strings.Replace(string(file),"##css##",css,1)
	base = strings.Replace(base,"##css##",css,1)
	//fmt.Println(base)
	py,create:= os.Create(name)
	if create != nil{
		fmt.Println(create)
		return
	}
	py.WriteString(base)
	//os.Chmod(name,7)
	exec.Command("chmod","+x", name).Run()
}





