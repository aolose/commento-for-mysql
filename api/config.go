package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func configParse() error {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("cannot load binary path: %v", err)
		return err
	}
	logger.Infof("%s", binPath)
	defaults := map[string]string{
		"CONFIG_FILE":   "",
		"DATABASE_TYPE": "mysql",
		"DATABASE_URL":  "commento:123@/commento?charset=utf8mb4&parseTime=True&loc=Local",

		// is 100, so let's use half that and leave the other half for other services.
		// Ideally, you'd be setting this to a much higher number (for example, at the
		// time of writing, commento.io uses 600). See https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
		"MAX_IDLE_DB_CONNECTIONS": "50",

		"BIND_ADDRESS": "127.0.0.1",
		"PORT":         "8080",
		"ORIGIN":       "",

		"CDN_PREFIX": "",

		"FORBID_NEW_OWNERS": "true",

		//"STATIC": "Z:\\go\\src\\commento\\build\\devel",
		"STATIC": binPath,

		"GZIP_STATIC": "false",

		"SMTP_USERNAME":     "",
		"SMTP_PASSWORD":     "",
		"SMTP_HOST":         "",
		"SMTP_PORT":         "",
		"SMTP_FROM_ADDRESS": "",

		"AKISMET_KEY": "",

		"GOOGLE_KEY":    "",
		"GOOGLE_SECRET": "",

		"GITHUB_KEY":    "",
		"GITHUB_SECRET": "",

		"TWITTER_KEY":    "",
		"TWITTER_SECRET": "",

		"GITLAB_KEY":    "",
		"GITLAB_SECRET": "",
	}

	for key, value := range defaults {
		if os.Getenv("COMMENTO_"+key) == "" {
			os.Setenv(key, value)
		} else {
			os.Setenv(key, os.Getenv("COMMENTO_"+key))
		}
	}

	if os.Getenv("CONFIG_FILE") != "" {
		if err := configFileLoad(os.Getenv("CONFIG_FILE")); err != nil {
			return err
		}
	}

	// Mandatory config parameters
	for _, env := range []string{"DATABASE_URL", "DATABASE_TYPE", "PORT", "ORIGIN", "FORBID_NEW_OWNERS", "MAX_IDLE_DB_CONNECTIONS"} {
		if os.Getenv(env) == "" {
			logger.Errorf("missing COMMENTO_%s environment variable", env)
			return errorMissingConfig
		}
	}

	os.Setenv("ORIGIN", strings.TrimSuffix(os.Getenv("ORIGIN"), "/"))
	os.Setenv("ORIGIN", addHttpIfAbsent(os.Getenv("ORIGIN")))

	if os.Getenv("CDN_PREFIX") == "" {
		os.Setenv("CDN_PREFIX", os.Getenv("ORIGIN"))
	}

	os.Setenv("CDN_PREFIX", strings.TrimSuffix(os.Getenv("CDN_PREFIX"), "/"))
	os.Setenv("CDN_PREFIX", addHttpIfAbsent(os.Getenv("CDN_PREFIX")))

	if os.Getenv("FORBID_NEW_OWNERS") != "true" && os.Getenv("FORBID_NEW_OWNERS") != "false" {
		logger.Errorf("COMMENTO_FORBID_NEW_OWNERS neither 'true' nor 'false'")
		return errorInvalidConfigValue
	}

	static := os.Getenv("STATIC")
	for strings.HasSuffix(static, "/") {
		static = static[0 : len(static)-1]
	}

	file, err := os.Stat(static)
	if err != nil {
		logger.Errorf("cannot load %s: %v", static, err)
		return err
	}

	if !file.IsDir() {
		logger.Errorf("COMMENTO_STATIC=%s is not a directory", static)
		return errorNotADirectory
	}

	os.Setenv("STATIC", static)

	if num, err := strconv.Atoi(os.Getenv("MAX_IDLE_DB_CONNECTIONS")); err != nil {
		logger.Errorf("invalid COMMENTO_MAX_IDLE_DB_CONNECTIONS: %v", err)
		return errorInvalidConfigValue
	} else if num <= 0 {
		logger.Errorf("COMMENTO_MAX_IDLE_DB_CONNECTIONS should be a positive integer")
		return errorInvalidConfigValue
	}

	return nil
}
