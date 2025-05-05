package main

type Metric struct {
	Intervals []Intervals `json:"intervals"`
}

type Intervals struct {
	Streams []Streams `json:"streams"`
}

type Streams struct {
	Socket        int     `json:"socket"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int     `json:"retransmits"`
	SndCwd        int     `json:"snd_cwnd"`
	SndWnd        int     `json:"snd_wnd"`
	RTT           int     `json:"rtt"`
	PMTU          int     `json:"pmtu"`
	Omitted       bool    `json:"omitted"`
	Sender        bool    `json:"sender"`
}
