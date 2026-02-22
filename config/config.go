package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/urfave/cli/v2"
	starcomm "starbase.ag/liquidity/liquid/common"
)

const PoolFeeTopic = "poolFee"

var ConfigFlag = &cli.StringFlag{
	Name:     "config",
	Value:    "./config.toml",
	Aliases:  []string{"c"},
	Usage:    "Path to config file",
	Required: false,
}
var _conf Config

// DBConfig configures the postgres database.
type DBConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Retry    int    `toml:"retry"`
}

// RedisConfig.
type RedisConfig struct {
	DB       int    `toml:"db"`
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

// Configures the server.
type ServerConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	WriteTimeout int    `toml:"timeout"`
}

type KafkaConfig struct {
	Topic   string `toml:"topic"`
	Brokers string `toml:"brokers"`
	Group   string `toml:"group"`
}

// ProviderConfig chain provider RPC and interval.
type ProviderConfig struct {
	RPC     string `toml:"rpc"`
	Tps     uint64 `toml:"tps"`
	Timeout uint64 `toml:"timeout"`
	V3      bool   `toml:"v3"` // only used by v3
}

type ArbConfig struct {
	NativeStablePair string   `toml:"native_stable_pair"`
	StablePools      []string `toml:"stable_pools"`
	BlackListTokens  []string `toml:"blacklist_tokens"`
	MinPoolETH       string   `toml:"min_pool_eth"`
	MinMidPoolETH    string   `toml:"min_mid_pool_eth"`
	Sender           string   `toml:"sender"`
	ArbContract      string   `toml:"arb_contract"`
	CalcAllTokens    bool     `toml:"calc_all_tokens"`
}

// ChainConfig configures of the chain being indexed.
type ChainConfig struct {
	ChainID        string                 `toml:"chain_id"`
	Network        string                 `toml:"network"`
	URL            string                 `toml:"url"`
	WSS            string                 `toml:"wss"`
	Subscribe      bool                   `toml:"subscribe"`
	StartingHeight uint64                 `toml:"starting_height"`
	Factory        []starcomm.SwapFactory `toml:"factory"`
	PairsQuery     string                 `toml:"pairs_query"`
	PoolQuery      map[string]uint        `toml:"pool_query"`
	Providers      []ProviderConfig       `toml:"providers"`
	// These configuration options will be removed once
	// native reorg handling is implemented
	ConfirmationDepth uint   `toml:"confirmation_depth"`
	PollingInterval   uint   `toml:"polling_interval"`
	BlockInterval     uint64 `toml:"block_interval"`
	PollingSteps      uint   `toml:"polling_steps"`
}

type Log struct {
	Level              string `toml:"level"`
	FileLoggingEnabled bool   `toml:"file_logging_enabled"`
	// Directory to log to to when filelogging is enabled
	Directory string `toml:"directory"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `toml:"filename"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `toml:"max_size"`
	// MaxAge the max age in days to keep a logfile
	MaxAge     int `toml:"max_age"`
	MaxBackups int `toml:"max_backups"`
}

// Config represents the `indexer.toml` file used to configure the indexer.
type Config struct {
	Env string `toml:"env"`
	// Mode          string       `toml:"mode"`
	Log           Log          `toml:"log"`
	Chain         ChainConfig  `toml:"chain"`
	Arb           ArbConfig    `toml:"arb"`
	Kafka         KafkaConfig  `toml:"kafka"`
	DB            DBConfig     `toml:"db"`
	Redis         RedisConfig  `toml:"redis"`
	HTTPServer    ServerConfig `toml:"http"`
	MetricsServer ServerConfig `toml:"metrics"`
	SentryOpts    SentryOpts   `toml:"sentry"`
}

func (c *Config) IsProd() bool {
	return strings.ToLower(c.Env) == "prod"
}

func (c *Config) IsLocal() bool {
	return strings.ToLower(c.Env) == "local"
}

// Contracts contracts address.
// type Contracts struct {
// 	UniswapQueryAddress common.Address `toml:"uniswap_query_address"`
// 	PancakeQueryAddress common.Address `toml:"pancake_query_address"`
// 	AeroQueryAddress    common.Address `toml:"areo_query_address"`
// }

// SentryOpts sentry config options.
type SentryOpts struct {
	DSN   string `toml:"dsn"`
	Env   string `toml:"env"`
	Debug bool   `toml:"debug"`
}

