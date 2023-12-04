from flask import Flask, request
from flask_restful import Resource, Api
import argparse

from bsky import create_session, post_text, post_image

app = Flask(__name__)
api = Api(app)

img_dir = ''
db_file = ''

try:
    create_session('lmao', 'haha')
except Exception:
    print('!!!WARNING: couldn\'t make a session')

class HelloWorld(Resource):
    def get(self):
        return {'hello':'world'}

api.add_resource(HelloWorld, '/')

if __name__ == '__main__':
    app.run(debug=True)
