/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package system

import (
    "github.com/gogo/protobuf/proto"
    "github.com/stretchr/testify/assert"
    "testing"
)

func Test_Nodes_AddNodes(t *testing.T) {
    addNodes := []NNode{
        {
            Database: "postgres",
        },
        {
            Database: "postgres1",
        },
    }

    for _, node := range addNodes {
        newNode, err := SystemCtx.Nodes.AddNode(node)
        if err != nil {
            t.Error(err)
            t.FailNow()
        }

        node.NodeId = newNode.NodeId
        assert.True(t, proto.Equal(&node, newNode), "returned node does not equal the expected result")
    }


    n, err := SystemCtx.Nodes.GetNodes()
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    assert.NotNil(t, n, "nodes should not be null")

    assert.NotEmpty(t, n, "no nodes found, there should be at least 1")

    for _, node := range n {
        assert.True(t, node.NodeId > 0, "node ID is not greater than 0, this means that its possible that the node wasn't created")
    }
}

func Test_Nodes_GetNodes(t *testing.T) {
    n, err := SystemCtx.Nodes.GetNodes()
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    assert.NotNil(t, n, "nodes should not be null")

    assert.NotEmpty(t, n, "no nodes found, there should be at least 1")

    for _, node := range n {
        assert.True(t, node.NodeId > 0, "node ID is not greater than 0, this means that its possible that the node wasn't created")
    }
}
