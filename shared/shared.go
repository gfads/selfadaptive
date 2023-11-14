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
const IpPortRabbitMQ = "192.168.0.20:5672" // Home Recife

// Training/Experiment parameters
const L = 1.0

//const Tau = 1.0  // original
//const T = 0.1 // original
const Tau = 0.1 // modified
const T = 0.01  // modified

//var RandomGoal = []float64{363, 1042, 1871, 2063, 1436, 585, 318, 888, 1754, 2094, 1585, 710, 300, 744, 1621, 2098, 1722}
//var RandomGoal = []float64{500, 1000, 750}
//var RandomGoals = []float64{866, 1440, 866}
var RandomGoals = []float64{866, 1440, 1000}

//var RandomGoals = []float64{1000}

//var RandomGoals = []float64{1000}

//var RandomGoals = []float64{1440}
//var InputSteps = []int{2, 1} // for Ziegler/Cohen/AMIGO
//var InputSteps = []int{4, 1, 11, 33, 60, 84, 98, 98, 84, 60, 33, 11, 1, 5, 22, 48, 74, 93, 100, 92, 72, 45, 20, 3, 1, 13, 35, 62, 86, 99, 97, 82, 57, 31, 9, 1, 6, 24, 50, 76, 94, 100, 91, 69, 43, 18, 3, 1, 14, 38, 65, 87, 99, 97, 80, 55, 28, 8, 1, 7, 26, 52, 78, 95, 100, 89, 67, 40, 16, 2, 2, 16, 40, 67, 89, 100, 96, 79, 53, 26, 7, 1, 8, 28, 54, 80, 96, 99, 88, 65, 38, 15, 1, 2, 17, 42, 69, 90, 100, 95, 77, 51}
var InputSteps = []int{1, 1, 1, 3, 5, 7, 11, 14, 18, 23, 27, 32, 38, 43, 48, 54, 59, 65, 70, 75, 79, 84, 87, 91, 94, 96, 98, 99, 100, 100, 99, 98, 97, 94, 92, 88, 84, 80, 76, 71, 66, 61, 55, 50, 44, 39, 34, 28, 24, 19, 15, 11, 8, 5, 3, 2, 1, 1, 1, 1, 2, 4, 7, 10, 13, 17, 21, 26, 31, 36, 41, 47, 52, 58, 63, 68, 73, 78, 82, 86, 90, 93, 95, 97, 99, 100, 100, 100, 99, 97, 95, 92, 89, 86, 82, 77, 72, 67, 62, 57, 51, 46, 40, 35, 30, 25, 20, 16, 12, 9, 6, 4, 2, 1, 1, 1, 1, 2, 4, 6, 9, 12, 16, 20, 24, 29, 34, 40, 45, 51, 56, 61, 67, 72, 77, 81, 85, 89, 92, 95, 97, 99, 100, 100, 100, 99, 98, 96, 93, 90, 87, 83, 79, 74, 69, 64, 58, 53, 47, 42, 37, 31, 26, 22, 17, 13, 10, 7, 4, 2, 1, 1, 1, 1, 1, 3, 5, 8, 11, 15, 19, 23, 28, 33, 38, 44, 49, 54, 60, 65, 70, 75, 80, 84, 88, 91, 94, 96, 98, 99, 100, 100, 99, 98, 96, 94, 91, 88, 84, 80, 75, 70, 65, 60, 55, 49, 44, 38, 33, 28, 23, 19, 15, 11, 8, 5, 3, 1, 1, 1, 1, 1, 2, 4, 7, 10, 13, 17, 22, 26, 31, 37, 42, 47, 53, 58, 64, 69, 74, 78, 83, 87, 90, 93, 96, 98, 99, 100, 100, 100, 99, 97, 95, 92, 89, 85, 81, 77, 72, 67, 62, 56, 51, 45, 40, 35, 29, 25, 20, 16, 12, 9, 6, 4, 2, 1, 1, 1, 1, 2, 4, 6, 9, 12, 16, 20, 25, 30, 35}
var Kp = map[string]string{
	BasicP + RootLocus:   "0.00448961", // -kp=0.00448961", "-ki=0.00000000", "-kd=0.00000000"
	BasicP + Ziegler:     "0.00022294", // "-kp=0.00022294", "-ki=0.00000000", "-kd=0.00000000"
	BasicP + Cohen:       "0.00100321", // "-kp=0.00100321", "-ki=0.00000000", "-kd=0.00000000"
	BasicP + Amigo:       "0.0",
	BasicPi + RootLocus:  "-0.00111867", //-kp=-0.00111867", "-ki=0.00148840", "-kd=0.00000000"
	BasicPi + Ziegler:    "0.00019285",  // "-kp=0.00019285", "-ki=0.00064284", "-kd=0.00000000"
	BasicPi + Cohen:      "0.00196709",  // "-kp=0.00196709", "-ki=0.07181427", "-kd=0.00000000"
	BasicPi + Amigo:      "0.00037871",  // "-kp=0.00037871", "-ki=0.01035190", "-kd=0.00000000"
	BasicPid + RootLocus: "0.00806",     //-kp=-0.00088937", "-ki=0.00141098", "-kd=0.00032814"
	BasicPid + Ziegler:   "0.00026446",  // "-kp=0.00026446", "-ki=0.00132228", "-kd=0.00001322"
	BasicPid + Cohen:     "0.00083304",  // "-kp=0.00083304", "-ki=0.00788611", "-kd=0.00001063"
	BasicPid + Amigo:     "0.00053993",  // "-kp=0.00053993", "-ki=0.01136109", "-kd=0.00000675"
}
var Ki = map[string]string{
	BasicP + RootLocus:   "0.0",
	BasicP + Ziegler:     "0.0",
	BasicP + Cohen:       "0.0",
	BasicP + Amigo:       "0.0",
	BasicPi + RootLocus:  "0.000153",
	BasicPi + Ziegler:    "0.00064284", //"-kp=0.00019285", "-ki=0.00064284", "-kd=0.00000000"
	BasicPi + Cohen:      "0.07181427",
	BasicPi + Amigo:      "0.01035190", // "-kp=0.00037871", "-ki=0.01035190", "-kd=0.00000000"
	BasicPid + RootLocus: "0.000493",
	BasicPid + Ziegler:   "0.00132228",
	BasicPid + Cohen:     "0.00788611",
	BasicPid + Amigo:     "0.01136109"}
