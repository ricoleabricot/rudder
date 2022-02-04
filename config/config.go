package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DevEnv     = "dev"
	PreprodEnv = "preprod"
	ProdEnv    = "prod"

	defaultEnvMode   = ProdEnv
	defaultDebugMode = false

	defaultMetricsAddr      = "0.0.0.0:8080"
	defaultProbeAddr        = "0.0.0.0:8081"
	enableElectionByDefault = false
	defaultManagerPort      = 9443
	enableWebhooksByDefault = true
)

type Config struct {
	// The current environment (dev/preprod/prod).
	Env string `json:"env" validate:"required"`

	// Debug mode enabled.
	Debug bool `json:"debug"`

	// The path to the kubeconfig.
	Kubeconfig string `json:"kubeconfig"`

	// The address the metric endpoint is bound to.
	MetricsAddr string `json:"metricsAddr" validate:"required"`

	// The address the probe endpoint binds to.
	ProbeAddr string `json:"probeAddr" validate:"required"`

	// Enable leader election for controller manager.
	// Enabling this will ensure there is only one active controller manager.
	LeaderElectionEnable bool `json:"leaderElection"`

	// The port on which to expose the controller manager.
	Port int `json:"port" validate:"required"`

	// Enable defaulter and validation on resources via Webhooks.
	EnableWebhooks bool `json:"enableWebhooks"`
}

// Load method will return configuration needed to properly run this operator.
// It allows multiple configuration mediums with the following order of importance:
// 1. The lowest: the configuration from file, all keywords must be defined in a config.yaml file.
// 2. In between: the env variables, will override any keys from config file.
// 3. The highest: the flags, used as arguments, will override any existing configuration key from the config.
func Load() (*Config, error) {
	cfg := Config{}

	// Setup defaults vars
	viper.SetDefault("env", defaultEnvMode)
	viper.SetDefault("debug", defaultDebugMode)
	viper.SetDefault("metricsAddr", defaultMetricsAddr)
	viper.SetDefault("probeAddr", defaultProbeAddr)
	viper.SetDefault("leaderElection", enableElectionByDefault)
	viper.SetDefault("enableWebhooks", enableWebhooksByDefault)

	// Enable configuration from env
	SetupConfigFromEnv()

	// Enable configuration from flags
	err := SetupConfigFromFlags()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config from flags: %w", err)
	}

	// Enable configuration from file
	err = SetupConfigFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config from file: %w", err)
	}

	// Bind viper configuration
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get viper configuration: %w", err)
	}

	// Init config validator
	validate := validator.New()
	err = validate.RegisterValidation("required_unless_dev", RequiredUnlessDevMode)
	if err != nil {
		return nil, fmt.Errorf("failed to register validation method: %w", err)
	}

	// Validate configuration
	err = validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return &cfg, nil
}

func SetupConfigFromEnv() {
	// Setup automatic env binding
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func SetupConfigFromFlags() error {
	// Define all flags
	pflag.String("env", defaultEnvMode, "The environment mode (dev/preprod/prod).")
	pflag.Bool("debug", defaultDebugMode, "If debug mode activated, more logs to print.")

	pflag.String("kubeconfig", "", "The path to the kubeconfig.")

	pflag.String("metricsAddr", defaultMetricsAddr, "The address the metric endpoint is bound to.")
	pflag.String("probeAddr", defaultProbeAddr, "The address the probe endpoint binds to.")
	pflag.Bool("leaderElection", enableElectionByDefault, "Enable leader election for controller manager.")
	pflag.Int("port", defaultManagerPort, "The port on which to expose the controller manager.")
	pflag.Bool("enableWebhooks", enableWebhooksByDefault, "Enable the webhook validator/defaulter.")

	pflag.Parse()

	return viper.BindPFlags(pflag.CommandLine)
}

func SetupConfigFromFile() error {
	// Define config file path
	viper.AddConfigPath("/etc/rudder")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read config file
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to bind file data to configuration: %w", err)
		}
	}

	return nil
}

// IsDevMod returns whether debug mode is active
func (conf *Config) IsDevMod() bool {
	return conf.Env == DevEnv
}

// RequiredUnlessDevMode is a custom validator for checking if env is in dev mode
// Returning true means the validation is OK
func RequiredUnlessDevMode(fl validator.FieldLevel) bool {
	if fl.Top().FieldByName("Env").String() == DevEnv {
		return true
	}

	return !fl.Field().IsZero()
}
