package config

// Option config 构建选项（函数选项模式）
type Option func(*options)

type options struct {
	configPath  string
	configName  string
	configType  string
	envPrefix   string
	autoEnv     bool
	watchConfig bool
}

func defaultOptions() *options {
	return &options{
		configPath:  ".",
		configName:  "config",
		configType:  "yaml",
		envPrefix:   "ANTBLOG",
		autoEnv:     true,
		watchConfig: false,
	}
}

// WithConfigPath 设置配置文件目录
func WithConfigPath(path string) Option {
	return func(o *options) { o.configPath = path }
}

// WithConfigName 设置配置文件名（不含扩展名）
func WithConfigName(name string) Option {
	return func(o *options) { o.configName = name }
}

// WithConfigType 设置配置文件类型（yaml/json/toml）
func WithConfigType(t string) Option {
	return func(o *options) { o.configType = t }
}

// WithEnvPrefix 设置环境变量前缀
func WithEnvPrefix(prefix string) Option {
	return func(o *options) { o.envPrefix = prefix }
}

// WithAutoEnv 是否自动绑定环境变量
func WithAutoEnv(auto bool) Option {
	return func(o *options) { o.autoEnv = auto }
}

// WithWatchConfig 是否启用配置热重载
func WithWatchConfig(watch bool) Option {
	return func(o *options) { o.watchConfig = watch }
}