var Kd = map[string]string{
	BasicP + RootLocus:   "0.0",
	BasicP + Ziegler:     "0.0",
	BasicP + Cohen:       "0.0",
	BasicP + Amigo:       "0.0",
	BasicPi + RootLocus:  "0.0",
	BasicPi + Ziegler:    "0.0",
	BasicPi + Cohen:      "0.0",
	BasicPi + Amigo:      "0.0",
	BasicPid + RootLocus: "0.0171",
	BasicPid + Ziegler:   "0.00831",
	BasicPid + Cohen:     "0.00001063",
	BasicPid + Amigo:     "0.00000675"}

const MinPC = "1"
const MaxPC = "100"
const MonitorInterval = "5"
const Adaptability = "true"
const InitialPC = "1"
const SetPoint = "500"
const Direction = "1.0"
const DeadZone = "200.0"
const HysteresisBand = "200.0"
const Alfa = "1.0"
const Beta = "0.9"
const ZieglerRepetitions = 2 * 10

// Experiments parameters
const TrainingSampleSize = 100
const TimeBetweenAdjustments = 1200 // seconds
const MaximumNrmse = 0.30
const WarmupTime = 30 // seconds
const TrainingAttempts = 100

const SizeOfSameLevel = 60 // used in the experiments
//const SizeOfSameLevel = 100 // used in the experiments

