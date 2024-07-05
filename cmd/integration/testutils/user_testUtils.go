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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-mi-tooling/cmd/integration/base"
	"github.com/wso2/product-mi-tooling/cmd/utils"
	artifactutils "github.com/wso2/product-mi-tooling/cmd/utils/artifactUtils"
)

type newUserRequestBody struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
	IsAdmin  string `json:"isAdmin"`
}

// ValidateUserList validate ctl output with list of users from the Management API
func ValidateUserList(t *testing.T, userCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, userCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementUserResource, &artifactutils.UserList{})
	validateUserListEqual(t, output, (artifactList.(*artifactutils.UserList)))
}

func validateUserListEqual(t *testing.T, userListFromCtl string, userList *artifactutils.UserList) {
	unmatchedCount := userList.Count
	for _, user := range userList.Users {
		assert.Truef(t, strings.Contains(userListFromCtl, user.UserId), "userListFromCtl: "+userListFromCtl+
			" , does not contain user.UserId: "+user.UserId)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "user lists are not equal")
}

// ValidateUser validate ctl output with the user from the Management API
func ValidateUser(t *testing.T, userCmd string, config *MiConfig, userName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, userCmd, userName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementUserResource+"/"+userName, nil, &artifactutils.UserSummary{})
	validateUserEqual(t, output, (artifact.(*artifactutils.UserSummary)))
}

func validateUserEqual(t *testing.T, userFromCtl string, user *artifactutils.UserSummary) {
	assert.Contains(t, userFromCtl, user.UserId)
	assert.Contains(t, userFromCtl, fmt.Sprint(user.IsAdmin))
	for _, role := range user.Roles {
		assert.Contains(t, userFromCtl, role)
	}
}

// AddNewUserFromAPI : Adds a new user using the MI Management API
func AddNewUserFromAPI(t *testing.T, config *MiConfig, userName, password, isAdmin string, cleanUp bool) {
	body := createAddUserRequestPayload(userName, password, isAdmin)
	addNewUserURL := getResourceURLWithQueryParam(config.MIClient.GetMiURL(), utils.MiManagementUserResource, nil)

	request := base.CreatePost(addNewUserURL, body)
	base.SetDefaultRestAPIHeaders(config.MIClient.accessToken, request)
	base.LogRequest("mi.AddNewUserFromAPI()", request)
	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("mi.AddNewUserFromAPI()", response, 200)
	if cleanUp {
		t.Cleanup(func() {
			removeUserFromAPI(config, userName)
		})
	}
}

func createAddUserRequestPayload(userName, password, isAdmin string) *bytes.Buffer {
	body := newUserRequestBody{
		UserID:   userName,
		Password: password,
		IsAdmin:  isAdmin,
	}
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

func removeUserFromAPI(config *MiConfig, userName string) {
	deleteUserURL := getResourceURLWithQueryParam(config.MIClient.GetMiURL(), utils.MiManagementUserResource+"/"+userName, nil)
	request := base.CreateDelete(deleteUserURL)
	base.SetDefaultRestAPIHeaders(config.MIClient.accessToken, request)
	base.LogRequest("mi.removeUserFromAPI()", request)
	response := base.SendHTTPRequest(request)
	defer response.Body.Close()
	base.ValidateAndLogResponse("mi.removeUserFromAPI()", response, 200)
}
