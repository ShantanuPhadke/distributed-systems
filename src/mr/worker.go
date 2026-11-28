package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "os"
import "io/ioutil"


// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

var coordSockName string // socket for coordinator


// main/mrworker.go calls this function.
func Worker(sockname string, mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	coordSockName = sockname

	// Your worker implementation here.

	// One way to get started is to modify mr/worker.go's Worker() to send an RPC
	// to the coordinator asking for a task. Then modify the coordinator to respond
	// with the file name of an as-yet-unstarted map task.
	taskReply := CallGetMapTask()
	// Then modify the worker to read that file and call the application Map function,
	// as in mrsequential.go.
	file, err := os.Open(taskReply.Filename)
	if err != nil {
		log.Fatalf("cannot open %v", taskReply.Filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", taskReply.Filename)
	}
	file.Close()
	kva := mapf(taskReply.Filename, string(content))
	fmt.Printf("%v\n", kva)

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()

}

func CallGetMapTask() GetMapTaskReply{
	args := GetMapTaskArgs{}
	reply := GetMapTaskReply{}
	ok := call("Coordinator.GetMapTask", &args, &reply)
	if ok {
		fmt.Printf("reply.Y %v\n", reply)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("unix", coordSockName)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	if err := c.Call(rpcname, args, reply); err == nil {
		return true
	}
	log.Printf("%d: call failed err %v", os.Getpid(), err)
	return false
}
