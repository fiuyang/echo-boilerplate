package config

type Config struct {
	Server   Server
	Database Database
	Swagger  Swagger
	Kong     Kong
	Obs      ObsHuawei
}

type Server struct {
	Port string
}

type Database struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type Swagger struct {
	Host string
	Url  string
	Mode string
}

type Kong struct {
	Url string
}

type ObsHuawei struct {
	Ak       string
	Sk       string
	Endpoint string
	Bucket   string
}
