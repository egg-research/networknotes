from http import HTTPStatus
import json
import openai
from flask import Flask, request, Response
import os

from lib.gpt import set_openai_key, Example
from lib.gcp_cnl import top_keywords
from lib.gpt_related_keywords import getRelatedKeywordGPT

import requests, PyPDF2, io

CONFIG_VAR = "OPENAI_CONFIG"
KEY_NAME = "OPENAI_KEY"
HEADER_REQUIREMENT = "NN_HEADER_CONFIG"

"""Creates Flask app to serve the React app."""
app = Flask(__name__)


def list_files(startpath):
    for root, dirs, files in os.walk(startpath):
        level = root.replace(startpath, '').count(os.sep)
        indent = ' ' * 4 * (level)
        print('{}{}/'.format(indent, os.path.basename(root)))
        subindent = ' ' * 4 * (level + 1)
        for f in files:
            print('{}{}'.format(subindent, f))
#list_files('./')

header_req = HEADER_REQUIREMENT
header_dir = os.environ[HEADER_REQUIREMENT]
with open(header_dir) as f: 
    header_req = f.readlines()
header_req = str(header_req[0])[:-1]

gpt = getRelatedKeywordGPT()

app.config.from_envvar(CONFIG_VAR)
set_openai_key(app.config[KEY_NAME])

@app.route('/')
def hello():
    return "EGG Networked Notes: Language Model endpoints!"

@app.route('/gpt_examples', methods=['GET'])
def get_example(example_id=None):
    """Gets a single example or all the examples."""
    # This is all local, requires no external GPT3 calls
    # Return all examples
    if not example_id:
        return json.dumps(gpt.get_all_examples())

    example = gpt.get_example(example_id)
    if not example:
        return error("id not found", HTTPStatus.NOT_FOUND)
    return json.dumps(example.as_dict())

@app.route("/related", methods=["GET", "POST"])
def related():
    # This will call GPT3 service with gpt.submit_request(prompt)i
    if request.json is not None and 'header_req' in request.json.keys():
        if request.json['header_req'] == header_req:
            prompt = request.json["prompt"]
            response = 'a, nice, list, of, related, keywords' # gpt.submit_request(prompt) 
            offset = 0
            if not gpt.append_output_prefix_to_query:
                offset = len(gpt.output_prefix)
            return {'text': response} # response['choices'][0]['text'][offset:]} #
    return 'Unable to reach endpoint'

@app.route("/keywords", methods=["GET", "POST"])
def keywords():
    # This will call Google Cloud language service with top_keywords(prompt)
    if request.json is not None and 'header_req' in request.json.keys():
        if request.json['header_req'] == header_req:
            prompt = request.json["prompt"]
            response = str(['tomato', 'pear', 'cherry', 'ice cream', 'sweet', 'icing', 'sundae', 'sprinkles', 'brownie', 'matcha'])# str(top_keywords(prompt))
            return {'text': response}
    return 'Unable to reach endpoint'


@app.route("/pdf", methods=["GET", "POST"])
def pdf():
    # This will call Google Cloud language service with top_keywords(prompt)
    if request.json is not None and 'header_req' in request.json.keys():
        if request.json['header_req'] == header_req:
            prompt = request.json["url"]
            return {'text': parse_pdf(prompt)}
    return 'Unable to reach endpoint'

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


if __name__=='__main__':
    app.run()
