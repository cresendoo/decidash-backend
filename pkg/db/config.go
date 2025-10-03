package db

type Config struct {
	DBName          string `yaml:"database"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
	Timezone        string `yaml:"timezone"`

	Addresses []Address `yaml:"addresses"`
}

type Address struct {
	User   string `yaml:"user"`
	Passwd string `yaml:"password"`
	Addr   string `yaml:"address"`
}
