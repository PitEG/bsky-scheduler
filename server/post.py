import requests
import argparse
import json
from datetime import datetime
import pytz

parser = argparse.ArgumentParser(description='id and password')
parser.add_argument('-u',action='store',type=str,help='user id')
parser.add_argument('-p',action='store',type=str,help='user password')
parser.add_argument('-i',action='store',type=str,help='img address')
args = parser.parse_args()

print(args)
# exit()

# CREATE SESSION
BSKY_HANDLE = args.u
BSKY_APP_PASS = args.p
IMAGE_PATH = args.i

res = requests.post(
    'https://bsky.social/xrpc/com.atproto.server.createSession',
    json={'identifier':BSKY_HANDLE, 'password':BSKY_APP_PASS}
)
res.raise_for_status()
session = res.json()
print(session['accessJwt'])

# CREATE IMAGE BLOB
with open(IMAGE_PATH, 'rb') as f:
    img_bytes = f.read()

# if len(img_bytes) > 1000000:
#     raise Exception(
#         f"image file size too large. 1000000 bytes maximum, got: {len(img_bytes)}"
#     )

res = requests.post(
    "https://bsky.social/xrpc/com.atproto.repo.uploadBlob",
    headers={
        "Content-Type": 'image/jpg', # TODO this should be interpreted
        "Authorization": "Bearer " + session["accessJwt"],
    },
    data=img_bytes,
)
res.raise_for_status()
blob = res.json()['blob']
print(blob)

# CREATE POST
now = datetime.now(pytz.utc).isoformat().replace("+00:00", "Z")

post = {
    '$type': 'app.bsky.feed.post',
    'text': 'a sword!',
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

res = requests.post(
    'https://bsky.social/xrpc/com.atproto.repo.createRecord',
    headers={'Authorization':'Bearer ' + session['accessJwt']},
    json={
        'repo': session['did'],
        'collection': 'app.bsky.feed.post',
        'record': post,
    }
)

print(json.dumps(res.json(), indent = 2))
res.raise_for_status()
