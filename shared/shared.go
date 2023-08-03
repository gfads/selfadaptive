package shared

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
	"time"
)

// Configuration RabbitMQ
//const IpPortRabbitMQ = "172.22.40.187:5672" // CIN
const IpPortRabbitMQ = "192.168.0.20:5672" // Home Recife

// Astar
const SV = 2.7           // Shutoff voltage (page 17) = 2.7 V
const OV = 3.7           // Optimum voltage (page 17) = 3.7 V
const HYSTERISIS = 0.001 // 10 mV

// Prefetch counter
const PcDefaultLimitMin = 1
const PcDefaultLimitMax = 1200 // TODO ASTAR

// Execution Types
const StaticGoal = "StaticGoal"
const Experiment = "Experiment"
const RootLocusTraining = "RootLocusTraining"
const ZieglerTraining = "ZieglerTraining"
const CohenTraining = "CohenTraining"
const AMIGOTraining = "AMIGOTraining"
const OnLineTraining = "OnlineTraining"
const WebTraining = "WebTraining"

// Controller type names
const BasicOnoff = "OnOff"
const DeadZoneOnoff = "OnOffwithDeadZone"
const HysteresisOnoff = "OnOffwithHysteresis"
const BasicP = "BasicP"
const BasicPi = "BasicPI"
const BasicPid = "BasicPID"
const SmoothingPid = "SmoothingDerivativePID"
const IncrementalFormPid = "IncrementalFormPID"
const ErrorSquarePidFull = "ErrorSquarePIDFull"
const ErrorSquarePidProportional = "ErrorSquarePIDProportional"
const DeadZonePid = "DeadZonePID"
const GainScheduling = "GainScheduling"
const PIwithTwoDegreesOfFreedom = "PIWithTwoDegreesOfFreedom"
const WindUp = "WindUp"
const SetpointWeighting = "SetpointWeighting"

const RootLocus = "RootLocus"
const Ziegler = "Ziegler"
const Cohen = "Cohen"
const Amigo = "AMIGO"
const None = "None"

const AsTAR = "AsTAR"
const HPA = "HPA"

// Experiments parameters
const TrainingSampleSize = 60
const TimeBetweenAdjustments = 1200 // seconds
const MaximumNrmse = 0.30
const WarmupTime = 30 // seconds
const TrainingAttempts = 30
const SizeOfSameLevel = 60

// Controller parameters
const MinOnoff = 10
const MaxOnoff = 600

const MonitorTime = 10 // seconds
const NumberOfColors = 7
const ColorReset = "\033[0m"

// Base Directories
const SourcesDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/envrnment/plugins/source"
const ExecutablesDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/envrnment/plugins/executable"
const DockerDir = "/app/data" // it is mapped into windows dir "C:\Users\user\go\selfadaptive\rabbitmq\data" (see execute-old.bat)

// const DataDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/rabbitmq/data/" // macos
const DataDir = "C:\\Users\\user\\go\\selfadaptive\\rabbitmq\\data" // macos
const DockerfilesDir = "C:\\Users\\user\\go\\selfadaptive"          // macos
const BatchfilesDir = "C:\\Users\\user\\go\\selfadaptive"           // macos
const BatchFileExperiments = "execute-all-experiments.bat"

//const SourcesDir = "C:\\Users\\user\\go\\selfadaptive\\example-plugin\\envrnment\\sources"
//const ExecutablesDir = "C:\\Users\\user\\go\\selfadaptive\\example-plugin\\envrnment\\executables"

// Goals
const AlwaysSecure = "Always Secure"
const NoWorry = "Do not take care"
const AlwaysUpdated = "Always Updated"
const BestEffort = "Best Effort"

// Environments security level
var EnvironmentSecurityLevels = []string{Secure, Suspicious, Unsecure}

// Symptom
type Symptoms struct {
	PluginSymptom   int
	SecuritySymptom string
}

// plugin symptoms
const NewPluginvAvailable = 0
const NoNewPluginAvailable = 1

// security symptoms
const Secure = "SecureEnvironment"
const Suspicious = "SuspiciousEnvironment"
const Unsecure = "UnsecureEnvironment"

// Request types
const NoChange = "No Update Needed"
const UsePlainText = "Use Plain Text"
const UseWeakCryptography = "Use Weak Cryptography"
const UseMediumCryptography = "Use Medium Cryptography"
const UseStrongCryptography = "Use Strong Cryptography"

// security
const PlainText = "This is the sent message"                                                                                      // This must be of 16 byte length!!
var Keys32 = []string{"WeakXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "MediumXXXXXXXXXXXXXXXXXXXXXXXXXXX", "StrongXXXXXXXXXXXXXXXXXXXXXXXXX"} // This must be of 32 byte length!!

// dir configurations
var DirGo = LocalizegGo() + "/bin"

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
			if attempts >= 3 {
				fmt.Printf("Shared:: Error on trying open plugin '%v' \n", pluginFile)
				os.Exit(0)
			} else {
				attempts++
				time.Sleep(10 * time.Millisecond)
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
		plgin := sources[i]
		name := plgin[strings.LastIndex(plgin, "/")+1:]
		pOut := dir + "/" + name[:strings.LastIndex(name, ".")]
		pIn := plgin

		_, err := exec.Command(DirGo+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
		if err != nil {
			ErrorHandler(GetFunction(), "Something wrong in generating plugin '"+pIn+"' in "+pOut+" "+err.Error())
		}
	}
}

func LoadPlugins(sourcesDir, executablesDir string) []func() {
	sourcePluginFiles := LoadSources(sourcesDir)
	GenerateExecutable(executablesDir, sourcePluginFiles)

	r := []func(){}
	for i := 0; i < len(sourcePluginFiles); i++ {
		p := sourcePluginFiles[i]
		plgin := LoadPlugin(executablesDir, p[strings.LastIndex(p, "/")+1:strings.LastIndex(p, ".")])
		f, err := plgin.Lookup("Behaviour")
		if err != nil {
			ErrorHandler(GetFunction(), "Function not found in plugin!!")
		}
		r = append(r, f.(func()))
	}
	return r

}

func RemoveContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func EncryptMessage(key string, message string) string {
	c, err := aes.NewCipher([]byte(key))
	//c, err := des.NewCipher([]byte(key))

	if err != nil {
		fmt.Println(err)
	}

	msgByte := make([]byte, len(message))
	c.Encrypt(msgByte, []byte(message))

	return hex.EncodeToString(msgByte)
}

func DecryptMessage(key string, message string) string {
	txt, _ := hex.DecodeString(message)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(txt))
	c.Decrypt(msgByte, []byte(txt))

	msg := string(msgByte[:])
	return msg
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
