import time
import bsky
import db
import argparse
import os

print('launching poster')

parser = argparse.ArgumentParser(
        prog='This thing posts the stuff'
        )
parser.add_argument('-db',default='schedule.db')
parser.add_argument('-imgs',default='.')
parser.add_argument('-u',help='should be the whole tag like user.bsky.social',metavar='username',required=True)
parser.add_argument('-p',metavar='password',required=True)
args = parser.parse_args()

db_file = args.db
img_dir = args.imgs
user = args.u
password = args.p

def check_post():
    conn = db.Connection(db_file)
    post = conn.get_earliest_post()
    if not post:
        return
    # print(post)
    img = conn.get_img(post['img_id'])
    if not img:
        return
    # print(img)
    session = bsky.Session(user, password)
    img_path = os.path.join(img_dir, img['path'])
    if not os.path.isfile(img_path):
        return
    # session.post_image(post['text'],img_path)
    print('successfully posted:',img_path,'\nfrom:',img['path'])
    conn.remove_post(post['id'])

while True:
   check_post()
   time.sleep(15)
