import requests
import argparse
import json
from datetime import datetime
import pytz

BSKY_CHAR_LIMIT = 300

class Session:
    def __init__(self, user : str, password : str):
        print('creating session...')
        res = requests.post(
            'https://bsky.social/xrpc/com.atproto.server.createSession',
            json={'identifier':user, 'password':password}
        )
        res.raise_for_status()

        self.session = res.json()
        pass

    def post_text(self, text : str):
        print('posting text...')
        now = datetime.now(pytz.utc).isoformat().replace("+00:00", "Z")
        post = {
            "$type": "app.bsky.feed.post",
            "text": text,
            "createdAt": now,
            'langs': ['en-US'],
        }
        res = requests.post(
               'https://bsky.social/xrpc/com.atproto.repo.createRecord',
               headers={'Authorization':'Bearer ' + self.session['accessJwt']},
               json={
                   'repo': self.session['did'],
                   'collection': 'app.bsky.feed.post',
                   'record': post,
               }
           )
        print('posted:', text)
        return
    
    
    def post_image(self, text : str, img_path : str):
        print('posting image...')
    
        # handle inputs
        if len(text) > BSKY_CHAR_LIMIT:
            raise Exception('Character count must not exceed 300')
        
        # read image
        with open(img_path, 'rb') as f:
            img_bytes = f.read()
    
    
        # if too big, just don't post it
        if len(img_bytes) > 1000000:
            raise Exception(
                f"image file size too large. 1000000 bytes maximum, got: {len(img_bytes)}"
            )
    
        # upload image blob
        res = requests.post(
            "https://bsky.social/xrpc/com.atproto.repo.uploadBlob",
            headers={
                "Content-Type": 'image/png', # TODO ideally, this should be interpreted from the image itself
                "Authorization": "Bearer " + self.session["accessJwt"],
            },
            data=img_bytes,
        )
        res.raise_for_status()
        blob = res.json()['blob']
        print(blob)
    
        # make post
        now = datetime.now(pytz.utc).isoformat().replace("+00:00", "Z")
        post = {
            '$type': 'app.bsky.feed.post',
            'text': text,
            'createdAt': now,
            'langs': ['en-US'],
            'embed': {
                '$type': 'app.bsky.embed.images',
                'images': [
                    {
                        "alt": "brief alt text description of the first image",
                        "image": {
                            "$type": "blob",
                            "ref": {
                                '$link': blob['ref']['$link'],
                            },
                            "mimeType": blob['mimeType'],
                            "size": blob['size'],
                        }
                    },
                ]
            }
        }
    
        # post the post
        res = requests.post(
            'https://bsky.social/xrpc/com.atproto.repo.createRecord',
            headers={'Authorization':'Bearer ' + self.session['accessJwt']},
            json={
                'repo': self.session['did'],
                'collection': 'app.bsky.feed.post',
                'record': post,
            }
        )
    
        print(json.dumps(res.json(), indent = 2))
        res.raise_for_status()
        return

# deprecated
def create_session(user, password):
    print('creating session...')
    res = requests.post(
        'https://bsky.social/xrpc/com.atproto.server.createSession',
        json={'identifier':user, 'password':password}
    )
    res.raise_for_status()
    session = res.json()
    print(session['accessJwt'])
    return session

# for testing
# parser = argparse.ArgumentParser(description='id and password')
# parser.add_argument('-u',action='store',type=str,help='user id')
# parser.add_argument('-p',action='store',type=str,help='user password')
# parser.add_argument('-i',action='store',type=str,help='img address')
# args = parser.parse_args()

# create session
# bsky_handle = args.u
# bsky_app_pass = args.p
# image_path = args.i

if __name__ == "__main__":
    print('hi')

    # parse stuff
    parser = argparse.ArgumentParser(description='id and password')
    parser.add_argument('-u',action='store',type=str,help='user id')
    parser.add_argument('-p',action='store',type=str,help='user password')
    parser.add_argument('-i',action='store',type=str,help='img address')
    args = parser.parse_args()
    
    # create session
    bsky_handle = args.u
    bsky_app_pass = args.p
    image_path = args.i
