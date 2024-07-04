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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-mi-tooling/cmd/integration/base"
	"github.com/wso2/product-mi-tooling/cmd/utils"
	artifactutils "github.com/wso2/product-mi-tooling/cmd/utils/artifactUtils"
)

// ValidateTransaction validate ctl output with transaction from the Management API
func ValidateTransaction(t *testing.T, transactionCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, transactionCmd, config)
	var transactionCountResource = utils.MiManagementTransactionResource + "/" + utils.MiManagementTransactionCountResource
	artifact := config.MIClient.GetArtifactListFromAPI(transactionCountResource, &artifactutils.TransactionCount{})
	validateTransactionEqual(t, output, (artifact.(*artifactutils.TransactionCount)))
}

func validateTransactionEqual(t *testing.T, transactionFromCtl string, transaction *artifactutils.TransactionCount) {
	assert.Contains(t, transactionFromCtl, fmt.Sprint(transaction.Year))
	assert.Contains(t, transactionFromCtl, fmt.Sprint(transaction.Month))
	assert.Contains(t, transactionFromCtl, fmt.Sprint(transaction.TransactionCount))
}

// ExecGetTransactionCountWithInvalidArgCount run get transaction-counts with invalid number of args
func ExecGetTransactionCountWithInvalidArgCount(t *testing.T, config *MiConfig, passed int, args ...string) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Password)
	getCmdArgs := []string{"get", "-e", config.MIClient.GetEnvName(), "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	expected := fmt.Sprintf("accepts exactly 0 or 2 arg(s), received %v", passed)
	assert.Contains(t, response, expected)
}
