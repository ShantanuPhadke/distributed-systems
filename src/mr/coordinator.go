package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"

type Coordinator struct {
	// Your definitions here.

	// Map of filenames to booleans keeping track of if each filename has been processed
	processed map[string]bool

}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) GetMapTask(args *GetMapTaskArgs, reply *GetMapTaskReply) error {
	// Loop over the processed map, find the first filename which hasn't been processed or is processing
	// mark that file name as processing, and set reply to be its information.
	for filename, fileProcessed := range c.processed {
		if !fileProcessed {
			c.processed[filename] = true
			reply.Filename = filename
			break
		}
	}
	return nil
}


// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server(sockname string) {
	rpc.Register(c)
	rpc.HandleHTTP()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatalf("listen error %s: %v", sockname, e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.


	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(sockname string, files []string, nReduce int) *Coordinator {
	c := Coordinator{
		processed: make(map[string]bool),
	}

	// Your code here.

	// Pick up the filenames
	for _, filename := range os.Args[2:] {
		// For the mapping phase we'll just have a map between filename and whether it is being processed
		c.processed[filename] = false
	}


	c.server(sockname)
	return &c
}
