package driver

import (
	"fmt"
	"github.com/hashicorp/nomad/client/allocdir"
	"github.com/hashicorp/nomad/client/testutil"
	"github.com/hashicorp/nomad/nomad/mock"
	"github.com/hashicorp/nomad/nomad/structs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mockAllocDir(t *testing.T) (*structs.Task, *allocdir.AllocDir) {
	alloc := mock.Alloc()
	task := alloc.Job.TaskGroups[0].Tasks[0]

	allocDir := allocdir.NewAllocDir(filepath.Join(os.TempDir(), alloc.ID))
	if err := allocDir.Build([]*structs.Task{task}); err != nil {
		log.Panicf("allocDir.Build() failed: %v", err)
	}

	return task, allocDir
}

func TestExec_IsolationWithJobChroot(t *testing.T) {
	testutil.ExecCompatible(t)

	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	task.ChrootEnv = map[string]string{
		"/etc/ld.so.cache":  "/etc/ld.so.cache",
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}

	driverCtx, execCtx := testDriverContexts(task)

	execCtx.AllocDir = allocDir

	//defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	handle, err := d.Start(execCtx, task)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle == nil {
		t.Fatalf("missing handle")
	}
	
	/*
	// Attempt to open
	handle2, err := d.Open(execCtx, handle.ID())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle2 == nil {
		t.Fatalf("missing handle")
	}
	*/
	
	expected :=
		`ld.so.cache
ld.so.conf
ld.so.conf.d/`
	fmt.Printf("%v\n", execCtx.AllocDir.LogDir())
	file := filepath.Join(execCtx.AllocDir.LogDir(), "web.stdout.0")
	fmt.Printf("%v\n", file)
	output, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Couldn't read file %v", file)
	}
	fmt.Printf("%v\n", string(output))
	fmt.Printf("%d\n", len(output))
	act := strings.TrimSpace(string(output))
	if act != expected {
		t.Fatalf("Command output incorrectly: want %v\n got: \n%v", expected, act)
	}
	handle.Kill()
	//handle2.Kill()
}
/*
func TestExec_IsolationWithJobChrootCompatibleToAgentChroot(t *testing.T) {
	testutil.ExecCompatible(t)

	agentChroot := map[string]string{
		"/etc/ld.so.cache":  "/etc/ld.so.cache",
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}
	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	task.ChrootEnv = map[string]string{
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}

	driverCtx, execCtx := testDriverContextsWithChrootEnv(task, agentChroot)

	execCtx.AllocDir = allocDir

	//defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	handle, err := d.Start(execCtx, task)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle == nil {
		t.Fatalf("missing handle")
	}

	// Attempt to open
	handle2, err := d.Open(execCtx, handle.ID())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle2 == nil {
		t.Fatalf("missing handle")
	}

	expected :=
		`ld.so.conf
ld.so.conf.d/`

	file := filepath.Join(execCtx.AllocDir.LogDir(), "web.stdout.0")
	output, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Couldn't read file %v", file)
	}

	act := strings.TrimSpace(string(output))
	if act != expected {
		t.Fatalf("Command output incorrectly: want %v\n got: \n%v", expected, act)
	}
	handle.Kill()
	handle2.Kill()
}

func TestExec_IsolationWithJobChrootIncompatibleToAgentChroot(t *testing.T) {
	testutil.ExecCompatible(t)

	agentChroot := map[string]string{
		"/etc/ld.so.cache":  "/etc/ld.so.cache",
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}
	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	task.ChrootEnv = map[string]string{
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/etc":              "/etc",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}

	driverCtx, execCtx := testDriverContextsWithChrootEnv(task, agentChroot)

	execCtx.AllocDir = allocDir

	//defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	_, err := d.Start(execCtx, task)
	if err == nil {
		t.Fatalf("Job Chroot was not subset of Agent Chroot, Command should not have started.")
	}

}

func TestExec_IsolationWithJobChrootCompatibleToDefaultChroot(t *testing.T) {
	testutil.ExecCompatible(t)

	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	task.ChrootEnv = map[string]string{
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}

	driverCtx, execCtx := testDriverContexts(task)

	execCtx.AllocDir = allocDir

	defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	handle, err := d.Start(execCtx, task)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle == nil {
		t.Fatalf("missing handle")
	}

	// Attempt to open
	handle2, err := d.Open(execCtx, handle.ID())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle2 == nil {
		t.Fatalf("missing handle")
	}

	expected :=
		`ld.so.conf
ld.so.conf.d/`

	file := filepath.Join(execCtx.AllocDir.LogDir(), "web.stdout.0")
	output, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Couldn't read file %v", file)
	}

	act := strings.TrimSpace(string(output))
	if act != expected {
		t.Fatalf("Command output incorrectly: want %v\n got: \n%v", expected, act)
	}
	handle.Kill()
	handle2.Kill()
}

func TestExec_IsolationWithJobChrootIncompatibleToDefaultChroot(t *testing.T) {
	testutil.ExecCompatible(t)
	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	task.ChrootEnv = map[string]string{
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/etc":              "/etc",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
		"/foo/bar":          "/foo/bar",
	}

	driverCtx, execCtx := testDriverContexts(task)

	execCtx.AllocDir = allocDir

	defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	_, err := d.Start(execCtx, task)
	if err == nil {
		t.Fatalf("Job Chroot was not subset of Default Chroot when no Agent Chroot given, Command should not have started.")
	}

}

func TestExec_IsolationWithAgentChrootCompatibleToDefaultChroot(t *testing.T) {
	testutil.ExecCompatible(t)

	agentChroot := map[string]string{
		"/etc/ld.so.cache":  "/etc/ld.so.cache",
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
	}
	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	driverCtx, execCtx := testDriverContextsWithChrootEnv(task, agentChroot)

	execCtx.AllocDir = allocDir

	//defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	handle, err := d.Start(execCtx, task)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle == nil {
		t.Fatalf("missing handle")
	}

	// Attempt to open
	handle2, err := d.Open(execCtx, handle.ID())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if handle2 == nil {
		t.Fatalf("missing handle")
	}

	expected :=
		`ld.so.cache
ld.so.conf
ld.so.conf.d/`
	file := filepath.Join(execCtx.AllocDir.LogDir(), "web.stdout.0")
	output, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Couldn't read file %v", file)
	}

	act := strings.TrimSpace(string(output))
	if act != expected {
		t.Fatalf("Command output incorrectly: want %v\n got: \n%v", expected, act)
	}
	handle.Kill()
	handle2.Kill()
}

func TestExec_IsolationWithAgentChrootIncompatibleToDefaultChroot(t *testing.T) {
	testutil.ExecCompatible(t)

	agentChroot := map[string]string{
		"/etc/ld.so.cache":  "/etc/ld.so.cache",
		"/etc/ld.so.conf":   "/etc/ld.so.conf",
		"/etc/ld.so.conf.d": "/etc/ld.so.conf.d",
		"/lib":              "/lib",
		"/lib64":            "/lib64",
		"/usr/lib":          "/usr/lib",
		"/bin/ls":           "/bin/ls",
		"/foo/bar":          "/foo/bar",
	}
	task, allocDir := mockAllocDir(t)

	task.Config = map[string]interface{}{
		"command": "/bin/ls",
		"args":    []string{"-F", "/etc/"},
	}

	driverCtx, execCtx := testDriverContextsWithChrootEnv(task, agentChroot)

	execCtx.AllocDir = allocDir

	defer execCtx.AllocDir.Destroy()

	d := NewExecDriver(driverCtx)

	_, err := d.Start(execCtx, task)
	if err == nil {
		t.Fatalf("Agent Chroot was not subset of Default Chroot, Command should not have started.")
	}

}
*/