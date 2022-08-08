package shell

import (
	"fmt"
	"strconv"
	"strings"

	execute "github.com/alexellis/go-execute/pkg/v1"
	log "github.com/sirupsen/logrus"
)

func RunAsSudo(myCommand string) (stdOut string, err error) {
	sudoPrefix := ""
	uid := execute.ExecTask{
		Command: "id",
		Args:    []string{"-u"},
		// Shell:       true,
		StreamStdio: false, // если true то выводит вконсоль и в Stdout
	}
	res, err := uid.Execute()
	if err != nil {
		return myCommand, err
	}

	if res.ExitCode != 0 {
		return myCommand, fmt.Errorf("%s", res.Stderr)
	}

	uidInt, err := strconv.Atoi(strings.Trim(res.Stdout, "\n"))
	if err != nil {
		return myCommand, err
	}
	if uidInt > 0 {
		sudoPrefix = "sudo "
	}
	runCommand := fmt.Sprintf("%s%s", sudoPrefix, myCommand)
	return runCommand, nil
}

// Run executes the local command don't output the result
func Run(myCommand string, stream bool, sudo bool, dryRunResponce string) (stdOut string, err error) {
	runCommand := myCommand
	if sudo {
		runCommand, err = RunAsSudo(myCommand)
		if err != nil {
			return "", err
		}
	}
	if len(dryRunResponce) > 0 {
		log.Infof("DRY-RUN %s", runCommand)
		return dryRunResponce, nil
	}
	log.Debugf("Executing: %s\n", runCommand)
	// log.Infof("Executing: %s", runCommand)
	cmd := execute.ExecTask{
		Command:     runCommand,
		StreamStdio: stream,
	}

	res, err := cmd.Execute()
	if err != nil {
		return "", err
	}

	if res.ExitCode != 0 {
		return "", fmt.Errorf("%s", res.Stderr)
	}

	return strings.Trim(res.Stdout, "\n"), nil
}

// RunLocalCommand runs the local command
// func RunLocalCommand(myCommand string, sudo bool, dryRun bool) (stdOut string, err error) {

// 	cmd := execute.ExecTask{
// 		Command: myCommand,
// 		// Args:        []string{"version"},
// 		StreamStdio: true,
// 	}

// 	res, err := cmd.Execute()
// 	if err != nil {
// 		// panic(err)
// 		log.Fatalln(err.Error())
// 	}

// 	if res.ExitCode != 0 {
// 		// log.Fatalln("Command failed with exit code: " + string(res.Stderr))
// 		return "", fmt.Errorf("%s", res.Stderr)
// 		// panic("Non-zero exit code: " + res.Stderr)
// 	}

// 	return res.Stdout, nil
// }
