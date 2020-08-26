package merger

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type KubeConfigCluster struct {
	Cluster struct {
		Server                string
		InsecureSkipTLSVerify bool `yaml:"insecure-skip-tls-verify"`
	}
	Name string
}

type KubeConfigContext struct {
	Context struct {
		Cluster   string
		Namespace string
		User      string
	}
	Name string
}

type KubeConfigUser struct {
	User struct {
		Password              string `yaml:"password,omitempty"`
		Username              string `yaml:"username,omitempty"`
		Token                 string `yaml:"token,omitempty"`
		ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
		ClientKeyData         string `yaml:"client-key-data,omitempty"`
	}
	Name string
}

type KubeConfig struct {
	APIVersion     string `yaml:"apiVersion"`
	Kind           string
	Preferences    interface{}
	CurrentContext string `yaml:"current-context"`
	Clusters       []KubeConfigCluster
	Contexts       []KubeConfigContext
	Users          []KubeConfigUser

	clustersMap map[string]KubeConfigCluster `yaml:"-"`
	contextsMap map[string]KubeConfigContext `yaml:"-"`
	usersMap    map[string]KubeConfigUser    `yaml:"-"`

	ToAddClusters []KubeConfigCluster `yaml:"-"`
	ToAddContexts []KubeConfigContext `yaml:"-"`
	ToAddUsers    []KubeConfigUser    `yaml:"-"`

	IsChanged bool `yaml:"-"`
}

func NewKubeConfig(kubeconfigPath string) *KubeConfig {
	kubeconfig := readFileToYaml(kubeconfigPath)

	kubeconfig.clustersMap = make(map[string]KubeConfigCluster)
	kubeconfig.contextsMap = make(map[string]KubeConfigContext)
	kubeconfig.usersMap = make(map[string]KubeConfigUser)

	for _, cluster := range kubeconfig.Clusters {
		kubeconfig.clustersMap[cluster.Name] = cluster
	}

	for _, context := range kubeconfig.Contexts {
		kubeconfig.contextsMap[context.Name] = context
	}

	for _, user := range kubeconfig.Users {
		kubeconfig.usersMap[user.Name] = user
	}

	return &kubeconfig
}

func (k *KubeConfig) MergeNewConfig(newConfig KubeConfig) {

	for _, cluster := range newConfig.Clusters {
		if _, cExists := k.clustersMap[cluster.Name]; !cExists {
			k.clustersMap[cluster.Name] = cluster
			k.Clusters = append(k.Clusters, cluster)
			k.ToAddClusters = append(k.ToAddClusters, cluster)
			k.IsChanged = true
		} else {
			fmt.Fprintf(os.Stdout, "Already exists cluster with same name %s\n", cluster.Name)
		}
	}

	for _, context := range newConfig.Contexts {
		if _, cExists := k.contextsMap[context.Name]; !cExists {
			k.contextsMap[context.Name] = context
			k.Contexts = append(k.Contexts, context)
			k.ToAddContexts = append(k.ToAddContexts, context)
			k.IsChanged = true
		} else {
			fmt.Fprintf(os.Stdout, "Already exists context with same name %s\n", context.Name)
		}
	}

	for _, user := range newConfig.Users {
		if _, uExists := k.usersMap[user.Name]; !uExists {
			k.usersMap[user.Name] = user
			k.Users = append(k.Users, user)
			k.ToAddUsers = append(k.ToAddUsers, user)
			k.IsChanged = true
		} else {
			fmt.Fprintf(os.Stdout, "Already exists user with same name %s\n", user.Name)
		}
	}
}

func readFileToYaml(filePath string) KubeConfig {
	configFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var config KubeConfig
	err = yaml.Unmarshal([]byte(configBytes), &config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return config
}

func OverrideKubeconfig(config KubeConfig) {

}
