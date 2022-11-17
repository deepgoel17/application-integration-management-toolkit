// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"regexp"

	"github.com/apigee/apigeecli/clilog"
)

// condition for Bindings
type condition struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Expression  string `json:"expression,omitempty"`
}

// binding for IAM Roles
type roleBinding struct {
	Role      string     `json:"role,omitempty"`
	Members   []string   `json:"members,omitempty"`
	Condition *condition `json:"condition,omitempty"`
}

// IamPolicy holds the response
type iamPolicy struct {
	Version  int           `json:"version,omitempty"`
	Etag     string        `json:"etag,omitempty"`
	Bindings []roleBinding `json:"bindings,omitempty"`
}

// SetIamPolicy holds the request to set IAM
type setIamPolicy struct {
	Policy iamPolicy `json:"policy,omitempty"`
}

// setIAMPermission set permissions for a member
func setIAMPermission(endpoint string, name string, memberName string, role string, memberType string) (err error) {

	u, _ := url.Parse(endpoint)
	u.Path = path.Join(u.Path, name+":getIamPolicy")
	getIamPolicyBody, err := HttpClient(false, u.String())
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	getIamPolicy := iamPolicy{}

	err = json.Unmarshal(getIamPolicyBody, &getIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	foundRole := false
	for i, binding := range getIamPolicy.Bindings {
		if binding.Role == role {
			//found members with the role already, add the new SA to the role
			getIamPolicy.Bindings[i].Members = append(binding.Members, memberType+":"+memberName)
			foundRole = true
		}
	}

	//no members with the role, add a new one
	if !foundRole {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, memberType+":"+memberName)
		getIamPolicy.Bindings = append(getIamPolicy.Bindings, binding)
	}

	u, _ = url.Parse(endpoint)
	u.Path = path.Join(u.Path, name+":setIamPolicy")

	setIamPolicy := setIamPolicy{}
	setIamPolicy.Policy = getIamPolicy

	setIamPolicyBody, err := json.Marshal(setIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	_, err = HttpClient(false, u.String(), string(setIamPolicyBody))

	return err
}

// SetConnectorIAMPermission set permissions for a member on a connection
func SetConnectorIAMPermission(name string, memberName string, iamRole string, memberType string) (err error) {
	var role string

	switch iamRole {
	case "admin":
		role = "roles/connectors.admin"
	case "invoker":
		role = "roles/connectors.invoker"
	case "viewer":
		role = "roles/connectors.viewer"
	default: //assume this is a custom role definition
		re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/roles\/([a-zA-Z0-9_-]+)`)
		result := re.FindString(iamRole)
		if result == "" {
			return fmt.Errorf("custom role must be of the format projects/{project-id}/roles/{role-name}")
		}
		role = iamRole
	}

	return setIAMPermission(GetBaseConnectorURL(), name, memberName, role, memberType)
}

// SetPubSubIAMPermission set permissions for a SA on a topic
func SetPubSubIAMPermission(project string, topic string, memberName string) (err error) {
	var endpoint = fmt.Sprintf("https://pubsub.googleapis.com/v1/projects/%s/topics", project)
	const memberType = "serviceAccount"
	const role = "roles/pubsub.publisher"
	return setIAMPermission(endpoint, topic, memberName, role, memberType)
}

func SetBigQueryIAMPermission(project string, datasetid string, memberName string) (err error) {
	var endpoint = fmt.Sprintf("https://bigquery.googleapis.com/bigquery/v2/projects/%s/datasets/%s", project, datasetid)
	const role = "WRITER"
	var content []byte

	//first fetch the information
	respBody, err := HttpClient(false, endpoint)
	if err != nil {
		return err
	}

	type accessType struct {
		Role         string  `json:"role,omitempty"`
		IamMember    *string `json:"iamMember,omitempty"`
		UserByEmail  *string `json:"userByEmail,omitempty"`
		SpecialGroup *string `json:"specialGroup,omitempty"`
		GroupByEmail *string `json:"groupByEmail,omitempty"`
	}

	type datasetType struct {
		Access []accessType `json:"access,omitempty"`
	}

	dataset := datasetType{}
	if err = json.Unmarshal(respBody, &dataset); err != nil {
		return err
	}

	access := accessType{}
	access.Role = role
	access.UserByEmail = new(string)
	*access.UserByEmail = memberName

	//merge the updates
	dataset.Access = append(dataset.Access, access)

	if content, err = json.Marshal(dataset); err != nil {
		return err
	}

	//patch the update
	if _, err = HttpClient(false, endpoint, string(content), "PATCH"); err != nil {
		return err
	}

	return nil
}
