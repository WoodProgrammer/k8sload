package lib

type TopologySpec struct {
	Port         int                 `yaml:"port"`
	Replicas     int                 `yaml:"replicas"`
	AntiAffinity bool                `yaml:"antiAffinity"`
	TopologyKeys []map[string]string `yaml:"topologyKeys"`
	Args         []map[string]string `yaml:"args"`
	Commands     []map[string]string `yaml:"commands"`
	HostNetwork  bool                `yaml:"hostNetwork"`
}

type Producer struct {
	Name      string       `yaml:"name"`
	Namespace string       `yaml:"namespace"`
	Spec      TopologySpec `yaml:"spec"`
}

type Topology struct {
	Producer Producer `yaml:"producer"`
}

type Config struct {
	Topology Topology `yaml:"topology"`
}
