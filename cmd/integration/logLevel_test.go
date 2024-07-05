/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-mi-tooling/cmd/integration/base"
	"github.com/wso2/product-mi-tooling/cmd/integration/testutils"
)

const validLoggerName = "org-apache-coyote"
const invalidLoggerName = "abc-logger"
const logLevelCmd = "log-levels"
const newLoggerName = "synapse-api"
const newLoggerClass = "org.apache.synapse.rest.API"

var validAddLoggerCmd = []string{"add", "log-level", newLoggerName, newLoggerClass, "DEBUG", "-e", "testing", "-k"}
var validUpdateLoggerCmd = []string{"update", "log-level", newLoggerName, "INFO", "-e", "testing", "-k"}

func TestGetLoggerByName(t *testing.T) {
	testutils.ValidateLogger(t, logLevelCmd, config, validLoggerName)
}

func TestGetNonExistingLoggerByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, config, logLevelCmd, invalidLoggerName)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of logger [ "+invalidLoggerName+" ]  Logger name ('"+invalidLoggerName+"') not found")
}

func TestGetLoggersWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, logLevelCmd, validLoggerName)
}

func TestGetLoggersWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, logLevelCmd, config, validLoggerName)
}

func TestGetLoggersWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, logLevelCmd, config, validLoggerName)
}

func TestGetLoggersWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, true, logLevelCmd, validLoggerName, invalidLoggerName)
}

func TestAddNewLoggerWithInvalidLogLevel(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "add", "log-level", newLoggerName, newLoggerClass, "ABC", "-e", miClient.GetEnvName(), "-k")
	base.Log(response)
	expected := "[ERROR]: Adding new logger [ " + newLoggerName + " ]  Invalid log level ABC"
	assert.Contains(t, response, expected)
}

func TestAddNewLoggerWithoutEnvFlag(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "add", "log-level", newLoggerName, newLoggerClass, "DEBUG", "-k")
	base.Log(response)
	expected := `required flag(s) "environment" not set`
	assert.Contains(t, response, expected)
}

func TestAddNewLoggerWithInvalidArgs(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "add", "log-level", newLoggerName, newLoggerClass, "-e", miClient.GetEnvName(), "-k")
	base.Log(response)
	expected := "accepts 3 arg(s), received 2"
	assert.Contains(t, response, expected)
}

func TestAddNewLoggerWithoutSettingUpEnv(t *testing.T) {
	response, _ := base.Execute(t, validAddLoggerCmd...)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

func TestAddLoggerWithoutLogin(t *testing.T) {
	base.SetupMIEnv(t, miClient.GetEnvName(), config.MIClient.GetMiURL())
	response, _ := base.Execute(t, validAddLoggerCmd...)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

func TestAddNewLogger(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, validAddLoggerCmd...)
	base.Log(response)
	expected := "Successfully added logger"
	assert.Contains(t, response, expected)
}

func TestAddExistingLogger(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, validAddLoggerCmd...)
	base.Log(response)
	expected := "[ERROR]: Adding new logger [ " + newLoggerName + " ]  Specified logger name ('" + newLoggerName + "') already exists, try updating the level instead"
	assert.Contains(t, response, expected)
}

func TestUpdateLoggerWithInvalidLogLevel(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "update", "log-level", newLoggerName, "ABC", "-e", miClient.GetEnvName(), "-k")
	base.Log(response)
	expected := "[ERROR]: updating logger [ " + newLoggerName + " ]  Invalid log level ABC"
	assert.Contains(t, response, expected)
}

func TestUpdateLoggerWithoutEnvFlag(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "update", "log-level", newLoggerName, "INFO", "-k")
	base.Log(response)
	expected := `required flag(s) "environment" not set`
	assert.Contains(t, response, expected)
}

func TestUpdateLoggerWithInvalidArgs(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "update", "log-level", newLoggerName, "-e", miClient.GetEnvName(), "-k")
	base.Log(response)
	expected := "accepts 2 arg(s), received 1"
	assert.Contains(t, response, expected)
}

func TestUpdateLoggerWithoutSettingUpEnv(t *testing.T) {
	response, _ := base.Execute(t, validUpdateLoggerCmd...)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

func TestUpdateLoggerWithoutLogin(t *testing.T) {
	base.SetupMIEnv(t, miClient.GetEnvName(), config.MIClient.GetMiURL())
	response, _ := base.Execute(t, validUpdateLoggerCmd...)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

func TestUpdateLogger(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, validUpdateLoggerCmd...)
	base.Log(response)
	expected := "Successfully added logger for ('" + newLoggerName + "') with level INFO"
	assert.Contains(t, response, expected)
}

func TestUpdateLoggerNonExistingLogger(t *testing.T) {
	testutils.SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "update", "log-level", invalidLoggerName, "INFO", "-e", miClient.GetEnvName(), "-k")
	base.Log(response)
	expected := "[ERROR]: updating logger [ " + invalidLoggerName + " ]  Specified logger ('" + invalidLoggerName + "') not found"
	assert.Contains(t, response, expected)
}
