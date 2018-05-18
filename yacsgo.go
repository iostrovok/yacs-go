package main

/*

Run YACS processing via the command line.

*/

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"diff"
	"helper"
	"loader"
	"utils"
)

type container struct {
	command, outDIR, inDIR           string
	sourceFile, copmareFile, outFile string
	verbose, quiet                   bool
	help                             bool
	needResolution                   bool
	needInheritance                  bool
	needValidation                   bool
	countCUPs                        int
	mode                             os.FileMode
	wg                               *sync.WaitGroup
	workFiles                        chan utils.FileForProcess
}

func (con *container) checkOutDir() {
	if con.outDIR == "" {
		con.printSimple("Need to set outDIR params\n")
		os.Exit(0)
	}

	err := utils.CreateDirIfNotExist(con.outDIR, con.mode)
	if err != nil {
		panic(err)
	}
}

func main() {

	con := &container{
		wg:        &sync.WaitGroup{},
		mode:      os.FileMode(0777),
		countCUPs: 4,
	}

	if runtime.NumCPU() > 1 {
		con.countCUPs = runtime.NumCPU()
	}

	var skipResolution, skipInheritance, skipValidation bool

	flag.BoolVar(&con.help, "help", false, "View help message.")
	flag.StringVar(&con.command, "command", "", "What are we doing?")
	flag.StringVar(&con.sourceFile, "file", "", "File which will be processed.")
	flag.StringVar(&con.copmareFile, "copmarefile", "", "File for copmare with 'file'. It's used with 'file' in the same time.")
	flag.StringVar(&con.outFile, "outfile", "", "File for storing result. It's used with 'file' in the same time.")

	flag.StringVar(&con.outDIR, "outdir", "", "Dir for storing result. Dir will be created if it doesn't exist.")
	flag.StringVar(&con.inDIR, "indir", "", "Dir (and all subdirs) which will be processed.")

	flag.BoolVar(&skipResolution, "skip-resolution", false, "Skip reference resolution step. (default \"false\")")
	flag.BoolVar(&skipInheritance, "skip-inheritance", false, "Skip inheritance step. (default \"false\")")
	flag.BoolVar(&skipValidation, "skip-validation", false, "Skip schema validation step. (default \"false\")")

	flag.BoolVar(&con.verbose, "verbose", false, "Shows details about the results of running. (default \"false\")")
	flag.BoolVar(&con.quiet, "quiet", false, "Silent operation. (default \"false\")")

	flag.Parse()

	con.needResolution = !skipResolution
	con.needInheritance = !skipInheritance
	con.needValidation = !skipValidation

	if con.help {
		con.viewhelp()
		return
	}

	switch con.command {
	case "batchdir":
		con.batchdir()
	case "onefile":
		con.onefile()
	case "compare":
		con.compare()
	default:
		con.viewhelp()
	}
}

func (con *container) viewhelp() {

	fmt.Println(`
  -help
        View help message.
  -command string
        What are we doing? May by "batchdir", "onefile", "compare"
  -copmarefile string
        File for copmare with 'file'. It's used with 'file' in the same time.
  -file string
        File which will be processed.
  -indir string
        Dir (and all subdirs) which will be processed.
  -outdir string
        Dir for storing result. Dir will be created if it doesn't exist.
  -outfile string
        File for storing result. It's used with 'file' in the same time.
  -quiet
        Silent operation. (default "false")
  -skip-inheritance
        Skip inheritance step. (default "false")
  -skip-resolution
        Skip reference resolution step. (default "false")
  -skip-validation
        Skip schema validation step. (default "false")
  -verbose
        Shows details about the results of running. (default "false")


Example:
> ./bin/yacsgo -help
> ./bin/yacsgo -verbose=t -command=batchdir -indir=./json-files/ -outdir=./test-out/
> ./bin/yacsgo -verbose=t -command=onefile --file=./mine.json -outfile=./out.json
> ./bin/yacsgo -verbose=t -command=compare -file=./mine.json -copmarefile=./yours.json

`)

}

