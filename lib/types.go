package lib

type TopologySpec struct {
	Port          int                 `yaml:"port"`
	Replicas      int                 `yaml:"replicas"`
	AntiAffinity  bool                `yaml:"antiAffinity"`
	TopologyKeys  []map[string]string `yaml:"topologyKeys"`
	Args          []string            `yaml:"args"`
	Commands      []string            `yaml:"commands"`
	HostNetwork   bool                `yaml:"hostNetwork"`
	ExporterImage string              `yaml:"exporterImage"`
	Image         string              `yaml:"image"`
}

type Producer struct {
	Name      string       `yaml:"name"`
	Namespace string       `yaml:"namespace"`
	Spec      TopologySpec `yaml:"spec"`
}

type Topology struct {
	Producer Producer `yaml:"producer"`
	Consumer Consumer `yaml:"consumer"`
}

type Config struct {
	Topology Topology `yaml:"topology"`
}

type Consumer struct {
	Name      string       `yaml:"name"`
	Namespace string       `yaml:"namespace"`
	Spec      TopologySpec `yaml:"spec"`
}
