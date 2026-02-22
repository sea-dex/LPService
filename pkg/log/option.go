package log

const (
	defaultFileMaxSizeMB  = 10
	defaultFileMaxAgeDays = 30
)

type Option func(*Options)

type Options struct {
	Console OptionsConsole `yaml:"console"`
	File    OptionsFile    `yaml:"file"`
}

type OptionsConsole struct {
	Disable    bool   `yaml:"disable"`
	Encoder    string `yaml:"encoder"`
	Level      string `yaml:"level"`
	TimeFormat string `yaml:"timeFormat"`
}

type OptionsFile struct {
	Encoder    string `yaml:"encoder"`
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	TimeFormat string `yaml:"timeFormat"`
	MaxSizeMB  int    `yaml:"maxSizeMBytes"`
	MaxAgeDays int    `yaml:"maxAgeDays"`
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		Console: OptionsConsole{
			Disable:    false,
			Encoder:    EncoderConsole,
			Level:      LevelDebug,
			TimeFormat: "RFC3339Nano",
		},
		File: OptionsFile{
			Encoder:    EncoderJSON,
			Level:      LevelInfo,
			TimeFormat: "RFC3339Nano",
			Path:       "",
			MaxSizeMB:  defaultFileMaxSizeMB,
			MaxAgeDays: defaultFileMaxAgeDays,
		},
	}

	for _, o := range opts {
		o(&opt)
	}

	return &opt
}

func ConsoleDisable(x bool) Option {
	return func(opts *Options) {
		opts.Console.Disable = x
	}
}

func ConsoleEncoder(enc string) Option {
	return func(opts *Options) {
		opts.Console.Encoder = enc
	}
}

func ConsoleLevel(l string) Option {
	return func(opts *Options) {
		opts.Console.Level = l
	}
}

func ConsoleTimeFormat(f string) Option {
	return func(opts *Options) {
		opts.Console.TimeFormat = f
	}
}

func FileEncoder(enc string) Option {
	return func(opts *Options) {
		opts.File.Encoder = enc
	}
}

func FileLevel(l string) Option {
	return func(opts *Options) {
		opts.File.Level = l
	}
}

func FileTimeFormat(f string) Option {
	return func(opts *Options) {
		opts.File.TimeFormat = f
	}
}

func FilePath(path string) Option {
	return func(opts *Options) {
		opts.File.Path = path
	}
}

func FileMaxSizeMB(s int) Option {
	return func(opts *Options) {
		opts.File.MaxSizeMB = s
	}
}

func FileMaxAgeDays(age int) Option {
	return func(opts *Options) {
		opts.File.MaxAgeDays = age
	}
}

func MergeOptions(logOpts Options) []Option {
	opts := make([]Option, 0)

	opts = mergeConsoleDisable(opts, logOpts.Console.Disable)
	opts = mergeConsoleEncoder(opts, logOpts.Console.Encoder)
	opts = mergeConsoleLevel(opts, logOpts.Console.Level)
	opts = mergeConsoleTimeFormat(opts, logOpts.Console.TimeFormat)

	opts = mergeFileEncoder(opts, logOpts.File.Encoder)
	opts = mergeFileLevel(opts, logOpts.File.Level)
	opts = mergeFileTimeFormat(opts, logOpts.File.TimeFormat)
	opts = mergeFilePath(opts, logOpts.File.Path)
	opts = mergeFileMaxSizeMB(opts, logOpts.File.MaxSizeMB)
	opts = mergeFileMaxAgeDays(opts, logOpts.File.MaxAgeDays)

	return opts
}

func mergeConsoleDisable(opts []Option, v bool) []Option {
	if v {
		return append(opts, ConsoleDisable(v))
	}

	return opts
}

func mergeConsoleEncoder(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, ConsoleEncoder(v))
	}

	return opts
}

func mergeConsoleLevel(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, ConsoleLevel(v))
	}

	return opts
}

func mergeConsoleTimeFormat(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, ConsoleTimeFormat(v))
	}

	return opts
}

func mergeFileEncoder(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, FileEncoder(v))
	}

	return opts
}

func mergeFileLevel(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, FileLevel(v))
	}

	return opts
}

func mergeFileTimeFormat(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, FileTimeFormat(v))
	}

	return opts
}

func mergeFilePath(opts []Option, v string) []Option {
	if v != "" {
		return append(opts, FilePath(v))
	}

	return opts
}

func mergeFileMaxSizeMB(opts []Option, v int) []Option {
	if v > 0 {
		return append(opts, FileMaxSizeMB(v))
	}

	return opts
}

func mergeFileMaxAgeDays(opts []Option, v int) []Option {
	if v > 0 {
		return append(opts, FileMaxAgeDays(v))
	}

	return opts
}
