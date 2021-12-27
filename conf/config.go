package conf

import (
	"git.solutions.im/Solutions.IM/goUtils/env"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DbHost             string
	DbUser             string
	DbPassword         string
	DbDatabase         string
	DbLog              bool
	DbInit             bool
	ListenAddr         string
	BaseUrl            string
	Version            string
	OpenIdURL          string
	OpenIdClientID     string
	OpenIdClientSecret string
	OpenIdLogoutPath   string
	S3Config           S3Config
}

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Region    string
	Bucket    string
}

func (c *Config) Load() {
	c.DbHost = env.GetEnvOrElse("DB_HOST", "postgres.agritracking.svc.cluster.local:5432")
	c.DbUser = env.GetEnvOrElse("DB_USER", "agritracking")
	c.DbPassword = env.GetEnvOrElse("DB_PASSWORD", "li7keegh4aexiToo")
	c.DbDatabase = env.GetEnvOrElse("DB_DATABASE", "agritracking")
	dbLog, err := env.GetBoolEnvOrElse("DB_LOG", false)
	if err != nil {
		log.Fatal(err)
	}
	c.DbLog = dbLog
	dbInit, err := env.GetBoolEnvOrElse("DB_INIT", true)
	if err != nil {
		log.Fatal(err)
	}
	c.DbInit = dbInit
	c.ListenAddr = env.GetEnvOrElse("LISTEN_ADDR", "0.0.0.0:4000")
	c.BaseUrl = env.GetEnvOrElse("BASE_URL", "http://localhost:4000/")
	c.OpenIdURL = env.GetEnvOrElse("OPENID_URL", "https://auth.solutions.im/auth/realms/solutions")
	c.OpenIdClientID = env.GetEnvOrElse("OPENID_CLIENT_ID", "test")
	c.OpenIdClientSecret = env.GetEnvOrElse("OPENID_CLIENT_SECRET", "0c862cea-64ca-4b07-b50f-2dca81a7a0b2")
	c.OpenIdLogoutPath = env.GetEnvOrElse("OPENID_LOGOUT_PATH", "auth/realms/solutions/protocol/openid-connect/logout")

	c.S3Config.Endpoint = env.GetEnvOrElse("S3_ENDPOINT", "minio.solutions.im")
	c.S3Config.AccessKey = env.GetEnvOrElse("S3_ACCESS_KEY", "U65Z81EH9S39NTNYZ71U")
	c.S3Config.SecretKey = env.GetEnvOrElse("S3_SECRET_KEY", "tOuTEgkp3gkAWsf3acExkxpBn+EYLbfJxJmHephF")
	c.S3Config.Region = env.GetEnvOrElse("S3_REGION", "us-east-1")
	c.S3Config.Bucket = env.GetEnvOrElse("S3_BUCKET", "agri")
	ssl, err := env.GetBoolEnvOrElse("S3_SSL", true)
	if err != nil {
		log.Fatal(err)
	}
	c.S3Config.UseSSL = ssl
	setupLogger()

	log.Infof(`
starting server with the following configuration :
- Database Host : %s
- Database Name: %s
- Database User : %s
- Listen to : %s
- Base URL : %s`, c.DbHost, c.DbDatabase, c.DbUser, c.ListenAddr, c.BaseUrl)
}

func setupLogger() {
	log.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "",
		HideKeys:        true,
		NoColors:        false,
		NoFieldsColors:  false,
		ShowFullLevel:   true,
		TrimMessages:    false,
	})
	log.SetLevel(log.DebugLevel)
}
