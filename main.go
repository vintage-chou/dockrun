package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall"

	docker "github.com/fsouza/go-dockerclient"
)

const TempFile = "/home/codes/repo/dockrun/tmp"

type DriverSet struct {
	Name string `json:"name"`
	Id   string `json:"-"`
}

type IQuery interface {
	Query() string
}

type Query struct {
	Content string
}

func NewQuery(c string) *Query {
	return &Query{Content: c}
}

func (q *Query) Query() string {
	return q.Content
}

func connect() (*docker.Client, error) {

	// grab directly from docker daemon
	endpoint := "unix:///var/run/docker.sock"

	var client *docker.Client
	var err error

	client, err = docker.NewClient(endpoint)

	env, _ := client.Version()

	fmt.Println(env.Get("ApiVersion"))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func listContainers() {

	cli, _ := connect()

	// stat, _ := os.Stdin.Stat()
	// fmt.Println(stat.Mode() & os.ModeCharDevice)

	containers, err := cli.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Println("error", err)
	}

	for i, c := range containers {
		fmt.Println(i, c.Image, c.Names)
	}
}
func listImages() {

	cli, _ := connect()

	clientImages, _ := cli.ListImages(docker.ListImagesOptions{All: true})
	image := clientImages[0]

	var previous string
	var newID string
	history, _ := cli.ImageHistory(image.ID)
	for i := len(history) - 1; i >= 0; i-- {

		h := sha256.New()
		h.Write([]byte(previous))
		h.Write([]byte(history[i].CreatedBy))
		h.Write([]byte(strconv.FormatInt(history[i].Created, 10)))
		h.Write([]byte(strconv.FormatInt(history[i].Size, 10)))
		newID = fmt.Sprintf("synth:%s", hex.EncodeToString(h.Sum(nil)))
		//fmt.Println(history[i].Tags)
	}
	fmt.Println("newID", newID)

	// for _, im := range clientImages {
	// 	fmt.Println(im.ID, im.ParentID)
	// }

	var imagesByParent = make(map[string][]string)

	for _, image := range clientImages {
		var s_id string
		if len(image.ParentID) == 0 {
			s_id = image.ParentID
		} else {
			s_id = image.ParentID[7:17]
		}

		if children, exists := imagesByParent[s_id]; exists {
			imagesByParent[s_id] = append(children, image.ID[7:17])
		} else {
			imagesByParent[s_id] = []string{image.ID[7:17]}
		}
	}

	fmt.Println(imagesByParent)
}

func collectRoots(images *[]docker.APIImages) []docker.APIImages {
	var roots []docker.APIImages
	for _, image := range *images {
		if image.ParentID == "" {
			roots = append(roots, image)
		}
	}

	return roots
}

type Dvil struct {
	Type bool
	Flag int
}

type DevilIf interface {
}

func NewDvil() DevilIf {
	return &Dvil{Type: true, Flag: 0}
}

func main() {
	//	listImages()
	var buff syscall.Statfs_t

	d := NewDvil()

	fmt.Printf("%T\n", d)

	str := "18 62 0:17 / /sys rw,nosuid,nodev,noexec,relatime shared:6 - sysfs sysfs rw,seclabel"
	//	str := "6 - sysfs sysfs rw,seclabel"

	index := strings.Index(str, " - ")
	fmt.Println(strings.Fields(str[index+3:]))

	if err := syscall.Statfs("/root", &buff); err != nil {
		fmt.Println(err)
	}

	ds := DriverSet{Name: "dm", Id: "123456"}

	cd, _ := json.Marshal(ds)

	fmt.Println(string(cd))

	fmt.Println("fs type:", buff.Type)

	q := NewQuery("hello")
	var iq IQuery

	iq = q

	if qr, ok := iq.(*Query); ok {
		fmt.Println(qr.Query())
	} else {
		fmt.Println(ok)
	}
}
