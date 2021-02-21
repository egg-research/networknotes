import requests, PyPDF2, io

# base_url = 'https://networknotes2.wl.r.appspot.com'
base_url = 'http://localhost:8080'

def parse_pdf(url):
    response = requests.get(url)
    my_raw_data = response.content

    doc_text = ""
    with io.BytesIO(response.content) as open_pdf_file:
        read_pdf = PyPDF2.PdfFileReader(open_pdf_file)
        num_pages = read_pdf.getNumPages()
        for p in range(num_pages):
            doc_text += read_pdf.getPage(p).extractText()
    return doc_text
    

def make_req(url, data):
    return requests.post(url, headers = {'content-type':'application/json'}, json=data)

def clear(authuid: str):
    url = base_url+'/clear'
    data = {"AuthUid": authuid}
    x = make_req(url, data)
    print(x.text)


def signup(authuid: str):
    url = base_url+'/signup'
    data = {"AuthUid": authuid}
    x = make_req(url, data)
    if not x.text.isnumeric():
        return login(authuid)
    return int(x.text)

def login(authuid: str):
    url = base_url+'/login'
    data = {"AuthUid": authuid}
    x = make_req(url, data)
    return int(x.text)


def writeDoc(uid: int, doc_name:str, doc_text:str):
    url = base_url+'/writeDoc'
    data = {"Uid": uid, "Doc": {"DocName":doc_name, "DocText":doc_text}}
    x = make_req(url, data)
    return int(x.text)

def writeDocWithId(uid: int, doc_id:int, doc_name:str, doc_text:str):
    url = base_url+'/writeDoc'
    data = {"Uid": uid, "Doc": {"DocId":doc_id, "DocName":doc_name, "DocText":doc_text}}
    x = make_req(url, data)
    return int(x.text)

if __name__ == "__main__":
    uid = signup('user0')

    doc_urls = [
        ('https://www.usenix.org/system/files/nsdi20-paper-liu-ming.pdf', 44)
    ]

    doc_ids = []
    for url, doc_id in doc_urls:
        doc_text = parse_pdf('https://www.usenix.org/system/files/nsdi20-paper-liu-ming.pdf')
        if doc_id >= 0:
            doc_id = writeDocWithId(uid, doc_id, url, doc_text)
        else:
            doc_id = writeDoc(uid, url, doc_text)
        doc_ids.append(doc_id)
    print(doc_ids)
