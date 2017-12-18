package main

import (
	"io"
	"os"
	"encoding/json"
	"fmt"
	"encoding/base64"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"context"
	"log"
)

func DockerTask(imageName string) {

	log.Printf("docker task for imageName=%s\n", imageName)

	cli := getDockerClient()

	loginDocker(cli)
	buildImage(imageName, cli)
	pushImage(cli, imageName)

	cleanContext()

}

func cleanContext() {
	if err := os.Remove("repo.tar"); err != nil {
		log.Fatal(err)
	}
}

func pushImage(cli *client.Client, imageName string) {
	encodedJSON, err := json.Marshal(getAuthConfig())
	if err != nil {
		panic(err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	r, err := cli.ImagePush(context.Background(), imageName, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, r)
	defer r.Close()
}

func buildImage(imageName string, cli *client.Client) {
	buildOptions := types.ImageBuildOptions{
		Tags:           []string{imageName},
		Dockerfile:     "Dockerfile",
		SuppressOutput: true,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
	}

	Tar(fmt.Sprintf("%s/%s", GetCurrentDirectory(), "output"), "repo.tar")
	dockerBuildContext, err := os.Open("repo.tar")
	defer dockerBuildContext.Close()

	if err != nil {
		log.Fatalf("build context failed:%v", err)
	}

	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatalf("buildImage=%s failed:%v", imageName, err)
	}

	defer buildResponse.Body.Close()
	io.Copy(os.Stdout, buildResponse.Body)
}

func loginDocker(cli *client.Client) {
	auth, err := cli.RegistryLogin(context.Background(), getAuthConfig())
	if err != nil {
		panic(err)
	}
	log.Printf("login response=%s", auth.Status)
}

func getAuthConfig() types.AuthConfig {
	authConfig := types.AuthConfig{
		Username:      "developer@baidao.com",
		Password:      os.Getenv("DOCKER_PASSWORD"),
		ServerAddress: "daocloud.io",
	}
	return authConfig
}

func getDockerClient() *client.Client {
	defaultHeaders := map[string]string{"User-Agent": "ego-v-0.0.1"}
	cli, _ := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	return cli
}
