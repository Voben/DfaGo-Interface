// *****************************************************************************
// Useful functions for the Python interface
// *****************************************************************************

package python

import (
	"errors"
	"github.com/go-python/cpy3" // creates bindings between python and go (github.com/go-python/cpy3)
)

// ---- PYBYTES/PYSTRING FUNCTIONS FOR PYOBJECTS -------------------------------------------

// Functions to change Strings to different PyObjects:
// 		PyBytes
func PyBytesFromString(str string) *python3.PyObject{
	return python3.PyBytes_FromString(str)
}

// 		PyByteArray
func PyByteArrayFromString(str string) *python3.PyObject{
	return python3.PyByteArray_FromStringAndSize(str)
}

// 		A PyString from Go String
// 	Uses a function to change a bytearray to string
func PyStringFromString(str string) (*python3.PyObject, error){
	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return nil, err
	}

	pyString, err := GetModuleFunc(PyFolder, "BytesToString")
	
	if err != nil {
		return nil, err
	}

	PyObj := PyByteArrayFromString(str)
	
	PyTuple := WrapObjectIntoTuple(PyObj)

	PyString := pyString.CallObject(PyTuple)

	return PyString, nil
}

// These functions will change PyObjects to Go Strings:
// 		PyBytes
func StringFromPyBytes(PyObj *python3.PyObject) string{
	return python3.PyBytes_AsString(PyObj)
}

// 		PyByteArray
func StringFromPyByteArray(PyObj *python3.PyObject) string{
	return python3.PyByteArray_AsString(PyObj)
}

// 		Python String
func StringFromPyString(PyObj *python3.PyObject) (string, error){
	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return "", err
	}

	pyString, err := GetModuleFunc(PyFolder, "StringToBytes")
	
	if err != nil {
		return "", err
	}
	
	PyTuple := WrapObjectIntoTuple(PyObj)

	ret := pyString.CallObject(PyTuple)

	str_Ret := StringFromPyBytes(ret)

	return str_Ret, nil
}

// ---- PYDICT FUNCTIONS FOR PYOBJECTS -------------------------------------------
// Creates an empty Python dictionary object
func CreateEmptyPyDict() *python3.PyObject{
	ret := python3.PyDict_New()

	return ret
}

// Will set a key, value pair given the PyObjects of a dictionary
// and the pairs' objects
func SetKeyValue(dict, k, v *python3.PyObject) (*python3.PyObject, error){
	err := python3.PyDict_SetItem(dict, k, v)

	if err != 0{
		return nil, errors.New("Error when setting key value pair")
	}

	return dict, nil
}

// This function gets a value from a dictionary given a key
func GetVal(dict, key *python3.PyObject) *python3.PyObject{
	valObj := python3.PyDict_GetItem(dict, key)

	return valObj
}

// ---- PYLONG FUNCTIONS FOR PYOBJECTS -------------------------------------------
// Different functions to create Python number variables
// depending on the type of Go variable

func PyObjFromInt(num int) *python3.PyObject{
	return python3.PyLong_FromGoInt(num)
}

func PyObjFromInt64(num int64) *python3.PyObject{
	return python3.PyLong_FromGoInt64(num)
}

func PyObjFromUint(num uint) *python3.PyObject{
	return python3.PyLong_FromGoUint(num)
}

func PyObjFromUint64(num uint64) *python3.PyObject{
	return python3.PyLong_FromGoUint64(num)
}

func PyObjFromFloat64(num float64) *python3.PyObject{
	return python3.PyLong_FromGoFloat64(num)
}

// And functions to obtain Go variables from PyObjects
func IntFromPyObj(PyObj *python3.PyObject)  int {
	return python3.PyLong_AsLong(PyObj)
}

func Int64FromPyObj(PyObj *python3.PyObject)  int64 {
	return python3.PyLong_AsLongLong(PyObj)
}

func UintFromPyObj(PyObj *python3.PyObject)  uint {
	return python3.PyLong_AsUnsignedLong(PyObj)
}

func Uint64FromPyObj(PyObj *python3.PyObject)  uint64 {
	return python3.PyLong_AsUnsignedLongLong(PyObj)
}

func Float64FromPyObj(PyObj *python3.PyObject)  float64 {
	return python3.PyLong_AsDouble(PyObj)
}

// ---- PYTUPLE FUNCTIONS FOR PYOBJECTS -------------------------------------------

// These functions wrap either wrap one Python object in a PyTuple or
// multiple in a larger sized PyTuple
func WrapObjectIntoTuple(object *python3.PyObject) *python3.PyObject {
	PyTuple := python3.PyTuple_New(1)
	python3.PyTuple_SetItem(PyTuple, 0, object)

	return PyTuple
}

func WrapObjectsIntoTuple(objects []*python3.PyObject) *python3.PyObject {
	PyTuple := python3.PyTuple_New(len(objects))

	for i := 0; i < len(objects); i++{
		python3.PyTuple_SetItem(PyTuple, 0, objects[i])		
	}

	return PyTuple
}

// Create PyTuple given a size
func NewPyTuple(size int) *python3.PyObject{
	PyTuple := python3.PyTuple_New(size)

	return PyTuple
}

// ---- PYLIST FUNCTIONS FOR PYOBJECTS -------------------------------------------

// Create a new PyList with the provided size
func NewPyList(size int) *python3.PyObject {
	PyList := python3.PyList_New(size)

	return PyList
}

// Append a PyObject to a PyList
func AppendToPyList(lst, obj *python3.PyObject) *python3.PyObject {
	python3.PyList_Append(lst, obj)

	return lst
}

// Change a slice of PyObjects to a PyList
func SliceToPyList(objs []*python3.PyObject) *python3.PyObject {
	PyList := python3.PyList_New(len(objs))

	for i, obj := range objs{
		python3.PyList_SetItem(PyList, i, obj)
	}

	return PyList
}

// Set PyObject to index in PyObject
func PyObjToIndex(lst *python3.PyObject, obj *python3.PyObject, idx int) *python3.PyObject {
	python3.PyList_SetItem(lst, idx, obj)

	return lst
}

// ---- REFERENCE FUNCTIONS FOR PYOBJECTS -------------------------------------------

// Increment reference count of a Python object
func IncrementRef(PyObj *python3.PyObject) *python3.PyObject {
	PyObj.IncRef()

	return PyObj
}

// Decrement reference count of a Python object
func DecrementRef(PyObj *python3.PyObject) *python3.PyObject {

	PyObj.DecRef()

	return PyObj
}

// Decrement reference count of all the Python objects in a slice
func DecrementRefSlice(PyObjs []*python3.PyObject) []*python3.PyObject{
	for _, PyObj := range PyObjs{
		PyObj.DecRef()
	}

	return PyObjs
}

// ---- MISC FUNCTIONS FOR PYOBJECTS -------------------------------------------

// Will print a PyObject from Python
func PrintFromPython(PyObj *python3.PyObject) {
	PyFolder, err := GetDfaGoPyModule("PyFolder")

	if err != nil {
		return 
	}

	pyPrint, err := GetModuleFunc(PyFolder, "printPyObject")
	
	if err != nil {
		return 
	}

	PyTuple := WrapObjectIntoTuple(PyObj)

	pyPrint.CallObject(PyTuple)
}

// Check if the Python interpreter is initilized
func CheckPython() bool {
	return python3.Py_IsInitialized()
}