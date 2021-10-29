package stats

import (
	"reflect"
	"testing"
	"time"

	"github.com/Microsoft/hcsshim"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	statsapi "k8s.io/kubelet/pkg/apis/stats/v1alpha1"
)

type fakehcsshim struct {
	containers       []hcsshim.ContainerProperties
	hcsStats         []hcsshim.NetworkStats
	statsToContainer map[string]string
}

type fakeConatiner struct {
	stat hcsshim.Statistics
}

func (f fakeConatiner) Start() error {
	return nil
}

func (f fakeConatiner) Shutdown() error {
	return nil
}

func (f fakeConatiner) Terminate() error {
	return nil
}

func (f fakeConatiner) Wait() error {
	return nil
}

func (f fakeConatiner) WaitTimeout(duration time.Duration) error {
	return nil
}

func (f fakeConatiner) Pause() error {
	return nil
}

func (f fakeConatiner) Resume() error {
	return nil
}

func (f fakeConatiner) HasPendingUpdates() (bool, error) {
	return false, nil
}

func (f fakeConatiner) Statistics() (hcsshim.Statistics, error) {
	return f.stat, nil
}

func (f fakeConatiner) ProcessList() ([]hcsshim.ProcessListItem, error) {
	return []hcsshim.ProcessListItem{}, nil
}

func (f fakeConatiner) MappedVirtualDisks() (map[int]hcsshim.MappedVirtualDiskController, error) {
	return map[int]hcsshim.MappedVirtualDiskController{}, nil
}

func (f fakeConatiner) CreateProcess(c *hcsshim.ProcessConfig) (hcsshim.Process, error) {
	return nil, nil
}

func (f fakeConatiner) OpenProcess(pid int) (hcsshim.Process, error) {
	return nil, nil
}

func (f fakeConatiner) Close() error {
	return nil
}

func (f fakeConatiner) Modify(config *hcsshim.ResourceModificationRequestResponse) error {
	return nil
}

func (s fakehcsshim) GetContainers(q hcsshim.ComputeSystemQuery) ([]hcsshim.ContainerProperties, error) {
	return s.containers, nil
}

func (s fakehcsshim) GetHNSEndpointByID(endpointID string) (*hcsshim.HNSEndpoint, error) {
	e := hcsshim.HNSEndpoint{
		Name: endpointID,
	}
	return &e, nil
}

func (s fakehcsshim) OpenContainer(id string) (hcsshim.Container, error) {
	sId := s.statsToContainer[id]
	c := fakeConatiner{}
	for _, stat := range s.hcsStats {
		if stat.InstanceId == sId {
			c.stat.Network = append(c.stat.Network, stat)
		}
	}

	return c, nil
}

func Test_criStatsProvider_listContainerNetworkStats(t *testing.T) {
	tests := []struct {
		name    string
		fields  fakehcsshim
		want    map[string]*statsapi.NetworkStats
		wantErr bool
	}{
		{
			name: "basic example",
			fields: fakehcsshim{
				containers: []hcsshim.ContainerProperties{
					{ID: "c1"},
					{ID: "c2"},
				},
				hcsStats: []hcsshim.NetworkStats{
					{
						BytesReceived: 1,
						BytesSent:     10,
						EndpointId:    "test",
						InstanceId:    "",
					},
					{
						BytesReceived: 2,
						BytesSent:     20,
						EndpointId:    "test2",
						InstanceId:    "",
					},
				},
				statsToContainer: map[string]string{
					"c1": "test",
					"c2": "test2",
				},
			},
			want: map[string]*statsapi.NetworkStats{
				"c1": &statsapi.NetworkStats{
					Time: v1.Time{},
					InterfaceStats: statsapi.InterfaceStats{
						Name:    "test",
						RxBytes: toP(1),
						TxBytes: toP(10),
					},
					Interfaces: []statsapi.InterfaceStats{
						{
							Name:    "test",
							RxBytes: toP(1),
							TxBytes: toP(10),
						},
					},
				},
				"c2": &statsapi.NetworkStats{
					Time: v1.Time{},
					InterfaceStats: statsapi.InterfaceStats{
						Name:    "test2",
						RxBytes: toP(2),
						TxBytes: toP(2),
					},
					Interfaces: []statsapi.InterfaceStats{
						{
							Name:    "test2",
							RxBytes: toP(2),
							TxBytes: toP(20),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &criStatsProvider{
				hcsshimInterface: fakehcsshim{
					containers:       tt.fields.containers,
					hcsStats:         tt.fields.hcsStats,
					statsToContainer: tt.fields.statsToContainer,
				},
			}
			got, err := p.listContainerNetworkStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("listContainerNetworkStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listContainerNetworkStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func toP(i uint64) *uint64 {
	return &i
}
