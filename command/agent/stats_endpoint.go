package agent

import (
	"net/http"
	"strings"

	"github.com/hashicorp/nomad/client/stats"
	cstructs "github.com/hashicorp/nomad/client/structs"
	"github.com/hashicorp/nomad/nomad/structs"
)

func (s *HTTPServer) ClientStatsRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Get the requested Node ID
	requestedNode := req.URL.Query().Get("node_id")

	// Build the request and parse the ACL token
	args := structs.NodeSpecificRequest{
		NodeID: requestedNode,
	}
	s.parse(resp, req, &args.QueryOptions.Region, &args.QueryOptions)

	// Determine the handler to use
	useLocalClient, useClientRPC, useServerRPC := s.rpcHandlerForNode(requestedNode)

	// Make the RPC
	var reply cstructs.ClientStatsResponse
	var rpcErr error
	if useLocalClient {
		rpcErr = s.agent.Client().ClientRPC("ClientStats.Stats", &args, &reply)
	} else if useClientRPC {
		rpcErr = s.agent.Client().RPC("ClientStats.Stats", &args, &reply)
	} else if useServerRPC {
		rpcErr = s.agent.Server().RPC("ClientStats.Stats", &args, &reply)
	} else {
		rpcErr = CodedError(400, "No local Node and node_id not provided")
	}

	if rpcErr != nil {
		if structs.IsErrNoNodeConn(rpcErr) {
			rpcErr = CodedError(404, rpcErr.Error())
		} else if strings.Contains(rpcErr.Error(), "Unknown node") {
			rpcErr = CodedError(404, rpcErr.Error())
		}

		return nil, rpcErr
	}

	return reply.HostStats, nil
}

func (s *HTTPServer) ClientsStatsRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	args := structs.NodeListRequest{}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var nodes structs.NodeListResponse
	if err := s.agent.RPC("Node.List", &args, &nodes); err != nil {
		return nil, err
	}

	out := make(map[string]*stats.HostStats)
	for _, node := range nodes.Nodes {
		args := structs.NodeSpecificRequest{
			NodeID: node.ID,
		}
		s.parse(resp, req, &args.QueryOptions.Region, &args.QueryOptions)

		// Determine the handler to use
		useLocalClient, useClientRPC, useServerRPC := s.rpcHandlerForNode(node.ID)

		// Make the RPC
		var reply cstructs.ClientStatsResponse
		var rpcErr error
		if useLocalClient {
			rpcErr = s.agent.Client().ClientRPC("ClientStats.Stats", &args, &reply)
		} else if useClientRPC {
			rpcErr = s.agent.Client().RPC("ClientStats.Stats", &args, &reply)
		} else if useServerRPC {
			rpcErr = s.agent.Server().RPC("ClientStats.Stats", &args, &reply)
		} else {
			rpcErr = CodedError(400, "No local Node and node_id not provided")
		}

		if rpcErr != nil {
			if structs.IsErrNoNodeConn(rpcErr) {
				rpcErr = CodedError(404, rpcErr.Error())
			} else if strings.Contains(rpcErr.Error(), "Unknown node") {
				rpcErr = CodedError(404, rpcErr.Error())
			}

			return nil, rpcErr
		}

		out[node.ID] = reply.HostStats
	}

	return out, nil
}
