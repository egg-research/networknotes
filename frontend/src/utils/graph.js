// const data = {
//   nodes: [...Array(N).keys()].map((i) => ({
//     id: i,
//     name: [
//       'Intro to Machine Learning',
//       'Computer Vision',
//       'Advanced Data Structures',
//     ][i],
//   })),
//   links: [
//     {
//       source: 0,
//       target: 1,
//       name: 'RNN',
//     },
//     {
//       source: 0,
//       target: 1,
//       name: 'CNN',
//     },
//   ],
// };

export function processGraph(graph) {
  console.log(graph)

  const nodeMap = {};
  for (const node of graph.nodes) {
    nodeMap[node.id] = node;
  }

  graph.nodes.forEach((node) => {
    node.neighbors = [];
    node.links = [];
  });

  graph.links.forEach((link) => {
    const a = nodeMap[link.source];
    const b = nodeMap[link.target];
    if (!a || !b) return

    a.neighbors.push(b);
    b.neighbors.push(a);

    a.links.push(link);
    b.links.push(link);
  });

  return graph;
}

function applyNodeFilter(graph, nodeFilter) {
  const newGraph = {};

  if (!nodeFilter.size) {
    Object.assign(newGraph, graph);
    return newGraph;
  }

  // get og
  const nodeFilterIdSet = new Set(Array.from(nodeFilter).map((x) => x.id));
  const nodes = graph.nodes.filter((node) => nodeFilterIdSet.has(node.id));

  const neighbors = graph.nodes.filter((currNode) => {
    for (const node of nodes) {
      if (node.neighbors.includes(currNode)) {
        return true;
      }
    }
    return false;
  });

  const nodesAndNeighbors = nodes.concat(neighbors);
  const nodeSet = new Set(nodesAndNeighbors.map((node) => node.id));

  const links = graph.links.filter(
    (link) => nodeSet.has(link.source) && nodeSet.has(link.target)
  );

  newGraph.nodes = nodesAndNeighbors;
  newGraph.links = links;

  const lol = processGraph(newGraph);
  console.log('after process');
  console.log(lol);

  return lol;
}

function applyLinkFilter(graph, linkFilter) {
  const newGraph = {};

  if (!linkFilter.size) {
    Object.assign(newGraph, graph);
    return newGraph;
  }

  const nodeMap = {};
  for (const node of graph.nodes) {
    nodeMap[node.id] = node;
  }

  const linkFilterIdSet = new Set(Array.from(linkFilter).map((x) => x.id));
  const links = graph.links.filter((link) => linkFilterIdSet.has(link.id));
  const nodeSet = new Set();
  for (const link of links) {
    nodeSet.add(nodeMap[link.source]);
    nodeSet.add(nodeMap[link.target]);
  }

  newGraph.nodes = Array.from(nodeSet);
  newGraph.links = links;

  return processGraph(newGraph);
}

// given a graph of the above format, and a set of nodeFilters, and a set of linkFilters return a
// new graph that only contains the filtered things
export function applyGraphFilter(graph, nodeFilter, linkFilter) {
  // only show nodes that are in nodeFilter and their neighbors
  // only show links that are in linkFilter
  //
  let newGraph = {};
  if (nodeFilter.length + linkFilter.length === 0) {
    Object.assign(newGraph, graph);
    return graph;
  }

  newGraph = applyNodeFilter(graph, nodeFilter);
  newGraph = applyLinkFilter(newGraph, linkFilter);

  return newGraph;
}
