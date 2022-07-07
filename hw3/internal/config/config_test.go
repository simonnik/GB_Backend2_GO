//go:build unit
// +build unit

package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	logLevel      = "error"
	logPath       = "/dev/console"
	port          = 8083
	serverTimeout = 60
	jwtSecret     = "secret"
	jwtTTL        = 60
	jwtIssuer     = "Issuer"
	jwtSubject    = "Subject"
	host          = "localhost"
	dbHost        = "postgresql"
	dbUser        = "price"
	dbPass        = "123123"
	dbPort        = 5432
	dbName        = "price"
	sslMode       = "disable"
	migrationsDir = "migrations"
	isShowEnvInfo = true
)

func requiredMsgByEnv(env string) string {
	return "required key " + env + " missing value"
}

func TestBuildConfigMissingPort(t *testing.T) {
	os.Clearenv()
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_PORT"))
	assert.Nil(t, c)

	os.Setenv("APP_PORT", strconv.Itoa(port))
}

func TestBuildConfigMissingHost(t *testing.T) {
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_HOST"))
	assert.Nil(t, c)

	os.Setenv("APP_HOST", host)
}

func TestBuildConfigMissingDBHost(t *testing.T) {
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_DB_HOST"))
	assert.Nil(t, c)

	os.Setenv("APP_DB_HOST", dbHost)
}

func TestBuildConfigMissingDBUser(t *testing.T) {
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_DB_USER"))
	assert.Nil(t, c)

	os.Setenv("APP_DB_USER", dbUser)
}

func TestBuildConfigMissingDBPassword(t *testing.T) {
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_DB_PASSWORD"))
	assert.Nil(t, c)

	os.Setenv("APP_DB_PASSWORD", dbPass)
}

func TestBuildConfigMissingDBName(t *testing.T) {
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_DB_NAME"))
	assert.Nil(t, c)

	os.Setenv("APP_DB_NAME", dbName)
}

func TestBuildConfigMissingJWTSecret(t *testing.T) {
	c, err := BuildConfig()
	require.EqualError(t, err, requiredMsgByEnv("APP_JWT_SECRET"))
	assert.Nil(t, c)

	os.Setenv("APP_JWT_SECRET", jwtSecret)
}

func TestBuildConfigInvalidJWTIssuer(t *testing.T) {
	c, err := BuildConfig()
	require.NoError(t, err)
	assert.IsType(t, new(Config), c)
	assert.NotEqual(t, jwtIssuer, c.JWT.Issuer)
}

func TestBuildConfigInvalidJWTSubject(t *testing.T) {
	c, err := BuildConfig()
	require.NoError(t, err)
	assert.IsType(t, new(Config), c)
	assert.NotEqual(t, jwtSubject, c.JWT.Subject)
}

func TestBuildConfigSuccessEnvProcess(t *testing.T) {
	c, err := BuildConfig()
	c.JWT.Issuer = jwtIssuer
	c.JWT.Subject = jwtSubject
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, host, c.Host)
	assert.Equal(t, port, c.Port)
	assert.Equal(t, serverTimeout, c.ServerTimeout)
	assert.Equal(t, jwtSecret, c.JWT.Secret)
	assert.Equal(t, jwtTTL, c.JWT.TTL)
	assert.Equal(t, jwtIssuer, c.JWT.Issuer)
	assert.Equal(t, jwtSubject, c.JWT.Subject)
	assert.Equal(t, port, c.Port)
	assert.Equal(t, isShowEnvInfo, c.IsShowEnvInfo)
	assert.Equal(t, dbHost, c.DB.Host)
	assert.Equal(t, dbPort, c.DB.Port)
	assert.Equal(t, dbName, c.DB.Name)
	assert.Equal(t, dbPass, c.DB.Password)
	assert.Equal(t, dbUser, c.DB.User)
	assert.Equal(t, sslMode, c.DB.SSLMode)
	assert.Equal(t, migrationsDir, c.DB.MigrationsDir)
	assert.Equal(t, logLevel, c.Log.Level)
	assert.Equal(t, logPath, c.Log.Path)
}
