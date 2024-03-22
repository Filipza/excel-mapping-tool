package settings

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var settings *viperWrapper

func GetSettings() Settings {
	return settings
}

type Settings interface {
	// Get can retrieve any value given the key to use.
	// Get is case-insensitive for a key.
	// Get has the behavior of returning the value associated with the first
	// place from where it is set. Viper will check in the following order:
	// override, flag, env, config file, key/value store, default
	//
	// Get returns an interface. For a specific value use one of the Get____ methods.
	Get(string) interface{}
	// GetString returns the value associated with the key as a string.
	GetString(string) string
	// GetBool returns the value associated with the key as a boolean.
	GetBool(string) bool
	// GetInt returns the value associated with the key as an integer.
	GetInt(string) int
	// GetInt32 returns the value associated with the key as an integer.
	GetInt32(string) int32
	// GetInt64 returns the value associated with the key as an integer.
	GetInt64(string) int64
	// GetUint returns the value associated with the key as an unsigned integer.
	GetUint(string) uint
	// GetUint32 returns the value associated with the key as an unsigned integer.
	GetUint32(string) uint32
	// GetUint64 returns the value associated with the key as an unsigned integer.
	GetUint64(string) uint64
	// GetFloat64 returns the value associated with the key as a float64.
	GetFloat64(string) float64
	// GetTime returns the value associated with the key as time.
	GetTime(string) time.Time
	// GetDuration returns the value associated with the key as a duration.
	GetDuration(string) time.Duration
	// GetIntSlice returns the value associated with the key as a slice of int values.
	GetIntSlice(string) []int
	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(string) []string
	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(string) map[string]interface{}
	// GetStringMapString returns the value associated with the key as a map of strings.
	GetStringMapString(string) map[string]string
	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	GetStringMapStringSlice(string) map[string][]string

	IsSet(string) bool

	GetDefault(string, interface{}) interface{}
	GetDefaultString(string, string) string
	GetDefaultBool(string, bool) bool
	GetDefaultInt(string, int) int
	GetDefaultInt32(string, int32) int32
	GetDefaultInt64(string, int64) int64
	GetDefaultUint(string, uint) uint
	GetDefaultUint32(string, uint32) uint32
	GetDefaultUint64(string, uint64) uint64
	GetDefaultFloat64(string, float64) float64
	GetDefaultTime(string, time.Time) time.Time
	GetDefaultDuration(string, time.Duration) time.Duration
	GetDefaultIntSlice(string, ...int) []int
	GetDefaultStringSlice(string, ...string) []string
	GetDefaultStringMap(string, map[string]interface{}) map[string]interface{}
	GetDefaultStringMapString(string, map[string]string) map[string]string
	GetDefaultStringMapStringSlice(string, map[string][]string) map[string][]string
}

type viperWrapper struct {
	*viper.Viper
}

func (w *viperWrapper) GetDefault(key string, v interface{}) interface{} {
	if !w.IsSet(key) {
		return v
	}
	return w.Get(key)
}

func (w *viperWrapper) GetDefaultString(key string, v string) string {
	if !w.IsSet(key) {
		return v
	}
	return w.GetString(key)
}

func (w *viperWrapper) GetDefaultBool(key string, v bool) bool {
	if !w.IsSet(key) {
		return v
	}
	return w.GetBool(key)
}

func (w *viperWrapper) GetDefaultInt(key string, v int) int {
	if !w.IsSet(key) {
		return v
	}
	return w.GetInt(key)
}

func (w *viperWrapper) GetDefaultInt32(key string, v int32) int32 {
	if !w.IsSet(key) {
		return v
	}
	return w.GetInt32(key)
}

func (w *viperWrapper) GetDefaultInt64(key string, v int64) int64 {
	if !w.IsSet(key) {
		return v
	}
	return w.GetInt64(key)
}

func (w *viperWrapper) GetDefaultUint(key string, v uint) uint {
	if !w.IsSet(key) {
		return v
	}
	return w.GetUint(key)
}

func (w *viperWrapper) GetDefaultUint32(key string, v uint32) uint32 {
	if !w.IsSet(key) {
		return v
	}
	return w.GetUint32(key)
}

func (w *viperWrapper) GetDefaultUint64(key string, v uint64) uint64 {
	if !w.IsSet(key) {
		return v
	}
	return w.GetUint64(key)
}

func (w *viperWrapper) GetDefaultFloat64(key string, v float64) float64 {
	if !w.IsSet(key) {
		return v
	}
	return w.GetFloat64(key)
}

func (w *viperWrapper) GetDefaultTime(key string, v time.Time) time.Time {
	if !w.IsSet(key) {
		return v
	}
	return w.GetTime(key)
}

func (w *viperWrapper) GetDefaultDuration(key string, v time.Duration) time.Duration {
	if !w.IsSet(key) {
		return v
	}
	return w.GetDuration(key)
}

func (w *viperWrapper) GetDefaultIntSlice(key string, v ...int) []int {
	if !w.IsSet(key) {
		return v
	}
	return w.GetIntSlice(key)
}

func (w *viperWrapper) GetDefaultStringSlice(key string, v ...string) []string {
	if !w.IsSet(key) {
		return v
	}
	return w.GetStringSlice(key)
}

func (w *viperWrapper) GetDefaultStringMap(key string, v map[string]interface{}) map[string]interface{} {
	if !w.IsSet(key) {
		return v
	}
	return w.GetStringMap(key)
}

func (w *viperWrapper) GetDefaultStringMapString(key string, v map[string]string) map[string]string {
	if !w.IsSet(key) {
		return v
	}
	return w.GetStringMapString(key)
}

func (w *viperWrapper) GetDefaultStringMapStringSlice(key string, v map[string][]string) map[string][]string {
	if !w.IsSet(key) {
		return v
	}
	return w.GetStringMapStringSlice(key)
}

func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("conf")
	viper.AddConfigPath("config")
	viper.AddConfigPath("configs")

	viper.SetConfigName("config")

	viper.ReadInConfig()

	venv := viper.New()
	for _, env := range os.Environ() {
		if strings.HasPrefix(strings.ToLower(env), "msd_") {
			pair := strings.SplitN(env, "=", 2)
			venv.Set(strings.ReplaceAll(strings.TrimPrefix(strings.ToLower(pair[0]), "msd_"), "_", "."), pair[1])
		}
	}

	viper.MergeConfigMap(venv.AllSettings())

	settings = &viperWrapper{
		viper.GetViper(),
	}
}