func (con *container) batchdir() {

	con.print("... command: %s\n    outdir: %s\n    indir: %s", con.command, con.outDIR, con.inDIR)

	startTime := time.Now()

	// checkOutDir(con.outDIR, con.mode)
	con.checkOutDir()
	con.workFiles = make(chan utils.FileForProcess, con.countCUPs*3)

	// Preparing...
	list, err := utils.FindAllFiles(con.inDIR, con.outDIR, "")
	if err != nil {
		con.printSimple("\n------------------------------------------\n ERROR:\ninDIR: %s\noutDIR: %s", con.inDIR, con.outDIR)
		panic(err)
	}

	con.print("Total %d files for processing", len(list))

	for i := 0; i < con.countCUPs; i++ {
		con.process(i)
	}

	// Pushs files for processing to goroutins
	con.pushFiles(list)

	// Closes our end of channel
	close(con.workFiles)

	// It's waiting for goroutines are working.
	con.wg.Wait()

	con.print("... command: %s\n    outdir: %s\n    indir: %s", con.command, con.outDIR, con.inDIR)
	con.print("Total %d files have been processed with %d threads in %.0f seconds", len(list), con.countCUPs, time.Now().Sub(startTime).Seconds())
}

func (con *container) pushFiles(list []utils.FileForProcess) {
	// Pushs files for processing to goroutins
	for i, bf := range list {
		bf.Num = i + 1
		con.workFiles <- bf
	}
}

func (con *container) print(text string, args ...interface{}) {
	if con.verbose && !con.quiet {
		fmt.Printf(text+"\n", args...)
	}
}

func (con *container) printSimple(text string, args ...interface{}) {
	if !con.quiet {
		fmt.Printf(text+"\n", args...)
	}
}

func (con *container) process(number int) {
	con.wg.Add(1)
	go con._process(number)
}

func (con *container) _process(number int) {
	defer con.wg.Done()
	for {
		select {
		case bf, ok := <-con.workFiles:
			if !ok {
				con.print("Thread %d has finished", number)
				return
			}

			err := con.processOneFile(bf.From, bf.To, false)
			if err != nil {
				// Any error should break process and stop all goroutines.
				// Other way we may not find error in tons of logs lines.
				panic(err)
			}
			con.print("%d] %s ===>>> %s", bf.Num, bf.From, bf.To)
		}
	}
}

func (con *container) onefile() {

	con.print("... command: %s\n    file: %s\n    outfile: %s", con.command, con.sourceFile, con.outFile)

	if err := con.processOneFile(con.sourceFile, con.outFile, con.verbose && !con.quiet); err != nil {
		panic(err)
	}

	con.printSimple("Out: %s ===>>> %s", con.sourceFile, con.outFile)
}

func (con *container) compare() {

	con.print("... command: %s\n    file: %s\n    copmarefile: %s", con.command, con.sourceFile, con.copmareFile)
	con.print("Start loading the %s...", con.copmareFile)

	comparebody, err := loader.GetURI(con.copmareFile)
	if err != nil {
		panic(err)
	}

	con.print("Start processing the %s...", con.copmareFile)

	processedDoc, err := helper.Process(con.sourceFile, con.needResolution, con.needInheritance, con.needValidation, false)
	if err != nil {
		panic(err)
	}

	diffres := diff.Diff(comparebody, processedDoc)

	if !con.verbose {
		// print here a short message
		con.printSimple("The files have %d differences\n", len(diffres))
		return
	}

	con.printSimple("\n\nComparison result...\n")

	if len(diffres) == 0 {
		con.printSimple("Congratulation! The files are equal.\n")
		return
	}

	for i, v := range diffres {
		con.printSimple("%d. %v\n", i+1, v)
	}
}

func (con *container) processOneFile(from, to string, verbose bool) error {
	processedDoc, err := helper.Process(from, con.needResolution, con.needInheritance, con.needValidation, verbose)
	if err != nil {
		return err
	}

	return utils.SaveJSONFile(to, processedDoc, con.mode)
}