// LoadConfig loads the `indexer.toml` config file from a given path.
func LoadConfig(path string) (Config, error) {
	log.Printf("loading config, path: %s", path)

	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	data = []byte(os.ExpandEnv(string(data)))
	// log.Debug().Msgf("parsed config file", "data", string(data))

	md, err := toml.Decode(string(data), &cfg)
	if err != nil {
		log.Printf("failed to decode config file: %v", err)
		return cfg, err
	}

	if len(md.Undecoded()) > 0 {
		log.Printf("unknown fields in config file, fields: %v", md.Undecoded())
		err = fmt.Errorf("unknown fields in config file: %v", md.Undecoded())

		return cfg, err
	}

	// Check to make sure some required properties are set
	var errs error
	if cfg.Chain.PollingInterval == 0 {
		errs = errors.Join(err, errors.New("`polling_interval` unset"))
	}

	if cfg.SentryOpts.Env == "" {
		cfg.SentryOpts.Env = cfg.Env
	}

	// log.Printf("loaded chain config: %v", cfg.Chain)

	return cfg, errs
}

var ErrConfigEmptyInAwsSecrets = errors.New("config is empty in AWS Secrets Manager")

// LoadConfigFromAWSSecrets load config from AWS secrets.
func LoadConfigFromAWSSecrets(secretName string) (Config, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	var cfg Config

	c, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return cfg, fmt.Errorf("conf.LoadConfigFromAWSSecrets LoadDefaultConfig: %w", err)
	}

	clt := sts.NewFromConfig(c)

	identity, err := clt.GetCallerIdentity(
		context.Background(),
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return cfg, fmt.Errorf("conf.LoadConfigFromAWSSecrets GetCallerIdentity: %w", err)
	}

	var (
		arn string
		uid string
	)

	if identity.Arn != nil {
		arn = *identity.Arn
	}

	if identity.UserId != nil {
		uid = *identity.UserId
	}

	fmt.Printf("AWS identity: arn=%s, uid=%s\n", arn, uid)

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(c)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.Background(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return cfg, fmt.Errorf("conf.LoadConfigFromAWSSecrets GetSecretValue: %w", err)
	}

	if result.SecretString == nil || *result.SecretString == "" {
		return cfg, fmt.Errorf("conf.LoadConfigFromAWSSecrets: %w", ErrConfigEmptyInAwsSecrets)
	}

	md, err := toml.Decode(*result.SecretString, &cfg)
	if err != nil {
		log.Printf("failed to decode config file: %v", err)
		return cfg, err
	}

	if len(md.Undecoded()) > 0 {
		log.Printf("unknown fields in config file, fields: %v %v", md.Undecoded(), err)
		err = fmt.Errorf("unknown fields in config file: %v", md.Undecoded())

		return cfg, err
	}

	if cfg.SentryOpts.Env == "" {
		cfg.SentryOpts.Env = cfg.Env
	}
	// err = yaml.Unmarshal([]byte(*result.SecretString), cfg)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Unmarshal config: %v\n", err)
	// }

	// cfg.setDefaults()

	return cfg, nil
}

const CONFIG_AWS_ENV_EVENTS = "CONFIG_AWS_ENV_EVENTS"

// LoadConfigAWS load config first from AWS, then config file.
func LoadConfigAWS(path string, envName string) (c Config, err error) {
	awsSecretsKey := os.Getenv(envName)

	if awsSecretsKey != "" {
		c, err = LoadConfigFromAWSSecrets(awsSecretsKey)
	} else {
		log.Printf("fallback to local config file, path: %v", path)
		// Path to config file can be passed in.
		c, err = LoadConfig(path)
	}

	SetDefaultLogConfig(&c.Log)

	if err == nil {
		_conf = c
	}

	return
}

func IsAWS(envName string) bool {
	return os.Getenv(envName) == ""
}

func IsLocal() bool {
	return _conf.Env == "local"
}

func SetDefaultLogConfig(logConfig *Log) {
	if !logConfig.FileLoggingEnabled {
		return
	}

	if logConfig.Directory == "" {
		logConfig.Directory = "./logs/"
	}

	if logConfig.Filename == "" {
		logConfig.Filename = "liquidity.log"
	}

	if logConfig.MaxAge == 0 {
		logConfig.MaxAge = 7 // days
	}

	if logConfig.MaxAge == 0 {
		logConfig.MaxSize = 500 // megabytes
	}

	if logConfig.MaxBackups == 0 {
		logConfig.MaxBackups = 10
	}
}
