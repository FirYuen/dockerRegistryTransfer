package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/types/ref"
	"strings"
)

var regIP = flag.String("reg", "docker.io", "registry ip")
var pullType = flag.String("type", "docker", "pull type")

func main() {
	flag.Parse()
	//registryIP := "http://192.168.30.157:5000"

	registryIP := *regIP
	if registryIP == "docker.io" {
		fmt.Println("reg not set ")
		return
	}
	if !strings.Contains(registryIP, "http") {
		registryIP = "http://" + registryIP
	}
	registryCfd := config.HostNewName(registryIP)

	registryCfd.Hostname = registryIP

	reg := regclient.New(regclient.WithConfigHost(*registryCfd))

	ctx := context.Background()
	resp, err := reg.RepoList(ctx, registryIP)
	if err != nil {
		panic(err)
	}
	for _, img := range resp.Repositories {
		//fmt.Println(img)
		reff, err := ref.New(fmt.Sprintf("%s", img))
		reff.Repository = img
		reff.Registry = registryIP
		removeHttp := strings.ReplaceAll(reff.Registry, "http://", "")
		if err != nil {
			//fmt.Println(err)
		}
		imgTags, err := reg.TagList(ctx, reff)
		if err != nil {
			//fmt.Println(err)
		} else {
			for _, tag := range imgTags.Tags {
				regName := fmt.Sprintf("%s/%s:%s", removeHttp, reff.Repository, tag)
				localName := fmt.Sprintf("%s:%s", reff.Repository, tag)
				//./skopeo copy --src-tls-verify=false --insecure-policy  --dest-tls-verify=false docker-daemon:gruebel/upx:latest dir:test
				if *pullType == "docker" {
					fmt.Printf("docker pull %s \ndocker tag %s %s\ndocker rmi %s\n", regName, regName, localName, regName)
				} else {
					fmt.Printf("./skopeo copy --src-tls-verify=false --insecure-policy  --dest-tls-verify=false docker://%s docker-daemon:%s\n", regName, localName)
				}
			}
		}

	}

}
