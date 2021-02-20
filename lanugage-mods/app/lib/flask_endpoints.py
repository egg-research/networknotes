from http import HTTPStatus
import json
#import subprocess
import openai
from flask import Flask, request, Response
from .gpt import set_openai_key, Example

CONFIG_VAR = "OPENAI_CONFIG"
KEY_NAME = "OPENAI_KEY"

def serve_app(gpt):
    """Creates Flask app to serve the React app."""
    app = Flask(__name__)

    app.config.from_envvar(CONFIG_VAR)
    set_openai_key(app.config[KEY_NAME])

    @app.route('/')
    def hello():
        return "NetworkedNotes Language Model Endpoints!"

    @app.route("/params", methods=["GET"])
    def get_params():
        # pylint: disable=unused-variable
        response = config.json()
        return response

    @app.route('/gpt_examples', methods=['GET'])
    def get_example(example_id=None):
        """Gets a single example or all the examples."""
        # Return all examples
        if not example_id:
            return json.dumps(gpt.get_all_examples())

        example = gpt.get_example(example_id)
        if not example:
            return error("id not found", HTTPStatus.NOT_FOUND)
        return json.dumps(example.as_dict())

    @app.route("/related", methods=["GET", "POST"])
    def related():
        # pylint: disable=unused-variable
        prompt = request.json["prompt"]
        response = gpt.submit_request(prompt)
        offset = 0
        if not gpt.append_output_prefix_to_query:
            offset = len(gpt.output_prefix)
        return {'text': response['choices'][0]['text'][offset:]}

    #subprocess.Popen(["yarn", "start"])
    app.run()
