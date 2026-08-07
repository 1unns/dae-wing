/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2023, daeuniverse Organization <team@v2raya.org>
 */

package general

import (
	"strconv"

	"github.com/daeuniverse/dae-wing/dae"
	"github.com/graph-gophers/graphql-go"
)

type RuntimeOverviewResolver struct {
	Overview *dae.RuntimeOverview
}

type RuntimeTrafficSampleResolver struct {
	Sample dae.RuntimeTrafficSample
}

type DeviceTrafficResolver struct {
	Device dae.DeviceTraffic
}

type ConnTrafficResolver struct {
	Conn dae.ConnTraffic
}

func (r *Resolver) RuntimeOverview(args *struct {
	WindowSec int32
	MaxPoints int32
}) (*RuntimeOverviewResolver, error) {
	overview, err := dae.GetRuntimeOverview(int(args.WindowSec), int(args.MaxPoints))
	if err != nil {
		return nil, err
	}
	return &RuntimeOverviewResolver{Overview: overview}, nil
}

func (r *RuntimeOverviewResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.Overview.UpdatedAt}
}

func (r *RuntimeOverviewResolver) UploadRate() float64 {
	return float64(r.Overview.UploadRate)
}

func (r *RuntimeOverviewResolver) DownloadRate() float64 {
	return float64(r.Overview.DownloadRate)
}

func (r *RuntimeOverviewResolver) UploadTotal() string {
	return strconv.FormatUint(r.Overview.UploadTotal, 10)
}

func (r *RuntimeOverviewResolver) DownloadTotal() string {
	return strconv.FormatUint(r.Overview.DownloadTotal, 10)
}

func (r *RuntimeOverviewResolver) ActiveConnections() int32 {
	return int32(r.Overview.ActiveConnections)
}

func (r *RuntimeOverviewResolver) UdpSessions() int32 {
	return int32(r.Overview.UDPSessions)
}

func (r *RuntimeOverviewResolver) Samples() []*RuntimeTrafficSampleResolver {
	resolvers := make([]*RuntimeTrafficSampleResolver, 0, len(r.Overview.Samples))
	for _, sample := range r.Overview.Samples {
		resolvers = append(resolvers, &RuntimeTrafficSampleResolver{Sample: sample})
	}
	return resolvers
}

func (r *RuntimeOverviewResolver) DeviceTraffics() []*DeviceTrafficResolver {
	resolvers := make([]*DeviceTrafficResolver, 0, len(r.Overview.DeviceTraffics))
	for _, dt := range r.Overview.DeviceTraffics {
		resolvers = append(resolvers, &DeviceTrafficResolver{Device: dt})
	}
	return resolvers
}

func (r *RuntimeOverviewResolver) ConnTraffics() []*ConnTrafficResolver {
	resolvers := make([]*ConnTrafficResolver, 0, len(r.Overview.ConnTraffics))
	for _, ct := range r.Overview.ConnTraffics {
		resolvers = append(resolvers, &ConnTrafficResolver{Conn: ct})
	}
	return resolvers
}

func (r *RuntimeTrafficSampleResolver) Timestamp() graphql.Time {
	return graphql.Time{Time: r.Sample.Timestamp}
}

func (r *RuntimeTrafficSampleResolver) UploadRate() float64 {
	return float64(r.Sample.UploadRate)
}

func (r *RuntimeTrafficSampleResolver) DownloadRate() float64 {
	return float64(r.Sample.DownloadRate)
}

func (r *DeviceTrafficResolver) Ip() string {
	return r.Device.IP
}

func (r *DeviceTrafficResolver) Mac() string {
	return ""
}

func (r *DeviceTrafficResolver) Name() string {
	return ""
}

func (r *DeviceTrafficResolver) ProxyUploadTotal() string {
	return strconv.FormatUint(r.Device.ProxyUploadTotal, 10)
}

func (r *DeviceTrafficResolver) ProxyDownloadTotal() string {
	return strconv.FormatUint(r.Device.ProxyDownloadTotal, 10)
}

func (r *DeviceTrafficResolver) DirectUploadTotal() string {
	return strconv.FormatUint(r.Device.DirectUploadTotal, 10)
}

func (r *DeviceTrafficResolver) DirectDownloadTotal() string {
	return strconv.FormatUint(r.Device.DirectDownloadTotal, 10)
}

func (r *ConnTrafficResolver) SrcIp() string {
	return r.Conn.SrcIP
}

func (r *ConnTrafficResolver) DstIp() string {
	return r.Conn.DstIP
}

func (r *ConnTrafficResolver) DstPort() int32 {
	return int32(r.Conn.DstPort)
}

func (r *ConnTrafficResolver) Id() graphql.ID {
	return graphql.ID(fmt.Sprintf("%s:%d-%s:%d", r.Conn.SrcIP, r.Conn.SrcPort, r.Conn.DstIP, r.Conn.DstPort))
}

func (r *ConnTrafficResolver) Domain() string {
	return ""
}

func (r *ConnTrafficResolver) Ip() string {
	return r.Conn.DstIP
}

func (r *ConnTrafficResolver) UploadTotal() string {
	return strconv.FormatUint(r.Conn.UploadTotal, 10)
}

func (r *ConnTrafficResolver) DownloadTotal() string {
	return strconv.FormatUint(r.Conn.DownloadTotal, 10)
}

func (r *ConnTrafficResolver) State() string {
	return "ESTABLISHED"
}
