import db
from datetime import datetime

conn = db.Connection('schedule.db')
conn.add_img('filesize.png', 'filesize lmao', datetime.now())
print(conn.get_earliest_post())
conn.schedule_post(1,'filesize lmao','bsky',datetime.now())
print(conn.get_earliest_post())
print(conn.get_imgs())
