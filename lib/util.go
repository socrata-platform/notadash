package lib

import (
    "fmt"
    "net"
    "os"
    "errors"
    "golang.org/x/crypto/ssh/terminal"
)

const CLR_Y = "\x1b[33;1m"
const CLR_N = "\x1b[0m"
const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"


func isTTY() (bool) {
    return terminal.IsTerminal(int(os.Stdin.Fd()))
}

func PrintBool(b bool) (string) {
    if b { return PrintGreen("true") }
    return PrintRed("false")
}


func PrintRed(txt string) (string) {
    if isTTY() {
        return fmt.Sprintf("%s%s%s", CLR_R, txt, CLR_N)
    } else {
        return txt
    }
}


func PrintGreen(txt string) (string) {
    if isTTY() {
        return fmt.Sprintf("%s%s%s", CLR_G, txt, CLR_N)
    } else {
        return txt
    }
}


func PrintYellow(txt string) (string) {
    if isTTY() {
        return fmt.Sprintf("%s%s%s", CLR_Y, txt, CLR_N)
    } else {
        return txt
    }
}


func GetExternalIP() (string, error) {
    if ifaces, err := net.Interfaces(); err != nil {
		return "", err
	} else {
        for _, iface := range ifaces {
            if iface.Flags&net.FlagUp == 0 { // interface down
                continue
            }
            if iface.Flags&net.FlagLoopback != 0 { // loopback interface
                continue
            }
            if addrs, err := iface.Addrs(); err != nil {
                return "", err
            } else {
                for _, addr := range addrs {
                    var ip net.IP
                    switch v := addr.(type) {
                    case *net.IPNet:
                        ip = v.IP
                    case *net.IPAddr:
                        ip = v.IP
                    }
                    if ip == nil || ip.IsLoopback() {
                        continue
                    }
                    ip = ip.To4()
                    if ip == nil { // not an ipv4 address
                        continue
                    }
                    return ip.String(), nil
                }
            }
        }
    }
	return "", errors.New("are you connected to the network?")
}
