package api

type configData struct {
	DockerRemoteAddress string
	DockerRemotePort    string
	DockerID            string
}

type imageInspectInfo struct {
	Id              string `json:"id"`
	Container       string `json:""container"`
	Comment         string
	Os              string
	Architecture    string
	Parent          string
	ContainerConfig ContainerConfig
	DockerVersion   string
	VirtualSize     int64
	Size            int64
	RootFS          RootFS
	//----- (possibly) invalid info: -----
	//Author
	//Created
	//GraphDriver
	//RepoDigests
	//RepoTags
	//Config
}

type ContainerConfig struct {
	Tty             bool
	Hostname        string
	Domainname      string
	AttachStdout    bool
	PublishService  string
	AttachStdin     bool
	OpenStdin       bool
	StdinOnce       bool
	NetworkDisabled bool
	OnBuild         []string
	Image           string
	User            string
	WorkingDir      string
	MacAddress      string
	AttachStderr    bool
	Lables          []string
	Env             []string
	Cmd             []string
}

type RootFS struct {
	Type   string
	Layers []string
}

type manifest struct {
	Config   string
	Layers   []string
	RepoTags string
}

type dpkgInfo struct {
	Package string
	Source  string
	Version string
}
