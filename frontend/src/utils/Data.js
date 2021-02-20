export function genRandomTree(N = 300, reverse = false) {
  const data = {
    nodes: [...Array(N).keys()].map((i) => ({
      id: i,
      name: ['Intro to Machine Learning', 'Advanced Data Structures'][
        Math.round(Math.random() * 2) - 1
      ],
    })),
    links: [...Array(N).keys()]
      .filter((id) => id)
      .map((id) => ({
        [reverse ? 'target' : 'source']: id,
        [reverse ? 'source' : 'target']: Math.round(Math.random() * (id - 1)),
        name: ['Intro to Machine Learning', 'Advanced Data Structures'][
          Math.round(Math.random() * 2) - 1
        ],
      })),
  };

  data.links.forEach((link) => {
    const a = data.nodes[link.source];
    const b = data.nodes[link.target];
    if (!a.neighbors) {
      a.neighbors = [];
    }

    if (!b.neighbors) {
      b.neighbors = [];
    }

    a.neighbors.push(b);
    b.neighbors.push(a);

    if (!a.links) {
      a.links = [];
    }

    if (!b.links) {
      b.links = [];
    }
    a.links.push(link);
    b.links.push(link);
  });

  return data;
}
