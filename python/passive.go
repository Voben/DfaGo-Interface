// *****************************************************************************
// DfaGo-to-Python bridge for passive learning algorithms (e.g., RPNI).
// *****************************************************************************

package python

import (
	"fmt"
	"time"

	"kguil.com/dfago"
)

// Calls Python RPNI learning function
func RunPythonRPNI(TrainingSet dfago.Dataset) {
	module, err := GetDfaGoPyModule("PyFolder")
	defer module.DecRef()

	if err != nil {
		fmt.Println(err)
	}

	rpni_Func, err := GetModuleFunc(module, "RPNIlearn")
	defer rpni_Func.DecRef()

	if err != nil {
		fmt.Println(err)
	}
	
	retSet, err := SetToPyObject(TrainingSet)
	defer retSet.DecRef()

	if err != nil {
		fmt.Println(err)
	}

	PyTuple := WrapObjectIntoTuple(retSet)
	defer PyTuple.DecRef()

	start := time.Now()
	rpni_Func.CallObject(PyTuple)
	end := time.Since(start)

	fmt.Println("Time Taken:", end)
}