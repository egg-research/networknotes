from lib.gpt_related_keywords import getRelatedKeywordGPT
from lib.flask_endpoints import serve_app

if __name__=='__main__':
    serve_app(getRelatedKeywordGPT())
