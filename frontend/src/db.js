const baseUrl = 'http://networknotes2.wl.r.appspot.com';
const userId = 28;

// return as int
export async function getUserId(username) {
  const body = JSON.stringify({ AuthUid: username });
  const res = await fetch(`${baseUrl}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}

// returns doc id as int
export async function makeNewDoc(uid, docTitle) {
  const body = JSON.stringify({ Uid: uid, Doc: { DocName: docTitle } });
  const res = await fetch(`${baseUrl}/writeDoc`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}

export async function updateDoc(uid, docId, docTitle, docText, rawDocText) {
  const body = JSON.stringify({
    Uid: uid,
    Doc: {
      DocId: docId,
      DocName: docTitle,
      DocText: docText,
      RawDocText: rawDocText,
    },
  });

  const res = await fetch(`${baseUrl}/writeDocFS`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}

export async function readDoc(uid, docId) {
  const body = JSON.stringify({
    Uid: uid,
    Doc: {
      DocId: docId,
    },
  });

  const res = await fetch(`${baseUrl}/readDoc`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  const resJSON = await res.json();
  return {
    keywords: resJSON[0] === null ? [] : resJSON[0][0],
    document: {
      id: docId,
      title: resJSON[1].docName,
      text: resJSON[1].docText,
      rawText: resJSON[1].rawDocText === null ? '' : resJSON[1].rawDocText,
    },
  };
}

// return array of {id: "33", title: "title", text: "", rawText: ""}
export async function getAllDocs(uid) {
  const body = JSON.stringify({
    Uid: uid,
  });

  const res = await fetch(`${baseUrl}/getAllDocs`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  const resJSON = await res.json();
  const keys = Object.keys(resJSON);

  return keys.map((x) => {
    const doc = resJSON[x];
    return {
      id: parseInt(x, 10),
      title: doc.docName,
      text: doc.docText,
      rawText: doc.rawDocText === null ? '' : doc.rawDocText,
    };
  });
}

export async function getAllKeywords(uid) {
  const body = JSON.stringify({
    Uid: uid,
  });

  const res = await fetch(`${baseUrl}/getAllKws`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  const resJSON = await res.json();
  const keys = Object.keys(resJSON);

  return keys.map((x) => {
    const kw = resJSON[x];
    return {
      id: parseInt(x, 10),
      name: kw.kw,
      text: kw.kwText,
    };
  });
}

export async function getDocGraph(uid) {
  const body = JSON.stringify({
    Uid: uid,
  });

  const res = await fetch(`${baseUrl}/getAll?q=doc`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}

export async function getKeywordGraph(uid) {
  const body = JSON.stringify({
    Uid: uid,
  });

  const res = await fetch(`${baseUrl}/getAll?q=kw`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}

export async function updateDocKeyword(uid, docId, keywords) {
  const body = JSON.stringify({
    Uid: uid,
    Doc: {
      DocId: docId,
    },
    Kws: keywords.map((x) => ({ Kw: x, KwText: '' })),
  });

  const res = await fetch(`${baseUrl}/writeKw`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}


export async function getRelatedKwds(uid, docId) {
  const body = JSON.stringify({
    Uid: uid,
    Doc: {
      DocId: docId,
    }
  });

  const res = await fetch(`${baseUrl}/related`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}
