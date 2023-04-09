package app

const defaultConfigFile = "arc.conf"

type Config map[string]any

func (conf Config) Read() error {
	if err := conf.readFromFile(defaultConfigFile); err != nil {
		return err
	}
	if err := conf.readFromEnv(); err != nil {
		return err
	}
	if err := conf.readFromCli(); err != nil {
		return err
	}
	return conf.Validate()
}

func (conf Config) readFromCli() error {
	// TODO parse arg
	return nil
}

func (conf Config) readFromFile(filePath string) error {
	// TODO read config file
	return nil
}

func (conf Config) readFromEnv() error {
	// TODO read from ENV
	conf["port"] = 6378
	conf["buffersize"] = 128 * 1024 // 128kb
	return nil
}

func (conf Config) Validate() error {
	return nil
}

func (conf Config) Get(key string) string {
	return conf[key].(string)
}

func (conf Config) GetInt(key string) int {
	return conf[key].(int)
}

func (conf Config) GetFloat(key string) float64 {
	return conf[key].(float64)
}
