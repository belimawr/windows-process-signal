A simple POC to gracefully shutdown process on Windows without
affecting the parent process and being able to also test this with `go
test`.

Linux and MacOS are pretty simple, because both support sending SIGINT
and SIGTERM directly to the child process. Windows is a challenge
because `GenerateConsoleCtrlEvent` can send the "CTRL + C" event,
however it is sent to a whole process group, which, by default in Go,
includes the parent process. When writing tests this includes the
process running the tests.

The solution is to create the process setting the creation flag
`CREATE_NEW_PROCESS_GROUP`, which disables sending `CTRL_C_EVENT` to
the child process, there fore the only way to signal a graceful
shutdown is to use `CTRL_BREAK_EVENT`.

## Examples
Creating a new process on Windows with `CREATE_NEW_PROCESS_GROUP` set

```go
cmd = exec.Command("./foo.exe")
cmd.Stderr = os.Stdout
cmd.Stdout = os.Stdout
cmd.SysProcAttr = &syscall.SysProcAttr{
	// This disables the child from receiveing CTRL_C events
	// But isolates other siganls from us, aka the parent
	CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
}
```

Sending the `CTRL_BREAK_EVENT`:
```go
windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(cmd.Process.Pid))
```

## References
 - https://learn.microsoft.com/en-us/windows/console/generateconsolectrlevent
 - https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
 - https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createprocessa
 - https://pkg.go.dev/golang.org/x/sys/windows
 - https://pkg.go.dev/os/exec
 - https://cs.opensource.google/go/go/+/refs/tags/go1.24.2:src/syscall/exec_windows.go;l=245
