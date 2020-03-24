package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/secsy/goftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var mapMutex = sync.RWMutex{}

func downiamgeFromFtp() {
	_ = execcmd()

	imap := jsonToMap()

	client, err, done := getftpClient()
	if done {
		return
	}
	b2 := new(bytes.Buffer)
	err = client.Retrieve("/root/file.txt", b2)
	if err != nil {
		fmt.Println(err)
		return
	}
	reg, _ := regexp.Compile(`/[\dabcdef]{2}/[\dabcdef]{2}/([\w]+-[\d]+-[\d]+-[\d]+-[\w]+)`)
	submatch := reg.FindAllStringSubmatch(b2.String(), -1)
	fmt.Println(" compile complete")
	group := sync.WaitGroup{}
	group.Add(len(submatch))
	for i, ss := range submatch {
		fileid := ss[1]
		split := strings.Split(fileid, "-")
		mapMutex.RLock()
		i2 := imap[split[0]]
		mapMutex.RUnlock()
		if i2 == 1 {
			group.Done()
			continue
		}

		fmt.Printf("第%d个\n", i)
		time.Sleep(time.Millisecond * 10)
		go func() {
			buf := new(bytes.Buffer)
			sprintf := fmt.Sprintf("/root/cache/cache/%s/%s/"+fileid, split[0][0:2], split[0][2:4])
			fmt.Println(sprintf)
			client.Retrieve(sprintf, buf)
			base := cachePath()
			path01 := filepath.Join(base, split[0][0:2])
			notExistorCreat(path01)
			path02 := filepath.Join(path01, split[0][2:4])
			notExistorCreat(path02)
			//join := filepath.Join(base, fileid)
			join := filepath.Join(path02, fileid)
			tfile, err := os.Create(join)
			if err != nil {
				fmt.Println(err)
				return
			}
			n, _ := tfile.Write(buf.Bytes())
			defer func() {
				if err := tfile.Close(); err != nil {
					return
				}
			}()
			mapMutex.Lock()
			imap[split[0]] = 1
			mapMutex.Unlock()
			fmt.Println("写入", n)
			group.Done()
		}()

	}
	group.Wait()
	map2json(imap)

}

func map2json(imap map[string]int) {
	file, _ := os.Create(savedPath())
	defer func() {
		_ = file.Close()
	}()
	bytes, _ := json.Marshal(imap)
	file.Write(bytes)

}

func jsonToMap() map[string]int {
	file, _ := os.OpenFile(savedPath(), os.O_RDWR, os.ModePerm)
	all, _ := ioutil.ReadAll(file)
	imap := make(map[string]int)
	json.Unmarshal(all, &imap)
	return imap
}

func getftpClient() (*goftp.Client, error, bool) {
	config := goftp.Config{
		User:               "root",
		Password:           "225116",
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
	}
	client, err := goftp.DialConfig(config, "router")
	if err != nil {
		fmt.Println(err)
		return nil, nil, true
	}
	return client, err, false
}

func execcmd() error {
	var stdOut, stdErr bytes.Buffer

	session, err := SSHConnect("root", "225116", "192.168.5.1", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Stdout = &stdOut
	session.Stderr = &stdErr

	session.Run(">file;find /opt/docker/volumes/hdata/_data/cache/ -type f >file.txt")
	return err
}

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}
