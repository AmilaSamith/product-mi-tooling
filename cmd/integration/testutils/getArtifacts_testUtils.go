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

package testutils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-mi-tooling/cmd/integration/base"
)

// ListArtifacts return ctl out from the command get artifactType
func ListArtifacts(t *testing.T, artifactType string, config *MiConfig) (string, error) {
	t.Helper()
	SetupAndLoginToMI(t, config)
	output, err := base.Execute(t, "get", artifactType, "-e", config.MIClient.GetEnvName(), "-k")
	return output, err
}

// GetArtifact return ctl out from the command get artifactType artifactName
func GetArtifact(t *testing.T, config *MiConfig, args ...string) (string, error) {
	t.Helper()
	SetupAndLoginToMI(t, config)
	getCmdArgs := []string{"get", "-e", config.MIClient.GetEnvName(), "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	output, err := base.Execute(t, getCmdArgs...)
	return output, err
}

// GetArtifactListFromAPI : Get Artifact Lists from MI Management API
func (instance *MiRESTClient) GetArtifactListFromAPI(resource string, artifactListType interface{}) interface{} {
	apisURL := getResourceURL(instance.GetMiURL(), resource)

	request := base.CreateGet(apisURL)
	base.SetDefaultRestAPIHeaders(instance.accessToken, request)
	base.LogRequest("mi.GetArtifactList()", request)
	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("mi.GetArtifactList()", response, 200)

	artifactListResponse := artifactListType
	json.NewDecoder(response.Body).Decode(&artifactListResponse)
	return artifactListResponse
}

// GetArtifactFromAPI : Get Artifacts from MI Management API
func (instance *MiRESTClient) GetArtifactFromAPI(resource string, params map[string]string, artifactType interface{}) interface{} {
	apisURL := getResourceURLWithQueryParam(instance.GetMiURL(), resource, params)

	request := base.CreateGet(apisURL)
	base.SetDefaultRestAPIHeaders(instance.accessToken, request)
	base.LogRequest("mi.GetArtifact()", request)
	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("mi.GetArtifact()", response, 200)

	artifactListResponse := artifactType
	json.NewDecoder(response.Body).Decode(&artifactListResponse)
	return artifactListResponse
}

// ExecGetCommandWithoutSettingEnv run get artifactType without setting up an environment
func ExecGetCommandWithoutSettingEnv(t *testing.T, args ...string) {
	t.Helper()
	getCmdArgs := []string{"get", "-e", "testing", "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

// ExecGetCommandWithoutLogin run get artifactType without login to MI
func ExecGetCommandWithoutLogin(t *testing.T, artifactType string, config *MiConfig, args ...string) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	getCmdArgs := []string{"get", artifactType, "-e", config.MIClient.GetEnvName(), "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

// ExecGetCommandWithoutEnvFlag run get artifactType without -e flag
func ExecGetCommandWithoutEnvFlag(t *testing.T, artifactType string, config *MiConfig, args ...string) {
	t.Helper()
	SetupAndLoginToMI(t, config)
	getCmdArgs := []string{"get", artifactType, "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	assert.Contains(t, response, `required flag(s) "environment" not set`)
}

// ExecGetCommandWithInvalidArgCount run get artifactType with invalid number of args
func ExecGetCommandWithInvalidArgCount(t *testing.T, config *MiConfig, required, passed int, fixedArgCout bool, args ...string) {
	t.Helper()
	SetupAndLoginToMI(t, config)
	getCmdArgs := []string{"get", "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	expected := fmt.Sprintf("accepts at most %v arg(s), received %v", required, passed)
	if fixedArgCout {
		expected = fmt.Sprintf("accepts %v arg(s), received %v", required, passed)
	}
	assert.Contains(t, response, expected)
}
