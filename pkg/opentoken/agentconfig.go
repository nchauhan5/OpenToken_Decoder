package opentoken

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// AgentConfiguration is
type AgentConfiguration struct {
	tokenName    string
	cookieDomain string
	cookiePath   string
	Password     string

	tokenLifetime      int
	renewUntilLifetime int
	notBeforeTolerance int
	Ciphersuite        int

	useCookie               bool
	sessionCookie           bool
	secureCookie            bool
	httpOnly                bool
	useVerboseErrorMessages bool
	obfuscatePassword       bool
}

// AgentConfig of type agentConfiguration struct with default values of the opentoken configs. These values will be updated from the current default values if they are present in the opentoken congig txts
var AgentConfig = AgentConfiguration{
	tokenName:    "opentoken",
	cookieDomain: "",
	cookiePath:   "/",
	Password:     "testPassword",

	tokenLifetime:      300,
	renewUntilLifetime: 43200,
	notBeforeTolerance: 120,
	Ciphersuite:        0,

	useCookie:               false,
	sessionCookie:           false,
	secureCookie:            false,
	useVerboseErrorMessages: false,
	httpOnly:                true,
	obfuscatePassword:       true,
}

func init() {
	dabsPath, _ := filepath.Abs("../../agent-config.txt")
	fmt.Println(dabsPath)
	AgentConfig = ReadConfigFromFile(dabsPath)
	fmt.Println(AgentConfig)
}

// ReadConfigFromFile is used to read config from file
func ReadConfigFromFile(filepath string) AgentConfiguration {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Returning the default Agent Config. Could not find a file at the path mentioned : " + filepath)
		return AgentConfig
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	keyValueConfigs := make(map[string]string)
	for index := 0; index < len(lines); index++ {
		splitSinglePair := strings.Split(lines[index], "=")
		keyValueConfigs[splitSinglePair[0]] = splitSinglePair[1]
	}

	if (len(keyValueConfigs["token-name"])) == 0 {
		fmt.Println("Returning the default Agent Config. Required field token name is missing from the agent config file")
		return AgentConfig
	}

	AgentConfig.tokenName = keyValueConfigs["token-name"]
	useCookie, _ := strconv.ParseBool(keyValueConfigs["use-cookie"])
	if useCookie {
		AgentConfig.useCookie = useCookie
		if cookieDomain, ok := keyValueConfigs["cookie-domain"]; ok {
			AgentConfig.cookieDomain = cookieDomain
		}
		if cookiePath, ok := keyValueConfigs["cookie-path"]; ok {
			AgentConfig.cookiePath = cookiePath
		}
		isSessionCookie, _ := strconv.ParseBool(keyValueConfigs["session-cookie"])
		if isSessionCookie {
			AgentConfig.sessionCookie = isSessionCookie
		}
		isSecureCookie, _ := strconv.ParseBool(keyValueConfigs["secure-cookie"])
		if isSecureCookie {
			AgentConfig.secureCookie = isSecureCookie
		}
		isHTTPOnly, _ := strconv.ParseBool(keyValueConfigs["http-only"])
		if !isHTTPOnly {
			AgentConfig.httpOnly = isHTTPOnly
		}
	}
	isVerboseErrorMessage, _ := strconv.ParseBool(keyValueConfigs["use-verbose-error-messages"])
	if isVerboseErrorMessage {
		AgentConfig.useVerboseErrorMessages = isVerboseErrorMessage
	}

	lifetime, _ := strconv.Atoi(keyValueConfigs["token-lifetime"])
	AgentConfig.tokenLifetime = lifetime

	renewLifetime, _ := strconv.Atoi(keyValueConfigs["token-renewuntil"])
	AgentConfig.renewUntilLifetime = renewLifetime

	notBeforeOffset, _ := strconv.Atoi(keyValueConfigs["token-notbefore-tolerance"])
	AgentConfig.notBeforeTolerance = notBeforeOffset

	obfuscatePassword, _ := strconv.ParseBool(keyValueConfigs["obfuscate-password"])
	if !obfuscatePassword {
		AgentConfig.obfuscatePassword = obfuscatePassword
	}

	password := removeReturnChar(keyValueConfigs["password"])
	if len(password) > 0 {
		ciphersuite, _ := strconv.Atoi(keyValueConfigs["cipher-suite"])
		AgentConfig.Ciphersuite = ciphersuite
		if obfuscatePassword {
			AgentConfig.Password = deObfuscate(password)
		} else {
			decodedPass, _ := base64.URLEncoding.DecodeString(password)
			AgentConfig.Password = string(decodedPass)
		}
	}

	return AgentConfig
}

func deObfuscate(password string) string {
	return password
}

func removeReturnChar(inputString string) string {
	re := regexp.MustCompile(`\r?\n`)
	input := re.ReplaceAllString(inputString, "")
	return input
}
