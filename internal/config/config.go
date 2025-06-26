package config

// Config представляет конфигурацию приложения.
type Config struct {
	Global  GlobalConfig  `yaml:"global"`
	Log     LogConfig     `yaml:"log"`
	Servers ServersConfig `yaml:"servers"`
	Clients ClientsConfig `yaml:"clients"`
	Storage StorageConfig `yaml:"storage"`
}

// GlobalConfig представляет глобальные настройки.
type GlobalConfig struct {
	// добавляем валидацию: обязательное поле, значения из {"local", "dev", "stage", "prod"}.
	Env string `yaml:"env" validate:"required,oneof=local dev stage prod"`
}

// LogConfig представляет настройки логирования.
type LogConfig struct {
	// добавляем валидацию: обязательное поле, значения из {"debug", "info", "warn", "error"}.
	Level string `yaml:"level" validate:"required,oneof=debug info warn error"`
}

// ServersConfig представляет настройки серверов.
type ServersConfig struct {
	Debug  DebugServerConfig  `yaml:"debug"`
	Client ClientServerConfig `yaml:"client"`
}

// DebugServerConfig представляет настройки отладочного сервера.
type DebugServerConfig struct {
	// добавляем валидацию: обязательное поле, значение должно быть в формате "host:port".
	Addr string `yaml:"addr" validate:"required,hostname_port"`
}

// ClientServerConfig представляет настройки клиентского API сервера.
type ClientServerConfig struct {
	Addr         string   `yaml:"addr" validate:"required,hostname_port"`
	AllowOrigins []string `yaml:"allow_origins"`
}

// ClientsConfig представляет настройки для внешних клиентов.
type ClientsConfig struct {
	Keycloak      KeycloakConfig `yaml:"keycloak"`       // back-end
	KeycloakAdmin KeycloakConfig `yaml:"keycloak_admin"` // woman-app-admin
}

// KeycloakConfig представляет настройки для Keycloak.
type KeycloakConfig struct {
	BasePath     string `yaml:"base_path" validate:"required,url"`
	Realm        string `yaml:"realm" validate:"required"`
	ClientID     string `yaml:"client_id" validate:"required"`
	ClientSecret string `yaml:"client_secret" validate:"required"`
	DebugMode    bool   `yaml:"debug_mode"`
}

// StorageConfig представляет настройки для хранения данных.
type StorageConfig struct {
	DBName        string `yaml:"db_name" validate:"required"`
	DBUser        string `yaml:"db_user" validate:"required"`
	DBPassword    string `yaml:"db_password" validate:"required"`
	DBHost        string `yaml:"db_host" validate:"required"`
	DBPort        string `yaml:"db_port" validate:"required"`
	DBSSLMode     string `yaml:"db_ssl_mode"`
	DBSSLRootCert string `yaml:"db_ssl_root_cert"`
	DBSSLKey      string `yaml:"db_ssl_key"`
}
