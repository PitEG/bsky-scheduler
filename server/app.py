from flask import Flask, send_file, request
from flask_restful import Resource, Api, reqparse
from markupsafe import escape
import argparse
import werkzeug

from bsky import create_session, post_text, post_image

parser = argparse.ArgumentParser(
        prog='Scheduler Backend',
        description='what the title says',
        )
parser.add_argument('-imgs',default='.') # not a good default
parser.add_argument('-db',default='schedule.db')
args = parser.parse_args()

img_dir = args.imgs
db_file = args.db
print(f'image directory:{img_dir}')
print(f'db file:{db_file}')

app = Flask(__name__)
api = Api(app)

try:
    create_session('lmao', 'haha')
except Exception:
    print('!!!WARNING: couldn\'t make a session')

# REST api
class HelloWorld(Resource):
    def get(self):
        return {'hello':'world'}

class Schedule(Resource):
    def get(self):
        return {'post1':'the stuff'}


class Upload(Resource):
    def post(self):
        parse = reqparse.RequestParser()
        parse.add_argument('image',type=werkzeug.datastructures.FileStorage,location='files')
        args = parse.parse_args()
        image_file = args['image']
        image_file.save(img_dir + '/' + 'image.png')

class Images(Resource):
    def get(self):
        return {'images':'the stuff lmao'}

api.add_resource(HelloWorld, '/')
api.add_resource(Schedule, '/schedule')
api.add_resource(Upload, '/upload')
api.add_resource(Images, '/imgs')

# images
@app.route('/img/<name>')
def get_img(name):
    print(f'requsted image:{escape(name)}')
    return send_file(f'{escape(img_dir+name)}',mimetype='image/png')

if __name__ == '__main__':
    app.run(debug=True)
