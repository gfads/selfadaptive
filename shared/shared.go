package shared

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"runtime"
	"strings"
	"time"
)

const MinOnoff = 10
const MaxOnoff = 600
const MonitorTime = 5000
const NumberOfColors = 7
const ColorReset = "\033[0m"

var DirGo = LocalizegGo() + "/bin"

const DirPluginsV1 = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/v1/env/plugins"

var ColorBehaviours = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}

func ErrorHandler(f string, msg string) {
	fmt.Println(f + "::" + msg)
	os.Exit(0)
}

func GetFunction() string {
	fpcs := make([]uintptr, 1)

	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}

	// Print the file name and line number
	//fmt.Println(caller.FileLine(fpcs[0]-1))

	// Print the name of the function
	//fmt.Println(caller.Name())

	return caller.Name()
}

func LocalizegGo() string {
	r := ""
	found := false

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "GOROOT" {
			r = pair[1]
			found = true
		}
	}

	if !found {
		fmt.Println("Shared:: Error:: OS Environment variable 'GOROOT' not configured\n")
		os.Exit(1)
	}
	return r
}

func LoadPlugin(dir string, pluginName string) plugin.Plugin {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := dir + "/" + pluginName
	attempts := 0
	for {
		plg, err = plugin.Open(pluginFile)

		if err != nil {
			if attempts >= 3 { // TODO
				fmt.Printf("Shared:: Error on trying open plugin '%v' \n", pluginFile)
				os.Exit(0)
			} else {
				attempts++
				time.Sleep(10 * time.Millisecond) // TODO
			}
		} else {
			break
		}
	}

	return *plg
}

func LoadSources(dir string) []string {
	r := []string{}

	folders, err1 := ioutil.ReadDir(dir)
	if err1 != nil {
		ErrorHandler(GetFunction(), err1.Error())
	}

	temp := []os.FileInfo{}

	for folder := range folders {
		temp, err1 = ioutil.ReadDir(dir + "/" + folders[folder].Name())
		if err1 != nil {
			ErrorHandler(GetFunction(), err1.Error())
		}

		for file := range temp {
			fullPathName := dir + "/" + folders[folder].Name() + "/" + temp[file].Name()
			r = append(r, fullPathName)
		}
	}
	return r
}

func GenerateExecutable(dir string, sources []string) {

	for i := range sources {
		plugin := sources[i]
		name := plugin[strings.LastIndex(plugin, "/")+1:]
		pOut := dir + "/" + name[:strings.LastIndex(name, ".")]
		pIn := plugin

		_, err := exec.Command(DirGo+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
		if err != nil {
			ErrorHandler(GetFunction(), "Something wrong in generating plugin '"+pIn+"' in "+pOut+" "+err.Error())
		}
	}
}

func LoadFuncs(sourcesDir, executablesDir string) []func() {
	sourcePluginFiles := LoadSources(sourcesDir)
	GenerateExecutable(executablesDir, sourcePluginFiles)

	r := []func(){}
	for i := 0; i < len(sourcePluginFiles); i++ {
		p := sourcePluginFiles[i]
		plugin := LoadPlugin(executablesDir, p[strings.LastIndex(p, "/")+1:strings.LastIndex(p, ".")])
		f, err := plugin.Lookup("Behaviour")
		if err != nil {
			ErrorHandler(GetFunction(), "Function not found in plugin!!")
		}
		r = append(r, f.(func()))
	}
	return r

}
