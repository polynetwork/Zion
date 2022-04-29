/*
 * Copyright (C) 2022 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */


package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/tool"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/urfave/cli.v1"
	"strconv"
)

var (
	genesisToolCommand = cli.Command{
		Name: "genesisTool",
		Usage: "A set of commands facilitating generating genesis configuration of maas chain",
		Category: "MISCELLANEOUS COMMANDS",
		Description: "",
		Subcommands: []cli.Command{
			{
				Name: "generate",
				Flags: []cli.Flag{
					basePathFlag,
				},
				Action: utils.MigrateFlags(generateMaasGensis),
			},
		},
	}
)

var (
	basePathFlag = cli.StringFlag{
		Name:  "basePath",
		Usage: "The path to store genesis configuration files",
	}
)

func generateMaasGensis(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}
	basePath := ctx.String(basePathFlag.Name)
	if basePath == ""{
		basePath = utils.DefaultBasePath()
	} else if basePath[len(basePath)-1:len(basePath)] != "/" {
		basePath += "/"
	}
	nodeNum, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		utils.Fatalf("parse node number error %s", err.Error())
	}

	if nodeNum < 4 {
		utils.Fatalf("got %v nodes, but hotstuff BFT requires at least 4 nodes", nodeNum)
	}

	genesis := new(utils.MaasGenesis)
	genesis.Default()

	staticNodes := make([]string, 0)

	nodes := make([]*tool.Node, 0)
	metaNodes := make([]*utils.Node, 0)

	for i := 0; i < nodeNum; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)

		nodekey := hexutil.Encode(crypto.FromECDSA(key))
		nodeInf, _ := tool.NodeKey2NodeInfo(nodekey)
		staticInf := tool.NodeStaticInfoTemp(nodeInf)

		node := &tool.Node{
			Address: addr.Hex(),
			NodeKey: nodekey,
			Static:  staticInf,
		}
		nodes = append(nodes, node)
	}

	sortedNodes := tool.SortNodes(nodes)
	genesisExtra, err := tool.Encode(tool.NodesAddress(sortedNodes))
	if err != nil {
		utils.Fatalf(err.Error())
	}

	genesis.Alloc = make(map[string]struct {
		PublicKey string `json:"publicKey"`
		Balance   string `json:"balance"`
	}, 0)

	for _, v := range sortedNodes {
		nodeInf, err := tool.NodeKey2NodeInfo(v.NodeKey)
		if err != nil {
			utils.Fatalf(err.Error())
		}

		pubInf, err := tool.NodeKey2PublicInfo(v.NodeKey)
		if err != nil {
			utils.Fatalf(err.Error())
		}

		genesis.Alloc[v.Address] = struct {
			PublicKey string `json:"publicKey"`
			Balance   string `json:"balance"`
		}{
			pubInf,
			"100000000000000000000000000000",
		}

		metaNode := &utils.Node{
			Address: v.Address,
			NodeKey: v.NodeKey,
			PubKey: pubInf,
			Static: v.Static,
		}

		metaNodes = append(metaNodes, metaNode)
		staticNodes = append(staticNodes, tool.NodeStaticInfoTemp(nodeInf))
	}
	filePaths := [3]string{}
	contents := [3]string{}

	genesis.ExtraData = genesisExtra
	geneJson, _ := genesis.Encode()
	contents[0] = geneJson

	filePaths[0] = basePath + "genesis.json"
	staticNodesEnc, err := json.MarshalIndent(staticNodes, "", "\t")
	if err != nil {
		utils.Fatalf(err.Error())
	}
	contents[1] = string(staticNodesEnc)
	filePaths[1] = basePath + "static-nodes.json"

	sortedNodesEnc, _ := json.MarshalIndent(metaNodes,"", "\t")
	contents[2] = string(sortedNodesEnc)
	filePaths[2] = basePath + "nodes.json"
	err = utils.DumpGenesis(filePaths, contents)
	if err != nil {
		utils.Fatalf(err.Error())
	}
	return nil
}
