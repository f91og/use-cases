package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Cluster struct {
	ID           json.Number `json:"id,omitempty"`
	Name         string      `json:"name"`
	ZkHosts      string      `json:"zkHosts,omitempty"`
	ModifiedTime string      `json:"modifiedTime,omitempty"`
}

type Config struct {
	ID           json.Number `json:"id,omitempty"`
	Name         string      `json:"name"`
	Content      string      `json:"content"`
	ClusterId    json.Number `json:"clusterId,omitempty"`
	ContentMd5   string      `json:"contentMd5,omitempty"`
	Status       string      `json:"status,omitempty"`
	ModifiedTime string      `json:"modifiedTime,omitempty"`
}

type ClusterAndConfig struct {
	CanalCluster Cluster `json:"canalCluster"`
	CanalConfig  Config  `json:"canalConfig"`
}

type NodeServer struct {
}

type Instance struct {
	ID              json.Number `json:"id,omitempty"`
	ClusterID       json.Number `json:"clusterId,omitempty"`
	Cluster         Cluster     `json:"canalCluster,omitempty"`
	ServerId        json.Number `json:"serverId,omitempty"`
	NodeServer      NodeServer  `json:"nodeServer,omitempty"`
	Name            string      `json:"name"`
	Content         string      `json:"content"`
	ContentMd5      string      `json:"contentMd5,omitempty"`
	Status          string      `json:"status,omitempty"`
	ModifiedTime    string      `json:"modifiedTime,omitempty"`
	ClusterServerId string      `json:"clusterServerId,omitempty"`
	RunningStatus   string      `json:"runningStatus"`
}

type InstanceConfig struct {
	ID           json.Number `json:"id,omitempty"`
	ClusterId    string      `json:"clusterServerId"`
	Content      string      `json:"content"`
	InstanceName string      `json:"name"`
	Status       string      `json:"status,omitempty"`
}

const instanceOldApiPath = "api/v1/canal/instance"
const clusterApiPath = "v2/api/canal/clusters"
const instanceApiPath = "v2/api/canal/instances"

func (c *Client) GetInstance(instanceId string) (*Instance, error) {
	log.Printf("GetInstance: %+v", instanceId)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, instanceOldApiPath, instanceId), nil)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	instance := Instance{}
	err = json.Unmarshal(data, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) CreateInstance(instanceConfig InstanceConfig) (string, error) {
	log.Printf("CreateInstance: %+v", instanceConfig)
	jsonEncodedBody, err := json.Marshal(instanceConfig)
	if err != nil {
		return "", err
	}
	log.Printf("requestBody: %s", string(jsonEncodedBody))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, instanceApiPath), strings.NewReader(string(jsonEncodedBody)))
	if err != nil {
		return "", err
	}

	data, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	instance := Instance{}
	err = json.Unmarshal(data, &instance)
	if err != nil {
		return "", err
	}

	return string(instance.ID), nil
}

func (c *Client) UpdateInstance(instanceConfig InstanceConfig) error {
	log.Printf("UpdateInstance: %+v", instanceConfig)
	jsonEncodedBody, err := json.Marshal(instanceConfig)
	if err != nil {
		return err
	}
	log.Printf("requestBody: %s", string(jsonEncodedBody))

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", c.HostURL, instanceApiPath), strings.NewReader(string(jsonEncodedBody)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteInstance(instanceId string) error {
	log.Printf("DeleteInstance: %+v", instanceId)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, instanceOldApiPath, instanceId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCluster(clusterId string) (*ClusterAndConfig, error) {
	log.Printf("GetCluster: %+v", clusterId)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, clusterApiPath, clusterId), nil)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	clusterAndConfig := ClusterAndConfig{}
	err = json.Unmarshal(data, &clusterAndConfig)
	if err != nil {
		return nil, err
	}

	return &clusterAndConfig, nil
}

func (c *Client) GetClusterByName(clusterName string) (*Cluster, error) {
	log.Printf("GetCluster: %+v", clusterName)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?name=%s", c.HostURL, clusterApiPath, clusterName), nil)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	clusters := []Cluster{}
	err = json.Unmarshal(data, &clusters)
	if err != nil {
		return nil, err
	}
	if len(clusters) != 1 {
		return nil, fmt.Errorf("failed to fetch the cluster from remote. name=%s, len=%d ", clusterName, len(clusters))
	}
	return &clusters[0], nil
}

func (c *Client) CreateCluster(clusterAndConfig ClusterAndConfig) (string, error) {
	log.Printf("CreateCluster: %+v", clusterAndConfig)
	jsonEncodedBody, err := json.Marshal(clusterAndConfig)
	if err != nil {
		return "", err
	}
	log.Printf("requestBody: %s", string(jsonEncodedBody))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, clusterApiPath), strings.NewReader(string(jsonEncodedBody)))
	if err != nil {
		return "", err
	}

	data, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	resClusterAndConfig := ClusterAndConfig{}
	err = json.Unmarshal(data, &resClusterAndConfig)
	if err != nil {
		return "", err
	}

	return string(resClusterAndConfig.CanalCluster.ID), nil
}

func (c *Client) UpdateCluster(clusterAndConfig ClusterAndConfig) error {
	log.Printf("UpdateCluster: obj: %+v", clusterAndConfig)
	jsonEncodedBody, err := json.Marshal(clusterAndConfig)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", c.HostURL, clusterApiPath, clusterAndConfig.CanalCluster.ID), strings.NewReader(string(jsonEncodedBody)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteCluster(clusterId string) error {
	log.Printf("DeleteCluster: %+v", clusterId)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, clusterApiPath, clusterId), nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return err
	}

	log.Printf("DeleteCluster Response: %+v", res)
	return nil
}
