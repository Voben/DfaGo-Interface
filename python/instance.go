// *****************************************************************************
// A problem instance created in DFA to be learned by some algorithm in Python.
// *****************************************************************************

package python

import (
	// "fmt"
	"os"
	"path/filepath"
	"strings"

	"errors"
	"path"
	"runtime"

	"github.com/go-python/cpy3" // creates bindings between python and go (github.com/go-python/cpy3)
	"kguil.com/dfago"
)

// Instance is an instance in a DfaGo-Python problem.
type Instance struct {
	TargetDfa   dfago.Dfa     // The target DFA that Python must find.
	TrainingSet dfago.Dataset // The training set that will be sent to Python for learning a hypothesis.
	TestingSet  dfago.Dataset // The testing set that will be used to evaluate the accuacy of a learnt DFA.
}

type RunInfo map[string]string

// NewAbbadingInstance creates a new Abbadingo problem instance to be learned
// by some algorithm in Python.
//
//   - nominalSize:          nominal size of the target DFA (e.g, 32 states).
//   - exact:                whether the target DFA created will have exactly the nominal size.
//   - trainingSetSize:      the size of the training set.
//   - testingSetSize:       the size of the testing set used to evaluate a hypothesis.
//   - symStructComplete:    whether the training set is symetrically structuarlly complete.
//   - balance:              the propostion of positive-to-negative strings in the training set
//     for it to be considered balanced. Set to < 0.0 for unrestricted.

func NewAbbadingoInstance(nominalSize int, exact bool, trainingSetSize, testingSetSize int, symStructComplete bool, balance float64) Instance {
	// Prep the result.
	instance := Instance{}

retry:
	// Create the Abbadingo instance.
	instance.TargetDfa, instance.TrainingSet, instance.TestingSet = dfago.NewAbbadingoInstance(nominalSize, exact, trainingSetSize, testingSetSize)

	// Check symetrically structurally complete.
	if symStructComplete && !instance.TrainingSet.SymStructurallyComplete(instance.TargetDfa) {
		goto retry
	}

	// Check balanced.
	if balance > 0.0 && !instance.TrainingSet.Balanced(balance) {
		goto retry
	}

	// Done.
	return instance
}

// ---- TESTS FOR HYPOTHESIS DFA -------------------------------------------

// Will return the accuracy of the hypothesis DFA on the instance's testing set
func (instance Instance) GetAccuracy(hypothesis dfago.Dfa) float64 {
	return instance.TestingSet.Accuracy(hypothesis)
}

// Will check if the hypothesis is the same as the instance target DFA
func (instance Instance) FoundTarget(hypothesis dfago.Dfa) bool {
	return instance.TargetDfa.SameAs(&hypothesis)
}

// ---- PYTHON INITILISATION FUNCTIONS -------------------------------------------

// InitPython initialised the Python subsystem.
func InitPython() {
	// initialize a Python interpreter
	python3.Py_Initialize()
}

// Finalise the Python 
func FinalPython() {
	// finalize and initializer a Python interpreter
	python3.Py_Finalize()
}

// Import third party package
func GetPyModule(module string){
	python3.PyRun_SimpleString("import "+ module)
}

// Get Python module from CWD
func GetCWDModule(module string) (*python3.PyObject, error) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// setting the sys variable in python to the current directory to import the python module
	python3.PyRun_SimpleString("import sys\nsys.path.append(\"" + dir + "\")")

	pyModule := python3.PyImport_ImportModule(module)

	if pyModule == nil {
		defer pyModule.DecRef()
		return nil, errors.New("Python Module Not Found")
	}

	return pyModule, nil
}

// Get Python Module object from DfaGo's python sub-directory
func GetDfaGoPyModule(module string) (*python3.PyObject, error) {
	_, filename, _, _ := runtime.Caller(0)

	dir :=  path.Dir(filename)

	// setting the sys variable in python to the current directory to import the python module
	python3.PyRun_SimpleString("import random\nimport re\nimport time\nimport sys\nsys.path.append(\"" + dir + "\")")

	pyModule := python3.PyImport_ImportModule(module)

	if pyModule == nil {
		defer pyModule.DecRef()
		return nil, errors.New("Python Module Not Found")
	}

	return pyModule, nil
}

// Get Function from Module Object
func GetModuleFunc(module *python3.PyObject, function string) (*python3.PyObject, error) {

	method_Dict := python3.PyModule_GetDict(module)
	functionObj := python3.PyDict_GetItemString(method_Dict, function)
	
	if !(functionObj != nil && python3.PyCallable_Check(functionObj)) {  
		defer functionObj.DecRef()
		return nil, errors.New("PyFunction Not Found")
	}

	return functionObj, nil
}

// ---- CONVERT TO PYOBJECTS -------------------------------------------

// Create DFA from DfaGo DFA
func DfaToPyObject(dfa dfago.Dfa) (*python3.PyObject, error) {
	dfaJSON, err := dfago.SerialiseToString(dfa.Serialise())

	if err != nil {
		return nil, err
	}

	Py_Bytes := python3.PyBytes_FromString(dfaJSON)

	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return nil, err
	}

	func_dfa_String, err := GetModuleFunc(PyFolder, "DfaGo_String")

	if err != nil {
		return nil, err
	}

	PyTuple := WrapObjectIntoTuple(Py_Bytes)
	PyDFA := func_dfa_String.CallObject(PyTuple)

	return PyDFA, nil
}

// Change a training set of an instance to a Python Object using a running server
func TrainingSetToPyObjectServer(port int) (*python3.PyObject, error) {

	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return nil, err
	}

	trainSet_Func, err := GetModuleFunc(PyFolder, "getTrainingSetServer")
	
	if err != nil {
		return nil, err
	}

	PyPort := PyObjFromInt(port)

	PyTuple := WrapObjectIntoTuple(PyPort)
	set_list := trainSet_Func.CallObject(PyTuple)

	return set_list, nil
}

// Change a testing set of an instance to a Python Object using a running server
func TestingSetToPyObjectServer(port int) (*python3.PyObject, error) {

	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return nil, err
	}

	testSet_Func, err := GetModuleFunc(PyFolder, "getTestingSetServer")
	
	if err != nil {
		return nil, err
	}

	PyPort := PyObjFromInt(port)
	PyTuple := WrapObjectIntoTuple(PyPort)
	set_list := testSet_Func.CallObject(PyTuple)

	return set_list, nil
}

// Change a provided dataset to a PyObject without using a server
func SetToPyObject(set dfago.Dataset) (*python3.PyObject, error) {
	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return nil, err
	}
	
	pySet_Func, err := GetModuleFunc(PyFolder, "getPySet")
	
	if err != nil {
		return nil, err
	}

	pos_list := []string{"positive"}
	neg_list := []string{"negative"}

	for _, sample := range set.Positive {
		pos_list = append(pos_list, strings.ReplaceAll(sample.String(), " ", ""))
	}

	for _, sample := range set.Negative {
		neg_list = append(neg_list, strings.ReplaceAll(sample.String(), " ", ""))
	}

	set_list_str := strings.Join(pos_list, ",") + "," +strings.Join(neg_list, ",")

	PyBytes := python3.PyBytes_FromString(set_list_str)

	PyTuple := WrapObjectIntoTuple(PyBytes)
	set_list := pySet_Func.CallObject(PyTuple)

	return set_list, nil
}