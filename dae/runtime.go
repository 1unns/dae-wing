/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2023, daeuniverse Organization <team@v2raya.org>
 */

package dae

import (
	"errors"
	"time"

	"github.com/daeuniverse/dae/control"
)

type RuntimeTrafficSample struct {
	Timestamp    time.Time
	UploadRate   uint64
	DownloadRate uint64
}

type DeviceTraffic struct {
	IP                  string
	ProxyUploadTotal    uint64
	ProxyDownloadTotal  uint64
	DirectUploadTotal   uint64
	DirectDownloadTotal uint64
}

type ConnTraffic struct {
	SrcIP         string
	DstIP         string
	DstPort       uint16
	UploadTotal   uint64
	DownloadTotal uint64
}

type RuntimeOverview struct {
	UpdatedAt         time.Time
	UploadRate        uint64
	DownloadRate      uint64
	UploadTotal       uint64
	DownloadTotal     uint64
	ActiveConnections int
	UDPSessions       int
	Samples           []RuntimeTrafficSample
	DeviceTraffics    []DeviceTraffic
	ConnTraffics      []ConnTraffic
}

func GetRuntimeOverview(windowSec int, maxPoints int) (*RuntimeOverview, error) {
	ctl, err := ControlPlane()
	if err != nil && !errors.Is(err, ErrControlPlaneNotInit) {
		return nil, err
	}

	// Use ctl.SnapshotRuntimeStats() so it reads from the correct ControlPlane
	// runtimeStats store (which records actual proxy traffic) and eBPF maps.
	// When ControlPlane is not yet initialized, fall back to the global store.
	var snapshot control.RuntimeStatsSnapshot
	if ctl != nil {
		snapshot = ctl.SnapshotRuntimeStats(windowSec, maxPoints)
	} else {
		activeTCPConnections := 0
		snapshot = control.SnapshotRuntimeStats(activeTCPConnections, control.DefaultUdpEndpointPool.Count(), windowSec, maxPoints)
	}

	samples := make([]RuntimeTrafficSample, 0, len(snapshot.Samples))
	for _, sample := range snapshot.Samples {
		samples = append(samples, RuntimeTrafficSample{
			Timestamp:    sample.Timestamp,
			UploadRate:   sample.UploadRate,
			DownloadRate: sample.DownloadRate,
		})
	}

	deviceTraffics := make([]DeviceTraffic, 0, len(snapshot.DeviceTraffics))
	for _, dt := range snapshot.DeviceTraffics {
		deviceTraffics = append(deviceTraffics, DeviceTraffic{
			IP:                  dt.IP,
			ProxyUploadTotal:    dt.ProxyUploadTotal,
			ProxyDownloadTotal:  dt.ProxyDownloadTotal,
			DirectUploadTotal:   dt.DirectUploadTotal,
			DirectDownloadTotal: dt.DirectDownloadTotal,
		})
	}

	connTraffics := make([]ConnTraffic, 0, len(snapshot.ConnTraffics))
	for _, ct := range snapshot.ConnTraffics {
		connTraffics = append(connTraffics, ConnTraffic{
			SrcIP:         ct.SrcIP,
			DstIP:         ct.DstIP,
			DstPort:       ct.DstPort,
			UploadTotal:   ct.UploadTotal,
			DownloadTotal: ct.DownloadTotal,
		})
	}

	return &RuntimeOverview{
		UpdatedAt:         snapshot.UpdatedAt,
		UploadRate:        snapshot.UploadRate,
		DownloadRate:      snapshot.DownloadRate,
		UploadTotal:       snapshot.UploadTotal,
		DownloadTotal:     snapshot.DownloadTotal,
		ActiveConnections: snapshot.ActiveConnections,
		UDPSessions:       snapshot.UDPSessions,
		Samples:           samples,
		DeviceTraffics:    deviceTraffics,
		ConnTraffics:      connTraffics,
	}, nil
}
