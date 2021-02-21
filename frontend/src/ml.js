const baseUrl = 'https://egg-network-notes.wl.r.appspot.com';

// return as int
export async function getText(pdfURL) {
  const body = JSON.stringify({ header_req:'french-scrambled-eggs', url: pdfURL });
  const res = await fetch(`${baseUrl}/pdf`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });

  return res.json();
}