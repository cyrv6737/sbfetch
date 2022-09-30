package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func main() {

	rpm_ostree_status := rpmostree_status()

	//Get Hostname
	out, err := os.ReadFile("/etc/hostname")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.TrimSuffix(string(out), "\n"))
	fmt.Println("--------------- General ----------------")

	//Get distro and version
	fmt.Println("Distro:", distro_name()+" "+version_string())

	//Get kernel
	out, err = os.ReadFile("/proc/version")
	if err != nil {
		log.Fatal(err)
	}
	proc_version := strings.Split(string(out), " ")
	fmt.Println("Kernel:", proc_version[2])

	//Get uptime
	out, err = os.ReadFile("/proc/uptime")
	if err != nil {
		log.Fatal(err)
	}
	seconds := strings.Split(string(out), ".")
	real_seconds, err := strconv.Atoi(seconds[0])
	fmt.Println("Uptime:", uptime_calc(real_seconds))

	//Get current shell
	fmt.Println("Shell:", get_shell())

	fmt.Println("--------------- RPM-OStree ----------------")
	fmt.Println("BaseCommit:", basecommit(rpm_ostree_status))
}

func uptime_calc(uptime_seconds int) string {
	days := uptime_seconds / 86400
	hours := (uptime_seconds - days*86400) / 3600
	minutes := (uptime_seconds - days*86400 - hours*3600) / 60
	seconds := uptime_seconds - days*86400 - hours*3600 - minutes*60

	result := strconv.Itoa(days) + " days " + strconv.Itoa(hours) + " hours " + strconv.Itoa(minutes) + " minutes " + strconv.Itoa(seconds) + " seconds "
	return result
}

//Determine distro name from os-release file. No args, returns string.
func distro_name() string {
	out, err := os.ReadFile("/etc/os-release")
	if err != nil {
		log.Fatal(err)
	}
	os_release_string := strings.Split(string(out), "\n")
	for i := 0; i < len(os_release_string); i++ {
		if strings.Contains(os_release_string[i], "NAME=\"") {
			distro_name := strings.TrimPrefix(os_release_string[i], "NAME=\"")
			distro_name = strings.TrimSuffix(distro_name, "\"")
			return distro_name
		}
	}
	return "Unknown Distro"
}

//Same thing as distro_name() but for the version string
func version_string() string {
	out, err := os.ReadFile("/etc/os-release")
	if err != nil {
		log.Fatal(err)
	}
	os_release_string := strings.Split(string(out), "\n")
	for i := 0; i < len(os_release_string); i++ {
		if strings.Contains(os_release_string[i], "VERSION=\"") {
			version_string := strings.TrimPrefix(os_release_string[i], "VERSION=\"")
			version_string = strings.TrimSuffix(version_string, "\"")
			return version_string
		}
	}
	return ""
}

func get_shell() string {

	out, err := os.ReadFile("/etc/passwd")
	if err != nil {
		log.Fatal(err)
	}
	passwd_file := strings.Split(string(out), "\n")
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	username := user.Username
	for i := 0; i < len(passwd_file); i++ {
		if strings.Contains(passwd_file[i], username) {
			user_line := strings.Split(passwd_file[i], ":")
			return user_line[len(user_line)-1]
		}
	}
	return ""
}

//Runs rpm-ostree status, returns output as string
func rpmostree_status() string {
	status, err := exec.Command("rpm-ostree", "status").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(status)
}

func basecommit(status string) string {
	split_status := strings.Split(status, "\n")
	for i := 0; i < len(split_status); i++ {
		if strings.Contains(split_status[i], "BaseCommit") {
			split_line := strings.Split(split_status[i], ":")
			return strings.TrimPrefix(split_line[1], " ")
		}
	}
	return ""
}
