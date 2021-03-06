package core

import (
	"os/exec"
	"github.com/pkg/errors"
	"fmt"
	"strings"
	"strconv"
)

func latestImageName(repoName string) string {
	return fmt.Sprintf("%s:latest", repoName)
}

func pullDockerImage(imageName string) error {
	cmd := exec.Command("docker", "pull", imageName)
	err := cmd.Run()
	return errors.Wrapf(err, "pullImage(%s) failed", imageName)
}

func runDockerContainer(imageName string, runArgs ...string) error {

	args := []string{"run", "-d"}
	args = append(args, runArgs...)
	args = append(args, imageName)

	cmd := exec.Command("docker", args...)
	return errors.Wrapf(cmd.Run(), "run(%s) failed", imageName)
}

func stopDockerContainer(containerName string) error {
	cmd := exec.Command("docker", "stop", containerName)
	return errors.Wrapf(cmd.Run(), "stop(%s) failed", containerName)
}

func getContainerIds(imageName string) ([]string, error) {
	argStr := fmt.Sprintf("ps -a -q --filter ancestor=%s", imageName)
	args := strings.Split(argStr, " ")

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "getContainerIds(%s) failed : %s", imageName, output)
	}

	containerIdStr := string(output)
	if containerIdStr == "" {
		return nil, nil
	}

	containerIdStr = strings.TrimSpace(containerIdStr)
	return strings.Split(containerIdStr, "\n"), nil
}

func removeDockerContainers(imageName string) error {
	const acceptableError = "exit status 1"

	containerIds, err := getContainerIds(imageName)

	if err != nil {
		return errors.Wrapf(err, "removeDockerContainers(%s) failed", imageName)
	}

	containerIdStr := strings.Join(containerIds, " ")
	argStr := fmt.Sprintf("rm %s --force", containerIdStr)
	args := strings.Split(argStr, " ")

	// todo: collect system return code=1
	cmd := exec.Command("docker", args...)
	err = cmd.Run()
	if err != nil {
		switch err.Error() {
		case acceptableError:
			return nil
		default:
			return errors.Wrapf(cmd.Run(), "removeContainers(%s) failed", imageName)
		}
	}
	return nil
}

func isContainerRunning(containerName string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "-f", "'{{.State.Running}}'", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrapf(err, "inspect(%s) failed", containerName)
	}

	outputStr := string(output)
	outputStr = strings.TrimSpace(outputStr)
	outputStr = strings.Trim(outputStr, "'")

	isRunning, err := strconv.ParseBool(outputStr)

	if err != nil {
		return false, errors.Wrapf(err, "inspect(%s) failed", containerName)
	}

	return isRunning, nil
}