// Astar
const SV = 2.7           // Shutoff voltage (page 17) = 2.7 V
const OV = 3.7           // Optimum voltage (page 17) = 3.7 V
const HYSTERISIS = 0.001 // 10 mV

// Prefetch counter
const PcDefaultLimitMin = 1
const PcDefaultLimitMax = 1200 // TODO ASTAR

// Execution Types
const OpenLoop = "OpenLoop"
const StaticGoal = "StaticGoal"
const Experiment = "Experiment"
const InputStep = "InputStep"
const RootTraining = "RootTraining"
const ZieglerTraining = "ZieglerTraining"
const CohenTraining = "CohenTraining"
const AMIGOTraining = "AMIGOTraining"
const OnLineTraining = "OnlineTraining"
const WebTraining = "WebTraining"
const ExperimentalDesign = "ExperimentalDesign"

const DockerFileStatic = "Dockerfile-static"
const DockerFileRoot = "Dockerfile-root"
const DockerFileZiegler = "Dockerfile-ziegler"

// Controller type names

const AsTAR = "AsTAR"
const HPA = "HPA"
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

var Controllers = []string{AsTAR, HPA, BasicOnoff, DeadZoneOnoff, HysteresisOnoff,
	BasicP, BasicPi, BasicPid, SmoothingPid, IncrementalFormPid, ErrorSquarePidProportional, ErrorSquarePidFull,
	DeadZonePid, GainScheduling, PIwithTwoDegreesOfFreedom, SetpointWeighting}

var TunningTypes = []string{
	RootLocus,
	Ziegler,
	Cohen,
	Amigo,
}

//const ExperimentFileBase = "raw-sin-36-static-"
const ZieglerInput = "ziegler-04.csv"
const RootInput = "root-01.csv"
const RootOutput = "root-01-output.csv"
const ExperimentInput = "experiment-4"
const ExperimentOutput = "data-all.csv"
const TrainingInput = "training-experiment-03-75-publishers.csv"
const TrainingOutput = "training-experiment-04-75-publishers-mean.csv"

const RootLocus = "RootLocus"
const Ziegler = "Ziegler"
const Cohen = "Cohen"
const Amigo = "AMIGO"
const None = "None"

// Controller parameters
const MinOnoff = 10
const MaxOnoff = 600

//const MonitorTime = 30 // seconds
const NumberOfColors = 7
const ColorReset = "\033[0m"
const DeltaTime = 5 // see page 103

// Base Directories
const SourcesDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/envrnment/plugins/source"
const ExecutablesDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/envrnment/plugins/executable"
const DockerDir = "/app/data" // it is mapped into windows dir "C:\Users\user\go\selfadaptive\rabbitmq\data" (see execute-old.bat)

// const DataDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/rabbitmq/data/" // macos
const DataDir = "C:\\Users\\user\\go\\selfadaptive\\rabbitmq\\data\\november2023" // macos
const DockerfilesDir = "C:\\Users\\user\\go\\selfadaptive\\temp"                  // macos
const BatchfilesDir = "C:\\Users\\user\\go\\selfadaptive"                         // macos
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

//creating a function to add zeroes to a string
func PadLeft(str string, length int) string {
	for len(str) < length {
		str = "0" + str
	}
	return str
}

// recreate Dockfile dir

func ConfigureDockerfileDir(dir string) {

	// check if directory exist
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		// remove old docker files
		err := os.RemoveAll(dir)
		if err != nil {
			ErrorHandler(GetFunction(), err.Error())
		}
	}

	// recreate docker files folder
	err := os.MkdirAll(dir, 0750)
	if err != nil && !os.IsExist(err) {
		ErrorHandler(GetFunction(), err.Error())
	}

}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsTuned(c string) bool {
	r := true
	if c == BasicOnoff || c == HysteresisOnoff || c == DeadZoneOnoff || c == HPA || c == AsTAR {
		r = false
	}
	return r
}
