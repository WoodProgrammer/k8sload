package lib

const ConsumerSvcTemplate = `
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Topology.Consumer.Name }}
  namespace: {{ .Topology.Consumer.Namespace }}
  labels:
    app: {{ .Topology.Consumer.Name }}
  annotations:
    prometheus.io/scrape: "true"            # Enable scraping
    prometheus.io/port: "{{ .Topology.Consumer.Spec.ExporterPort }}"              # Port the exporter listens on
    prometheus.io/path: "/metrics"          # Metrics endpoint (default for most exporters)
spec:
  selector:
    app: {{ .Topology.Consumer.Name }}
  ports:
    - name: http
      port: {{ .Topology.Consumer.Spec.ExporterPort }}
      targetPort: {{ .Topology.Consumer.Spec.ExporterPort }}
  type: ClusterIP
`

const ConsumerDeploymentTemplate = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Topology.Consumer.Name }}
  namespace: {{ .Topology.Consumer.Namespace }}
spec:
  replicas: {{ .Topology.Consumer.Spec.Replicas }}
  selector:
    matchLabels:
      app: {{ .Topology.Consumer.Name }}
  template:
    metadata:
      labels:
        app: {{ .Topology.Consumer.Name }}
    spec:
      hostNetwork:  {{ .Topology.Consumer.Spec.HostNetwork }}
      volumes:
      - name: metrics-data
        emptyDir: {}
      containers:
      - name: {{ .Topology.Consumer.Name }}-exporter
        image: {{ .Topology.Consumer.Spec.ExporterImage }}
        env:
        - name: "LOAD_TEST_NAME"
          value: "{{ .Topology.Name }}"
        volumeMounts:
        - name: metrics-data
          mountPath: /opt/metrics/
        ports:
        - containerPort: 9001
      - name: {{ .Topology.Consumer.Name }}
        image: {{ .Topology.Consumer.Spec.Image }}
        volumeMounts:
        - name: metrics-data
          mountPath: /opt/metrics/
        command: {{ .Topology.Consumer.Spec.Commands }}
        args:
        {{- range $keyVal := .Topology.Consumer.Spec.Args }}
        - "{{ $keyVal }}"
        {{- end }}
        ports:
        - containerPort: {{ .Topology.Consumer.Spec.Port }}
{{- if .Topology.Consumer.Spec.AntiAffinity }}
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
{{- range $keyVal := .Topology.Consumer.Spec.TopologyKeys }}
                {{- range $key, $val := $keyVal }}
              - key: {{ $key }}
                operator: In
                values:
                - {{ $val }}
                {{- end }}
{{- end }}
            topologyKey: "kubernetes.io/hostname"
{{- end }}

`

const ProducerDeploymentTemplate = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Topology.Producer.Name }}
  namespace: {{ .Topology.Producer.Namespace }}
spec:
  replicas: {{ .Topology.Producer.Spec.Replicas }}
  selector:
    matchLabels:
      app: {{ .Topology.Producer.Name }}
  template:
    metadata:
      labels:
        app: {{ .Topology.Producer.Name }}
    spec:
      hostNetwork:  {{ .Topology.Producer.Spec.HostNetwork }}
      containers:
      - name: {{ .Topology.Producer.Name }}
        image: {{ .Topology.Producer.Spec.Image }}
        command: {{ .Topology.Producer.Spec.Commands }}
        args:
        {{- range $keyVal := .Topology.Producer.Spec.Args }}
        - "{{ $keyVal }}"
        {{- end }}
        ports:
        - containerPort: {{ .Topology.Producer.Spec.Port }}
{{- if .Topology.Producer.Spec.AntiAffinity }}
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
{{- range $keyVal := .Topology.Producer.Spec.TopologyKeys }}
                {{- range $key, $val := $keyVal }}
              - key: {{ $key }}
                operator: In
                values:
                - {{ $val }}
                {{- end }}
{{- end }}
            topologyKey: "kubernetes.io/hostname"
{{- end }}
`

const ProducerSvcTemplate = `
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Topology.Producer.Name }}
  namespace: {{ .Topology.Producer.Namespace }}
  labels:
    app: {{ .Topology.Producer.Name }}
spec:
  selector:
    app: {{ .Topology.Producer.Name }}
  ports:
    - name: http
      port: {{ .Topology.Producer.Spec.Port }}
      targetPort: {{ .Topology.Producer.Spec.Port }}
  type: ClusterIP
`
