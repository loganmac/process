// +build windows

package main

var (
	taskSpinner       = []string{"|", "/", "-", "\\"}
	taskSpinnerLength = len(taskSpinner)
	taskSuccess       = "√"
	taskFailure       = "×"
	taskPause         = "*"
)

func getSpinner(pos int) string {
	return taskSpinner[pos%taskSpinnerLength]
}
