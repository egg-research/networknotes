import React, { useState, useCallback } from 'react';

import { ForceGraph2D } from 'react-force-graph';

const NODE_R = 8;

function Graph({ data, height, width }) {
  const [highlights, setHighlight] = useState({
    nodes: new Set(),
    edges: new Set(),
    hoverNode: null,
  });

  const handleNodeHover = (node) => {
    const newHighlights = {
      nodes: new Set(),
      edges: new Set(),
      hoverNode: null,
    };

    if (node) {
      node.neighbors.forEach((neighbor) => newHighlights.nodes.add(neighbor));
      node.links.forEach((edge) => newHighlights.edges.add(edge));
      newHighlights.nodes.add(node);
      newHighlights.hoverNode = node;
    }

    setHighlight(newHighlights);
  };

  const paintRing = useCallback(
    (node, ctx) => {
      ctx.beginPath();
      ctx.arc(node.x, node.y, NODE_R * 1.4, 0, 2 * Math.PI, false);
      ctx.fillStyle = node === highlights.hoverNode ? 'red' : 'orange';
      ctx.fill();
    },
    [highlights]
  );

  return (
    <ForceGraph2D
      graphData={data}
      nodeRelSize={NODE_R}
      height={height}
      width={width}
      linkWidth={(link) => (highlights.edges.has(link) ? 5 : 1)}
      nodeCanvasObjectMode={(node) =>
        highlights.nodes.has(node) ? 'before' : undefined
      }
      nodeCanvasObject={paintRing}
      onNodeHover={handleNodeHover}
    />
  );
}

export default Graph;
