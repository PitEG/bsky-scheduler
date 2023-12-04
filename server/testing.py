import db
from datetime import datetime

conn = db.Connection('schedule.db')
conn.add_img('somewhere.txt', 'example_name.png', datetime.now())
print(conn.get_earliest_post())
conn.schedule_post(7,'hi there!',datetime.now())
print(conn.get_earliest_post())
print(conn.get_imgs())
