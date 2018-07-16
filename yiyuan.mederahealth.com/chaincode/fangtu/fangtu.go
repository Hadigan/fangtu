/*
This chaincode is implemented to put the md5 value of
consensual letter into blockchain, namely , the ledger.
*/

package main

import (
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Consensual Letter Chaincode implementation
type ConsensualLetterChaincode struct {
}

func (t *ConsensualLetterChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *ConsensualLetterChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "update" {
		// Make payment of X units from A to B
		return t.update(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *ConsensualLetterChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//   0       1
	// "md5", "md5"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	//Input sanitation
	fmt.Println("- start init consensual letter")
	if len(args[0]) != 32 {
		return shim.Error("argument must be a 32bits string")
	} else {
		// regular expression check if the input parameter is a md5 value
		reg := regexp.MustCompile("[[:xdigit:]]{32}")
		if !reg.MatchString(args[0]) {
			return shim.Error("argument must be a hexadecimal string")
		}
	}

	kMd5Letter := args[0]
	vMd5Letter := args[0]

	//Check if consensual letter already exists
	letterAsBytes, err := stub.GetState(kMd5Letter)
	if err != nil {
		return shim.Error("Failed to get letter: " + err.Error())
	} else if letterAsBytes != nil {
		fmt.Println("This letter already exists: " + kMd5Letter)
		return shim.Error("This letter already exists: " + kMd5Letter)
	}

	//Save letter to state
	err = stub.PutState(kMd5Letter, []byte(vMd5Letter))
	if err != nil {
		return shim.Error(err.Error())
	}

	//Letter saved and indexed. Return success
	fmt.Println("- end init letter")
	return shim.Success(nil)
}

// query the consensual letter from the ledger
func (t *ConsensualLetterChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var kMd5Letter string // md5 of the letter, input parameter
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	kMd5Letter = args[0]

	// Get the state from the ledger
	letterValBytes, err := stub.GetState(kMd5Letter)
	if err != nil {
		md5Resp := "{\"Error\":\"Failed to get state for " + kMd5Letter + "\"}"
		return shim.Error(md5Resp)
	}

	if letterValBytes == nil {
		md5Resp := "{\"Error\":\"Nil amount for " + kMd5Letter + "\"}"
		return shim.Error(md5Resp)
	}

	md5Resp := "{\"Name\":\"" + kMd5Letter + "\",\"Amount\":\"" + string(letterValBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", md5Resp)
	return shim.Success(letterValBytes)
}

func main() {
	err := shim.Start(new(ConsensualLetterChaincode))
	if err != nil {
		fmt.Printf("Error starting Consensual Letter chaincode: %s", err)
	}
}
