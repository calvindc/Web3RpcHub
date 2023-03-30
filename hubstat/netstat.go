package hubstat

import (
	"sync"

	"fmt"
	"sort"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/calvindc/Web3RpcHub/internal/broadcasts"
	"github.com/calvindc/Web3RpcHub/refs"
	"github.com/go-kit/kit/log"
)

//hubStatMap 如果管理图谱结构，修改muxrpc.Endpoint为map[ref.feed]muxrpc.Endpoint
type hubStatMap map[string]cmuxrpc.Endpoint

func (rsm hubStatMap) AsList() []string {
	memberList := make([]string, 0, len(rsm))
	for m := range rsm {
		memberList = append(memberList, m)
	}
	sort.Strings(memberList)
	//按照时间排序
	return memberList
}

//---------------------------------------------------------------------------

// HubNetManager
type HubNetManager struct {
	hubMu    *sync.Mutex
	hubStats hubStatMap
	logger   log.Logger

	endpointsUpdater        broadcasts.EndpointsEmitter
	endpointsbroadcaster    *broadcasts.EndpointsBroadcast
	participantsUpdater     broadcasts.AttendantsEmitter
	participantsbroadcaster *broadcasts.AttendantsBroadcast
}

// NewHubNetManager
func NewHubNetManager(log log.Logger) *HubNetManager {
	ee, eb := broadcasts.NewEndpointsEmitter()
	pe, pb := broadcasts.NewAttendantsEmitter()
	return &HubNetManager{
		hubMu:                   new(sync.Mutex),
		hubStats:                make(hubStatMap),
		logger:                  log,
		endpointsUpdater:        ee,
		endpointsbroadcaster:    eb,
		participantsUpdater:     pe,
		participantsbroadcaster: pb,
	}
}

func (hm *HubNetManager) RegisterLegacyEndpoints(sink broadcasts.EndpointsEmitter) {
	hm.endpointsbroadcaster.Register(sink)
}

func (hm *HubNetManager) RegisterAttendantsUpdates(sink broadcasts.AttendantsEmitter) {
	hm.participantsbroadcaster.Register(sink)
}

func (hm *HubNetManager) List() []string {
	hm.hubMu.Lock()
	defer hm.hubMu.Unlock()
	return hm.hubStats.AsList()
}

func (m *HubNetManager) ListAsRefs() []refs.FeedRef {
	m.hubMu.Lock()
	lst := m.hubStats.AsList()
	m.hubMu.Unlock()

	rlst := make([]refs.FeedRef, len(lst))
	for i, s := range lst {
		fr, err := refs.ParseFeedRef(s)
		if err != nil {
			panic(fmt.Errorf("invalid feed ref in hub state: %d: %s", i, err))
		}
		rlst[i] = fr
	}
	return rlst
}

// AddEndpoint adds the endpoint to the hub
func (hm *HubNetManager) AddEndpoint(who refs.FeedRef, edp cmuxrpc.Endpoint) {
	hm.hubMu.Lock()
	// add ref to to the hub map
	hm.hubStats[who.String()] = edp
	currentMembers := hm.hubStats.AsList()
	hm.hubMu.Unlock()
	// update all the connected tunnel.endpoints calls
	hm.endpointsUpdater.Update(currentMembers)
	// update all the connected hub.patricipants calls
	hm.participantsUpdater.Joined(who)
}

// Remove removes the peer from the hub
func (hm *HubNetManager) Remove(who refs.FeedRef) {
	hm.hubMu.Lock()
	// remove ref from lobby
	delete(hm.hubStats, who.String())
	currentMembers := hm.hubStats.AsList()
	hm.hubMu.Unlock()
	// update all the connected tunnel.endpoints calls
	hm.endpointsUpdater.Update(currentMembers)
	// update all the connected hub.attendants calls
	hm.participantsUpdater.Left(who)
}

// AlreadyAdded returns true if the peer was already added to the hub.
// if it isn't it will be added.
func (hm *HubNetManager) AlreadyAdded(who refs.FeedRef, edp cmuxrpc.Endpoint) bool {
	hm.hubMu.Lock()
	var currentMembers []string
	// if the peer didn't call tunnel.announce()
	_, has := hm.hubStats[who.String()]
	if !has {
		// register them as if they didnt
		hm.hubStats[who.String()] = edp
		currentMembers = hm.hubStats.AsList()
	}
	hm.hubMu.Unlock()

	if !has {
		// update everyone
		hm.endpointsUpdater.Update(currentMembers)
		hm.participantsUpdater.Joined(who)
	}

	return has
}

// Has returns true and the endpoint if the peer is in the hub
func (hm *HubNetManager) Has(who refs.FeedRef) (cmuxrpc.Endpoint, bool) {
	hm.hubMu.Lock()
	// add ref to to the hub map
	edp, has := hm.hubStats[who.String()]
	hm.hubMu.Unlock()
	return edp, has
}
