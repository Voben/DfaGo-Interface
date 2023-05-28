// *****************************************************************************
// DfaGo-to-Python bridge for active learning algorithms (e.g., L*).
// *****************************************************************************

package python

import (
	"fmt"
	"net"
	"strconv"

	"log"

	python3 "github.com/go-python/cpy3"
	"kguil.com/dfago"
)

// ---- MEMBERSHIP OR EQUIVALENCE QUERIES FOR HYPOTHESIS DFA -------------------------------------------

// Will perform a modified BFS to return the shortest counter-example of 2
// DFAs (the target DFA in the instance and a hypothesis DFA)
func (instance Instance) EquivalenceQuery(hypothesis dfago.Dfa) string {
	visited := map[string]int{}

	queue := [][]string{
		{"", strconv.Itoa(instance.TargetDfa.StartingState)},
	}

	node_slice := []string{}

	var idx int

	if instance.TargetDfa.Parse(dfago.BinaryStringToSample("")) != hypothesis.Parse(dfago.BinaryStringToSample("")) {
		return ""
	}

	for {
		if len(queue) <= 0 {
			return "None"
		}

		node_slice, queue = queue[0], queue[1:]

		val, exists := visited[node_slice[1]]

		if exists {
			visited[node_slice[1]]++
		} else {
			visited[node_slice[1]] = 1
		}

		if val < 2 {
			idx, _ = strconv.Atoi(node_slice[1])

			for char, next_idx := range instance.TargetDfa.States[idx].Next {

				if instance.TargetDfa.Parse(dfago.BinaryStringToSample(node_slice[0]+strconv.Itoa(char))) != hypothesis.Parse(dfago.BinaryStringToSample(node_slice[0]+strconv.Itoa(char))) {
					return node_slice[0] + strconv.Itoa(char)
				}

				queue = append(queue, []string{node_slice[0] + strconv.Itoa(char), strconv.Itoa(next_idx)})
			}
		}
	}
}

// Will parase the word and return the result when parsed by the target DFA
func (instance Instance) MembershipQuery(word string) bool {
	sample_word := dfago.BinaryStringToSample(word)
	return instance.TargetDfa.Parse(sample_word)
}

// ---- SERVER RELATED FUNCTIONS -------------------------------------------

// Start TCP server on the given port number and use the buffer size to be 
// used for the size of the messages being sent between the server and client
func (instance Instance) TCPServer(portNum, bufferSize int) {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(portNum))

	if err != nil {
		log.Fatal(err)	
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		
		if err != nil {
			log.Fatal(err)
		}

		go instance.HandlerFunction(conn, bufferSize)
	}
}

// Handler function for the TCP server
func (instance Instance) HandlerFunction(conn net.Conn, bufferSize int) {
	for {
		buffer := make([]byte, bufferSize)
		bufLen, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Read returns error...")
			log.Fatal(err)
		}

		query := string(buffer[:bufLen])

		if query[:2] == "mq" {
			word := query[3:]

			if instance.MembershipQuery(word) == true {
				conn.Write([]byte("true"))
			} else {
				conn.Write([]byte("false"))
			}
		} else if query[:2] == "eq" {
			desObj, _ := dfago.DeserialiseFromBytes([]byte(query[3:]))
			hypo_dfa := dfago.DeserialiseDfa(desObj)
			ret := instance.EquivalenceQuery(hypo_dfa)
			conn.Write([]byte(ret))
		} else if query[:3] == "end" {
			conn.Close()
			break
		} else {
			conn.Write([]byte("Not a query"))
		}
	}
}

// ---- RUN PYTHON LSTAR ALGORITHM -------------------------------------------

func RunLstar(portNum, bufferSize int) dfago.Dfa {
	PyFolder, err := GetDfaGoPyModule("PyFolder")
	defer PyFolder.DecRef()

	if err != nil {
		fmt.Println("Error not nil after GetPyModule")
	}

	init_func, err := GetModuleFunc(PyFolder, "LearnLstar")
	defer init_func.DecRef()
	
	if err != nil{
		fmt.Println("Error in init")
	}
	PyTuple_s3 := python3.PyTuple_New(3)
	defer PyTuple_s3.DecRef()

	ipPy := python3.PyBytes_FromString("localhost")
	portPy := python3.PyLong_FromLong(8000)
	bufferSizePy := python3.PyLong_FromLong(8192)
	defer ipPy.DecRef()
	defer portPy.DecRef()
	defer bufferSizePy.DecRef()

	python3.PyTuple_SetItem(PyTuple_s3, 0, ipPy)
	python3.PyTuple_SetItem(PyTuple_s3, 1, portPy)
	python3.PyTuple_SetItem(PyTuple_s3, 2, bufferSizePy)

	ret_DFA_py := init_func.CallObject(PyTuple_s3)
	defer ret_DFA_py.DecRef()

	JSON_string := python3.PyBytes_AsString(ret_DFA_py)

	deserObj, _ := dfago.DeserialiseFromBytes([]byte(JSON_string))

	ret_dfa := dfago.DeserialiseDfa(deserObj)
	return ret_dfa
}