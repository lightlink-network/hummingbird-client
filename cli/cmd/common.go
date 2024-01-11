package cmd

import (
	"crypto/ecdsa"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

func ConsoleLogger() *slog.Logger {
	w := os.Stderr
	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level:      parseLogLevel(viper.GetString("log-level")),
		TimeFormat: time.Kitchen,
		AddSource:  viper.GetBool("log-source"),
	}))

	return logger
}

func JSONLogger(w io.Writer) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     parseLogLevel(viper.GetString("log-level")),
		AddSource: viper.GetBool("log-source"),
	}))

	return logger
}

func GetLogger(logType string) *slog.Logger {
	switch logType {
	case "console":
		return ConsoleLogger()
	case "json":
		return JSONLogger(os.Stderr)
	default:
		panic("log type must be 'console' or 'json' got: " + logType)
	}
}

func parseLogLevel(lvl string) slog.Level {
	switch lvl {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		panic("log level not known: " + lvl)
	}
}

func getEthKey() *ecdsa.PrivateKey {
	key := os.Getenv("ETH_KEY")
	if key == "" {
		panic("env ETH_KEY not set")
	}

	ethKey, err := crypto.ToECDSA(hexutil.MustDecode(key))
	if err != nil {
		panic("Failed to decode ETH_KEY: " + err.Error())
	}

	return ethKey
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
